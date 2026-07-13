package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	paddle "github.com/PaddleHQ/paddle-go-sdk/v5"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/webhook"
	"vessl.dev/vessl/internal/cloud/repositories"
	"vessl.dev/vessl/internal/utils"
)

type BillingHandler struct {
	cloudRepo           repos.CloudRepo
	stripeWebhookSecret string
	paddleWebhookSecret string
}

func NewBillingHandler(cloudRepo repos.CloudRepo) *BillingHandler {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &BillingHandler{
		cloudRepo:           cloudRepo,
		stripeWebhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
		paddleWebhookSecret: os.Getenv("PADDLE_WEBHOOK_SECRET"),
	}
}

func planFromStripePriceID(priceID string) string {
	if priceID == os.Getenv("STRIPE_PRICE_PRO") {
		return "pro"
	}
	if priceID == os.Getenv("STRIPE_PRICE_TEAM") {
		return "team"
	}
	return "hobby"
}

func planFromPaddlePriceID(priceID string) string {
	if priceID == os.Getenv("PADDLE_PRICE_PRO") {
		return "pro"
	}
	if priceID == os.Getenv("PADDLE_PRICE_TEAM") {
		return "team"
	}
	return "hobby"
}

// @Summary Stripe Webhook
// @Description Receives billing events from Stripe to update subscription status
// @Tags Cloud-Billing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /billing/stripe/webhook [post]
func (h *BillingHandler) HandleStripeWebhook(c echo.Context) error {
	const maxBodyBytes = int64(65536)
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, maxBodyBytes)
	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return utils.Error(c, http.StatusServiceUnavailable, "Error reading request body")
	}

	sigHeader := c.Request().Header.Get("Stripe-Signature")
	event, err := h.parseStripeEvent(payload, sigHeader)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	log.Printf("Received Stripe webhook event: %s", event.Type)

	switch event.Type {
	case "checkout.session.completed":
		if err := h.handleStripeCheckoutCompleted(event); err != nil {
			return utils.Error(c, http.StatusBadRequest, err.Error())
		}
	case "customer.subscription.created", "customer.subscription.updated":
		if err := h.handleStripeSubscriptionUpdate(event); err != nil {
			return utils.Error(c, http.StatusBadRequest, err.Error())
		}
	case "customer.subscription.deleted":
		if err := h.handleStripeSubscriptionDeleted(event); err != nil {
			return utils.Error(c, http.StatusBadRequest, err.Error())
		}
	}

	return utils.Success(c, "received", nil)
}

func (h *BillingHandler) parseStripeEvent(payload []byte, sigHeader string) (stripe.Event, error) {
	var event stripe.Event
	var err error

	if h.stripeWebhookSecret != "" {
		event, err = webhook.ConstructEvent(payload, sigHeader, h.stripeWebhookSecret)
		if err != nil {
			log.Printf("Stripe signature verification failed: %v", err)
			return event, err
		}
	} else {
		if err := json.Unmarshal(payload, &event); err != nil {
			log.Printf("Failed to parse webhook body json: %v\n", err)
			return event, err
		}
	}
	return event, nil
}

func (h *BillingHandler) handleStripeCheckoutCompleted(event stripe.Event) error {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		return err
	}

	workspaceIDStr := session.ClientReferenceID
	if workspaceIDStr == "" || session.Customer == nil {
		log.Println("Checkout completed, but no workspaceID or customer found in session")
		return nil
	}

	workspaceID, err := strconv.ParseUint(workspaceIDStr, 10, 32)
	if err != nil {
		log.Printf("Invalid workspace ID format: %s\n", workspaceIDStr)
		return nil
	}

	team, err := h.cloudRepo.GetTeamByID(uint(workspaceID))
	if err != nil || team == nil {
		log.Printf("Checkout completed, but workspace %d not found\n", workspaceID)
		return nil
	}

	team.StripeCustomerID = session.Customer.ID
	h.cloudRepo.UpdateTeam(team)

	log.Printf("Checkout Completed | Workspace: %d | Customer: %s", workspaceID, session.Customer.ID)
	return nil
}

func (h *BillingHandler) handleStripeSubscriptionUpdate(event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		return err
	}

	team, err := h.cloudRepo.GetTeamByStripeCustomerID(sub.Customer.ID)
	if err == nil && team != nil && len(sub.Items.Data) > 0 {
		priceID := sub.Items.Data[0].Price.ID
		newPlan := planFromStripePriceID(priceID)
		if sub.Status == stripe.SubscriptionStatusActive || sub.Status == stripe.SubscriptionStatusTrialing {
			team.Plan = newPlan
			h.cloudRepo.UpdateTeam(team)
		}
	}

	log.Printf("Subscription Update | Customer: %s | Status: %s | Plan: %s",
		sub.Customer.ID, sub.Status, sub.Items.Data[0].Price.ID)
	return nil
}

func (h *BillingHandler) handleStripeSubscriptionDeleted(event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		return err
	}

	team, err := h.cloudRepo.GetTeamByStripeCustomerID(sub.Customer.ID)
	if err == nil && team != nil {
		team.Plan = "hobby"
		h.cloudRepo.UpdateTeam(team)
	}

	log.Printf("Subscription Deleted | Customer: %s", sub.Customer.ID)
	return nil
}

// @Summary Paddle Webhook
// @Description Receives billing events from Paddle
// @Tags Cloud-Billing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /billing/paddle/webhook [post]
func (h *BillingHandler) HandlePaddleWebhook(c echo.Context) error {
	const maxBodyBytes = int64(65536)
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, maxBodyBytes)

	if err := h.verifyPaddleSignature(c.Request()); err != nil {
		return utils.Error(c, http.StatusForbidden, "Invalid signature")
	}

	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return utils.Error(c, http.StatusServiceUnavailable, "Error reading request body")
	}

	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid payload")
	}

	eventType, _ := event["event_type"].(string)
	log.Printf("Received Paddle webhook event: %s", eventType)

	switch eventType {
	case "subscription.created", "subscription.updated":
		h.handlePaddleSubscriptionUpdate(event)
	case "subscription.canceled":
		h.handlePaddleSubscriptionCanceled(event)
	}

	return utils.Success(c, "received", nil)
}

func (h *BillingHandler) verifyPaddleSignature(r *http.Request) error {
	if h.paddleWebhookSecret == "" {
		return nil
	}
	verifier := paddle.NewWebhookVerifier(h.paddleWebhookSecret)
	ok, err := verifier.Verify(r)
	if err != nil {
		log.Printf("Paddle verification error: %v", err)
		return err
	}
	if !ok {
		return os.ErrPermission
	}
	return nil
}

func (h *BillingHandler) handlePaddleSubscriptionUpdate(event map[string]interface{}) {
	log.Println("Handling Paddle subscription update...")
	if data, ok := event["data"].(map[string]interface{}); ok {
		customerID, _ := data["customer_id"].(string)
		status, _ := data["status"].(string)
		if items, ok := data["items"].([]interface{}); ok && len(items) > 0 {
			if item, ok := items[0].(map[string]interface{}); ok {
				if price, ok := item["price"].(map[string]interface{}); ok {
					priceID, _ := price["id"].(string)
					team, err := h.cloudRepo.GetTeamByPaddleCustomerID(customerID)
					if err == nil && team != nil && (status == "active" || status == "trialing") {
						team.Plan = planFromPaddlePriceID(priceID)
						h.cloudRepo.UpdateTeam(team)
					}
				}
			}
		}
	}
}

func (h *BillingHandler) handlePaddleSubscriptionCanceled(event map[string]interface{}) {
	log.Println("Handling Paddle subscription canceled...")
	if data, ok := event["data"].(map[string]interface{}); ok {
		customerID, _ := data["customer_id"].(string)
		team, err := h.cloudRepo.GetTeamByPaddleCustomerID(customerID)
		if err == nil && team != nil {
			team.Plan = "hobby"
			h.cloudRepo.UpdateTeam(team)
		}
	}
}

type CheckoutRequest struct {
	PlanID      string `json:"plan_id"`
	ReturnURL   string `json:"return_url"`
	WorkspaceID string `json:"workspace_id"`
}

// @Summary Create Stripe Checkout
// @Description Generates a checkout URL for subscriptions
// @Tags Cloud-Billing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /billing/stripe/checkout [post]
func (h *BillingHandler) CreateStripeCheckout(c echo.Context) error {
	var req CheckoutRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid request")
	}

	priceID := stripePriceID(req.PlanID)
	if priceID == "" {
		return utils.Error(c, http.StatusBadRequest, "Invalid plan ID")
	}

	params := &stripe.CheckoutSessionParams{
		ClientReferenceID:  stripe.String(req.WorkspaceID),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(req.ReturnURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(req.ReturnURL + "?canceled=true"),
	}

	s, err := checkoutsession.New(params)
	if err != nil {
		log.Printf("Stripe checkout error: %v", err)
		return utils.Error(c, http.StatusInternalServerError, "Failed to create checkout session")
	}

	return utils.Success(c, "Session created", map[string]string{"url": s.URL})
}

// @Summary Create Paddle Checkout
// @Description Generates a checkout URL/Transaction for subscriptions
// @Tags Cloud-Billing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /billing/paddle/checkout [post]
func (h *BillingHandler) CreatePaddleCheckout(c echo.Context) error {
	var req CheckoutRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid request")
	}

	priceID := paddlePriceID(req.PlanID)
	if priceID == "" {
		return utils.Error(c, http.StatusBadRequest, "Invalid plan ID")
	}

	return utils.Success(c, "Ready for Paddle", map[string]string{
		"price_id": priceID,
		"status":   "ready_for_frontend_paddle_js",
	})
}

func stripePriceID(plan string) string {
	switch plan {
	case "hobby":
		return os.Getenv("STRIPE_PRICE_HOBBY")
	case "pro":
		return os.Getenv("STRIPE_PRICE_PRO")
	case "team":
		return os.Getenv("STRIPE_PRICE_TEAM")
	}
	return ""
}

func paddlePriceID(plan string) string {
	switch plan {
	case "hobby":
		return os.Getenv("PADDLE_PRICE_HOBBY")
	case "pro":
		return os.Getenv("PADDLE_PRICE_PRO")
	case "team":
		return os.Getenv("PADDLE_PRICE_TEAM")
	}
	return ""
}

// @Summary Create Stripe Customer Portal
// @Description Generates a portal URL for subscription management
// @Tags Cloud-Billing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /billing/stripe/portal [post]
func (h *BillingHandler) CreateStripePortal(c echo.Context) error {
	var req struct {
		WorkspaceID uint   `json:"team_id"`
		ReturnURL   string `json:"return_url"`
	}
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid request")
	}

	team, err := h.cloudRepo.GetTeamByID(req.WorkspaceID)
	if err != nil || team == nil || team.StripeCustomerID == "" {
		return utils.Error(c, http.StatusBadRequest, "Team has no Stripe customer")
	}

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(team.StripeCustomerID),
		ReturnURL: stripe.String(req.ReturnURL),
	}
	s, err := session.New(params)
	if err != nil {
		log.Printf("Stripe portal error: %v", err)
		return utils.Error(c, http.StatusInternalServerError, "Failed to create portal session")
	}

	return utils.Success(c, "Session created", map[string]string{"url": s.URL})
}

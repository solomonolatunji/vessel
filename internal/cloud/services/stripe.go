package services

type StripeService struct{}

func NewStripeService() *StripeService {
	return &StripeService{}
}

func (s *StripeService) CreateCheckoutSession(teamID string) (string, error) {
	return "chk_stub", nil
}

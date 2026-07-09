package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"vessel.dev/vessel/internal/api"
	"vessel.dev/vessel/internal/store"
	"vessel.dev/vessel/internal/types"
)

func TestTeamsAndSettingsEndpoints(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "vessel_teams_settings_test")
	_ = os.RemoveAll(tempDir)
	_ = os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "vessel.db")
	s, err := store.NewStore(dbPath)
	if err != nil {
		t.Fatalf("failed to init store: %v", err)
	}
	defer s.Close()

	srv := api.NewServer(s, nil, nil, nil)

	// 1. Register admin user (User A)
	regAdmin := []byte(`{"email":"owner@vessel.dev","password":"securepassword123","role":"admin"}`)
	reqRegA := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(regAdmin))
	reqRegA.Header.Set("Content-Type", "application/json")
	recRegA := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recRegA, reqRegA)

	var respA map[string]any
	_ = json.NewDecoder(recRegA.Body).Decode(&respA)
	tokenA, _ := respA["token"].(string)

	// 2. Register second user (User B)
	regMember := []byte(`{"email":"collaborator@vessel.dev","password":"securepassword123","role":"user"}`)
	reqRegB := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(regMember))
	reqRegB.Header.Set("Content-Type", "application/json")
	recRegB := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recRegB, reqRegB)

	var respB map[string]any
	_ = json.NewDecoder(recRegB.Body).Decode(&respB)
	tokenB, _ := respB["token"].(string)
	userBMap, _ := respB["user"].(map[string]any)
	userBID, _ := userBMap["id"].(string)

	// --- Teams & Organizations Tests ---

	// Create Team with User A
	teamPayload := []byte(`{"name":"Vessel Core Team"}`)
	reqCreateTeam := httptest.NewRequest(http.MethodPost, "/api/teams", bytes.NewReader(teamPayload))
	reqCreateTeam.Header.Set("Content-Type", "application/json")
	reqCreateTeam.Header.Set("Authorization", "Bearer "+tokenA)
	recCreateTeam := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recCreateTeam, reqCreateTeam)

	if recCreateTeam.Code != http.StatusCreated {
		t.Fatalf("expected status %d on create team, got %d. Body: %s", http.StatusCreated, recCreateTeam.Code, recCreateTeam.Body.String())
	}
	var createdTeam types.Team
	_ = json.NewDecoder(recCreateTeam.Body).Decode(&createdTeam)

	// List Teams for User A
	reqListTeams := httptest.NewRequest(http.MethodGet, "/api/teams", nil)
	reqListTeams.Header.Set("Authorization", "Bearer "+tokenA)
	recListTeams := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recListTeams, reqListTeams)
	if recListTeams.Code != http.StatusOK {
		t.Fatalf("expected status %d on list teams, got %d", http.StatusOK, recListTeams.Code)
	}

	// Invite an unregistered email to the Team
	invitePayload := []byte(`{"email":"external@vessel.dev","role":"Member"}`)
	reqInvite := httptest.NewRequest(http.MethodPost, "/api/teams/"+createdTeam.ID+"/invite", bytes.NewReader(invitePayload))
	reqInvite.Header.Set("Content-Type", "application/json")
	reqInvite.Header.Set("Authorization", "Bearer "+tokenA)
	recInvite := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recInvite, reqInvite)

	if recInvite.Code != http.StatusCreated {
		t.Fatalf("expected status %d on invite member, got %d. Body: %s", http.StatusCreated, recInvite.Code, recInvite.Body.String())
	}
	var inviteResp map[string]any
	_ = json.NewDecoder(recInvite.Body).Decode(&inviteResp)
	invMap, _ := inviteResp["invite"].(map[string]any)
	invToken, _ := invMap["token"].(string)

	// Inspect Invite via Token
	reqGetInv := httptest.NewRequest(http.MethodGet, "/api/team-invites/"+invToken, nil)
	recGetInv := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recGetInv, reqGetInv)
	if recGetInv.Code != http.StatusOK {
		t.Fatalf("expected status %d on get invite, got %d", http.StatusOK, recGetInv.Code)
	}

	// Invite existing User B directly by email
	inviteBPayload := []byte(`{"email":"collaborator@vessel.dev","role":"Admin"}`)
	reqInviteB := httptest.NewRequest(http.MethodPost, "/api/teams/"+createdTeam.ID+"/invite", bytes.NewReader(inviteBPayload))
	reqInviteB.Header.Set("Content-Type", "application/json")
	reqInviteB.Header.Set("Authorization", "Bearer "+tokenA)
	recInviteB := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recInviteB, reqInviteB)

	if recInviteB.Code != http.StatusCreated {
		t.Fatalf("expected status %d on invite existing user, got %d", http.StatusCreated, recInviteB.Code)
	}

	// Verify User B is now listed in team members
	reqMembers := httptest.NewRequest(http.MethodGet, "/api/teams/"+createdTeam.ID+"/members", nil)
	reqMembers.Header.Set("Authorization", "Bearer "+tokenA)
	recMembers := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recMembers, reqMembers)
	if recMembers.Code != http.StatusOK {
		t.Fatalf("expected status %d on list members, got %d", http.StatusOK, recMembers.Code)
	}

	// Remove User B from Team
	reqRemove := httptest.NewRequest(http.MethodDelete, "/api/teams/"+createdTeam.ID+"/members/"+userBID, nil)
	reqRemove.Header.Set("Authorization", "Bearer "+tokenA)
	recRemove := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recRemove, reqRemove)
	if recRemove.Code != http.StatusNoContent {
		t.Fatalf("expected status %d on remove member, got %d", http.StatusNoContent, recRemove.Code)
	}

	// --- Global Server Settings & System Prune Tests ---

	reqGetSettings := httptest.NewRequest(http.MethodGet, "/api/settings", nil)
	reqGetSettings.Header.Set("Authorization", "Bearer "+tokenA)
	recGetSettings := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recGetSettings, reqGetSettings)
	if recGetSettings.Code != http.StatusOK {
		t.Fatalf("expected status %d on get server settings, got %d", http.StatusOK, recGetSettings.Code)
	}

	updateSettingsPayload := []byte(`{"id":"global","caddyWildcardIp":"198.51.100.1","discordWebhookUrl":"https://discord.com/api/webhooks/test","notificationAlerts":true}`)
	reqUpdateSettings := httptest.NewRequest(http.MethodPut, "/api/settings", bytes.NewReader(updateSettingsPayload))
	reqUpdateSettings.Header.Set("Content-Type", "application/json")
	reqUpdateSettings.Header.Set("Authorization", "Bearer "+tokenA)
	recUpdateSettings := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recUpdateSettings, reqUpdateSettings)
	if recUpdateSettings.Code != http.StatusOK {
		t.Fatalf("expected status %d on update server settings, got %d. Body: %s", http.StatusOK, recUpdateSettings.Code, recUpdateSettings.Body.String())
	}

	// Trigger Docker System Prune
	reqPrune := httptest.NewRequest(http.MethodPost, "/api/settings/prune", nil)
	reqPrune.Header.Set("Authorization", "Bearer "+tokenA)
	recPrune := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recPrune, reqPrune)
	if recPrune.Code != http.StatusOK {
		t.Fatalf("expected status %d on system prune, got %d", http.StatusOK, recPrune.Code)
	}

	// --- Profile & Personal Access Tokens (PATs) Tests ---

	reqGetProfile := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	reqGetProfile.Header.Set("Authorization", "Bearer "+tokenB)
	recGetProfile := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recGetProfile, reqGetProfile)
	if recGetProfile.Code != http.StatusOK {
		t.Fatalf("expected status %d on get profile, got %d", http.StatusOK, recGetProfile.Code)
	}

	patPayload := []byte(`{"name":"Vessel CLI Dev"}`)
	reqCreatePAT := httptest.NewRequest(http.MethodPost, "/api/profile/tokens", bytes.NewReader(patPayload))
	reqCreatePAT.Header.Set("Content-Type", "application/json")
	reqCreatePAT.Header.Set("Authorization", "Bearer "+tokenB)
	recCreatePAT := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recCreatePAT, reqCreatePAT)
	if recCreatePAT.Code != http.StatusCreated {
		t.Fatalf("expected status %d on create PAT, got %d. Body: %s", http.StatusCreated, recCreatePAT.Code, recCreatePAT.Body.String())
	}

	var patResp map[string]any
	_ = json.NewDecoder(recCreatePAT.Body).Decode(&patResp)
	patTokenStr, _ := patResp["token"].(string)
	if patTokenStr == "" {
		t.Fatal("expected raw access token to be generated")
	}
	patObj, _ := patResp["pat"].(map[string]any)
	patID, _ := patObj["id"].(string)

	reqListPATs := httptest.NewRequest(http.MethodGet, "/api/profile/tokens", nil)
	reqListPATs.Header.Set("Authorization", "Bearer "+tokenB)
	recListPATs := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recListPATs, reqListPATs)
	if recListPATs.Code != http.StatusOK {
		t.Fatalf("expected status %d on list PATs, got %d", http.StatusOK, recListPATs.Code)
	}

	reqDelPAT := httptest.NewRequest(http.MethodDelete, "/api/profile/tokens/"+patID, nil)
	reqDelPAT.Header.Set("Authorization", "Bearer "+tokenB)
	recDelPAT := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recDelPAT, reqDelPAT)
	if recDelPAT.Code != http.StatusNoContent {
		t.Fatalf("expected status %d on delete PAT, got %d", http.StatusNoContent, recDelPAT.Code)
	}
}

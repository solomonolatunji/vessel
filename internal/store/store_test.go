package store_test

import (
	"os"
	"testing"

	"github.com/solomonolatunji/vessel/internal/store"
	"github.com/solomonolatunji/vessel/internal/types"
)

func TestStoreAndVault(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vessel-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	s, err := store.NewStore(tmpDir)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	defer s.Close()

	// 1. Test Project CRUD
	proj := &types.ProjectConfig{
		Name:          "Aeroplane-Inspired App",
		RepositoryURL: "https://github.com/solomonolatunji/sample-app",
		Branch:        "main",
		InternalPort:  3000,
		Domain:        "app.vessel.dev",
	}
	if err := s.CreateProject(proj); err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	fetchedProj, err := s.GetProject(proj.ID)
	if err != nil {
		t.Fatalf("GetProject failed: %v", err)
	}
	if fetchedProj.Name != proj.Name {
		t.Errorf("expected project name %s, got %s", proj.Name, fetchedProj.Name)
	}

	// 2. Test Domain CRUD
	domain := &types.DomainConfig{
		ProjectID:     proj.ID,
		DomainName:    "custom-app.vessel.dev",
		PathPrefix:    "/",
		SSLCertStatus: "pending",
	}
	if err := s.AddDomain(domain); err != nil {
		t.Fatalf("AddDomain failed: %v", err)
	}

	domains, err := s.ListDomains(proj.ID)
	if err != nil || len(domains) != 1 {
		t.Fatalf("expected 1 domain, got %d (err: %v)", len(domains), err)
	}

	// 3. Test AES Encrypted EnvVars
	if err := s.SetEnvVar(proj.ID, "DATABASE_URL", "postgres://user:secret@db:5432/app"); err != nil {
		t.Fatalf("SetEnvVar failed: %v", err)
	}
	envs, err := s.GetEnvVars(proj.ID)
	if err != nil {
		t.Fatalf("GetEnvVars failed: %v", err)
	}
	if envs["DATABASE_URL"] != "postgres://user:secret@db:5432/app" {
		t.Errorf("expected decrypted secret, got %s", envs["DATABASE_URL"])
	}

	// 4. Test User & Invite CRUD
	user := &types.User{
		Email:        "admin@vessel.dev",
		PasswordHash: "hashed-pass",
		Role:         "owner",
	}
	if err := s.CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	inv := &types.Invite{
		Email:     "colleague@vessel.dev",
		Role:      "developer",
		InvitedBy: user.ID,
	}
	if err := s.CreateInvite(inv); err != nil {
		t.Fatalf("CreateInvite failed: %v", err)
	}
}

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/repositories"
	"codedock.dev/codedock/internal/utils"
)

func runSetup() {
	fmt.Println("🛰️  Codedock Setup Wizard")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━")

	_, db, _ := initDataDir()

	repo := repositories.NewUserRepo(db)

	users, _, _ := repo.ListUsers(context.Background(), 1, 0)
	if len(users) > 0 {
		fmt.Println("✅ An admin account already exists. Use 'codedockd reset-password' if needed.")
		return
	}

	fmt.Println("\nLet's create your admin account.")
	email := prompt("Email: ")
	name := prompt("Name: ")
	password := promptPassword("Password: ")
	confirm := promptPassword("Confirm password: ")
	if password != confirm {
		exitError("Passwords do not match")
	}
	if len(password) < 8 {
		exitError("Password must be at least 8 characters")
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		exitError("Failed to hash password: %v", err)
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Email:        email,
		Name:         name,
		PasswordHash: hashed,
		Role:         "admin",
	}

	if err := repo.CreateUser(context.Background(), user); err != nil {
		exitError("Failed to create admin user: %v", err)
	}

	fmt.Printf("\n✅ Admin account created for %s (%s)\n", name, email)
	fmt.Println("You can now log in at the Codedock dashboard.")

	domain := promptOptional("Domain for apps (e.g. codedock.example.com, or press Enter to skip): ")
	tlsEmail := promptOptional("Let's Encrypt email (optional, for SSL): ")
	if domain != "" || tlsEmail != "" {
		envPath := filepath.Join(filepath.Dir(os.Getenv("CODEDOCK_DATA_DIR")), ".env")
		f, err := os.OpenFile(envPath, os.O_APPEND|os.O_WRONLY, 0o644)
		if err == nil {
			if domain != "" {
				fmt.Fprintf(f, "\nCODEDOCK_DOMAIN=%s\n", domain)
				fmt.Printf("✅ Domain set to %s. Restart Codedock to apply.\n", domain)
				fmt.Printf("   DNS: *.%s  A  <your-server-ip>\n", domain)
			}
			if tlsEmail != "" {
				fmt.Fprintf(f, "\nCODEDOCK_TLS_EMAIL=%s\n", tlsEmail)
				fmt.Println("✅ TLS email saved. Restart Codedock for SSL to take effect.")
			}
			f.Close()
		} else {
			fmt.Printf("⚠️  Could not update .env. Edit it manually: CODEDOCK_DOMAIN=%s CODEDOCK_TLS_EMAIL=%s\n", domain, tlsEmail)
		}
	}

	fmt.Println("\n📖 Next steps:")
	fmt.Println("   Dashboard: http://localhost:" + os.Getenv("PORT"))
	fmt.Println("   Docs: https://docs.codedock.dev")
}

func exitError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "❌ "+format+"\n", args...)
	os.Exit(1)
}

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
	"vessl.dev/vessl/internal/utils"
)

func runSetup() {
	fmt.Println("🛰️  Vessl Setup Wizard")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━")

	_, db, _ := initDataDir()

	repo := repositories.NewUserSQLiteRepository(db)

	users, _, _ := repo.ListUsers(context.Background(), 1, 0)
	if len(users) > 0 {
		fmt.Println("✅ An admin account already exists. Use 'vessld reset-password' if needed.")
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
	fmt.Println("You can now log in at the Vessl dashboard.")

	tlsEmail := promptOptional("Let's Encrypt email (optional, for SSL): ")
	if tlsEmail != "" {
		envPath := filepath.Join(filepath.Dir(os.Getenv("VESSL_DATA_DIR")), ".env")
		if f, err := os.OpenFile(envPath, os.O_APPEND|os.O_WRONLY, 0o644); err == nil {
			fmt.Fprintf(f, "\nVESSL_TLS_EMAIL=%s\n", tlsEmail)
			f.Close()
			fmt.Println("✅ TLS email saved. Restart Vessl for SSL to take effect.")
		} else {
			fmt.Printf("⚠️  Could not update .env. Set VESSL_TLS_EMAIL=%s manually and restart.\n", tlsEmail)
		}
	}

	fmt.Println("\n📖 Next steps:")
	fmt.Println("   Dashboard: http://localhost:" + os.Getenv("PORT"))
	fmt.Println("   Docs: https://docs.vessl.dev")
}

func exitError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "❌ "+format+"\n", args...)
	os.Exit(1)
}

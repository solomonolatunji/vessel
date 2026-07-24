package main

import (
	"context"
	"fmt"

	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/utils"
)

func runResetPassword() {
	_, db, _ := initDataDir()

	email := prompt("Admin email: ")
	repo := repositories.NewUserRepo(db)
	user, err := repo.GetUserByEmail(context.Background(), email)
	if err != nil {
		exitError("User with email %s not found", email)
	}

	password := promptOptional("New password: ")
	hashed, err := utils.HashPassword(password)
	if err != nil {
		exitError("Failed to hash password: %v", err)
	}

	user.PasswordHash = hashed
	if err := repo.UpdateUser(context.Background(), user); err != nil {
		exitError("Failed to update password: %v", err)
	}

	fmt.Printf("✅ Password reset for %s\n", email)
}

package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/occult/pagode/pkg/log"
	"github.com/occult/pagode/pkg/services"
)

// main creates a new admin user with the email passed in via the flag.
func main() {
	// Start a new container.
	c := services.NewContainer()
	defer func() {
		// Gracefully shutdown all services.
		if err := c.Shutdown(); err != nil {
			log.Default().Error("shutdown failed", "error", err)
		}
	}()

	var email string
	flag.StringVar(&email, "email", "", "email address for the admin user")
	flag.Parse()

	if len(email) == 0 {
		invalid("email is required")
	}

	// Generate a random placeholder password (Casdoor manages real passwords).
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		invalid("failed to generate a random password")
	}
	pw := hex.EncodeToString(b)

	// Create the admin user.
	err := c.ORM.User.
		Create().
		SetEmail(email).
		SetName("Admin").
		SetAdmin(true).
		SetVerified(true).
		SetPassword(pw).
		Exec(context.Background())

	if err != nil {
		invalid(err.Error())
	}

	fmt.Println("")
	fmt.Println("-- ADMIN USER CREATED --")
	fmt.Printf("Email: %s\n", email)
	fmt.Println("Password is managed by Casdoor SSO.")
	fmt.Println("----")
	fmt.Println("")
}

func invalid(msg string) {
	fmt.Printf("[ERROR] %s\n", msg)
	os.Exit(1)
}

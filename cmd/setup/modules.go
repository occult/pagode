package main

// Module describes an optional feature that can be removed from the project.
type Module struct {
	Name        string   // Display name in TUI (e.g., "Payment (Stripe)")
	FeatureName string   // Short name used in [feature:X] markers (e.g., "payment")
	Description string   // One-liner description
	Files       []string // Standalone files to delete
	Dirs        []string // Directories to delete recursively
	EntSchemas  []string // Ent schema files to remove (triggers ent-gen)
}

// Modules returns all optional feature modules.
func Modules() []Module {
	return []Module{
		paymentModule(),
		chatModule(),
		mailModule(),
		tasksModule(),
		filesModule(),
	}
}

func paymentModule() Module {
	return Module{
		Name:        "Payment (Stripe)",
		FeatureName: "payment",
		Description: "Subscription billing, one-time payments, premium access",
		Files: []string{
			"pkg/handlers/billing.go",
			"pkg/handlers/plans.go",
			"pkg/handlers/products.go",
			"pkg/handlers/premium.go",
			"pkg/services/payment.go",
			"pkg/services/payment_stripe.go",
			"resources/js/Pages/Billing.tsx",
			"resources/js/Pages/Plans.tsx",
			"resources/js/Pages/Products.tsx",
			"resources/js/Pages/Premium.tsx",
			"resources/js/components/PaymentForm.tsx",
			"resources/js/types/stripe.d.ts",
		},
		EntSchemas: []string{
			"ent/schema/paymentcustomer.go",
			"ent/schema/paymentmethod.go",
			"ent/schema/paymentintent.go",
			"ent/schema/subscription.go",
		},
	}
}

func chatModule() Module {
	return Module{
		Name:        "Chat (WebSocket)",
		FeatureName: "chat",
		Description: "Real-time community chat with rooms, voice messages",
		Files: []string{
			"pkg/handlers/chat.go",
			"resources/js/hooks/useChat.ts",
			"resources/js/types/chat.d.ts",
			"resources/js/hooks/useAudioRecorder.ts",
		},
		Dirs: []string{
			"static/chat-uploads",
			"resources/js/Pages/Chat",
			"resources/js/components/chat",
			"pkg/chat",
		},
		EntSchemas: []string{
			"ent/schema/chatroom.go",
			"ent/schema/chatmessage.go",
			"ent/schema/chatban.go",
		},
	}
}

func mailModule() Module {
	return Module{
		Name:        "Mail (Resend)",
		FeatureName: "mail",
		Description: "Transactional emails (verification, password reset, contact)",
		Files: []string{
			"pkg/handlers/contact.go",
			"pkg/services/mail.go",
			"pkg/services/mail_test.go",
			"pkg/ui/forms/contact.go",
			"pkg/ui/pages/contact.go",
			"pkg/ui/emails/auth.go",
		},
	}
}

func tasksModule() Module {
	return Module{
		Name:        "Background Tasks",
		FeatureName: "tasks",
		Description: "Async job queue with admin monitoring UI",
		Files: []string{
			"pkg/handlers/task.go",
			"pkg/ui/forms/task.go",
			"pkg/ui/pages/task.go",
		},
		Dirs: []string{"pkg/tasks"},
	}
}

func filesModule() Module {
	return Module{
		Name:        "File Upload",
		FeatureName: "files",
		Description: "File upload handling with filesystem abstraction",
		Files: []string{
			"pkg/handlers/files.go",
			"resources/js/Pages/UploadFile.tsx",
			"pkg/ui/forms/file.go",
			"pkg/ui/pages/file.go",
			"pkg/ui/models/file.go",
		},
	}
}

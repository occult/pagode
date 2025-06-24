package services

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/occult/pagode/config"
	"github.com/occult/pagode/pkg/log"
	"github.com/resend/resend-go/v2"
	"maragu.dev/gomponents"

	"github.com/labstack/echo/v4"
)

type (
	// MailClient provides a client for sending email
	// This is purposely not completed because there are many different methods and services
	// for sending email, many of which are very different. Choose what works best for you
	// and populate the methods below. For now, emails will just be logged.
	MailClient struct {
		// config stores application configuration.
		config *config.Config
		sender *resend.Client
	}

	// mail represents an email to be sent.
	mail struct {
		client    *MailClient
		from      string
		to        string
		subject   string
		body      string
		component gomponents.Node
	}
)

// NewMailClient creates a new MailClient.
func NewMailClient(cfg *config.Config) (*MailClient, error) {
	return &MailClient{
		config: cfg,
		sender: resend.NewClient(cfg.Mail.ResendApiKey),
	}, nil
}

// Compose creates a new email.
func (m *MailClient) Compose() *mail {
	return &mail{
		client: m,
		from:   m.config.Mail.FromAddress,
	}
}

// skipSend determines if mail sending should be skipped.
func (m *MailClient) skipSend() bool {
	return m.config.App.Environment != config.EnvProduction
}

// send attempts to send the email.
func (m *MailClient) send(email *mail, ctx echo.Context) error {
	switch {
	case email.to == "":
		return errors.New("email cannot be sent without a to address")
	case email.body == "" && email.component == nil:
		return errors.New("email cannot be sent without a body or component to render")
	}

	// Check if a component was supplied.
	if email.component != nil {
		// Render the component and use as the body.
		// TODO pool the buffers?
		buf := bytes.NewBuffer(nil)
		if err := email.component.Render(buf); err != nil {
			return err
		}

		email.body = buf.String()
	}

	// Check if mail sending should be skipped.
	if m.skipSend() {
		log.Ctx(ctx).Debug("skipping email delivery",
			"to", email.to,
		)
		return nil
	}

	// Build Resend request.
	params := &resend.SendEmailRequest{
		From:    email.from,
		To:      []string{email.to},
		Subject: email.subject,
		Html:    email.body,
	}

	// Fire away.
	if _, err := m.sender.Emails.Send(params); err != nil {
		return fmt.Errorf("resend: %w", err)
	}

	log.Ctx(ctx).Info("email sent", "to", email.to, "subject", email.subject)
	return nil
}

// From sets the email from address.
func (m *mail) From(from string) *mail {
	m.from = from
	return m
}

// To sets the email address this email will be sent to.
func (m *mail) To(to string) *mail {
	m.to = to
	return m
}

// Subject sets the subject line of the email.
func (m *mail) Subject(subject string) *mail {
	m.subject = subject
	return m
}

// Body sets the body of the email.
// This is not required and will be ignored if a component is set via Component().
func (m *mail) Body(body string) *mail {
	m.body = body
	return m
}

// Component sets a renderable component to use as the body of the email.
func (m *mail) Component(component gomponents.Node) *mail {
	m.component = component
	return m
}

// Send attempts to send the email.
func (m *mail) Send(ctx echo.Context) error {
	return m.client.send(m, ctx)
}

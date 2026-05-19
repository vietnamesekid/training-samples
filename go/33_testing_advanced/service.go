package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// EmailSender interface — easy to mock in tests
type EmailSender interface {
	Send(ctx context.Context, to, subject, body string) error
}

// UserStore interface
type UserStore interface {
	FindByID(ctx context.Context, id string) (*User, error)
	Save(ctx context.Context, user *User) error
}

type User struct {
	ID        string
	Name      string
	Email     string
	Active    bool
	CreatedAt time.Time
}

// NotificationService sends email when events occur
type NotificationService struct {
	store  UserStore
	mailer EmailSender
}

func NewNotificationService(store UserStore, mailer EmailSender) *NotificationService {
	return &NotificationService{store: store, mailer: mailer}
}

func (s *NotificationService) WelcomeUser(ctx context.Context, userID string) error {
	user, err := s.store.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}

	subject := fmt.Sprintf("Welcome, %s!", user.Name)
	body := fmt.Sprintf("Hi %s,\n\nWelcome to our platform!\n\nBest regards,\nTeam", user.Name)

	if err := s.mailer.Send(ctx, user.Email, subject, body); err != nil {
		return fmt.Errorf("send welcome email: %w", err)
	}

	user.Active = true
	return s.store.Save(ctx, user)
}

// Calculator — simple functions for testing demos
func Add(a, b int) int { return a + b }
func Sub(a, b int) int { return a - b }
func Mul(a, b int) int { return a * b }

func Div(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// Reverse reverses a string (UTF-8 safe)
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// FormatUser formats user information for display/snapshot testing
func FormatUser(u *User) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ID: %s\n", u.ID))
	sb.WriteString(fmt.Sprintf("Name: %s\n", u.Name))
	sb.WriteString(fmt.Sprintf("Email: %s\n", u.Email))
	sb.WriteString(fmt.Sprintf("Active: %v\n", u.Active))
	return sb.String()
}

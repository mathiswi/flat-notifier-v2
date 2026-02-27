package notify

import (
	"errors"
	"fmt"
	"log"
)

// Notifier can send a notification with a message and optional image URL.
type Notifier interface {
	Send(message string, imageURL string) error
}

// Dispatcher sends notifications to all configured channels.
type Dispatcher struct {
	notifiers []Notifier
}

// NewDispatcher creates a dispatcher with all configured notification channels.
func NewDispatcher() (*Dispatcher, error) {
	var notifiers []Notifier

	tg, err := NewTelegramNotifier()
	if err != nil {
		return nil, fmt.Errorf("telegram: %w", err)
	}
	if tg != nil {
		notifiers = append(notifiers, tg)
	}

	dc, err := NewDiscordNotifier()
	if err != nil {
		return nil, fmt.Errorf("discord: %w", err)
	}
	if dc != nil {
		notifiers = append(notifiers, dc)
	}

	if len(notifiers) == 0 {
		return nil, fmt.Errorf("no notification channels configured (set TELEGRAM_TOKEN or DISCORD_TOKEN)")
	}

	return &Dispatcher{notifiers: notifiers}, nil
}

// Send sends the message to all configured channels, collecting any errors.
func (disp *Dispatcher) Send(message string, imageURL string) error {
	var errs []error
	for _, notifier := range disp.notifiers {
		if err := notifier.Send(message, imageURL); err != nil {
			log.Printf("notification error: %v", err)
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

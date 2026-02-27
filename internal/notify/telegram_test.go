package notify

import "testing"

func TestNewTelegramNotifier_NoEnv(t *testing.T) {
	// With no env vars set, should return nil, nil
	t.Setenv("TELEGRAM_TOKEN", "")
	t.Setenv("TELEGRAM_CHAT_ID", "")

	n, err := NewTelegramNotifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != nil {
		t.Error("expected nil notifier when TELEGRAM_TOKEN is not set")
	}
}

func TestNewTelegramNotifier_MissingChatID(t *testing.T) {
	t.Setenv("TELEGRAM_TOKEN", "fake-token")
	t.Setenv("TELEGRAM_CHAT_ID", "")

	_, err := NewTelegramNotifier()
	if err == nil {
		t.Error("expected error when TELEGRAM_CHAT_ID is missing")
	}
}

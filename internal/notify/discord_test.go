package notify

import "testing"

func TestNewDiscordNotifier_NoEnv(t *testing.T) {
	t.Setenv("DISCORD_TOKEN", "")
	t.Setenv("DISCORD_USER_ID", "")

	n, err := NewDiscordNotifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != nil {
		t.Error("expected nil notifier when DISCORD_TOKEN is not set")
	}
}

func TestNewDiscordNotifier_MissingUserID(t *testing.T) {
	t.Setenv("DISCORD_TOKEN", "fake-token")
	t.Setenv("DISCORD_USER_ID", "")

	_, err := NewDiscordNotifier()
	if err == nil {
		t.Error("expected error when DISCORD_USER_ID is missing")
	}
}

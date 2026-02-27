package notify

import (
	"fmt"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const placeholderImage = "https://heuft.com/upload/image/400x267/no_image_placeholder.png"

// TelegramNotifier sends notifications via Telegram.
type TelegramNotifier struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

// NewTelegramNotifier creates a Telegram notifier from env vars.
// Returns nil if TELEGRAM_TOKEN is not set.
func NewTelegramNotifier() (*TelegramNotifier, error) {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		return nil, nil
	}

	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	if chatIDStr == "" {
		return nil, fmt.Errorf("TELEGRAM_TOKEN set but TELEGRAM_CHAT_ID missing")
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid TELEGRAM_CHAT_ID: %w", err)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("creating Telegram bot: %w", err)
	}

	return &TelegramNotifier{bot: bot, chatID: chatID}, nil
}

// Send sends a photo with caption to the configured Telegram chat.
func (t *TelegramNotifier) Send(message string, imageURL string) error {
	if imageURL == "" {
		imageURL = placeholderImage
	}

	photo := tgbotapi.NewPhoto(t.chatID, tgbotapi.FileURL(imageURL))
	photo.Caption = message
	photo.ParseMode = "Markdown"

	_, err := t.bot.Send(photo)
	if err != nil {
		return fmt.Errorf("sending Telegram photo: %w", err)
	}
	return nil
}

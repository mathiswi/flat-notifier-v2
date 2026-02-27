package notify

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

// DiscordNotifier sends notifications via Discord DM.
type DiscordNotifier struct {
	session *discordgo.Session
	userID  string
}

// NewDiscordNotifier creates a Discord notifier from env vars.
// Returns nil if DISCORD_TOKEN is not set.
func NewDiscordNotifier() (*DiscordNotifier, error) {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return nil, nil
	}

	userID := os.Getenv("DISCORD_USER_ID")
	if userID == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN set but DISCORD_USER_ID missing")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("creating Discord session: %w", err)
	}

	return &DiscordNotifier{session: session, userID: userID}, nil
}

// Send sends a message (and optional image) to the configured Discord user via DM.
func (d *DiscordNotifier) Send(message string, imageURL string) error {
	channel, err := d.session.UserChannelCreate(d.userID)
	if err != nil {
		return fmt.Errorf("creating DM channel: %w", err)
	}

	if imageURL != "" {
		_, err = d.session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
			Content: message,
			Embeds: []*discordgo.MessageEmbed{
				{
					Image: &discordgo.MessageEmbedImage{URL: imageURL},
				},
			},
		})
	} else {
		_, err = d.session.ChannelMessageSend(channel.ID, message)
	}

	if err != nil {
		return fmt.Errorf("sending Discord message: %w", err)
	}
	return nil
}

// Close closes the underlying Discord session.
func (d *DiscordNotifier) Close() error {
	return d.session.Close()
}

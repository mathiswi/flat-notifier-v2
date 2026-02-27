# flat-notifier-v2

Scrapes Kleinanzeigen apartment listings and sends notifications via Telegram/Discord when new flats appear. Runs as an AWS Lambda on a 5-minute schedule.

## Prerequisites

- Go 1.21+
- AWS CDK CLI (`npm install -g aws-cdk`)
- AWS credentials configured

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `OVERVIEW_URL` | Yes | Kleinanzeigen search URL |
| `TELEGRAM_TOKEN` | * | Telegram bot token |
| `TELEGRAM_CHAT_ID` | * | Telegram chat ID |
| `DISCORD_TOKEN` | * | Discord bot token |
| `DISCORD_USER_ID` | * | Discord user ID for DMs |

\* At least one notification channel (Telegram or Discord) must be configured.

## Commands

```sh
make test      # Run tests
make dev       # Run locally
make build     # Build Lambda binary
make deploy    # Build + CDK deploy
make clean     # Remove build artifacts
```

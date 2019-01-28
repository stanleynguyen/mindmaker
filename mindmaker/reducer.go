package mindmaker

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// HandleUpdates reduce updates to correct handler function
func HandleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		log.Printf("%+v\n", update.Message.Text)
	}
}

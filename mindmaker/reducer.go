package mindmaker

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stanleynguyen/mindmaker/persistence"
)

// Reducer reduce and handle updates accordingly
type Reducer struct {
	Bot         *tgbotapi.BotAPI
	Persistence *persistence.Persistence
}

// NewReducer get a new Reducer
func NewReducer(bot *tgbotapi.BotAPI, db *persistence.Persistence) *Reducer {
	return &Reducer{bot, db}
}

// HandleUpdates reduce updates to correct handler function
func (r *Reducer) HandleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		switch update.Message.Command() {
		case "create":

		}
	}
}

func getPrettyArgumentString(rawArgString string) string {
	return strings.Trim(rawArgString, " ")
}

func (r *Reducer) handleCreateCommand(update tgbotapi.Update) error {
	argStr := getPrettyArgumentString(update.Message.CommandArguments())
	bucketName := string(update.Message.Chat.ID) + " - " + argStr
	err := r.createABucket(bucketName)
	if err != nil {
		return err
	}
	err = r.setDefaultBucket(bucketName)
	if err != nil {
		return err
	}

	return nil
}

func (r *Reducer) createABucket(name string) error {
	return nil
}

func (r *Reducer) setDefaultBucket(name string) error {
	return nil
}

package mindmaker

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stanleynguyen/mindmaker/persistence"
)

// Reducer reduce and handle updates accordingly
type Reducer struct {
	Bot         *tgbotapi.BotAPI
	Persistence persistence.Persistence
}

// NewReducer get a new Reducer
func NewReducer(bot *tgbotapi.BotAPI, db persistence.Persistence) *Reducer {
	return &Reducer{bot, db}
}

// HandleUpdates reduce updates to correct handler function
func (r *Reducer) HandleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		switch update.Message.Command() {
		case "create":
			r.handleCreateCommand(update)
		}
	}
}

func getPrettyArgumentString(rawArgString string) string {
	return strings.Trim(rawArgString, " ")
}

func (r *Reducer) handleCreateCommand(update tgbotapi.Update) {
	argStr := getPrettyArgumentString(update.Message.CommandArguments())
	if argStr == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I'm confused 😵 My buckets can't have empty names\nPlease tell me in this format: /create <your bucket name>")
		r.Bot.Send(msg)
		return
	}

	bucketName := string(update.Message.Chat.ID) + " - " + argStr
	// TODO: check if bucket already exists
	err := r.Persistence.InsertBucket(bucketName)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}
	err = r.Persistence.UpdateDefaultBucket(update.Message.Chat.ID, bucketName)
	if err != nil {
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Bucket %v created for you, boss!", argStr))
	r.Bot.Send(msg)
}

func (r *Reducer) sendErrMessage(chatID int64) {
	errMsg := tgbotapi.NewMessage(chatID, "Sorry I'm not feeeling very well :( Please try again later")
	r.Bot.Send(errMsg)
}

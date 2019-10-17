package reducer

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (r *Reducer) handleCreateCommand(update tgbotapi.Update) {
	argStr := getPrettyArgumentString(update.Message.CommandArguments())
	chatID := update.Message.Chat.ID
	if argStr == "" {
		msg := tgbotapi.NewMessage(chatID, "I'm confused 😵 My buckets cant have empty names\nPlease tell me in this format: /create <your bucket name>")
		r.Bot.Send(msg)
		return
	}

	bucketDoesExist, err := r.Persistence.Exists(chatID, argStr)
	if err != nil {
		r.sendErrMessage(chatID)
		return
	} else if bucketDoesExist {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Sorry boss 🧐 bucket %v is already existing\nI'm setting it default bucket for you!", argStr))
		r.Bot.Send(msg)
		err := r.Persistence.UpdateDefaultBucket(chatID, argStr)
		if err != nil {
			r.sendErrMessage(chatID)
			return
		}
		msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("Updated bucket %v to be default!", argStr))
		r.Bot.Send(msg)
		return
	}
	err = r.Persistence.InsertBucket(chatID, argStr)
	if err != nil {
		r.sendErrMessage(chatID)
		return
	}
	err = r.Persistence.UpdateDefaultBucket(chatID, argStr)
	if err != nil {
		r.sendErrMessage(chatID)
		return
	}

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Bucket %v created for you, boss!", argStr))
	r.Bot.Send(msg)
}

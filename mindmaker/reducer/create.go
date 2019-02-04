package reducer

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (r *Reducer) handleCreateCommand(update tgbotapi.Update) {
	argStr := getPrettyArgumentString(update.Message.CommandArguments())
	if argStr == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I'm confused 😵 My buckets cant have empty names\nPlease tell me in this format: /create <your bucket name>")
		r.Bot.Send(msg)
		return
	}

	bucketName := getBucketNameFromChatID(update.Message.Chat.ID, argStr)
	bucketDoesExist, err := r.Persistence.Exists(bucketName)
	if err != nil {
		r.sendErrMessage(update.Message.Chat.ID)
		return
	} else if bucketDoesExist {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Sorry boss 🧐 bucket %v is already existing\nI'm setting it default bucket for you!", argStr))
		r.Bot.Send(msg)
		err := r.Persistence.UpdateDefaultBucket(update.Message.Chat.ID, getBucketNameFromChatID(update.Message.Chat.ID, argStr))
		if err != nil {
			r.sendErrMessage(update.Message.Chat.ID)
			return
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Updated bucket %v to be default!", argStr))
		r.Bot.Send(msg)
		return
	}
	err = r.Persistence.InsertBucket(bucketName)
	if err != nil {
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

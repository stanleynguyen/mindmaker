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
	// TODO: check if bucket already exists
	err := r.Persistence.InsertBucket(bucketName)
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

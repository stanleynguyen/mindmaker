package reducer

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (r *Reducer) handleSetCommand(update tgbotapi.Update) {
	argStr := getPrettyArgumentString(update.Message.CommandArguments())
	if argStr == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Boss 😱 Please gimme some name at least by telling me in this format: /set <bucket name>")
		r.Bot.Send(msg)
		return
	}

	bucketName := getBucketNameFromChatID(update.Message.Chat.ID, argStr)
	// TODO: check if bucket exists
	err := r.Persistence.UpdateDefaultBucket(update.Message.Chat.ID, bucketName)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("You have selected bucket %v to draw decision from 😉", argStr))
	r.Bot.Send(msg)
}

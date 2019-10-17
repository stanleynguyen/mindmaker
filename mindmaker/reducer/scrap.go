package reducer

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (r *Reducer) handleScrapCommand(update tgbotapi.Update) {
	argStr := getPrettyArgumentString(update.Message.CommandArguments())
	if argStr == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Boss 😱 Please gimme some name at least by telling me in this format: /scrap <bucket name>")
		r.Bot.Send(msg)
		return
	}

	bucketDoesExist, err := r.Persistence.Exists(update.Message.Chat.ID, argStr)
	if err != nil {
		r.sendErrMessage(update.Message.Chat.ID)
		return
	} else if !bucketDoesExist {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Oops 😶 I can't find bucket %v anywhere, boss\nAre you sure you you have created it with the create command before?", argStr))
		r.Bot.Send(msg)
		return
	}
	err = r.Persistence.DeleteBucket(update.Message.Chat.ID, argStr)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Bucket %v is out of the window 😶", argStr))
	r.Bot.Send(msg)
}

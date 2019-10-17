package reducer

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stanleynguyen/mindmaker/domain"
)

func (r *Reducer) handleAddCommand(update tgbotapi.Update) {
	defaultSet, err := r.Persistence.DefaultWasSet(update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}
	if !defaultSet {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "There's currently no bucket set to add this decision to 😕")
		r.Bot.Send(msg)
		return
	}
	bucketName, err := r.Persistence.GetDefaultBucket(update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}

	argStr := getPrettyArgumentString(update.Message.CommandArguments())
	if argStr == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Boss, your decision can't be an empty one 🤪\nPlease tell me in this format: /add <your decision>")
		r.Bot.Send(msg)
		return
	}
	err = r.Persistence.InsertOption(update.Message.Chat.ID, bucketName, domain.Option(argStr))
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Decision %v added to bucket %v", argStr, bucketName))
	r.Bot.Send(msg)
}

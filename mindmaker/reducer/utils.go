package reducer

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stanleynguyen/mindmaker/domain"
)

func getPrettyArgumentString(rawArgString string) string {
	return strings.Trim(rawArgString, " ")
}

func getFormattedListOfOptions(options []domain.Option) string {
	if len(options) == 0 {
		return "Opps! Your current bucket doesnt contain any options"
	}

	rtv := "Possible decisions in your current default bucket:\n"
	for i := 0; i < len(options); i++ {
		rtv += fmt.Sprintf("%v. %v", i+1, options[i])
		if i < len(options)-1 {
			rtv += "\n"
		}
	}

	return rtv
}

func (r *Reducer) sendErrMessage(chatID int64) {
	errMsg := tgbotapi.NewMessage(chatID, "Sorry I'm not feeeling very well 🤒 Please try again later")
	r.Bot.Send(errMsg)
}

package reducer

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const defaultResponse string = "Sorry boss 🙇 I didn't get what you just said\nBut I have a dad joke for you if you reply `yes` to this message 😎"

func (r *Reducer) handleAllOthersUpdate(update tgbotapi.Update) {
	isYesToDadJokeQn := update.Message.ReplyToMessage != nil &&
		update.Message.ReplyToMessage.Text == defaultResponse &&
		strings.EqualFold(update.Message.Text, "YES")
	if isYesToDadJokeQn {
		joke, err := getRandomDadJoke()
		var msg tgbotapi.MessageConfig
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "dad joke")
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, joke)
		}

		r.Bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, defaultResponse)
	r.Bot.Send(msg)
}

func getRandomDadJoke() (string, error) {
	dadJokeServiceErr := errors.New("Dont have one on top of my head right now 😶 Sorry boss")

	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com", nil)
	if err != nil {
		return "", dadJokeServiceErr
	}
	req.Header.Add("Accept", "text/plain")
	httpsCli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{},
		},
	}
	resp, err := httpsCli.Do(req)
	if err != nil {
		return "", dadJokeServiceErr
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", dadJokeServiceErr
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	body := string(bodyBytes)
	return body, nil
}

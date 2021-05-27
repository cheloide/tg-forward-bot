package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/hellowearemito/go-telegram-structs"
)

const uriBase = "https://api.telegram.org/bot%s"
const getUpdatesFormat = uriBase + "/getUpdates?offset=%d&timeout=%d&limit=1"
const forwardMessageFormat = uriBase + "/forwardMessage"

//Make persistent
var offset int64 = 0

func main() {

	var settings Settings = getSettings("settings.json")

	var botToken string = settings.Token
	var from int64 = settings.ForwardFrom
	var to int64 = settings.ForwardTo

	fmt.Printf("Forwarding from %d to %d\n", from, to)

	for {
		doStuff(botToken, from, to)
		time.Sleep(10 * time.Second)
	}
}

func getSettings(path string) Settings {

	f, err := ioutil.ReadFile(path)
	logFatalIfErr(err)

	var settings Settings
	err = json.Unmarshal([]byte(f), &settings)
	logFatalIfErr(err)
	return settings
}

func doStuff(botToken string, from int64, to int64) {
	var updates = getUpdates(botToken)

	for i := 0; i < len(updates); i++ {
		var update = updates[i]
		offset = update.UpdateID + 1
		if update.ChannelPost != nil && update.ChannelPost.Chat.ID == from {
			forward(botToken, from, to, update.ChannelPost.MessageID)
		}

	}
}

func getUpdates(botToken string) []telegram.Update {

	var timeout int = 60
	client := http.Client{
		Timeout: 60 * time.Second,
	}
	req, err := client.Get(fmt.Sprintf(getUpdatesFormat, botToken, offset, timeout))

	logFatalIfErr(err)

	responseBody, err := ioutil.ReadAll(req.Body)
	logFatalIfErr(err)

	var updates GetUpdatesResponse
	err = json.Unmarshal(responseBody, &updates)
	logFatalIfErr(err)
	return updates.Result
}

func logFatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func forward(botToken string, from int64, to int64, messageId int64) {

	var disableNotification bool = true

	var forwardMessage telegram.ForwardMessage = telegram.ForwardMessage{
		ChatID:              fmt.Sprint(to),
		FromChatID:          fmt.Sprint(from),
		DisableNotification: &disableNotification,
		MessageID:           messageId,
	}

	payload, err := json.Marshal(forwardMessage)
	logFatalIfErr(err)

	_, err = http.Post(fmt.Sprintf(forwardMessageFormat, botToken), "application/json", bytes.NewBuffer(payload))
	logFatalIfErr(err)
}

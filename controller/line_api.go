package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/shimokp/takizawa-garbage-bot/manager/config"
	"github.com/shimokp/takizawa-garbage-bot/manager/garbage"
	"github.com/shimokp/takizawa-garbage-bot/model"
)

func CallbackHandler(c *gin.Context) {
	var resp = ""

	bufBody := new(bytes.Buffer)
	bufBody.ReadFrom(c.Request.Body)
	defer c.Request.Body.Close()
	log.Println("RequestBuffer::", bufBody.String())

	var message = model.MessageText{}
	err := json.Unmarshal(bufBody.Bytes(), &message)
	if err != nil {
		log.Println(err)
		addString(&resp, err.Error())
	}

	//TODO: 一個でもエラーが出たら連鎖していろんなifに引っかかりそう
	if len(message.Events) > 0 {
		event := message.Events[0]

		err = sendToSlack(event)
		if err != nil {
			log.Println(err)
			addString(&resp, err.Error())
		}

		err = replyMessage(event)
		if err != nil {
			log.Println(err)
			addString(&resp, err.Error())
		}
	}

	log.Println("RESP::", resp)

	c.JSON(http.StatusOK, gin.H{
		"message": resp,
	})
}

func addString(base *string, text string) {
	*base = *base + "\n" + text
}

func replyMessage(event model.Event) error {
	sendMessage := switchMessage(event)

	bot, err := linebot.New(config.GetInstance().TGB_CHANNEL_SECRET, config.GetInstance().TGB_CHANNEL_ACCESS_TOKEN)
	if err != nil {
		return err
	}
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(sendMessage)).Do(); err != nil {
		return err
	}

	return nil
}

func switchMessage(event model.Event) string {
	switch event.Message.Text {
	case "今日":
		return garbage.GetMessage(model.Today, model.A)
	case "明日":
		return garbage.GetMessage(model.Tomorrow, model.A)
	}

	return ""
}

func sendToSlack(event model.Event) error {
	s := fmt.Sprintf(`
	{ 	"text" : " `+
		" ``` "+
		`replyToken: %s
type: %s
timeStamp: %s
---Source---
type: %s
userId: %s
iD: %s
---Message---
type: %s
text: %s`+
		" ``` "+` "}`, event.ReplyToken,
		event.Type,
		event.Timestamp,
		event.Source.Type,
		event.Source.UserID,
		event.Message.ID,
		event.Message.Type,
		event.Message.Text)

	body := strings.NewReader(s)

	url := os.Getenv("SlackToMe")
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	log.Println("Slack Responce Status::", resp.Status)

	bufBody := new(bytes.Buffer)
	_, err = bufBody.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	log.Println("Slack Responce::", bufBody.String())
	defer resp.Body.Close()
	return nil
}

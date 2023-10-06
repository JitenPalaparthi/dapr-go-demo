package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/golang/glog"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func Health(c *gin.Context) {
	c.String(200, "Ok")
}

func SubscribeTip(token, chatid string) func(*gin.Context) {
	return func(c *gin.Context) {
		result := make(map[string]any)
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			glog.Errorln(err)
			return
		}
		err = json.Unmarshal(data, &result)
		if err != nil {
			glog.Errorln(err)
			return
		}
		glog.Infoln(result["data"])
		// telegram code

		b, err := bot.New(token)
		if err != nil {
			glog.Errorln(err)
		}

		jsonString, err := json.Marshal(result["data"])

		if err != nil {
			glog.Errorln("XXXXXXXXXfdsfsdfsf___>>", err)
		}
		_, err = b.SendMessage(context.TODO(), &bot.SendMessageParams{
			ChatID: chatid,
			Text:   string(jsonString),
		})
		if err != nil {
			glog.Errorln("XXXXX___>>", err)
		}

		// end of telegram code
		obj, err := json.Marshal(data)
		if err != nil {
			glog.Errorln(err)
			return
		}

		_, err = c.Writer.Write(obj)
		if err != nil {
			glog.Errorln(err)
			return
		}
	}

}

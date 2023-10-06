package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/JitenPalaparthi/atipaday-telegram-subscriber/handlers"
	"github.com/JitenPalaparthi/dapr-go-http-wrapper/wrapper"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

var (
	appPort                   string
	telegram_Token            string
	telegram_ChatId           string
	err                       error
	daprObj                   *wrapper.Dapr
	DAPR_HOST, DAPR_HTTP_PORT string
	okHost, okPort            bool
)

const (
	DAPR_SECRET_STORE = "localsecretstore"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: example -stderrthreshold=[INFO|WARNING|ERROR|FATAL] -log_dir=[string]\n")
	flag.PrintDefaults()
	//os.Exit(2)
}

func init() {
	flag.Usage = usage
}

func init() {
	if DAPR_HOST, okHost = os.LookupEnv("DAPR_HOST"); !okHost {
		DAPR_HOST = "http://localhost"
	}
	if DAPR_HTTP_PORT, okPort = os.LookupEnv("DAPR_HTTP_PORT"); !okPort {
		DAPR_HTTP_PORT = "3500"
	}
	//daprObj = dapr.New(DAPR_HOST, DAPR_HTTP_PORT)
	daprObj = wrapper.New(DAPR_HOST, DAPR_HTTP_PORT)
}
func main() {
	appPort = os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "6002"
	}
	flag.Set("stderrthreshold", "INFO") // can set up the glog
	flag.Parse()
	defer glog.Flush()

	telegram_Token, err = daprObj.GetSecretFromFile(DAPR_SECRET_STORE, "telegram-token")
	if err != nil {
		glog.Fatalln("no telegram token;", err.Error())
	}

	telegram_ChatId, err = daprObj.GetSecretFromFile(DAPR_SECRET_STORE, "telegram-chatid")
	if err != nil {
		glog.Fatalln("no telegram chatid;", err.Error())
	}

	glog.Infoln("Telegram-Token:", telegram_Token)
	glog.Infoln("Telegram-ChatId:", telegram_ChatId)

	r := gin.Default()

	r.GET("/ping", handlers.Ping)
	r.GET("/health", handlers.Health)
	r.POST("/tips", handlers.SubscribeTip(telegram_Token, telegram_ChatId))

	glog.Infoln("Server started , up and run on port:" + appPort)
	// Start the server; this is a blocking call
	err := r.Run(":" + appPort) // http.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		glog.Fatalln(err)
	}
}

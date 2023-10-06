package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/JitenPalaparthi/atipaday/database"
	"github.com/JitenPalaparthi/atipaday/handlers"
	"github.com/JitenPalaparthi/dapr-go-http-wrapper/wrapper"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

var (
	appPort string
	DSN     string
	ctx     context.Context
	db      any
	err     error
)

var (
	DAPR_HOST, DAPR_HTTP_PORT string
	okHost, okPort            bool
	daprObj                   *wrapper.Dapr
)

const (
	DAPR_SECRET_STORE = "localsecretstore"
)

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

func usage() {
	fmt.Fprintf(os.Stderr, "usage: example -stderrthreshold=[INFO|WARNING|ERROR|FATAL] -log_dir=[string]\n")
	flag.PrintDefaults()
	//os.Exit(2)
}

func init() {
	flag.Usage = usage
	ctx = context.Background()
}

func main() {
	if appPort, okHost = os.LookupEnv("APP_PORT"); !okHost {
		appPort = "50090"
	}

	flag.Set("stderrthreshold", "INFO") // can set up the glog
	flag.Parse()
	defer glog.Flush()
	DSN, err := daprObj.GetSecretFromFile(DAPR_SECRET_STORE, "dsn")
	//DSN, err := apis.GetSecret(DAPR_HOST+":"+DAPR_HTTP_PORT+"/v1.0/secrets/"+DAPR_SECRET_STORE, "dsn")
	//glog.Infoln("dsn info for d", DSN)
	if err != nil || DSN == "" {
		glog.Info("Cannot connect to dapr secrets.Taking from local settings.")
		flag.StringVar(&DSN, "db", "host=127.0.0.1 user=admin password=admin123 dbname=atipadaydb port=55432 sslmode=disable TimeZone=Asia/Shanghai", "--db=host=localhost user=admin password=admin123 dbname=atipadaydb port=55432 sslmode=disable TimeZone=Asia/Shanghai")
		if os.Getenv("DB_CONN") != "" {
			DSN = os.Getenv("DB_CONN")
		}
	}

	ctx, cancel := context.WithCancel(ctx)

	glog.Info("Connecting to the database--")
	db, err = database.GetConnection(DSN)

	if err != nil {
		cancel()
		glog.Fatal(err)
	}

	r := gin.Default()

	r.GET("/ping", handlers.Ping)
	r.GET("/health", handlers.Health)

	glog.Info("creating new instance of TipDb object")
	cdb := new(database.Tip)
	cdb.DB = db
	cHandler := new(handlers.Tip)
	cHandler.ITip = cdb
	cHandler.Dapr = daprObj
	v1_tip := r.Group("v1/private/tip")
	v1_tip.POST("/", cHandler.Create(ctx))
	v1_tip.GET("/:id", cHandler.GetBy(ctx))
	v1_tip.GET("/all/:offset/:limit", cHandler.GetAllByOffset(ctx))
	v1_tip.GET("/all/:offset/:limit/:search", cHandler.Search(ctx))
	v1_tip.PUT("/:id", cHandler.UpdateBy(ctx))
	v1_tip.DELETE("/:id", cHandler.DeleteBy(ctx))

	r.Run(":" + appPort) // http.ListenAndServe()
}

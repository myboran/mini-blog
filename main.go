package main

import (
	"fmt"
	"gin-blog/models"
	"gin-blog/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"gin-blog/pkg/gredis"
	"gin-blog/pkg/setting"
	"gin-blog/routers"
)

func init() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()
	util.Setup()
}
func main() {
	gin.SetMode(setting.ServerSetting.RunMode)
	router := routers.InitRouter()

	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	s := &http.Server{
		Addr:           fmt.Sprintf(endPoint),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("[info] start http server listening %s", endPoint)
	if err := s.ListenAndServe(); err != nil {
		panic("启动失败")
	}
}

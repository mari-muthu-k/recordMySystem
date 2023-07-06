package main

import (
	"github.com/gin-gonic/gin"
	"github.com/recordMySystem/service/systeminfo"
)

func main(){
	r := gin.New()
	r.GET("/getSystemInfo",GetSystemInfo)
	r.POST("/startRecording",StartRecording)
	r.Run("localhost:8080")
}

func GetSystemInfo(ctx *gin.Context){
	sysInfo := systeminfo.GetSystemInfo()
	ctx.JSON(200,sysInfo)
}

func StartRecording(ctx *gin.Context){
	ctx.JSON(200,map[string]string{"message":"ok"})
}
package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/recordMySystem/handlers"
	"github.com/recordMySystem/server"
	"github.com/recordMySystem/service/database"
)

func appStartup(){
	fmt.Println("connecting influx client...")
	//Connect influx db
	database.ConnectDB(&server.InfluxClient)
	defer server.InfluxClient.Close()

	server.InfluxWriteAPI = server.InfluxClient.WriteAPIBlocking(server.ORGANIZATION,server.BUCKET)
	server.InfluxQueryAPI = server.InfluxClient.QueryAPI(server.ORGANIZATION)
	fmt.Println("influx client connected")
}

func main(){
	appStartup()
	r := gin.New()
	r.GET("/getCurrentSystemInfo",handlers.GetCurrentSystemInfo)
	r.GET("/getSystemInfo",handlers.GetSystemInfo)
	r.POST("/startRecording",handlers.StartRecording)
	r.Run("localhost:8080")
}
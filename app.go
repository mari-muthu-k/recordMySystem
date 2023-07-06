package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/recordMySystem/server"
	"github.com/recordMySystem/service/database"
	"github.com/recordMySystem/service/systeminfo"
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
	r.GET("/getSystemInfo",GetSystemInfo)
	r.POST("/startRecording",StartRecording)
	r.Run("localhost:8080")
}

func GetSystemInfo(ctx *gin.Context){
	sysInfo := systeminfo.GetSystemInfo()
	ctx.JSON(200,sysInfo)
}

func StartRecording(ctx *gin.Context){
	recordingChan := make(chan bool)
	//Recod the system
	go func(){
		var seconds int64 = 1
		for{
			fmt.Println("Recording started ",fmt.Sprint(seconds),"s ago")

			if seconds >= 60{
				break
			}
            
			//Get current system info
			sysData := systeminfo.GetSystemInfo()

			//Insert it into the bucket
			_,err := database.InsertData(&sysData); if err != nil {
				fmt.Println("unable to insert system data")
				panic(err)
			}

			seconds++
			time.Sleep(1*time.Second)
		}

		fmt.Println("recording finished")
		recordingChan <- true
	}()
	ctx.JSON(200,map[string]string{"message":"ok"})
}
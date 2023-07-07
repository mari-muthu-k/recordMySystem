package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/recordMySystem/service/database"
	"github.com/recordMySystem/service/systeminfo"
)

func GetCurrentSystemInfo(c *gin.Context){
	sysInfo := systeminfo.GetSystemInfo()
	c.JSON(200,sysInfo)
}

func GetSystemInfo(c *gin.Context){
	var id string
	var startTime ,endTime string
	var fields []string
	var sysInfo systeminfo.GetSystemInfoData

	queryParams := c.Request.URL.Query()
	for field,value := range queryParams {
		switch field {
		case "id":
			id = value[0]
		case "startTime":
			startTime = value[0]
		case "endTime":
			endTime = value[0]
		default:
			if value[0] == "true" {
				fields = append(fields, field)
			}
		}
	}

	if startTime == "" || endTime == "" {
		c.JSON(400,map[string]string{
			"message":"please select the start time and end time",
		})
		return
	}

	isDataExist,err := database.QueryData(fields,startTime,endTime,id,&sysInfo)
	if err != nil {
		log.Fatal(err)
	}

	if !isDataExist{
		c.JSON(404,map[string]string{
			"message":"no data found",
		})
		return
	}

	c.JSON(200,sysInfo)
}

func StartRecording(c *gin.Context){
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
	c.JSON(200,map[string]string{"message":"ok"})
}
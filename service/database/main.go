package database

import (
	"context"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/recordMySystem/server"
	"github.com/recordMySystem/service/systeminfo"
)

func ConnectDB(dbClient *influxdb2.Client){
	token := os.Getenv("INFLUXDB_TOKEN")
	url := "http://localhost:8086"
	*dbClient = 	influxdb2.NewClient(url, token)
}

func InsertData(sysData *systeminfo.SystemInfo)(bool,error){
	var tags = map[string]string{
		"id" : sysData.Id,
	}

	var fields = map[string]interface{}{
		"hostName":sysData.HostName,
		"batteryPercentage":sysData.BatteryPercentage,
		"memoryUsage":sysData.MemoryUsage,
		"temperature":sysData.Temperature,
		"cpuPercentage":sysData.CpuPercentage,
	}

	point := write.NewPoint(server.MEASUREMENT,tags,fields,time.Now())
	if err := server.InfluxWriteAPI.WritePoint(context.Background(),point); err != nil {
		return false,err
	}

	return true,nil
}
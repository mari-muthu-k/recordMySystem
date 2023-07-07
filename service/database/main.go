package database

import (
	"context"
	"fmt"
	"os"
	"strconv"
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

func QueryData(fields []string,startTime string,endTime string,id string,sysData *systeminfo.GetSystemInfoData)(bool,error){
	var isDataExist bool
	query := BuildQuery(fields,startTime,endTime,id)
	results, err := server.InfluxQueryAPI.Query(context.Background(), query)
	if err != nil {
		return isDataExist,err
	}

	for results.Next() {
		record := results.Record()
		val := record.Value()
        if !isDataExist || val != nil {
			isDataExist = true
		}

		switch record.Field() {
		case "cpuPercentage":
			sysData.CpuPercentage = append(sysData.CpuPercentage, val.(float64))
		case "temperature":
			sysData.Temperature = append(sysData.Temperature, val.(float64))
		case "memoryUsage":
			sysData.MemoryUsage = append(sysData.MemoryUsage, val.(float64))
		case "batteryPercentage":
			sysData.BatteryPercentage = append(sysData.BatteryPercentage, val.(float64))
		case "hostName":
			sysData.HostName = append(sysData.HostName, val.(string))
		}
	}
	
	if err := results.Err(); err != nil {
		return isDataExist,err
	}

	return isDataExist,nil
}

func BuildQuery(fields []string,startTime string,endTime string,id string)string{
	sTime,err := strconv.Atoi(startTime)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	eTime,err := strconv.Atoi(endTime)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Convert Unix timestamp to time.Time
	st := time.Unix(0,  int64(sTime)*int64(time.Millisecond))
	et := time.Unix(0,  int64(eTime)*int64(time.Millisecond))

	sISO := st.UTC().Format("2006-01-02T15:04:05.999Z")
	eISO := et.UTC().Format("2006-01-02T15:04:05.999Z")

	query := fmt.Sprintf(
		      `from(bucket: "%s")
			  |> range(start:%s,stop:%s)
			  |> filter(fn: (r) => r["_measurement"] == "%s")`,server.BUCKET,sISO,eISO,server.MEASUREMENT)

	for _,field := range fields {
		    query += fmt.Sprintf(`|> filter(fn: (r)=> r["_field"]=="%s")`,field)
	}

	query += fmt.Sprintf(`|> filter(fn:(r)=> r["id"]=="%s")
	                      |> keep(columns:["_field","_value"])`,id)
	return query
}
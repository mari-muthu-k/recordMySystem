package systeminfo

import (
	"context"
	"log"

	battery "github.com/distatus/battery"
	"github.com/google/uuid"
	cpu "github.com/shirou/gopsutil/v3/cpu"
	host "github.com/shirou/gopsutil/v3/host"
	mem "github.com/shirou/gopsutil/v3/mem"
)

func GetSystemInfo()SystemInfo{
	var sysInfo SystemInfo
	uuidObj :=  uuid.New()
	sysInfo.Id       = uuidObj.String()
	sysInfo.HostName = GetHostName()
	sysInfo.BatteryPercentage = GetBatteryPercentage()
	sysInfo.MemoryUsage = GetMemoryUsage()
	sysInfo.Temperature = GetAvgTemperature()
	sysInfo.CpuPercentage = GetCpuPercentage()
	return sysInfo
}

func GetHostName()string{
	hostInfo,err := host.InfoWithContext(context.Background())
	if err != nil {
		log.Fatalf("unable to read host name info : ",err)
	}
	return hostInfo.Hostname
}

func GetBatteryPercentage()float64{
	battery, err := battery.Get(0)
	if err != nil {
		log.Fatalf("unable to read battery info : ",err)
	}

	return ((battery.Current/battery.Full)*100)
}

func GetMemoryUsage()float64{
	gM,err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("unable to read memory info : ",err)
	}
	return gM.UsedPercent
}


func GetAvgTemperature()float64{
	var totalTemp float64
	tempArr,err := host.SensorsTemperatures()
	if err != nil {
		log.Fatalf("unable to read temperature info : ",err)
	}

	for _,curr := range tempArr {
		totalTemp += curr.Temperature
	}
	return totalTemp/float64(len(tempArr))

}

func GetCpuPercentage()float64{
	cpuPercentage,err := cpu.Percent(0,false) 
	if err != nil {
		log.Fatalf("unable to read cpu info : ",err)
	}
	return cpuPercentage[0]
}
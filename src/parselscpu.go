package main

import (
	"fmt"
	"strconv"
	"strings"
)

type CPUInfo struct {
	Architecture string `json:"Architecture"`
	CPUCount     string `json:"CPUCount"`
	ThreadCount  string `json:"ThreadCount"`
	CoreCount    string `json:"CoreCount"`
	SocketCount  string `json:"SocketCount"`
	ModelName    string `json:"ModelName"`
	Frequency    string `json:"Frequency"`
	MaxFrequency string `json:"MaxFrequency"`
}

func parseLsCpu(s string) {

	slcLsCpu := strings.Split(s, "\n")

	cpuMax := 0
	cpuSpeed := 0

	for _, line := range slcLsCpu {
		if strings.Contains(line, "Architecture") {
			fmt.Println("\t", line)
			cpuInfo.Architecture = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "CPU(s)") {
			fmt.Println("\t", line)
			cpuInfo.CPUCount = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "Thread(s)") {
			fmt.Println("\t", line)
			cpuInfo.ThreadCount = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "Core(s)") {
			fmt.Println("\t", line)
			cpuInfo.CoreCount = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "Socket(s)") {
			fmt.Println("\t", line)
			cpuInfo.SocketCount = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "Model name") {
			fmt.Println("\t", line)
			cpuInfo.ModelName = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Split(line, ":")[0] == "CPU MHz" {
			fmt.Println("\t", line)
			strCPUspeed := strings.TrimSpace(strings.Split(strings.Split(line, ":")[1], ".")[0])
			cpuSpeed, _ = strconv.Atoi(strCPUspeed)
			cpuInfo.Frequency = strCPUspeed
		}
		if strings.Split(line, ":")[0] == "CPU max MHz" {
			fmt.Println("\t", line)
			strCPUMaxSpeed := strings.TrimSpace(strings.Split(strings.Split(line, ":")[1], ".")[0])
			cpuMax, _ = strconv.Atoi(strCPUMaxSpeed)
			cpuInfo.MaxFrequency = strCPUMaxSpeed

			if cpuMax-cpuSpeed > 100 {
				sWarning := "Warning! The CPU runs on speeds below its capabilities. Please verify the settings and configuration as this may impact the NVMesh performance and user experience."
				fmt.Println(formatYellow("\t" + sWarning))
				mReport["CPU"] = sWarning
			}
		}

	}
}

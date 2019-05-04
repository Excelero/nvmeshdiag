package main

import (
	"fmt"
	"strings"
)

type OperatingSystem struct {
	LinuxDistribution string `json:"LinuxDistribution"`
	Kernel            string `json:"Kernel"`
}

func parseOSinfo(o string, k string) {

	kernel := k
	slcOS := strings.Split(o, "\n")

	for _, line := range slcOS {
		if strings.Contains(line, "Description:") {
			linuxDistribution := strings.TrimSpace(strings.Split(line, ":")[1])
			fmt.Println("\t", "Linux Distribution:", linuxDistribution)
			operatingSystem.LinuxDistribution = linuxDistribution
		}
	}
	if len(kernel) > 1 {
		fmt.Print("\t", " Kernel:", kernel)
		operatingSystem.Kernel = kernel
	}
}

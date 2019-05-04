package main

import (
	"fmt"
	"strings"
)

func parseLshw(s string) {

	slcLshw := strings.Split(strings.Split(s, "*")[0], "\n")

	serverName := slcLshw[0]
	fmt.Println(formatBoldWhite("Server/Hostname:"), serverName)
	nvmeshDiag.ServerName = serverName

	for _, line := range slcLshw {
		if strings.Contains(line, "product") {
			platformName := strings.TrimSpace(strings.Split(line, ":")[1])
			fmt.Println(formatBoldWhite("Platform:"), platformName)
			nvmeshDiag.Platform = platformName
		}
		if strings.Contains(line, "vendor") {
			vendorName := strings.TrimSpace(strings.Split(line, ":")[1])
			fmt.Println(formatBoldWhite("Manufacturer:"), vendorName)
			nvmeshDiag.Manufacturer = vendorName
		}
		if strings.Contains(line, "serial") {
			serialNumber := strings.TrimSpace(strings.Split(line, ":")[1])
			fmt.Println(formatBoldWhite("Serial #:"), serialNumber)
			nvmeshDiag.SerialNumber = serialNumber
		}
	}
}

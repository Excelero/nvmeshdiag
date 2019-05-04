package main

import (
	"fmt"
	"strings"
)

func parseDmiDecode(s string) {

	slcDmiInformation := strings.Split(s, "\n\n")

	for _, line := range strings.Split(slcDmiInformation[0], "\n") {
		if strings.Contains(line, "Product Name") {
			baseBoardType := strings.TrimSpace(strings.Split(line, ":")[1])
			fmt.Println(formatBoldWhite("Baseboard type:"), baseBoardType)
			nvmeshDiag.BaseboardType = baseBoardType
		}
		if strings.Contains(line, "Version") {
			baseBoardVersion := strings.TrimSpace(strings.Split(line, ":")[1])
			fmt.Println(formatBoldWhite("Baseboard version:"), baseBoardVersion)
			nvmeshDiag.BaseboardVersion = baseBoardVersion
		}
		if strings.Contains(line, "Serial Number") {
			baseBoardSerial := strings.TrimSpace(strings.Split(line, ":")[1])
			fmt.Println(formatBoldWhite("Baseboard serial #:"), baseBoardSerial)
			nvmeshDiag.BaseboardSerial = baseBoardSerial
		}
	}
}

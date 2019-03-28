package main

import (
	"fmt"
	"strings"
)

func parseOSinfo (o string, k string) {

	strKernel := k
	slcOS := strings.Split(o,"\n")
	for _, line := range slcOS{
		if strings.Contains(line, "Description:"){
			fmt.Println("\t", "Linux Distribution:", strings.TrimSpace(strings.Split(line, ":")[1]))
		}
	}
	if len(strKernel) > 1{
		fmt.Print("\t", " Kernel:", strKernel)
	}
	}

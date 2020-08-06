package main

import (
	"fmt"
	"strings"
)

func checkSystemTuning(){
	fmt.Println(formatBoldWhite("\nServer/Host System Tuning:"))
	checkTuneD()
	checkIRQBalance()
}

func checkTuneD() {
	var sWarning string
	tunedReturnCode := getCommandReturnCode(strings.Fields("systemctl status tuned"))
	if tunedReturnCode == 4 {
		sWarning = formatYellow("Not installed! To achieve best performance and user experience you should run tuned in the 'latency-performance' profile.")
	} else if tunedReturnCode != 0 {
		sWarning = formatYellow("Service failed! Please check the tuned service.")
	}
	if tunedReturnCode == 0 {
		tunedProfile, _ := runCommand(strings.Fields("tuned-adm active"))
		{
			if strings.Contains(tunedProfile, "latency-performance") {
				fmt.Println("\tTuneD service:", formatGreen("OK"))
				return
			} else {
				sWarning = formatYellow("Service is running but to achieve best user experience and performance, the 'latency-performance' is recommended.")
			}
		}
	}
	fmt.Println("\tTuneD service:", sWarning)
	mReport["TuneD"] = sWarning
	return
}

func checkIRQBalance(){
	var sWarning string
	irqBalance := getCommandReturnCode(strings.Fields("systemctl status irqbalance"))
	if irqBalance == 4 {
		sWarning = formatYellow("Not installed! To achieve best performance and user experience you should run IRQ Balance.")
	} else if irqBalance != 0 {
		sWarning = formatYellow("Service failed! Please check the IRQ Balance service.")
	}
	if irqBalance == 0{
		fmt.Println("\tIRQBalance:", formatGreen("OK"))
		return
	}
	fmt.Println("\tIRQ Balance:", sWarning)
	mReport["IRQ Balance"] = sWarning
	return

}
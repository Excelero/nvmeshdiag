package main

import (
	"fmt"
	"strings"
)

type SystemTuning struct {
	TuneDService string `json:"TuneDService"`
	IRQBalance   string `json:"IRQBalance"`
}

func checkSystemTuning() {
	fmt.Println(formatBoldWhite("\nServer/Host System Tuning:"))
	checkTuneD()
	checkIRQBalance()
}

func checkTuneD() {
	tunedReturnCode := getCommandReturnCode(strings.Fields("systemctl status tuned"))
	if tunedReturnCode == 4 {
		sWarning := "Not installed! To achieve best performance and user experience you should run tuned in the 'latency-performance' profile."
		fmt.Println("\tTuneD service:", formatYellow(sWarning))
		mReport["TuneD"] = sWarning
		systemTuning.TuneDService = sWarning
		return
	}
	if tunedReturnCode != 0 {
		sWarning := "Service failed! Please check the tuned service."
		fmt.Println("\tTuneD service:", formatYellow(sWarning))
		mReport["TuneD"] = sWarning
		systemTuning.TuneDService = sWarning
		return
	}
	if tunedReturnCode == 0 {
		tunedProfile, _ := runCommand(strings.Fields("tuned-adm active"))
		{
			if strings.Contains(tunedProfile, "latency-performance") {
				fmt.Println("\tTuneD service:", formatGreen("OK"))
				systemTuning.TuneDService = "OK"
				return
			} else {
				sWarning := " Service is running but to achieve best user experience and performance, the 'latency-performance' is recommended."
				fmt.Println("\tTuneD service:", formatYellow(sWarning))
				mReport["TuneD"] = sWarning
				systemTuning.TuneDService = sWarning
				return
			}
		}
	}
}

func checkIRQBalance() {
	irqBalance := getCommandReturnCode(strings.Fields("systemctl status irqbalance"))
	if irqBalance == 4 {
		sWarning := "Not installed! To achieve best performance and user experience you should run IRQ Balance."
		fmt.Println("\tIRQ Balance:", formatYellow(sWarning))
		mReport["IRQ Balance"] = sWarning
		systemTuning.IRQBalance = sWarning
		return
	}
	if irqBalance != 0 {
		sWarning := "Service failed! Please check the IRQ Balance service."
		fmt.Println("\tIRQ Balance:", formatYellow(sWarning))
		mReport["IRQ Balance"] = sWarning
		systemTuning.IRQBalance = sWarning
		return
	}
	if irqBalance == 0 {
		fmt.Println("\tIRQBalance:", formatGreen("OK"))
		systemTuning.IRQBalance = "OK"
		return
	}
}

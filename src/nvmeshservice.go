package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

func checkNvmeshInstalled() []string {

	var slcNvmeshService []string
	var strPackageQuery string
	r := regexp.MustCompile(`(?i)nvmesh\-([a-z]*)`)

	if checkExecutableExists("apt") {
		strPackageQuery, _ = runCommand(strings.Fields("apt list --installed"))
	} else if checkExecutableExists("rpm") {
		strPackageQuery, _ = runCommand(strings.Fields("rpm -qa"))
	}

	for _, match := range r.FindAllString(strPackageQuery, -1) {
		slcNvmeshService = append(slcNvmeshService, match)
	}
	return slcNvmeshService
}

func readTomaLeader() {
	dat, err := ioutil.ReadFile("/var/log/NVMesh/toma_leader_name")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("\tNVMesh Toma Leader:", strings.TrimRight(string(dat), "\n\r"))
}

func checkNvmeshStatus() {
	slcNvmeshService := checkNvmeshInstalled()

	if len(slcNvmeshService) > 0 {
		var status string
		fmt.Println(formatBoldWhite("\nNVMesh Service Information:"))
		for _, line := range slcNvmeshService {
			if strings.Contains(line, "core") {
				status = statusFormat(getCommandReturnCode(strings.Fields("systemctl status nvmeshtarget")))
				fmt.Println("\tNVMesh Target Service:", status)
				mReport["NVMesh Target Service: "] = status
				status = statusFormat(getCommandReturnCode(strings.Fields("systemctl status nvmeshclient")))
				fmt.Println("\tNVMesh Client Service:", status)
				mReport["NVMesh Client Service: "] = status
			}
			if strings.Contains(line, "management") {
				status = statusFormat(getCommandReturnCode(strings.Fields("systemctl status nvmeshmgr")))
				fmt.Println("\tNVMesh Management Service:", status)
				mReport["NVMesh Management Service: "] = status
			}
			if strings.Contains(line, "target") {
				status = statusFormat(getCommandReturnCode(strings.Fields("systemctl status nvmeshtarget")))
				fmt.Println("\tNVMesh Target Service:", status)
				mReport["NVMesh Target Service: "] = status
			}
			if strings.Contains(line, "client") {
				status = statusFormat(getCommandReturnCode(strings.Fields("systemctl status nvmeshtarget")))
				fmt.Println("\tNVMesh Client Service:", status)
				mReport["NVMesh Client Service: "] = status
			}
		}
		readTomaLeader()

	} else {
		fmt.Println(formatBoldWhite("\nNVMesh Service Information:"), "No NVMesh services found.")
	}
}

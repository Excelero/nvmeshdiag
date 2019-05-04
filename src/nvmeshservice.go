package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

type NVMeshInfo struct {
	TargetPackage           string `json:"TargetPackage"`
	TargetServiceStatus     string `json:"TargetServiceStatus"`
	ClientPackage           string `json:"ClientPackage"`
	ClientServiceStatus     string `json:"ClientServiceStatus"`
	ManagementPackage       string `json:"ManagementPackage"`
	ManagementServiceStatus string `json:"ManagementServiceStatus"`
	CorePackage             string `json:"CorePackage"`
	TomaLeader              string `json:"TomaLeader"`
}

func checkNvmeshInstalled() []string {

	var slcNvmeshService []string
	var strPackageQuery string
	r := regexp.MustCompile(`(?i)(nvmesh-[a-z]*-\d*\.\d*\.\d*-\d*)`)

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
	var tomaLeader string
	tomaLeader = strings.TrimRight(string(dat), "\n\r")
	fmt.Println("\tNVMesh Toma Leader:", tomaLeader)
	nvmeshInfo.TomaLeader = tomaLeader
}

func checkNvmeshStatus() {
	slcNvmeshService := checkNvmeshInstalled()

	if len(slcNvmeshService) > 0 {
		var status string
		fmt.Println(formatBoldWhite("\nNVMesh Service Information:"))
		for _, line := range slcNvmeshService {
			if strings.Contains(line, "core") {
				fmt.Println("\tNVMesh Core Package Information:", line)
				mReport["NVMesh Core Package: "] = line
				nvmeshInfo.CorePackage = line
				status = statusFormat(getCommandReturnCode(strings.Fields("systemctl status nvmeshtarget")))
				fmt.Println("\tNVMesh Target Service Status:", status)
				mReport["NVMesh Target Service: "] = status
				nvmeshInfo.TargetServiceStatus = status
				status = statusFormat(getCommandReturnCode(strings.Fields("systemctl status nvmeshclient")))
				fmt.Println("\tNVMesh Client Service Status:", status)
				mReport["NVMesh Client Service: "] = status
				nvmeshInfo.ClientServiceStatus = status
			}
			if strings.Contains(line, "management") {
				fmt.Println("\tNVMesh Management Package Information:", line)
				mReport["NVMesh Management Package: "] = line
				nvmeshInfo.ManagementPackage = line
				status = statusFormat(getCommandReturnCode(strings.Fields("systemctl status nvmeshmgr")))
				fmt.Println("\tNVMesh Management Service Status:", status)
				mReport["NVMesh Management Service: "] = status
				nvmeshInfo.ManagementServiceStatus = status
			}
			if strings.Contains(line, "target") {
				fmt.Println("\tNVMesh Target Package Information:", line)
				mReport["NVMesh Target Package: "] = line
				nvmeshInfo.TargetPackage = line
				status = statusFormat(getCommandReturnCode(strings.Fields("systemctl status nvmeshtarget")))
				fmt.Println("\tNVMesh Target Service Status:", status)
				mReport["NVMesh Target Service: "] = status
				nvmeshInfo.TargetServiceStatus = status
			}
			if strings.Contains(line, "client") {
				fmt.Println("\tNVMesh Client Package Information:", line)
				mReport["NVMesh Client Package: "] = line
				nvmeshInfo.ClientPackage = line
				status = statusFormat(getCommandReturnCode(strings.Fields("systemctl status nvmeshtarget")))
				fmt.Println("\tNVMesh Client Service Status:", status)
				mReport["NVMesh Client Service: "] = status
				nvmeshInfo.ClientServiceStatus = status
			}
		}
		readTomaLeader()

	} else {
		fmt.Println(formatBoldWhite("\nNVMesh Service Information:"), "No NVMesh services found.")
	}
}

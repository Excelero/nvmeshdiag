package main

import (
	"flag"
	"fmt"
	"strings"
)

//Define global variables to be used throughout the project
var (
	mReport                                                        = make(map[string]string)
	plainOutput                                                    = false
	jsonOutputFormat                                               = false
	nvmeshDiag, operatingSystem, nvmeshInfo, systemTuning, cpuInfo = NvmeshDiag{}, OperatingSystem{}, NVMeshInfo{}, SystemTuning{}, CPUInfo{}
)

//Define the struct which will be updated and populated by various routines.
//The content of this struct eventually will be exported into the JSON formatted string.
type NvmeshDiag struct {
	ServerName       string          `json:"ServerName"`
	Platform         string          `json:"Platform"`
	Manufacturer     string          `json:"Manufacturer"`
	SerialNumber     string          `json:"SerialNumber"`
	BaseboardType    string          `json:"BaseboardType"`
	BaseboardVersion string          `json:"BaseboardVersion"`
	BaseboardSerial  string          `json:"BaseboardSerial"`
	InstalledMemory  string          `json:"InstalledMemory"`
	OperatingSystem  OperatingSystem `json:"OperatingSystem"`
	NVMeshInfo       NVMeshInfo      `json:"NVMeshInfo"`
	SystemTuning     SystemTuning    `json:"SystemTuning"`
	CPUInfo          CPUInfo         `json:"CPUInfo"`
	OFEDInfo         string          `json:"OFEDInfo"`
	FirewallInfo     string          `json:"FirewallInfo"`
}

func init() {

}

func main() {
	//Parse the commandline arguments
	//valid arguments are '-p' for plain and unformatted text output and '-j' to save the information locally in the users home folder in JSON format
	flag.BoolVar(&plainOutput, "p", false, "No color and no formatted console output in plain text.")
	flag.BoolVar(&jsonOutputFormat, "j", false, "Print/output the result and information in JSON format.")
	flag.Parse()

	//Check the user level. If other than root or in sudo mode, some nvmeshdiag will not be able to pull all the information.
	//If all details are desired/required, run nvmesh diage using sudo or root
	checkUser()

	fmt.Println(formatBoldWhite("Starting System Scan. Please wait..."))

	//Running/executing the 'lshw' linus OS comamnd and passing the output on to parse the content
	if checkExecutableExists("lshw") {
		strLshwOutput, _ := runCommand(strings.Fields("lshw -c system -quiet"))
		parseLshw(strLshwOutput)
	}

	//The execution of 'dmidecode' requires elevated user privileges to cature all Information
	//If logged in as user root or runnning it in sudo, 'dmidecode' will be run and the output will be passed on and parsed.
	if rootUser() {
		strDmiDecodeOutput, _ := runCommand(strings.Fields("dmidecode --no-sysfs -t baseboard -q"))
		parseDmiDecode(strDmiDecodeOutput)
	}

	//
	if checkExecutableExists("free") {
		strFree, _ := runCommand(strings.Fields("free --si -h"))
		slcInstMem := strings.Split(strFree, "\n")
		installedMemory := strings.Fields(slcInstMem[1])[1]
		fmt.Println(formatBoldWhite("Installed Memory:"), installedMemory)
		nvmeshDiag.InstalledMemory = installedMemory
	}
	if checkExecutableExists("lsb_release") {
		fmt.Println(formatBoldWhite("\nOperating System:"))
		strOSinfo, _ := runCommand(strings.Fields("lsb_release -a"))
		strKernel, _ := runCommand(strings.Fields("uname -r"))
		parseOSinfo(strOSinfo, strKernel)
	}

	checkNvmeshStatus()

	checkSystemTuning()

	if checkExecutableExists("ofed_info") {
		ofedInfo, _ := runCommand(strings.Fields("ofed_info -n"))
		ofedVersion := strings.TrimRight(ofedInfo, "\n\r")
		println(formatBoldWhite("\nMellanox OFED:"), ofedVersion)
		nvmeshDiag.OFEDInfo = ofedVersion
	} else {
		fmt.Println(formatBoldWhite("\nMellanox OFED:"), "No OFED found.")
	}

	if !checkFirewall() == false {
		sWarning := "Warning. Firewall is running! Make sure that all necessary TCP/IP ports as listed in the Excelero NVMesh documentation are configured and open."
		fmt.Println(formatBoldWhite("\nFirewall:"), formatYellow(sWarning))
		mReport["Firewall"] = sWarning
		nvmeshDiag.FirewallInfo = sWarning
	} else {
		fmt.Println(formatBoldWhite("\nFirewall:"), "No firewall found.")

	}

	fmt.Println(formatBoldWhite("\nCPU Information:"))
	strLsCpuOutput, _ := runCommand(strings.Fields("lscpu"))
	parseLsCpu(strLsCpuOutput)

	nvmeshDiag.OperatingSystem = operatingSystem
	nvmeshDiag.NVMeshInfo = nvmeshInfo
	nvmeshDiag.SystemTuning = systemTuning
	nvmeshDiag.CPUInfo = cpuInfo

	strLspciOutput, _ := runCommand(strings.Fields("lspci -vvv"))
	parseLSPCI(strLspciOutput)

	if rootUser() {

		if checkExecutableExists("nvme") {
			fmt.Println(formatBoldWhite("\nMore NVMe drive details:"))
			strNvmeCliOutput, _ := runCommand(strings.Fields("nvme list"))
			slcNvmeCliOutput := strings.Split(strNvmeCliOutput, "\n")
			for _, line := range slcNvmeCliOutput {
				fmt.Println("\t", line)
			}
		} else {
			fmt.Println(formatBoldWhite("\nMore NVMe drive details:"),
				"The nvme utilities are not installed or missing. Cannot scan the NVMe drives for more details without the nvme utilities.")
		}
	} else {
		fmt.Println(formatYellow("Need root user or sudo to get lower level NVme device information"))
	}

	fmt.Println(formatBoldWhite("\nGeneric Block Device Information:"))
	strLsBlkOutput, _ := runCommand(strings.Fields("lsblk -l -d -f"))
	slcLsBlkOutput := strings.Split(strLsBlkOutput, "\n")
	for _, line := range slcLsBlkOutput {
		fmt.Println("\t", line)
	}

	fmt.Println(formatBoldWhite("\nIP Network Interface Information:"))

	for _, ipCommand := range []string{"ip -4 a s", "ip -s link"} {

		strIpOutput, _ := runCommand(strings.Fields(ipCommand))
		slcIpOutput := strings.Split(strIpOutput, "\n")
		for _, line := range slcIpOutput {
			fmt.Println("\t", line)
		}
	}

	parseIBDEVInfo()

	fmt.Println(formatBoldWhite("\nSummary:"))
	if len(mReport) < 1 {
		fmt.Println("\t No troubles found.")
	} else {
		for k, v := range mReport {
			fmt.Println("\t", k, v)
		}
	}
}

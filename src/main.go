package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	mReport          = make(map[string]string)
	plainOutput      = false
	jsonOutputFormat = false
)

func init() {

}

func main() {

	flag.BoolVar(&plainOutput, "p", false, "No color and no formatted console output in plain text.")
	flag.BoolVar(&jsonOutputFormat, "j", false, "Print/output the result and information in JSON format.")
	flag.Parse()

	checkUser()

	fmt.Println(formatBoldWhite("Starting System Scan. Please wait..."))
	strLshwOutput, _ := runCommand(strings.Fields("lshw -c system -quiet"))
	parseLshw(strLshwOutput)

	if rootUser() {
		strDmiDecodeOutput, _ := runCommand(strings.Fields("dmidecode --no-sysfs -t baseboard -q"))
		parseDmiDecode(strDmiDecodeOutput)
	}
	if checkExecutableExists("free") {
		strFree, _ := runCommand(strings.Fields("free --si -h"))
		slcInstMem := strings.Split(strFree, "\n")
		fmt.Println(formatBoldWhite("Installed Memory:"), strings.Fields(slcInstMem[1])[1])
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
		ofedVersion, _ := runCommand(strings.Fields("ofed_info -n"))
		println(formatBoldWhite("\nMellanox OFED:"), strings.TrimRight(ofedVersion, "\n\r"))
	} else {
		fmt.Println(formatBoldWhite("\nMellanox OFED:"), "No OFED found.")
	}

	if !checkFirewall() == false {
		sWarning := "Warning. Firewall is running! Make sure that all necessary TCP/IP ports as listed in the Excelero NVMesh documentation are configured and open."
		fmt.Println(formatBoldWhite("\nFirewall:"), formatYellow(sWarning))
		mReport["Firewall"] = sWarning
	} else {
		fmt.Println(formatBoldWhite("\nFirewall:"), "No firewall found.")
	}

	fmt.Println(formatBoldWhite("\nCPU Information:"))
	strLsCpuOutput, _ := runCommand(strings.Fields("lscpu"))
	parseLsCpu(strLsCpuOutput)

	strLspciOutput, _ := runCommand(strings.Fields("lspci -vvv"))
	parseLSPCI(strLspciOutput)

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
			fmt.Println("\t", k, ":", v)
		}
	}
}

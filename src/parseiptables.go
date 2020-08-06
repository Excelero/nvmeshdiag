package main

import (
	"fmt"
	"strings"
	"strconv"
)

func parseIPTables(s string) bool {
	// Default management ports
	fmt.Println(formatBoldWhite("\nFirewall:\n\tFirewall is running! Checking default ports (4000-4006):"))
	port := 4000
	all_ok := true
	for i := 0; i < 7; i++ {
		port_string := strconv.Itoa(port + i)
		if !strings.Contains(s, "tcp dpt:" + port_string) {
			fmt.Println(formatYellow("\tNVMesh Manangement Port tcp " + port_string + " must be set and open!"))
			all_ok = false
		}
	}
	return all_ok
}


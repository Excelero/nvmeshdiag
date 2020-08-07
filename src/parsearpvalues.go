package main

import (
	"fmt"
	"strings"
)

func parseArpValues() bool {
	success := true
	if checkExecutableExists("sysctl") {
		const all_prefix string = "net.ipv4.conf.all"
		const default_prefix string = "net.ipv4.conf.default"
		fmt.Println(formatBoldWhite("\nChecking arp values:"))
		ALL_ARP_FILTER,_ := runCommand(strings.Fields("sysctl -ar "+all_prefix+".arp_filter"))
		DEFAULT_ARP_FILTER,_ := runCommand(strings.Fields("sysctl -ar "+default_prefix+".arp_filter"))
		ALL_RP_FILTER,_ := runCommand(strings.Fields("sysctl -ar "+all_prefix+".rp_filter"))
		DEFAULT_RP_FILTER,_ := runCommand(strings.Fields("sysctl -ar "+default_prefix+".rp_filter"))
		ALL_ARP_IGNORE,_ := runCommand(strings.Fields("sysctl -ar "+all_prefix+".arp_ignore"))
		DEFAULT_ARP_IGNORE,_ := runCommand(strings.Fields("sysctl -ar "+default_prefix+".arp_ignore"))
		ALL_ARP_ANNOUNCE,_ := runCommand(strings.Fields("sysctl -ar "+all_prefix+".arp_announce"))
		DEFAULT_ARP_ANNOUNCE,_ := runCommand(strings.Fields("sysctl -ar "+default_prefix+".arp_announce"))
		if !strings.Contains(ALL_ARP_FILTER, "1") {
			fmt.Println(formatYellow("\tWARNING: " + strings.TrimSpace(ALL_ARP_FILTER) + ". Please set to 1, using 'sysctl -w "+all_prefix+".arp_filter=1'"))
			success = false
		}
		if !strings.Contains(DEFAULT_ARP_FILTER, "1") {
			fmt.Println(formatYellow("\tWARNING: " + strings.TrimSpace(DEFAULT_ARP_FILTER) + ". Please set to 1, using 'sysctl -w "+default_prefix+".arp_filter=1'"))
			success = false
		}
		if !strings.Contains(ALL_RP_FILTER, "2") {
			fmt.Println(formatYellow("\tWARNING: " + strings.TrimSpace(ALL_RP_FILTER) + ". Please set to 2, using 'sysctl -w "+all_prefix+".rp_filter=2'"))
			success = false
		}
		if !strings.Contains(DEFAULT_RP_FILTER, "2") {
			fmt.Println(formatYellow("\tWARNING: " + strings.TrimSpace(DEFAULT_RP_FILTER) + ". Please set to 2, using 'sysctl -w "+default_prefix+".rp_filter=2'"))
			success = false
		}
		if !strings.Contains(ALL_ARP_IGNORE, "2") {
			fmt.Println(formatYellow("\tWARNING: " + strings.TrimSpace(ALL_ARP_IGNORE) + ". Please set to 2, using 'sysctl -w "+all_prefix+".arp_ignore=2'"))
			success = false
		}
		if !strings.Contains(DEFAULT_ARP_IGNORE,"2") {
			fmt.Println(formatYellow("\tWARNING: " + strings.TrimSpace(DEFAULT_ARP_IGNORE) + ". Please set to 2, using 'sysctl -w "+default_prefix+".arp_ignore=2'"))
			success = false
		}
		if !strings.Contains(ALL_ARP_ANNOUNCE, "2") {
			fmt.Println(formatYellow("\tWARNING: " + strings.TrimSpace(ALL_ARP_ANNOUNCE) + ". Please set to 2, using 'sysctl -w "+all_prefix+".arp_announce=2'"))
			success = false
		}
		if !strings.Contains(DEFAULT_ARP_ANNOUNCE, "2") {
			fmt.Println(formatYellow("\tWARNING: " + strings.TrimSpace(DEFAULT_ARP_ANNOUNCE) + ". Please set to 2, using 'sysctl -w "+default_prefix+".arp_announce=2'"))
			success = false
		}
		if success {
			fmt.Println(formatBoldWhite("\tArp values correctly assigned"))
		}
	} else {
		fmt.Println(formatBoldWhite("\tCannot check arp values"))
	}
	return success
}

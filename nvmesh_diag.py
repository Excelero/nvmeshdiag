#!/usr/bin/env python
#
# Copyright (c) 2018 Excelero, Inc. All rights reserved.
#
# This Software is licensed under one of the following licenses:
#
# 1) under the terms of the "Common Public License 1.0" a copy of which is
#    available from the Open Source Initiative, see
#    http://www.opensource.org/licenses/cpl.php.
#
# 2) under the terms of the "The BSD License" a copy of which is
#    available from the Open Source Initiative, see
#    http://www.opensource.org/licenses/bsd-license.php.
#
# 3) under the terms of the "GNU General Public License (GPL) Version 2" a
#    copy of which is available from the Open Source Initiative, see
#    http://www.opensource.org/licenses/gpl-license.php.
#
# Licensee has the right to choose one of the above licenses.
#
# Redistributions of source code must retain the above copyright
# notice and one of the license notices.
#
# Redistributions in binary form must reproduce both the above copyright
# notice, one of the license notices in the documentation
# and/or other materials provided with the distribution.
#
# Author:        Andreas Krause
# Version:       1
# Maintainer:    Andreas Krause
# Email:         andreas@excelero.com

import argparse
import subprocess
from tempfile import TemporaryFile
import logging
import platform
import datetime
import re

parser = argparse.ArgumentParser()
parser.add_argument('-m', '--mode',
                    help='Please specify the run mode. Either preinstall, healthcheck or full. Usage nvmeshdiag '
                         '-m/--mode [preinstall|healthcheck|full]',
                    required=False)
parser.add_argument('-v', '--verbose', type=bool, nargs='?', const=True, default=False,
                    help="Activate verbose output mode.")
parser.add_argument('-s', '--set-parameters', type=bool, nargs='?', const=True, default=False,
                    help="Set the recommended parameters where possible.")
args = parser.parse_args()
exec_mode = getattr(args, "mode")
verbose_mode = getattr(args, "verbose")
set_parameters = getattr(args, "set_parameters")

REGEX_HCAID = r"(mlx5_\d*)"
REGEX_INSTALLED_MEMORY = r"\S*Mem:\s*(\d*[A-Za-z])"
REGEX_HCA_MAX = r"LnkCap:\s\S*\s\S*\s\S*\s([A-Za-z0-9]*/s),\s\S*\s(\S[0-9]*)"
REGEX_HCA_ACTUAL = r"LnkSta:\s\S*\S*\s([A-Za-z0-9]*/s),\s\S*\s(\S[0-9]*)"
EXCELERO_MANAGEMENT_PORTS = [("tcp", 4000), ("tcp", 4001)]
ROCEV2_TARGET_PORT = ("udp", 4791)
MONGODB_PORT = ("tcp", 27017)
MELLANOX_INBOX_DRIVERS = ["libibverbs", "librdmacm", "libibcm", "libibmad", "libibumad", "libmlx4", "libmlx5", "opensm",
                          "ibutils", "infiniband-diags", "perftest", "mstflint", "rdmacmutils", "ibverbs-utils",
                          "librdmacm-utils", "libibverbs-utils"]
CMD_SET_TUNED_PARAMETERS = "tuned-adm profile latency-performance"
CMD_SET_SELINUX = "sed -i 's/^SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config"
CMD_SET_ONE_QP = "mlxconfig -d %s -b /etc/opt/NVMesh/Excelero_mlxconfig.db set ONE_QP_PER_RECOVERY=1"
CMD_GET_ONE_QP = "mlxconfig -d %s -b ./Excelero_mlxconfig.db query ONE_QP_PER_RECOVERY | grep ONE_QP_PER_RECOVERY"
CMD_DISABLE_FIREWALL = ["systemctl stop firewalld", "systemctl disable firewalld"]
CMD_SET_FIREWALL_FOR_NVMESH_MGMT = ["firewall-cmd --permanent --direct --add-rule ipv4 filter INPUT 0 -p tcp --dport "
                                    "4000 -j ACCEPT -m comment --comment Excelero-Management", "firewall-cmd "
                                    "--permanent --direct --add-rule ipv4 filter INPUT 0 -p tcp --dport 4001 -j "
                                    "ACCEPT -m comment --comment Excelero-Management"]
CMD_SET_FIREWALL_FOR_ROCEV2 = "firewall-cmd --permanent --direct --add-rule ipv4 filter INPUT 0 -p udp --dport " \
                              "4791 -j ACCEPT -m comment --comment RoCEv2-Target"
CMD_SET_FIREWALL_FOR_MOGODB = "firewall-cmd --permanent --direct --add-rule ipv4 filter INPUT 0 -p tcp --dport 27017 " \
                              "-j ACCEPT -m comment --comment MongoDB"
CMD_RELOAD_FIREWALL_RULES = "firewall-cmd --reload"
CMD_GET_IRQ_BALANCER_STATUS = "systemctl status irqbalance"
CMD_START_IRQ_BALANCER = ["systemctl enable irqbalance", "systemctl start irqbalance"]
CMD_GET_FIREWALLD_STATUS = "systemctl status firewalld"

host_name = platform.node()
output = open(host_name + '_' + str(datetime.datetime.utcnow()).replace(" ", "_") + '_nvmesh_diag_output.txt', 'w')


def get_cmd_output(cmd_to_execute):
    with TemporaryFile() as t:
        try:
            logging.debug("Running Shell Command: " + cmd_to_execute)
            out = subprocess.check_output(cmd_to_execute, stderr=t, shell=True)
            logging.debug("Success")
            return 0, out
        except subprocess.CalledProcessError as e:
            t.seek(0)
            logging.error(t.read().strip("\n"))
            return e.returncode, cmd_to_execute, t.read()


def return_cmd_output(cmd):
    cmd_output = get_cmd_output(cmd)
    if cmd_output[0] == 0:
        return cmd_output[1]
    elif cmd_output[0] == 127:
        return print_yellow(cmd.strip().split(" ")[0] + " not found or not installed!")
    elif cmd_output[0] == 255:
        return print_yellow(cmd.strip().split(" ")[
                                0] + " shows no data as there is no IB transport layer. Looks like Ethernet "
                                     "connectivity.")
    elif cmd_output[0] != 0:
        return print_red("Error has occurred while executing " + cmd.strip().split(" ")[
            0] + "! Details can be found in the nvmesh_diag.log file.")


def run_cmd(cmd):
    cmd_output = get_cmd_output(cmd)
    if cmd_output[0] == 0:
        return print_green("Done. Successful.")
    elif cmd_output[0] == 127:
        return print_red("Error! " + cmd.strip().split(" ")[0] + " not found or not installed!")
    elif cmd_output[0] != 0:
        return print_red("Error has occurred while executing " + cmd.strip().split(" ")[
            0] + "! Details can be found in the nvmesh_diag.log file.")


def print_and_log_info(text):
    logging.info(strip_tabs((text.lstrip("\n")).rstrip(":")))
    print('\033[1m' + '\033[4m' + text + '\033[0m')
    output.write(text + "\n")


def strip_tabs(text):
    text = text.replace("\t", "")
    return text


def print_green(text):
    logging.info(strip_tabs(text))
    output.write(text + "\n")
    return '\033[92m' + text + '\033[0m'


def print_yellow(text):
    logging.warning(strip_tabs(text))
    output.write(text + "\n")
    return '\033[33m' + text + '\033[0m'


def print_red(text):
    logging.error(text + "\n")
    output.write(text)
    return '\033[31m' + text + '\033[0m'


# region Logging Configuration
logging.basicConfig(filename='nvmesh_diag.log', format='%(asctime)s\t%(levelname)s\t%(message)s',
                    datefmt='%m/%d/%Y %I:%M:%S %p', level=logging.DEBUG)
logging.debug("Execution Started -------------------------------------------------------------------------------------")
logging.debug("Execution mode: " + unicode(exec_mode))
logging.debug("Verbose output enabled: " + unicode(verbose_mode))
logging.debug("Setting recommended parameters: " + unicode(set_parameters))
logging.debug("Storing output to: " + output.name)
# endregion

# region Collecting Host Name Information - non critical for NVMesh install
print_and_log_info("Collecting Host Name Information:")
output.write(host_name + "\n\n")
if verbose_mode is True:
    print(host_name)
print("Done.")
# endregion

# region Collecting Hardware Vendor and System Information - non critical for NVMesh install
print_and_log_info("\nCollecting Hardware Vendor and System Information:")
system_information = return_cmd_output("dmidecode | grep -A 4 'System Information'").split("\n")
base_board_information = return_cmd_output("dmidecode | grep -A 5 'Base Board Information'").split("\n")
output.write(str(system_information[1].strip() + "\n" + system_information[2].strip() + "\n\n"))

if verbose_mode is True:
    print system_information[1].strip(), system_information[2]
print("Done.")
# endregion

# region Collecting Operating System Information - critical for NVMesh install
print_and_log_info("\nCollecting Operating System Information:")
os_details = platform.linux_distribution()[0:2]
kernel_version = str(platform.release())
output.write(os_details[0] + os_details[1] + kernel_version + "\n\n")
print os_details[0], os_details[1], kernel_version, "\nPlease verify this information with the latest support matrix!"
# endregion

# region Collecting and Verifying SELinux Information - critical for NVMesh install
print_and_log_info("\nCollecting and Verifying SELinux Information:")
selinux_status = return_cmd_output("sestatus").splitlines()
for line in selinux_status:
    output.write(line + "\n")
output.write("\n")

if not selinux_status[0].split(":")[1].strip() == "disabled":
    print print_red("SELinux active. Should be disabled!")
    if set_parameters is True:
        if "y" in raw_input("Do you want to Disable SELinux now?[Yes/No]: ").lower():
            print(run_cmd(CMD_SET_SELINUX))
        else:
            print("No it is. Going on...")
    if verbose_mode is True:
        for line in selinux_status:
            print(line)
    print("\n")
else:
    print print_green("Disabled - OK")
# endregion

# region Collecting and Verifying Firewall Information - critical for NVMesh install
print_and_log_info("\nCollecting and Verifying Firewall Information:")
if not return_cmd_output("systemctl status firewalld | grep Active").strip(" ").split(" ")[1] == "inactive":
    iptables_output = return_cmd_output("iptables -nL")
    output.write(iptables_output + "\n")
    print(print_yellow("Firewall active!"))
    if set_parameters is True:
        if "y" in raw_input("Do you want to disable the firewall now?[Yes/No]: ").lower():
            for command in CMD_DISABLE_FIREWALL:
                print(run_cmd(command))
            pass
        else:
            print("No it is. Going on...")
    if re.findall(r"%s dpt:%s" % (EXCELERO_MANAGEMENT_PORTS[0][0], EXCELERO_MANAGEMENT_PORTS[0][1]), iptables_output) and re.findall(r"%s dpt:%s" % (EXCELERO_MANAGEMENT_PORTS[1][0], EXCELERO_MANAGEMENT_PORTS[1][1]), iptables_output):
        print print_green("For NVMesh Client operations OK. Excelero Management ports are configured.")
    else:
        print print_red(
            "Not OK for NVMesh client operations. Excelero Management ports tcp 4000 and tcp 4001 must be set and "
            "open!")
        if set_parameters is True:
            if "y" in raw_input("Do you want to open port 4000 and 4001 now?[Yes/No]: ").lower():
                for command in CMD_SET_FIREWALL_FOR_NVMESH_MGMT:
                    print(run_cmd(command))
            else:
                print("No it is. Going on...")
    if re.findall(r"%s dpt:%s" % (ROCEV2_TARGET_PORT[0], ROCEV2_TARGET_PORT[1]), iptables_output):
        print print_green("For NVMesh Target operations OK if the Link Layer is Ethernet. RoCEv2 ports are set.")
    else:
        print print_red(
            "Not OK for NVMesh Target operations if the Link Layer is Ethernet. RoCEv2 udp port 4791 must be set and "
            "open!")
        if set_parameters is True:
            if "y" in raw_input("Do you want to open port 4791 now?[Yes/No]: ").lower():
                    print(run_cmd(CMD_SET_FIREWALL_FOR_ROCEV2))
        else:
            print("No it is. Going on...")
    if re.findall(r"%s dpt:%s" % (MONGODB_PORT[0], MONGODB_PORT[1]), iptables_output):
        print print_green("For MongoDB HA operations OK. MongoDB ports are set.")
    else:
        print print_red("Not OK for MongoDB HA operations. MongoDB tcp port 27017 must be set and open!")
        if set_parameters is True:
            if "y" in raw_input("Do you want to open port 4791 now?[Yes/No]: ").lower():
                print(run_cmd(CMD_SET_FIREWALL_FOR_MOGODB))
            else:
                print("No it is. Going on...")
    if set_parameters is True:
        print "Reloading the firewall rules to apply changes.\n", run_cmd(CMD_RELOAD_FIREWALL_RULES)
else:
    print print_green("Disabled - OK")
# endregion

# region Collecting and Verifying CPU Information - non critcal for NVMesh install
print_and_log_info("\nCollecting and Verifying CPU Information:")
cpu_info = return_cmd_output("lscpu").split("\n")
for line in cpu_info:
    output.write(line + "\n")
output.write("\n")
actual_cpu_frequency = None
max_cpu_frequency = None
for line in cpu_info:
    if "Socket(s)" in line:
        print line.split(":")[1].strip() + " Physical CPU"
    if "Model name" in line:
        print line.split(":")[1].lstrip()
    if "CPU MHz" in line:
        actual_cpu_frequency = float((line.split(":")[1]).strip())
        max_frequency_info = re.compile("\d+")
        for frequency in return_cmd_output("dmidecode -s processor-frequency").split("\n"):
            match = max_frequency_info.search(frequency)
            if match:
                max_cpu_frequency = float(match.group(0))
if actual_cpu_frequency and max_cpu_frequency is not None:
    if max_cpu_frequency - actual_cpu_frequency >= 100:
        print print_yellow(
            "Actual running CPU frequency is lower than the maximum CPU frequency. This might impact performance! "
            "Check BIOS settings and verify System Tuning settings as below.\n")
    else:
        print print_green("CPU frequency settings OK.\n")
# endregion

# region Collecting and Verifying System Tuning Information - not critical for NVMesh install
print_and_log_info("Collecting and Verifying System Tuning Information:")
tuned_adm_info = return_cmd_output("tuned-adm active")
output.write(tuned_adm_info + "\n")
if "latency-performance" in tuned_adm_info:
    print(print_green("Tuned profile settings are OK"))
else:
    print(print_yellow(
        "Tuned settings is not as recommended! Please run 'tuned-adm profile latency-performance' to set and enable "
        "the recommended Tuned profile and veryfy the Tuned service is running."))
    if set_parameters is True:
        if "y" in raw_input("Do you want to set the recommended tuned parameters now?[Yes/No]: ").lower():
            print(run_cmd(CMD_SET_TUNED_PARAMETERS))
        else:
            print("No it is. Going on...")
irq_balancer_info = return_cmd_output(CMD_GET_IRQ_BALANCER_STATUS)
if "active" in irq_balancer_info.lower():
    print print_green("IRQ Balancer is running - OK")
else:
    print(print_yellow("IRQ Balance is not running. This might severely impact the system performance."))
    if set_parameters is True:
        if "y" in raw_input("Do you want to Enable and Start the IRQ Balancer now?[Yes/No]: "):
            for command in CMD_START_IRQ_BALANCER:
                print(run_cmd(command))
        else:
            print("No it is. Going on...")
# endregion

# region Collecting Memory Information - not critical for NVMesh install
print_and_log_info("\nCollecting Memory Information:")
memory_info = return_cmd_output("free -h")
output.write(memory_info + "\n")
installed_memory = re.findall(REGEX_INSTALLED_MEMORY, memory_info)[0]
if verbose_mode is True:
    print(installed_memory + " installed memory.\n")
print("Done.")
# endregion

# region Collecting High Level Block Device Information - not crotical for NVMesh install
print_and_log_info("Collecting High Level Block Device Information:")
lsblk_output = return_cmd_output("lsblk")
output.write(lsblk_output + "\n")
if verbose_mode is True:
    print(lsblk_output)
print("Done.")
# endregion

# region Collecting NVMe Storage Device Information - not critical for NVMesh install
print_and_log_info("Collecting NVMe Storage Device Information:")
nvme_list_output = return_cmd_output("./nvme list").splitlines()
nvme_numa_output = None
if len(nvme_list_output) > 2:
    nvme_numa_output = return_cmd_output("lspci -vv | grep -A 10 Volatile | grep -e Volatile -e NUMA")
    if verbose_mode is True:
        for line in nvme_list_output:
            print(line)
        print("\n" + nvme_numa_output)
        output.write("\n" + nvme_numa_output + "\n")
else:
    print print_yellow("No NVMe SSD found on this server! This server can only be configured as a NVMesh Client.")
for line in nvme_list_output:
    output.write(line + "\n")
print("Done.")
# endregion

# region Collecting And Verifying R-NIC Information - critical for NVMesh install
print_and_log_info("\nCollecting and Verifying R-NIC information:")
rnics_output = return_cmd_output(
    "for i in `lspci | awk '/Mellanox/ {print $1}'`;do echo $i; echo \"FW level:\" | tr '\n' ' '; cat "
    "/sys/bus/pci/devices/0000:$i/infiniband/mlx*_*/fw_ver; lspci -s $i -vvv | egrep -e Connect-X -e "
    "\"Product Name:\" -e Subsystem -e NUMA -e \"LnkSta:\" -e \"LnkCap\" -e \"MaxPayload\"; echo ""; done")
output.write("\n" + rnics_output + "\n")
rnics = rnics_output.split("\n\n")

for rnic in rnics:
    max_rnic_speed_and_pcie_width = re.findall(REGEX_HCA_MAX, rnic)
    actual_rnic_speed_and_pcie_with = re.findall(REGEX_HCA_ACTUAL, rnic)
    rnic_details = rnic.splitlines()
    if len(rnic_details) > 0:
        print("Checking HCA at PCIe address: " + rnic_details[0])
        print("\tVendor/OEM information:" + rnic_details[2].split("Device")[0])
        try:
            if "Product Name" in rnic_details[8]:
                print "\tHCA Type: " + rnic_details[8].split(":")[1]
        except:
            pass
        print("\tFirmware level: " + rnic_details[1].split(":")[1].strip())

        if max_rnic_speed_and_pcie_width[0][0] == actual_rnic_speed_and_pcie_with[0][0]:
            print print_green("\tHCA PCIe speed settings OK. Running at " + actual_rnic_speed_and_pcie_with[0][0])
        else:
            print print_yellow("\tThe HCA is capable of ") + max_rnic_speed_and_pcie_width[0][
                0] + " but its running at " + actual_rnic_speed_and_pcie_with[0][
                      0] + "! Check BIOS and HW settings to ensure max performance and a stable environment!"
        if max_rnic_speed_and_pcie_width[0][1] == actual_rnic_speed_and_pcie_with[0][1]:
            print print_green(
                "\tHCA PCIe width settings OK. Running at " + actual_rnic_speed_and_pcie_with[0][1]) + "\n"
        else:
            print print_yellow("\tThe HCA is capable of ") + max_rnic_speed_and_pcie_width[0][
                1] + " but its running at " + actual_rnic_speed_and_pcie_with[0][
                      0] + "! Check BIOS and HW settings to ensure max performance and a stable environment!\n"
# endregion

# region Collecting And Verifying Mellanox Driver Information - critical for NVMesh install
print_and_log_info("Collecting And Verifying Mellanox Driver Information:")
ofed_version = return_cmd_output("ofed_info -n")
output.write("OFED: " + ofed_version + "\n")
if "not found or not installed" in ofed_version:
    print "OFED not installed! Checking for inbox drivers now."
    missing_inbox_drivers = []
    for rpm_package in MELLANOX_INBOX_DRIVERS:
        if "Error" in return_cmd_output("rpm -q %s" % rpm_package):
            missing_inbox_drivers.append(rpm_package)
            output.write("missing " + rpm_package + "!")
            print print_red("\t%s missing " % rpm_package + "!")
    if len(missing_inbox_drivers) != 0:
        print print_red("OFED is not installed and Inbox drivers are missing! \n"
                        "You must install either OFED or the missing Inbox drivers!\n")
    else:
        print print_green("Inbox drivers installed - OK ")
else:
    print print_green("OFED installed - OK "), "\nVersion: " + ofed_version + "Please verify this information " \
                                                                              "with the latest support matrix!"
# endregion

# region Collecting and Verifying RDMA HCA Information - critical for NVMesh install
print_and_log_info("\nCollecting And Verifying RDMA Specific Information:")
ibv_devinfo_output = return_cmd_output("ibv_devinfo | grep -e hca_id -e guid")
output.write(ibv_devinfo_output + "\n")
hca_list = re.findall("(mlx5_\\d*)\\s*node_guid:\\s*([A-Za-z0-9]*):([A-Za-z0-9]*):([A-Za-z0-9]*):([A-Za-z0-9]*)",
                      ibv_devinfo_output, re.MULTILINE)
for (hca, guid1, guid2, guid3, guid4) in hca_list:
    if guid1 == "0000":
        print print_red(hca + " guid information seems incorrect or missing. Please check!")
    else:
        print print_green(hca + " guid OK")
    one_qp_per_recovery = re.sub("\s\s+", " ", (return_cmd_output(CMD_GET_ONE_QP % hca).lstrip(" ")))
    if "True" in one_qp_per_recovery:
        print print_green(hca + " ready and configured for RDDA")
    elif "False" in one_qp_per_recovery:
        print print_yellow(
            hca + " will support RDDA but is not configured correctly. You have to enable ONE_QP_PER_RECOVERY in the "
                  "Mellanox firmware if you want to use RDDA")
    elif "-E-" in one_qp_per_recovery:
        print print_red(
            hca + " will not support RDDA due to firmware limitations on the HCA. If you intent to use RDDA, you have "
                  "to update the firmware and enable ONE_QP_PER_RECOVERY on the HCA.")
if verbose_mode is True:
    print(return_cmd_output("ibdev2netdev -v"))
print("\nDone.")
# endregion

# region Collect Infiniband Environment Specific Information - not critical for NVMesh install
print_and_log_info("\nCollecting Infiniband Specific Information:")
output.write(return_cmd_output("ibhosts") + "\n" + return_cmd_output("ibswitches") + "\n")
if verbose_mode is True:
    print return_cmd_output("ibhosts"), "\n" + return_cmd_output("ibswitches")
print("Done.")
# endregion

# region Collect IP Address Information - not critical for NVMesh install
print_and_log_info("\nCollecting IP Address Information:")
output.write(return_cmd_output("ip -4 a"))
if verbose_mode is True:
    print return_cmd_output("ip -4 a")
print("Done.")
# endregion

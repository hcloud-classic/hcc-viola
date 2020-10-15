package controlcli

import (
	"fmt"
	"hcc/viola/lib/logger"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

//VerifyClusrter : Verify Cluster system
func VerifyClusrter() bool {
	act := "krgadm cluster"
	cmd := exec.Command("/bin/bash", "-c", act)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return true
	// return strings.TrimSuffix(string(result), "\n")
}

//WriteNFSConfigObject : /etc/exports
func WriteNFSConfigObject() {
	filename := nfsFilePath
	// volume.Pool = config.volumeig.VOLUMEPOOL
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("err")
	}
	defer func() {
		_ = file.Close()
	}()

	input := nfsConfigBuilder()

	_, err = file.WriteString(input)
	if err != nil {
		logger.Logger.Println("WriteNFSConfigObject Faild")
	}
	logger.Logger.Println("WriteNFSConfigObject Success")

}

func nfsConfigBuilder() string {
	var retContents string
	// hostName := hostnameSet()
	clustersIP := hostIPSet()
	for _, args := range mountPointRenewal() {
		tempContents := nfsContents
		tempContents = strings.Replace(tempContents, "PATH", args, -1)
		parseClusterIP := strings.Split(clustersIP, ".")
		exportClusterIP := parseClusterIP[0] + "." + parseClusterIP[1] + "." + parseClusterIP[2] + ".0"

		tempContents = strings.Replace(tempContents, "CLUSTERADDRESS", exportClusterIP, -1)
		tempContents = strings.Replace(tempContents, "NETMASK", "255.255.255.0", -1)
		retContents += tempContents
	}

	return retContents
}

func mountPointRenewal() []string {
	act := "mount | grep /dev/sd | awk '{print $3}'|xargs -i echo {} "
	cmd := exec.Command("/bin/bash", "-c", act)

	result, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(result)
	}
	mounPath := strings.Split(strings.Trim(string(result), "\n"), "\n")
	return mounPath

}

//WriteHostsConfigObject : /etc/hosts
func WriteHostsConfigObject() {
	filename := hostFilePath
	// volume.Pool = config.volumeig.VOLUMEPOOL
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("err")
	}
	defer func() {
		_ = file.Close()
	}()

	input := hostConfigBuilder()

	_, err = file.WriteString(input)
	if err != nil {
		logger.Logger.Println("WriteHostsConfigObject Faild")

		// strerr := "create_volume action status=>iscsistatus " + fmt.Sprintln(err)
	}
	logger.Logger.Println("WriteHostsConfigObject Success")

}

func hostConfigBuilder() string {
	var computeList string
	tempContents := hostContents
	hostName := hostnameSet()
	clustersIP := hostIPSet()
	tempContents = strings.Replace(tempContents, "CLUSTERHOSTNAME", hostName, -1)
	tempContents = strings.Replace(tempContents, "LEADERNODEHOSTNAME", hostName+convComputeIP(strconv.Itoa(1)), -1)
	tempContents = strings.Replace(tempContents, "CLUSTERADDRESS", clustersIP, -1)
	parseClusterIP := strings.Split(clustersIP, ".")
	for i := 2; i <= 4; i++ {
		computeList += parseClusterIP[0] + "." + parseClusterIP[1] + "." + parseClusterIP[2] + "." + strconv.Itoa(i) + " " + hostName + convComputeIP(strconv.Itoa(i)) + "\n"
	}

	tempContents = strings.Replace(tempContents, "COMPUTE", computeList, -1)

	return tempContents
}

func convComputeIP(nodeIP string) string {
	var retString string
	act := "printf '\\%03d' " + nodeIP
	cmd := exec.Command("/bin/bash", "-c", act)

	result, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println(err)
	}
	re := regexp.MustCompile("[0-9]+")
	extractscope := re.FindAllString(string(result), -1)
	for _, args := range extractscope {
		retString += args
	}
	return retString
}

func hostIPSet() string {
	act := "ip -4 addr | grep \"eth\\|eno\" | grep -oP '(?<=inet\\s)\\d+(\\.\\d+){3}' "
	cmd := exec.Command("/bin/bash", "-c", act)

	result, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println(err)
	}
	return strings.TrimSuffix(string(result), "\n")
}

func hostnameSet() string {
	act := "cat /etc/sysconfig/network | grep HOSTNAME | cut -d '=' -f2 "
	cmd := exec.Command("/bin/bash", "-c", act)

	result, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println(err)

	}
	return strings.TrimSuffix(string(result), "\n")
}

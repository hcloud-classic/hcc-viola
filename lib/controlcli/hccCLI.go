package controlcli

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var nodemap map[string]string

// HccCli : Hcc integration Command line interface
func HccCli(action string) interface{} {
	parseaction := strings.Split(action, " ")
	area := parseaction[1]
	actclass := parseaction[2]
	actscope, _ := strconv.Atoi(parseaction[3])

	if kerrighedContainerVerify() {
		switch area {
		case "nodes":
			cmdNodes(actclass, actscope)
		case "cluster":
			cmdCluster(actclass, actscope)
		default:
			fmt.Println("Please choose the area {nodes or cluster}")

		}
	} else {
		fmt.Println("Please Continue in Kerrighed Container")
		return false
	}
	return false
}

func cmdNodes(actclass string, actscope int) interface{} {

	switch actclass {
	case "status":
		verbosenode := nodeStatus(actscope)
		fmt.Println("node", verbosenode)
	case "add":
		if checkNFS() {

		} else {
			restartNFS()
		}
	case "del":
	default:
		fmt.Println("Please Choose Operation {status, add, del}")
		return false
	}
	return false
}

func cmdCluster(actclass string, actscope int) {

}

func addNodes(action string) {

}

func checkNFS() bool {
	cmd := exec.Command("service", "nfs-common", "status")
	result, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("NFS Service error")

	} else {
		if strings.Contains(string(result), "all daemons running") {
			return true
		}
	}
	return false
}

func restartNFS() {
	cmd := exec.Command("service", "nfs-common", "restart")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("NFS Service Can't start")

	}

}

// @ N node nodeStatus index == n
// @ all node status index == 0
func nodeStatus(index int) interface{} {
	cmd := exec.Command("krgadm", "nodes", "status")
	result, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Node status error Occured!!")
	}
	fmt.Println("Compute ON? ", verifyNPort("192.168.110.2", "2222"))
	if index == 0 {
		fmt.Println("HCC All Nodes Status \nIP  status\n", result)
		nodeStatusRegister(string(result))
		return string(result)
	} else {
		tmpstr := strings.Split(string(result), "\n")
		for _, words := range tmpstr {
			retoken := strings.Split(string(words), ":")
			if string(words[0]) == strconv.Itoa(index) {
				fmt.Println(index, " th node status = > ", retoken[1])
				return string(retoken[1])
			}
		}
	}
	return nil
}

func nodeOnlineCheck() {

}
func nodeStatusRegister(status string) {
	tmpstr := strings.Split(string(status), "\n")
	for _, words := range tmpstr {
		retoken := strings.Split(string(words), ":")
		nodemap[string(words[0])] = retoken[1]
	}
}

func printOutput(outs string) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", outs)
	}
}

func kerrighedContainerVerify() bool {

	if fileExists("/proc/nodes/self/nodeid") {
		fmt.Println("Kerrighed Container load")
		return true
	}

	return false
}

func extractToken(srcstr string, delimiter string, index int) string {
	tmpstr := strings.Split(srcstr, delimiter)
	return tmpstr[index]
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// cmd := exec.Command("ls", "-al")
// result, err := cmd.CombinedOutput()
// qwe := strings.Split(string(result), "\n")
// for i, words := range qwe {
// 	fmt.Println(i, "= > ", words)
// }
// if err != nil {
// 	fmt.Println("Error occured!!")
// }

func verifyNPort(ip string, port string) bool {
	cmd := exec.Command("nmap", ip)
	result, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(ip, " has not configure the ", port)
	}
	//compute on?
	if strings.Contains(string(result), port) {
		return true
	}

	return false
}

func nodeInit() {
	nodemap = make(map[string]string)
}

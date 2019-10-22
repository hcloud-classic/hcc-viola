package controlcli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//ParedAction : action parser struct
type ParedAction struct {
	area        string
	class       string
	scope       []string
	rangeoption bool
}

var tokenaction ParedAction
var nodemap map[string]string

// ActionParser : Parcing Action
func ActionParser(action string) {

	tmpstr := strings.Split(action, " ")
	tokenaction.area = tmpstr[1]
	tokenaction.class = tmpstr[2]
	hasOption := strings.Contains(action, "-n")
	hasRangeOption := strings.Contains(action, ":")
	tokenaction.rangeoption = hasRangeOption
	//deliIndex : delimeter index
	var deliIndex int
	var endOfIndex int
	if hasOption {
		if hasRangeOption {
			for i, tmpact := range tmpstr {
				if tmpact == "-n" {
					deliIndex = i
				}
				endOfIndex++
			}
			tokenaction.scope = append(tokenaction.scope, tmpstr[deliIndex+1])
		} else {
			for i, tmpact := range tmpstr {
				if tmpact == "-n" {
					deliIndex = i
				}
				endOfIndex++
			}
			for t := deliIndex + 1; t < endOfIndex; t++ {
				tokenaction.scope = append(tokenaction.scope, tmpstr[t])
			}
		}

	}
	//Debug : tokenaction Structure
	// fmt.Println("area =>", tokenaction.area)
	// fmt.Println("class => ", tokenaction.class)
	// fmt.Println("scope => ", tokenaction.scope)

}

// HccCli : Hcc integration Command line interface
func HccCli(action string) interface{} {
	ActionParser(action)
	// fmt.Println(tokenaction.area, tokenaction.class, tokenaction.scope)
	if kerrighedContainerVerify() {
		switch tokenaction.area {
		case "nodes":
			cmdNodes(tokenaction.class, tokenaction.scope)
		case "cluster":
			cmdCluster(tokenaction.class, tokenaction.scope)
		default:
			fmt.Println("Please choose the area {nodes or cluster}")
		}
	} else {
		fmt.Println("Please Continue in Kerrighed Container")
		return false
	}
	clearAction()
	return false
}
func clearAction() {
	tokenaction.area = ""
	tokenaction.class = ""
	tokenaction.scope = nil
}
func cmdNodes(actclass string, actscope []string) interface{} {

	switch actclass {
	case "status":
		verbosenode := nodeStatus(actscope[0])
		fmt.Println("Node Status\n ", verbosenode, "\n++++++++++++++++")
	case "add":
		if checkNFS() {
		} else {
			restartNFS()
		}
		nodeStatus("0")

		if nNodeAvailable(actscope[0]) {
			addNodes(actscope[0])
		}

	case "del":
	default:
		fmt.Println("Please Choose Operation {status, add, del}")
		return false
	}
	return false
}

func cmdCluster(actclass string, actscope []string) {

}

func addNodes(actscope string) interface{} {
	cmd := exec.Command("krgadm", "nodes", "add", "-n", actscope)
	result, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Node Can't add the Num of [", actscope, "] Node")
		return err
	}
	return result
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
func nodeStatus(index string) interface{} {
	cmd := exec.Command("krgadm", "nodes", "status")
	result, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Node status error Occured!!")
	}
	fmt.Println("Compute ON? ", verifyNPort("192.168.110.2", "2222"))
	if index == "0" {
		fmt.Println("HCC All Nodes Status \nIP  status\n", string(result))
		nodeStatusRegister(string(result))
		return string(result)
	} else {
		tmpstr := strings.Split(string(result), "\n")
		for _, words := range tmpstr {
			retoken := strings.Split(string(words), ":")
			if string(words[0]) == index {
				fmt.Println(index, " th node status = > ", retoken[1])
				return string(retoken[1])
			}
		}
	}
	return nil
}

func nodeStatusRegister(status string) {
	nodemap = make(map[string]string)
	tmpstr := strings.Split(status, "\n")
	for _, words := range tmpstr {
		if strings.Contains(string(words), ":") {
			retoken := strings.Split(string(words), ":")
			if string(words[0]) != "1" {
				nodemap[string(words[0])] = retoken[1]
			}
			// fmt.Println("words => ", string(words[0]), "retoken => ", retoken[1])
			// fmt.Println("register => ", nodemap[])

		}
	}
}

// nNodeAvailable : check
func nNodeAvailable(actscope string) bool {
	// fmt.Println("qwe => ", nodemap["1"])
	if nodemap[actscope] == "present" {
		return true
	}
	// for key, val := range nodemap {
	// 	if val == "online" && key != "1" {
	// 		return false
	// 	}
	// 	fmt.Println("Codex => ", key, val)
	// }
	return true
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

// NodeInit :qwe
// func NodeInit() {
// 	nodemap = make(map[string]string)
// 	// nodemap = map[string]string{}
// }

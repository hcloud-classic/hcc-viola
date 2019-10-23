package controlcli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//AtomicAction : action parse struct
type AtomicAction struct {
	area        string
	class       string
	scope       []string
	rangeoption bool
	iprange     []string
}

var actiontype = []string{"area", "class", "scope"}
var tokenaction AtomicAction
var nodemap map[string]string

// HccCli : Hcc integration Command line interface
func HccCli(action string, iprange string) (bool, interface{}) {
	clearAction()
	fmt.Println("Receive : ", action)
	err := ActionParser(action, iprange)
	if err != nil {
		return false, errors.New("ActionParcer Faild")
	}
	fmt.Println(tokenaction.area, tokenaction.class, tokenaction.scope)
	iskerrighed, err := kerrighedContainerVerify()
	if iskerrighed {
		switch tokenaction.area {
		case "nodes":
			return cmdNodes(tokenaction.class, tokenaction.scope)
		case "cluster":
			cmdCluster(tokenaction.class, tokenaction.scope)
		default:
			fmt.Println("Please choose the area {nodes or cluster}")
		}
	} else {
		return false, errors.New("Please Continue in Kerrighed Container")
	}

	return false, nil
}

// ActionParser : Parcing Action
func ActionParser(action string, iprange string) interface{} {
	//Action parsing
	tmpstr := strings.Split(action, " ")
	tmplength := len(tmpstr)
	// fmt.Println("action : ", action, "\n", "Length =>", tmplength, "++++++", tmpstr[1])

	tokenaction.area = tmpstr[1]
	if tmplength <= 3 && tmplength >= 2 {
		tokenaction.class = "status"
		// tokenaction.scope = append(tokenaction.scope, "0")
	} else if tmplength < 2 {
		return errors.New("Hcc Command Line Invalid")
	} else {
		tokenaction.class = tmpstr[2]
		// tokenaction.scope = append(tokenaction.scope, tmpstr[3])
	}

	hasOption := strings.Contains(action, "-n")
	hasRangeOption := strings.Contains(action, ":")
	fmt.Println("hasOption=> ", hasOption, "] hasRangeOption => ", hasRangeOption)
	//deliIndex : delimeter index
	var deliIndex int = 0
	var endOfIndex int = 0
	if hasOption {
		if hasRangeOption {
			tokenaction.rangeoption = hasRangeOption
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

	} else {
		tokenaction.scope = append(tokenaction.scope, "0")
	}
	//ip parsing

	if len(iprange) > 0 {
		iptmp := strings.Split(iprange, " ")
		tokenaction.iprange = append(tokenaction.iprange, iptmp[1])
		tokenaction.iprange = append(tokenaction.iprange, iptmp[2])
		//Debug For iprange
		// for i, words := range iptmp {
		// 	fmt.Println(i, "=>", words)
		// }
	}

	//Debug : tokenaction Structure
	fmt.Println("area =>", tokenaction.area)
	fmt.Println("class => ", tokenaction.class)
	fmt.Println("scope => ", tokenaction.scope)
	fmt.Println("iprange => ", tokenaction.iprange)
	fmt.Println("rangeoption => ", tokenaction.rangeoption)

	return nil
}

func clearAction() {
	tokenaction.area = ""
	tokenaction.class = ""
	tokenaction.scope = nil
	tokenaction.rangeoption = false
	tokenaction.iprange = nil
}
func cmdNodes(actclass string, actscope []string) (bool, interface{}) {

	switch actclass {
	case "status":
		err, verbosenode := nodeStatus(actscope[0])
		if err != false {
			return true, nil
		} else {

			return false, verbosenode
		}
	case "add":
		if checkNFS() {
			fmt.Println("Leader Node NFS Service On")
		} else {
			restartNFS()
			fmt.Println("Leader Node NFS Service restart")

		}
		//For nodeMap renewal
		nodeStatus("0")
		if nAvailableNodeAdd(actscope[0]) {
			return true, errors.New("All Nodes is Preparing with online")
		} else {
			return false, errors.New("All Nodes is Not Preparing with online")
		}
	case "del":
	default:
		return false, errors.New("Please Choose Operation {status, add, del}")
	}
	return false, errors.New("Not Available Command")
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
func nodeStatus(index string) (bool, interface{}) {
	cmd := exec.Command("krgadm", "nodes", "status")
	result, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Node status error Occured!!")
	}

	if index == "0" {
		fmt.Println("HCC All Nodes Status \nIP  status\n", string(result))
		nodeStatusRegister(string(result))
		return true, string(result)
	} else {
		if nodeConnectCheck(index) {
			tmpstr := strings.Split(string(result), "\n")
			for _, words := range tmpstr {
				retoken := strings.Split(string(words), ":")
				if string(words[0]) == index {
					fmt.Println(index, " th node status = > ", retoken[1])
					return true, string(retoken[1])
				}
			}
		} else {
			result := "[" + index + "] th Node Is Not in Cluster Area"
			return false, errors.New(result)
		}

	}
	return false, nil
}

func nodeStatusRegister(status string) {
	nodemap = make(map[string]string)
	tmpstr := strings.Split(status, "\n")
	for _, words := range tmpstr {
		if strings.Contains(string(words), ":") {
			retoken := strings.Split(string(words), ":")
			// if string(words[0]) != "1" {}
			nodemap[string(words[0])] = retoken[1]

			fmt.Println("words => ", string(words[0]), "retoken => ", retoken[1])
			fmt.Println("register => ", nodemap[string(words[0])])

		}
	}
}
func isAllNodeOnline(startRange int, endRange int) bool {
	for i := startRange; i < endRange; i++ {
		if nodemap[string(i)] == "present" {
			return false
		}
	}
	return true
}
func nodeConnectCheck(actscope string) bool {
	for key := range nodemap {
		// fmt.Println(key, val)
		if key == actscope {
			return true
		}
	}
	return false
}

// nAvailableNodeAdd : check
func nAvailableNodeAdd(actscope string) bool {
	// fmt.Println("qwe => ", nodemap["1"])
	fmt.Println("Now actscope=>", actscope)
	subnet := strings.Split(tokenaction.iprange[0], ".")
	if tokenaction.rangeoption {
		parseScope := strings.Split(actscope, ":")
		startRange, err := strconv.Atoi(parseScope[0])
		endRange, err := strconv.Atoi(parseScope[1])
		if err == nil {
			fmt.Println("Available node Can't parse")
			return false
		}
		//Compute node Is available?
		retry := 0

		for !isAllNodeOnline(startRange, endRange) {
			fmt.Println("Availabe Node Add retry : [", retry+1, "/10]")
			if retry > 10 {
				return false
			}
			for i := startRange; i < endRange; i++ {
				subnet[3] = nodemap[string(i)]
				if nodemap[string(i)] == "present" && verifyNPort(strings.Join(subnet, "."), "2222") {
					addNodes(string(i))
				}
			}
			retry++
		}
		return true
	} else {
		if nodeConnectCheck(actscope) {
			subnet[3] = actscope
			if nodemap[actscope] == "present" && verifyNPort(strings.Join(subnet, "."), "2222") {
				result := addNodes(actscope)
				fmt.Println("Action Result : ", result)
				return true

			}
		}

		return false
	}

	//Debug For nodeMap
	// for key, val := range nodemap {
	// 	if val == "present" && key != "1" {
	// 		return false
	// 	}
	// 	fmt.Println("Codex => ", key, val)
	// }
	//
}

func printOutput(outs string) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", outs)
	}
}

func kerrighedContainerVerify() (bool, error) {

	if fileExists("/proc/nodes/self/nodeid") {
		fmt.Println("Kerrighed Container load")
		return true, nil
	}

	return false, errors.New("Not Kerrighed Container")
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

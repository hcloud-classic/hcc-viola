package controlcli

import (
	"errors"
	"hcc/viola/lib/logger"
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

// var actiontype = []string{"area", "class", "scope"}
var tokenaction AtomicAction
var nodemap map[string]string

// HccCli : Hcc integration Command line interface
func HccCli(action string, iprange string) (bool, interface{}) {
	clearAction()
	logger.Logger.Println("Receive : ", action)
	err := ActionParser(action, iprange)
	if err != nil {
		return false, errors.New("ActionParcer Faild")
	}

	logger.Logger.Println(tokenaction.area, tokenaction.class, tokenaction.scope)
	iskerrighed, _ := kerrighedContainerVerify()
	//if err != nil {
	//	return false, err
	//}

	if iskerrighed {
		switch tokenaction.area {
		case "nodes":
			return cmdNodes(tokenaction.class, tokenaction.scope)
		case "cluster":
			cmdCluster(tokenaction.class, tokenaction.scope)
			break
		default:
			logger.Logger.Println("Please choose the area {nodes or cluster}")
		}
	} else {
		return false, errors.New("please proceed in Kerrighed container")
	}

	return false, nil
}

// ActionParser : Parcing Action
func ActionParser(action string, iprange string) interface{} {
	//Action parsing
	tmpstr := strings.Split(action, " ")
	tmplength := len(tmpstr)
	// logger.Logger.Println("action : ", action, "\n", "Length =>", tmplength, "++++++", tmpstr[1])

	tokenaction.area = tmpstr[1]
	if tmplength <= 3 && tmplength >= 2 {
		tokenaction.class = "status"
		// tokenaction.scope = append(tokenaction.scope, "0")
	} else if tmplength < 2 {
		return errors.New("invalid Hcc command line")
	} else {
		tokenaction.class = tmpstr[2]
		// tokenaction.scope = append(tokenaction.scope, tmpstr[3])
	}

	hasOption := strings.Contains(action, "-n")
	hasRangeOption := strings.Contains(action, ":")
	logger.Logger.Println("hasOption=> ", hasOption, "] hasRangeOption => ", hasRangeOption)
	//deliIndex : delimeter index
	var deliIndex = 0
	var endOfIndex = 0
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
		// 	logger.Logger.Println(i, "=>", words)
		// }
	}

	//Debug : tokenaction Structure
	logger.Logger.Println("area =>", tokenaction.area)
	logger.Logger.Println("class => ", tokenaction.class)
	logger.Logger.Println("scope => ", tokenaction.scope)
	logger.Logger.Println("iprange => ", tokenaction.iprange)
	logger.Logger.Println("rangeoption => ", tokenaction.rangeoption)

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
		}
		return false, verbosenode
	case "add":
		if checkNFS() {
			logger.Logger.Println("Leader Node NFS Service On")
		} else {
			restartNFS()
			logger.Logger.Println("Leader Node NFS Service restart")

		}
		//For nodeMap renewal
		nodeStatus("0")
		if nAvailableNodeAdd(actscope[0]) {
			return true, errors.New("all nodes are preparing with online")
		}
		return false, errors.New("all nodes are not preparing with online")
	case "del":
		// TODO : Add del operation
	default:
		return false, errors.New("please choose operation {status, add, del}")
	}
	return false, errors.New("not available command")
}

func cmdCluster(actclass string, actscope []string) {

}

func addNodes(actscope string) interface{} {
	cmd := exec.Command("krgadm", "nodes", "add", "-n", actscope)
	result, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println("Node Can't add the Num of [", actscope, "] Node")
		return err
	}
	return result
}

func checkNFS() bool {
	cmd := exec.Command("service", "nfs-common", "status")
	result, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println("NFS Service error")

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
		logger.Logger.Println("NFS Service Can't start")

	}

}

// @ N node nodeStatus index == n
// @ all node status index == 0
func nodeStatus(index string) (bool, interface{}) {
	cmd := exec.Command("krgadm", "nodes", "status")
	result, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println("Node status error occurred!!")
	}

	if index == "0" {
		logger.Logger.Println("HCC All Nodes Status \nIP  status\n", string(result))
		nodeStatusRegister(string(result))
		return true, string(result)
	}

	if nodeConnectCheck(index) {
		tmpstr := strings.Split(string(result), "\n")
		for _, words := range tmpstr {
			retoken := strings.Split(words, ":")
			if string(words[0]) == index {
				logger.Logger.Println(index, " th node status = > ", retoken[1])
				return true, retoken[1]
			}
		}
	} else {
		result := "[" + index + "] th Node Is Not in Cluster Area"
		return false, errors.New(result)
	}

	return false, nil
}

func nodeStatusRegister(status string) {
	nodemap = make(map[string]string)
	tmpstr := strings.Split(status, "\n")
	for _, words := range tmpstr {
		if strings.Contains(words, ":") {
			retoken := strings.Split(words, ":")
			// if string(words[0]) != "1" {}
			nodemap[string(words[0])] = retoken[1]

			logger.Logger.Println("words => ", string(words[0]), "retoken => ", retoken[1])
			logger.Logger.Println("register => ", nodemap[string(words[0])])

		}
	}
}

func isAllNodeOnline(startRange int, endRange int) bool {
	for i := startRange; i < endRange; i++ {
		if nodemap[string(i)] == "present" {
			logger.Logger.Println(nodemap[string(i)])
			return false
		}
	}
	return true
}

func nodeConnectCheck(actscope string) bool {
	for key := range nodemap {
		// logger.Logger.Println(key, val)
		if key == actscope {
			return true
		}
	}
	return false
}

func checkAllNodeOnline(startRange int, endRange int, subnet []string) bool {
	retry := 0

	for !isAllNodeOnline(startRange, endRange) {
		logger.Logger.Println("Available Node Add retry : [", retry+1, "/10]")
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
}

// nAvailableNodeAdd : check
func nAvailableNodeAdd(actscope string) bool {
	// logger.Logger.Println("qwe => ", nodemap["1"])
	logger.Logger.Println("Now actscope=>", actscope)
	subnet := strings.Split(tokenaction.iprange[0], ".")
	if tokenaction.rangeoption {
		parseScope := strings.Split(actscope, ":")
		startRange, err := strconv.Atoi(parseScope[0])
		if err != nil {
			logger.Logger.Println("Can't parse available node")
			return false
		}
		endRange, err := strconv.Atoi(parseScope[1])
		if err != nil {
			logger.Logger.Println("Can't parse available node")
			return false
		}
		//Compute node Is available?
		return checkAllNodeOnline(startRange, endRange, subnet)
	}

	if nodeConnectCheck(actscope) && actscope != "0" {
		subnet[3] = actscope
		if nodemap[actscope] == "present" && verifyNPort(strings.Join(subnet, "."), "2222") {
			result := addNodes(actscope)
			logger.Logger.Println("Action Result : ", result)
			return true

		}
	} else {
		start := strings.Split(tokenaction.iprange[0], ".")
		end := strings.Split(tokenaction.iprange[1], ".")
		logger.Logger.Println(start, "   ", end)
		startip, err := strconv.Atoi(start[3])
		if err != nil {
			logger.Logger.Println("Can't parse IP range")
			return false
		}
		endip, err := strconv.Atoi(end[3])
		if err != nil {
			logger.Logger.Println("Can't parse IP range")
			return false
		}
		logger.Logger.Println(startip, "   ", endip)

		return checkAllNodeOnline(startip, endip, subnet)
	}

	return false

	//Debug For nodeMap
	// for key, val := range nodemap {
	// 	if val == "present" && key != "1" {
	// 		return false
	// 	}
	// 	logger.Logger.Println("Codex => ", key, val)
	// }
	//
}

func printOutput(outs string) {
	if len(outs) > 0 {
		logger.Logger.Printf("==> Output: %s\n", outs)
	}
}

func kerrighedContainerVerify() (bool, error) {
	if fileExists("/proc/nodes/self/nodeid") {
		logger.Logger.Println("Kerrighed Container load")
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
// 	logger.Logger.Println(i, "= > ", words)
// }
// if err != nil {
// 	logger.Logger.Println("Error occurred!!")
// }

func verifyNPort(ip string, port string) bool {
	cmd := exec.Command("nmap", ip)
	result, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println(ip, " has not configure the ", port)
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

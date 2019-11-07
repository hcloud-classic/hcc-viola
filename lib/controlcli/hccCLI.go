package controlcli

import (
	"errors"
	"fmt"
	"hcc/viola/model"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

//AtomicAction : Parsing tjuhe
type AtomicAction struct {
	area        string
	class       string
	scope       []string
	iprange     []string
	rangeoption bool
	publisher   string
	receiver    string
}

var actiontype = []string{"area", "class", "scope"}
var tokenaction AtomicAction
var nodemap map[string]string

// HccCli : Hcc integration Command line interface
func HccCli(parseaction model.Control) (bool, interface{}) {
	clearAction()
	ActionClassify(parseaction)
	//Debug Option
	// fmt.Println("Receive : ", action)

	fmt.Println(tokenaction.area, tokenaction.class, tokenaction.scope)
	ishcccluster, err := hccContainerVerify()
	if err != nil {
		return false, errors.New("ActionParcer Faild")
	}

	if !ishcccluster {
		switch tokenaction.area {
		case "nodes":
			return cmdNodes()
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

func normalActionparser() interface{} {
	return nil
}

func hccActionparser(parseaction model.HccAction) interface{} {
	tokenaction.area = parseaction.ActionArea
	tokenaction.class = parseaction.ActionClass
	//ip range parse
	if strings.Contains(parseaction.HccIPRange, "range") {
		splitip := strings.Split(parseaction.HccIPRange, " ")
		tokenaction.iprange = append(tokenaction.iprange, splitip[1])
		tokenaction.iprange = append(tokenaction.iprange, splitip[2])
	} else {
		return errors.New("[hccActionparser] Invaild Ip rande, Failed parse iprange")
	}

	//Action effective scope parsing
	if parseaction.ActionScope != "" {
		if strings.Contains(parseaction.ActionScope, ":") {
			tokenaction.rangeoption = true
		}
		re := regexp.MustCompile("[0-9]+")
		extractscope := re.FindAllString(parseaction.ActionScope, -1)
		tokenaction.scope = extractscope

	} else {
		return errors.New("[hccActionparser] Invaild scope, Failed parse scope")
	}
	//Debug : tokenaction Structure
	fmt.Println("area =>", tokenaction.area)
	fmt.Println("class => ", tokenaction.class)
	fmt.Println("scope => ", tokenaction.scope)
	fmt.Println("iprange => ", tokenaction.iprange)

	return nil
}

// ActionClassify : Parcing Action
func ActionClassify(parsingmsg model.Control) interface{} {
	tokenaction.publisher = parsingmsg.Publisher
	tokenaction.receiver = parsingmsg.Receiver
	//Classify Action Type
	switch parsingmsg.Control.ActionType {
	case "hcc":
		err := hccActionparser(parsingmsg.Control.HccType)
		if err != nil {
			errstr := fmt.Sprintf("%v", err)
			return errors.New("[Hcc Action Parsing] Can't parse hcc action (" + errstr + ")")
		}
	case "normal":
		err := normalActionparser()
		if err != nil {
			errstr := fmt.Sprintf("%v", err)
			return errors.New("[Normal Action Parsing] Can't parse normal action (" + errstr + ")")
		}
	default:
		return errors.New("[Parsing Error]Please Correct Action type")
	}

	return nil
}

func clearAction() {
	tokenaction.area = ""
	tokenaction.class = ""
	tokenaction.scope = nil
	tokenaction.iprange = nil
	tokenaction.rangeoption = false
}
func cmdNodes() (bool, interface{}) {

	switch tokenaction.class {
	case "status":
		err, verbosenode := nodeStatus(tokenaction.scope[0])
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
		if nAvailableNodeAdd() {
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
			// fmt.Println(nodemap[string(i)])
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
func nAvailableNodeAdd() bool {

	// 0 =>  all node add
	//  x:y =>  x~y nodes add

	subnetstart := strings.Split(tokenaction.iprange[0], ".")
	subnetend := strings.Split(tokenaction.iprange[1], ".")
	//N number of  nodes add
	if tokenaction.rangeoption {
		startRange, err := strconv.Atoi(tokenaction.scope[0])
		endRange, err := strconv.Atoi(tokenaction.scope[1])
		if err != nil {
			fmt.Println("Available node Can't parse")
			return false
		}
		//Compute node Is available?
		retry := 0
		fmt.Println("startRange : ", startRange, " | endRange  ", endRange, " | subnet : ", subnetstart)
		for i := 0; i < len(subnetstart); i++ {
			fmt.Println("subnet[", i, "]  ", subnetstart[i])
		}
		for !isAllNodeOnline(startRange, endRange) {
			fmt.Println("Availabe Node Add retry : [", retry+1, "/10]")
			if retry > 10 {
				return false
			}
			for i := startRange; i < endRange; i++ {
				subnetstart[3] = nodemap[string(i)]
				if nodemap[string(i)] == "present" && verifyNPort(strings.Join(subnetstart, "."), "2222") {
					addNodes(string(i))
				}
			}
			retry++
		}
		return true
	} else {
		// Specific the number node add
		if nodeConnectCheck(tokenaction.scope[0]) && tokenaction.scope[0] != "0" {
			subnetstart[3] = tokenaction.scope[0]
			if nodemap[tokenaction.scope[0]] == "present" && verifyNPort(strings.Join(subnetstart, "."), "2222") {
				result := addNodes(tokenaction.scope[0])
				fmt.Println("Action Result : ", result)
				return true

			}
		} else {
			// range option is zero, Add all nodes
			fmt.Println(subnetstart, "   ", subnetend)
			startip, err := strconv.Atoi(subnetstart[3])
			endip, err := strconv.Atoi(subnetend[3])
			fmt.Println(startip, " to ", endip)

			if err != nil {

			}
			retry := 0
			for !isAllNodeOnline(startip, endip) {
				fmt.Println("Availabe Node Add retry : [", retry+1, "/10]")
				if retry > 10 {
					return false
				}
				for i := startip; i < endip; i++ {
					subnetstart[3] = nodemap[string(i)]
					if nodemap[string(i)] == "present" && verifyNPort(strings.Join(subnetstart, "."), "2222") {
						addNodes(string(i))
					}
				}
				retry++
			}

			return true
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

func hccContainerVerify() (bool, error) {

	if fileExists("/proc/nodes/self/nodeid") {
		fmt.Println("Hcloud Container load")
		return true, nil
	}

	return false, errors.New("Not Hcloud Container")
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

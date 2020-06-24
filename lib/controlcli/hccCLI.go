package controlcli

import (
	"errors"
	"fmt"
	"hcc/viola/lib/logger"
	"hcc/viola/model"
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

var tokenaction AtomicAction

// var nodemap map[string]string
var nodemap = make(map[string]string)

// HccCli : Hcc integration Command line interface
func HccCli(parseaction model.Control) (bool, interface{}) {
	clearAction()
	ActionClassify(parseaction)
	logger.Logger.Println(tokenaction.area, tokenaction.class, tokenaction.scope)
	ishcccluster, err := hccContainerVerify()
	if err != nil {
		return false, errors.New("ActionParcer Faild")
	}
	istelegrafset, checker := telegrafSetting(parseaction)
	if !istelegrafset {
		logger.Logger.Println(checker)
	}

	if ishcccluster {
		switch tokenaction.area {
		case "nodes":
			return cmdNodes()
		case "cluster":
			cmdCluster(tokenaction.class, tokenaction.scope)
		default:
			logger.Logger.Println("Please choose the area {nodes or cluster}")
		}
	} else {
		return false, errors.New("Please Continue in Hcloud Container")
	}

	return false, nil
}

func normalActionparser() interface{} {
	return nil
}

func isipv4(host string) bool {
	parts := strings.Split(host, ".")

	if len(parts) < 4 {
		return false
	}
	for _, x := range parts {
		if i, err := strconv.Atoi(x); err == nil {
			if i < 0 || i > 255 {
				return false
			}
		} else {
			return false
		}

	}
	return true
}

func hccActionparser(parseaction model.HccAction) interface{} {
	tokenaction.area = parseaction.ActionArea
	tokenaction.class = parseaction.ActionClass
	splitip := strings.Split(parseaction.HccIPRange, " ")
	if isipv4(splitip[0]) && isipv4(splitip[1]) {

		tokenaction.iprange = append(tokenaction.iprange, splitip[0])
		tokenaction.iprange = append(tokenaction.iprange, splitip[1])
	} else {
		return errors.New("[hccActionparser] Invaild Ip range, Failed parse iprange")
	}

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
	return nil
}

// ActionClassify : Parcing Action
func ActionClassify(parsingmsg model.Control) interface{} {
	logger.Logger.Println("Receive : ", parsingmsg)
	tokenaction.publisher = parsingmsg.Publisher
	tokenaction.receiver = parsingmsg.Receiver
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
		}
		return false, verbosenode
	case "add":
		if checkNFS() {
			logger.Logger.Println("Leader Node NFS Service On")
		} else {
			restartService("nfs-common")
			logger.Logger.Println("Leader Node NFS Service restart")
		}
		//For nodeMap renewal
		nodeStatus("0")
		if nAvailableNodeAdd() {
			return true, errors.New("All Nodes is Preparing and online")
		}
		return false, errors.New("Some Nodes is Not Preparing")
	case "del":
	default:
		return false, errors.New("Please Choose Operation {status, add, del}")
	}
	return false, errors.New("Not Available Command")
}

func cmdCluster(actclass string, actscope []string) {

}

func addNodes(actscope string) interface{} {
	cmd := exec.Command("hccadm", "nodes", "add", "-n", actscope)
	result, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println("Node Can't add the Num of [", actscope, "] Node")
		return err
	}
	nodemap[actscope] = "online"

	return string(result)
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

func nodeStatus(index string) (bool, interface{}) {
	cmd := exec.Command("hccadm", "nodes", "status")
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
			retoken := strings.Split(string(words), ":")
			if string(words[0]) == index {
				logger.Logger.Println(index, " th node status = > ", retoken[1])
				return true, string(retoken[1])
			}
		}
	} else {
		result := "[" + index + "] th Node Is Not in Cluster Area"
		return false, errors.New(result)
	}

	return false, nil
}

func nodeStatusRegister(status string) {
	tmpstr := strings.Split(status, "\n")
	for _, words := range tmpstr {
		if strings.Contains(string(words), ":") {
			retoken := strings.Split(string(words), ":")
			nodemap[string(words[0])] = retoken[1]
		}
	}
}

func isAllNodeOnline(startRange int, endRange int) bool {
	var needednode = endRange - startRange + 1
	var count = 0
	for i := startRange; i <= endRange; i++ {
		if nodemap[strconv.Itoa(i)] == "present" {
			return false
		}
		if nodemap[strconv.Itoa(i)] == "online" {
			count++
		}
	}
	logger.Logger.Println("needednode : ", needednode, " || count : ", count)
	if needednode == count {
		return true
	}

	return false
}

func nodeConnectCheck(actscope string) bool {
	for key := range nodemap {
		if key == actscope {
			return true
		}
	}
	return false
}

func rangeAtoiParse(start string, end string) (int, int, error) {
	startInt, err := strconv.Atoi(start)
	if err != nil {
		return 0, 0, errors.New("failed to parse start")
	}

	endInt, err := strconv.Atoi(end)
	if err != nil {
		return 0, 0, errors.New("failed to parse end")
	}

	return startInt, endInt, nil
}

func nAvailableNodeAdd() bool {
	subnetstart := strings.Split(tokenaction.iprange[0], ".")
	subnetend := strings.Split(tokenaction.iprange[1], ".")
	if tokenaction.rangeoption {
		startRange, endRange, err := rangeAtoiParse(tokenaction.scope[0], tokenaction.scope[1])
		if err != nil {
			logger.Logger.Println("Available node Can't parse")
			return false
		}

		retry := 0
		logger.Logger.Println("startRange : ", startRange, " | endRange  ", endRange, " | subnet : ", subnetstart)

		for !isAllNodeOnline(startRange, endRange) {
			logger.Logger.Println("Available Node Add retry : [", retry+1, "/100]")
			nodeStatus("0")
			if retry > 100 {
				return false
			}
			for i := startRange; i <= endRange; i++ {
				subnetstart[3] = strconv.Itoa(i)
				logger.Logger.Println(nodeVerifyAdd(strconv.Itoa(i), subnetstart))
			}
			time.Sleep(4 * time.Second)
			retry++
		}
		return true
	}
	if nodeConnectCheck(tokenaction.scope[0]) && tokenaction.scope[0] != "0" {
		retry := 0
		specificnode, err := strconv.Atoi(tokenaction.scope[0])
		if err != nil {
			return false
		}
		for !isAllNodeOnline(specificnode, specificnode) {
			if retry > 100 {
				return false
			}
			logger.Logger.Println("Available Node Add retry : [", retry+1, "/100]")
			nodeStatus("0")
			subnetstart[3] = tokenaction.scope[0]
			logger.Logger.Println(nodeVerifyAdd(tokenaction.scope[0], subnetstart))
			retry++
			time.Sleep(4 * time.Second)
		}
		return true
	}
	startip, endip, err := rangeAtoiParse(subnetstart[3], subnetend[3])
	if err != nil {
		logger.Logger.Println("Available node Can't parse")
		return false
	}
	retry := 0
	for !isAllNodeOnline(startip, endip) {
		if retry > 100 {
			return false
		}
		logger.Logger.Println("Available Node Add retry : [", retry+1, "/100]")
		nodeStatus("0")
		for i := startip; i <= endip; i++ {
			subnetstart[3] = strconv.Itoa(i)
			logger.Logger.Println(nodeVerifyAdd(strconv.Itoa(i), subnetstart))
		}
		retry++
		time.Sleep(4 * time.Second)
	}
	return true
}
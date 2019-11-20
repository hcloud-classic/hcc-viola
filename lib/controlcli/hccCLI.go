package controlcli

import (
	"errors"
	"fmt"
	"hcc/viola/lib/config"
	"hcc/viola/lib/logger"
	"hcc/viola/model"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
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

// var actiontype = []string{"area", "class", "scope"}
var tokenaction AtomicAction

// var nodemap map[string]string
var nodemap = make(map[string]string)

// HccCli : Hcc integration Command line interface
func HccCli(parseaction model.Control) (bool, interface{}) {
	clearAction()
	ActionClassify(parseaction)
	//Debug Option
	// logger.Logger.Println("Receive : ", parseaction)

	logger.Logger.Println(tokenaction.area, tokenaction.class, tokenaction.scope)
	ishcccluster, err := hccContainerVerify()
	if err != nil {
		return false, errors.New("ActionParcer Faild")
	}
	istelegrafset, checkerr := telegrafSetting(parseaction)
	if !istelegrafset {
		logger.Logger.Println(checkerr)
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
	// logger.Logger.Println("hccActionparser : ", parseaction, "  parseaction.ActionArea : ", parseaction.ActionArea)
	tokenaction.area = parseaction.ActionArea
	tokenaction.class = parseaction.ActionClass
	//ip range parse
	splitip := strings.Split(parseaction.HccIPRange, " ")
	if isipv4(splitip[0]) && isipv4(splitip[1]) {

		tokenaction.iprange = append(tokenaction.iprange, splitip[0])
		tokenaction.iprange = append(tokenaction.iprange, splitip[1])
	} else {
		return errors.New("[hccActionparser] Invaild Ip range, Failed parse iprange")
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
	// logger.Logger.Println("area =>", tokenaction.area)
	// logger.Logger.Println("class => ", tokenaction.class)
	// logger.Logger.Println("scope => ", tokenaction.scope)
	// logger.Logger.Println("iprange => ", tokenaction.iprange)

	return nil
}

// ActionClassify : Parcing Action
func ActionClassify(parsingmsg model.Control) interface{} {
	logger.Logger.Println("Receive : ", parsingmsg)
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

// @ N node nodeStatus index == n
// @ all node status index == 0
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

			// logger.Logger.Println("words => ", string(words[0]), "retoken => ", retoken[1])
			// logger.Logger.Println("register => ", nodemap[string(words[0])])

		}
	}
}

func isAllNodeOnline(startRange int, endRange int) bool {
	var needednode = endRange - startRange + 1
	var count = 0
	for i := startRange; i <= endRange; i++ {
		if nodemap[strconv.Itoa(i)] == "present" {
			//For Debug
			// logger.Logger.Println("i : [", i, "] => ", nodemap[strconv.Itoa(i)])
			return false
		}
		if nodemap[strconv.Itoa(i)] == "online" {
			//For Debug
			// logger.Logger.Println("i : [", i, "] => ", nodemap[strconv.Itoa(i)])
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
		//For Debug
		// logger.Logger.Println(key, val)
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

// nAvailableNodeAdd : check
func nAvailableNodeAdd() bool {
	// 0 =>  all node add
	//  x:y =>  x~y nodes add

	subnetstart := strings.Split(tokenaction.iprange[0], ".")
	subnetend := strings.Split(tokenaction.iprange[1], ".")
	//N number of  nodes add
	if tokenaction.rangeoption {
		startRange, endRange, err := rangeAtoiParse(tokenaction.scope[0], tokenaction.scope[1])
		if err != nil {
			logger.Logger.Println("Available node Can't parse")
			return false
		}

		//Compute node Is available?
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
	// Specific the number node add
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

	// logger.Logger.Println(subnetstart, "   ", subnetend)
	startip, endip, err := rangeAtoiParse(subnetstart[3], subnetend[3])
	if err != nil {
		logger.Logger.Println("Available node Can't parse")
		return false
	}
	// logger.Logger.Println(startip, " to ", endip)

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
	// Debug For nodeMap
	// for key, val := range nodemap {
	// 	// if val == "present" && key != "1" {
	// 	// 	return false
	// 	// }
	// 	logger.Logger.Println("Codex => ", key, val)
	// }
	// range option is zero, Add all nodes
	return true
}

func hccContainerVerify() (bool, error) {

	if fileExists("/proc/nodes/self/nodeid") {
		logger.Logger.Println("Hcloud Container load")
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

func verifyNPort(ip string, port string) bool {
	cmd := exec.Command("nmap", ip)
	result, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println(ip, " has not configure the ", port)
	}
	//compute on?
	if strings.Contains(string(result), port) {
		logger.Logger.Println(ip, " : ", port, "Connect")

		return true
	}

	return false
}

func nodeVerifyAdd(mapnum string, subnetstart []string) interface{} {
	if nodemap[mapnum] == "present" && verifyNPort(strings.Join(subnetstart, "."), "2222") {
		result := addNodes(mapnum)
		//For Debug
		// logger.Logger.Println("Action Result : ", result)
		return result
	}
	return "Faild Add Node" + mapnum
}

//TelegrafCheck :telegraf config
func TelegrafCheck() (bool, interface{}) {
	if !fileExists(telegrafDir + "/telegraf.conf") {
		return false, errors.New("Telegraf setting is failed, Please check " + telegrafDir + "/telegraf.conf")
	}
	return true, "Telegraf Config Exist!\n"
}

func telegrafSetting(parseaction model.Control) (bool, interface{}) {
	state, err := TelegrafCheck()
	if !state {
		strtmp := fmt.Sprintf("%v", err)
		return false, errors.New(strtmp)
	}
	b, err := ioutil.ReadFile(telegrafDir + "/telegraf.conf") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	r, _ := regexp.Compile(parseaction.Control.HccType.ServerUUID)
	if r.MatchString(string(b)) {
		return true, "Already setting complete\n"
	}

	teleconf := agent + outputsInfluxdb + cpuInfo + inputsDisk + etcSet
	teleconf = strings.Replace(teleconf, "SERVER_UUID", parseaction.Control.HccType.ServerUUID, -1)
	teleconf = strings.Replace(teleconf, "INFLUX_DB_IP", config.InfluxDB.IP, -1)
	teleconf = strings.Replace(teleconf, "PORT", config.InfluxDB.Port, -1)
	err = ioutil.WriteFile(telegrafDir+"/telegraf.conf", []byte(teleconf), 644)
	if err != nil {
		return false, errors.New("failed to write to telegraf.conf file")
	}
	restartService("telegraf")
	return true, "Telegraf Setting is complete!!\n"
}

func restartService(servname string) {
	cmd := exec.Command("service", servname, "restart")
	_, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println(servname + " Service Can't start")

	}

}

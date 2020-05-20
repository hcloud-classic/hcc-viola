package controlcli

import (
	"errors"
	"fmt"
	"hcc/viola/lib/logger"
	"hcc/viola/model"
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

	return false, nil
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

	return nil
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

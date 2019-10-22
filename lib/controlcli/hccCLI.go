package controlcli

import (
	"fmt"
	"os/exec"
	"strings"
)

// HccCLI :
func HccCLI(area string, action string) {

	switch area {
	case "nodes":
		cmdNodes(action)
	case "cluster":
		cmdCluster(action)
	default:
		fmt.Println("Please choose the area {nodes or cluster}")

	}

}
func cmdNodes(action string) {
	nodeStatus()
	// switch action {
	// 	case "add"

	// }
}

func cmdCluster(action string) {

}

func nodeStatus() {
	cmd := exec.Command("krgadm", "nodes", "status")
	result, err := cmd.CombinedOutput()
	printOutput(string(result))
	if err != nil {
		fmt.Println("qweqweqwe")
	}
}

func printOutput(outs string) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", outs)
	}
}

func extractToken(str string, num int) {
	tmp := strings.Split(str, " ")

	fmt.Println(tmp[1])
}

package controlcli

var hostFilePath = "/etc/hosts"
var nfsFilePath = "/etc/exports"

var hostContents = "127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4\n" +
	"::1       CLUSTERHOSTNAME localhost localhost.localdomain localhost6 localhost6.localdomain6\n" +
	"CLUSTERADDRESS CLUSTERHOSTNAME LEADERNODEHOSTNAME\n" +
	"COMPUTE\n"

var nfsContents = "PATH CLUSTERADDRESS/NETMASK(rw,nohide,sync,no_subtree_check,no_root_squash,insecure)\n"

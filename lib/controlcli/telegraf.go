package controlcli

var telegrafDir = "/etc/telegraf"

var globalTags = "[global_tags]\n"

var agent = "[agent]\n" +
	"interval = \"10s\"\n" +
	"round_interval = true\n" +
	"metric_batch_size = 1000\n" +
	"metric_buffer_limit = 10000\n" +
	"collection_jitter = \"0s\"\n" +
	"flush_interval = \"10s\"\n" +
	"flush_jitter = \"0s\"\n" +
	"precision = \"\"\n" +
	"hostname = \"SERVER_UUID\"\n" +
	"omit_hostname = false\n"

var outputsInfluxdb = "[[outputs.influxdb]]\n" +
	"# Address of influxdb\n" +
	"urls = [\"http://INFLUX_DB_IP:PORT\"]\n" +
	"database = \"SERVER_UUID\"\n" +
	"skip_database_creation = false\n"

var cpuInfo = "[[inputs.cpu]]\n" +
	"percpu = true\n" +
	"totalcpu = true\n" +
	"collect_cpu_time = false\n" +
	"report_active = false\n"

var inputsDisk = "[[inputs.disk]]\n" +
	"ignore_fs = [\"tmpfs\", \"devtmpfs\", \"devfs\", \"iso9660\", \"overlay\", \"aufs\", \"squashfs\"]\n"

var etcSet = "# Info detail setting\n" +
	"[[inputs.diskio]]\n" +
	"[[inputs.kernel]]\n" +
	"[[inputs.mem]]\n" +
	"[[inputs.processes]]\n" +
	"[[inputs.swap]]\n" +
	"[[inputs.system]]"

package config

type viola struct {
	ServerAddress    string `goconf:"violin:violin_server_address"`     // ServerAddress : IP address of server which installed violin module
	ServerPort       int64  `goconf:"violin:violin_server_port"`        // ServerPort : Listening port number of violin module
	RequestTimeoutMs int64  `goconf:"violin:violin_request_timeout_ms"` // RequestTimeoutMs : HTTP timeout for GraphQL request to violin module
}

// Viola : violin config structure
var Viola viola

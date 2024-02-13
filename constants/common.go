package constants

const (
	TEZPEAK_VERSION  = "<VERSION>"
	TEZPEAK_CODENAME = "<CODENAME>"

	DEFAULT_LISTEN_ADDRESS = "localhost:8733"

	ENV_TEZBAKE_HOME        = "TEZBAKE_HOME"
	ENV_TEZPEAK_CONFIG_FILE = "TEZPEAK_CONFIG_FILE"

	MAX_SERVICES_REFRESH_INTERVAL = 300 // 5 minutes
	MIN_SERVICES_REFRESH_INTERVAL = 5   // 5 seconds
	DEFAULT_BAKER_NODE_URL        = "http://localhost:8732"

	DEFAULT_REFERENCE_NODE_URL                = "https://rpc.tzbeta.net/"
	DEFAULT_REFERENCE_NODE_IS_RIGHTS_PROVIDER = true
	DEFAULT_REFERENCE_NODE_IS_BLOCK_PROVIDER  = false

	DEFAULT_REFERENCE_NODE_2_URL                = "https://rpc.tzkt.io/mainnet/"
	DEFAULT_REFERENCE_NODE_2_IS_RIGHTS_PROVIDER = false
	DEFAULT_REFERENCE_NODE_2_IS_BLOCK_PROVIDER  = true
)

var (
	PRIVATE_NETWORK_HOSTS = []string{
		"localhost",
		"127.0.0.1",
		"::1",
	}
)

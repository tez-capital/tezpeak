package constants

const (
	TEZPEAK_VERSION  = "<VERSION>"
	TEZPEAK_CODENAME = "<CODENAME>"

	DEFAULT_LISTEN_ADDRESS       = "localhost:8733"
	DEFAULT_HTTP_TIMEOUT_SECONDS = 30

	// tezbake
	TEZBAKE_MODULE_ID             = "tezbake"
	ENV_TEZPEAK_CONFIG_FILE       = "TEZPEAK_CONFIG_FILE"
	MAX_SERVICES_REFRESH_INTERVAL = 300 // 5 minutes
	MIN_SERVICES_REFRESH_INTERVAL = 5   // 5 seconds
	DEFAULT_NODE_APP_PATH         = "node"
	DEFAULT_SIGNER_APP_PATH       = "signer"
	DEFAULT_BAKER_NODE_URL        = "http://localhost:8732"
	DEFAULT_BAKER_SIGNER_URL      = "http://localhost:20090"
	DEFAULT_RIGHTS_BLOCK_WINDOW   = 50

	// tezpay
	TEZPAY_MODULE_ID        = "tezpay"
	DEFAULT_TEZPAY_APP_PATH = "pay"

	// tx constants
	MAX_OPERATION_TTL         = 12
	MAX_WAIT_FOR_CONFIRMATION = 120
)

var (
	PRIVATE_NETWORK_HOSTS = []string{
		"localhost",
		"127.0.0.1",
		"::1",
	}
)

package constants

const (
	TEZPEAK_VERSION = "0.0.1"

	ENV_TEZBAKE_HOME = "TEZBAKE_HOME"

	MAX_SERVICES_REFRESH_INTERVAL = 300 // 5 minutes
	MIN_SERVICES_REFRESH_INTERVAL = 5   // 5 seconds
	DEFAULT_BAKER_NODE_URL        = "http://localhost:8732"
	BAKER_NODE_ID                 = "#baker"

	DEFAULT_REFERENCE_NODE_URL                = "https://rpc.tzbeta.net/"
	DEFAULT_REFERENCE_NODE_IS_RIGHTS_PROVIDER = true
	DEFAULT_REFERENCE_NODE_IS_BLOCK_PROVIDER  = false

	DEFAULT_REFERENCE_NODE_2_URL                = "https://rpc.tzkt.io/mainnet/"
	DEFAULT_REFERENCE_NODE_2_IS_RIGHTS_PROVIDER = false
	DEFAULT_REFERENCE_NODE_2_IS_BLOCK_PROVIDER  = true
)

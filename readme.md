## TEZPEAK 

Simple monitoring interface for tezos bakers.

### Setup

Right now there are two supported ways to run tezpeak:
- as module of tezbake
- as a standalone server with ami

#### tezbake

Since tezbake 0.13.0-alpha the tezpeak is natively supported module. You can setup it in 3 simple steps:
1. `tezbake setup --peak`
2. Adjust configuration as needed
3. `tezbake start --peak`

Sample minimal configuration:
```hjson
{
	listen: "0.0.0.0:8733"
	bakers: [
		tz1P6WKJu2rcbxKiKRZHKQKmKrpC9TfW1AwM
		tz1hZvgjekGo7DmQjWh7XnY5eLQD8wNYPczE
	]
}
```

#### Standalone

You can run tezpeak as a standalone server as a binary or on linux as a service with [ami-tezpeak](https://github.com/tez-capital/ami-tezpeak).

##### As binary

1. Download the latest release from the [releases page](https://github.com/tez-capital/tezpeak)
2. Add configuration file `config.hjson` to the same directory as the binary
- Sample standalone minimal configuration:
```hjson
{
	listen: "127.0.0.1:8733"
	node: http://localhost:8732
	bakers: [
		tz1P6WKJu2rcbxKiKRZHKQKmKrpC9TfW1AwM
		tz1hZvgjekGo7DmQjWh7XnY5eLQD8wNYPczE
	]
	providers: {
		// right now only tezbake is supported, so disable it altogether with 'none' to suppress the warnings
		services: none 
	}
}
```
3. Run the binary
4. Open the browser and navigate to `http://localhost:8733`
5. Enjoy

##### As an ami based service

Refer to the [ami-tezpeak readme](https://github.com/tez-capital/ami-tezpeak) for the installation and usage instructions.

### Configuration

The configuration is stored in `config.hjson` file. The default configuration is:

```hjson
{
	# Id to show in the header
  	id: ""
	# Address to listen on
  	listen: 127.0.0.1:8733
    modules: {
        tezbake: {
			# uncomment bellow to disable tezbake package monitoring
            # applications: null
            bakers: [
				# list of bakers to monitor for balances and rights
                tz1P6WKJu2rcbxKiKRZHKQKmKrpC9TfW1AwM
            ]
        }
        tezpay: {
			# can be null to disable tezpay package monitoring
			applications: {
				# path to tezpay ami package, either absolute or relative to parent directory peak
				tezpay: tezpay
			}
			payout_wallet: tz1X7U9XxVz6NDxL4DSZhijME61PW45bYUJE
            payout_wallet_preferences: {
                balance_warning_threshold: 100
                balance_error_threshold: 50
            }
			# forces all operations to be dry run
			force_dry_run: true
		}
    }
	
	# List of reference nodes to connect to
	# The reference nodes are used to get the rights and blocks if the baker's node is not available
	#nodes: {
	#	"Tezos Foundation": {
	#		address: https://rpc.tzbeta.net/
	#		is_rights_provider: true
	#		is_block_provider: false
	#	},
	#	"tzkt": {
	#		address: https://rpc.tzkt.io/mainnet/
	#		is_rights_provider: false
	#		is_block_provider: true
	#       # reports error if node not available
	#       is_essential: false
	#	}
	#}
	# The mode tezpeak should operate in
	# auto - if bound to localhost, it will operate in private mode if not, it will operate in public mode
	# public - assumes public environment, only readonly operations are allowed
	# private - assumes private environment, all operations are allowed
	mode: auto
}
``` 
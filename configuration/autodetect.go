package configuration

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"slices"

	"github.com/hjson/hjson-go/v4"
	"github.com/tez-capital/tezbake/ami"
	"github.com/tez-capital/tezpay/configuration"
	"github.com/tez-capital/tezpeak/constants"

	signer_engines "github.com/tez-capital/tezpay/engines/signer"
	"github.com/tez-capital/tezpay/state"
)

/*
	tezpay: {
		applications: {
				tezpay: tezpay
		}
		payout_wallet: tz1X7U9XxVz6NDxL4DSZhijME61PW45bYUJE
		payout_wallet_preferences: {
			balance_warning_threshold: 100
			balance_error_threshold: 50
		}
	}
*/
func autoDetectTezpayConfiguration(rootDir string) (json.RawMessage, error) {
	if rootDir == "" {
		return nil, errors.New("rootDir is empty")
	}

	if !ami.IsAppInstalled(path.Join(rootDir, constants.DEFAULT_TEZPAY_APP_PATH)) {
		return nil, errors.New("tezpay not found, skipping")
	}

	config, err := configuration.Load()
	if err != nil {
		return nil, errors.Join(errors.New("failed to load tezpay configuration"), err)
	}

	signerEngine := state.Global.SignerOverride
	if signerEngine == nil {
		signerEngine, err = signer_engines.Load(string(config.PayoutConfiguration.WalletMode))
		if err != nil {
			return nil, errors.Join(errors.New("failed to to load tezpay signer"), err)
		}
	}

	slog.Info("Checking PKH used for tezpay payouts...")
	pkh := signerEngine.GetPKH()

	tezpayModuleConfiguration := &TezpayModuleConfiguration{
		moduleConfigurationbase: moduleConfigurationbase{
			Applications: map[string]string{
				"tezpay": constants.DEFAULT_TEZPAY_APP_PATH,
			},
		},
		PayoutWallet: pkh.String(),
		PayoutWalletPreferences: PayoutWalletPreferences{
			BalanceWarningThreshold: 100,
			BalanceErrorThreshold:   50,
		},
	}
	return hjson.MarshalWithOptions(tezpayModuleConfiguration, hjson.DefaultOptions())
}

/*
	{
		"configuration": {
				...
				"additional_key_aliases": [ "key" ]
		},
	}
*/
type nodeAppJsonPartialConfiguration struct {
	RemoteSignerUrl      string   `json:"REMOTE_SIGNER_ADDR,omitempty"`
	AdditionalKeyAliases []string `json:"additional_key_aliases,omitempty"`
}

type nodeAppJsonPartial struct {
	Configuration nodeAppJsonPartialConfiguration
}

/*
	tezbake: {
		bakers: [
			tz1P6WKJu2rcbxKiKRZHKQKmKrpC9TfW1AwM
			tz1hZvgjekGo7DmQjWh7XnY5eLQD8wNYPczE
		]
	}
*/
func autoDetectTezbakeConfiguration(rootDir string) (json.RawMessage, error) {
	nodeAppPath := path.Join(rootDir, constants.DEFAULT_NODE_APP_PATH)
	signerAppPath := path.Join(rootDir, constants.DEFAULT_SIGNER_APP_PATH)

	applications := map[string]string{}
	bakers := []string{}
	remoteSignerUrl := constants.DEFAULT_BAKER_SIGNER_URL
	if ami.IsAppInstalled(nodeAppPath) {
		applications["node"] = constants.DEFAULT_NODE_APP_PATH

		keys := []string{"baker"}

		nodeAppConfig := nodeAppJsonPartial{
			Configuration: nodeAppJsonPartialConfiguration{
				AdditionalKeyAliases: []string{},
			},
		}

		nodeAppPath := path.Join(rootDir, constants.DEFAULT_NODE_APP_PATH, "app.json")
		fileContent := []byte("{}")
		if data, err := os.ReadFile(nodeAppPath); err == nil {
			fileContent = data
		} else {
			nodeAppPath = path.Join(rootDir, constants.DEFAULT_NODE_APP_PATH, "app.hjson")
			if data, err := os.ReadFile(nodeAppPath); err == nil {
				fileContent = data
			}
		}
		if err := hjson.Unmarshal(fileContent, &nodeAppConfig); err == nil {
			keys = append(keys, nodeAppConfig.Configuration.AdditionalKeyAliases...)
		}

		if nodeAppConfig.Configuration.RemoteSignerUrl != "" {
			remoteSignerUrl = nodeAppConfig.Configuration.RemoteSignerUrl
		}

		// read pkhs based on aliases
		pathToPkhs := path.Join(rootDir, constants.DEFAULT_NODE_APP_PATH, "data/.tezos-client/public_key_hashs")
		pkhs := nodePublicKeys{}
		fileContent = []byte{}
		if data, err := os.ReadFile(pathToPkhs); err == nil {
			fileContent = data
		}
		if err := hjson.Unmarshal(fileContent, &pkhs); err != nil {
			return nil, errors.New("failed to read public key hashes")
		}
		for _, pkh := range pkhs {
			if slices.Contains(keys, pkh.Name) {
				bakers = append(bakers, pkh.Hash)
			}
		}
	} else {
		slog.Warn("Node app is not found, skipping")
	}

	if ami.IsAppInstalled(signerAppPath) {
		applications["signer"] = constants.DEFAULT_SIGNER_APP_PATH
	} else {
		slog.Warn("Signer app is not found, skipping")
	}

	tezbakeModuleConfiguration := &TezbakeModuleConfiguration{
		moduleConfigurationbase: moduleConfigurationbase{
			Applications: applications,
		},
		SignerUrl:         remoteSignerUrl,
		RightsBlockWindow: constants.DEFAULT_RIGHTS_BLOCK_WINDOW,
		Bakers:            bakers,
	}

	return hjson.MarshalWithOptions(tezbakeModuleConfiguration, hjson.DefaultOptions())
}

func AutoDetect(rootDir string, destinationFile string) {
	modules := map[string]json.RawMessage{}
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		slog.Warn("Failed to get absolute path to root dir", "error", err.Error())
	} else {
		rootDir = absRootDir
	}

	tezpayConfig, err := autoDetectTezpayConfiguration(rootDir)
	if err != nil {
		slog.Warn("Failed to auto-detect tezpay configuration", "error", err.Error())
	} else {
		modules[constants.TEZPAY_MODULE_ID] = tezpayConfig
	}

	tezbakeConfig, err := autoDetectTezbakeConfiguration(rootDir)
	if err != nil {
		slog.Warn("Failed to auto-detect tezbake configuration", "error", err.Error())
	} else {
		modules[constants.TEZBAKE_MODULE_ID] = tezbakeConfig
	}

	config := Runtime{
		Id:      "",
		AppRoot: rootDir,
		Listen:  constants.DEFAULT_LISTEN_ADDRESS,
		Modules: modules,
		Mode:    AutoPeakMode,
	}

	if data, err := hjson.MarshalWithOptions(config, hjson.DefaultOptions()); err == nil {
		if err := os.WriteFile(destinationFile, data, 0644); err != nil {
			slog.Error("Failed to write autodetected configuration", "error", err.Error())
		}
	} else {
		slog.Error("Failed to marshal autodetected configuration", "error", err.Error())
	}

}

package configuration

import (
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hjson/hjson-go/v4"
	"github.com/tez-capital/tezpeak/constants"
	"github.com/trilitech/tzgo/tezos"
	"golang.org/x/mod/semver"
)

type TezbakeModuleConfiguration struct {
	moduleConfigurationbase

	SignerUrl         string   `json:"signer_url"`
	RightsBlockWindow int64    `json:"rights_block_window"`
	Bakers            []string `json:"bakers"`
	LedgerWallets     []string `json:"ledger_wallets"`
	ArcBinaryPath     string   `json:"arc_binary_path"`
}

func getDefaultTezbakeModuleConfiguration() *TezbakeModuleConfiguration {
	return &TezbakeModuleConfiguration{
		moduleConfigurationbase: moduleConfigurationbase{
			Applications: map[string]string{
				"node":   constants.DEFAULT_NODE_APP_PATH,
				"signer": constants.DEFAULT_SIGNER_APP_PATH,
			},
		},
		SignerUrl:         constants.DEFAULT_BAKER_SIGNER_URL,
		RightsBlockWindow: constants.DEFAULT_RIGHTS_BLOCK_WINDOW,
		// ledger
		LedgerWallets: constants.DEFAULT_LEDGER_WALLETS,
		ArcBinaryPath: constants.DEFAULT_ARC_BINARY_PATH,
	}
}

type nodeAppJsonConfiguration struct {
	AdditionalKeysAliases []string `json:"additional_keys_aliases,omitempty"`
}

type nodeAppJson struct {
	Configuration nodeAppJsonConfiguration `json:"configuration,omitempty"`
}

type nodePublicKeyHashAlias struct {
	Name string `json:"name,omitempty"`
	Hash string `json:"value,omitempty"`
}

type nodePublicKeys []nodePublicKeyHashAlias

func (configuration *TezbakeModuleConfiguration) GetNodeAppPath() string {
	path, ok := configuration.Applications["node"]
	if !ok {
		return ""
	}
	return path
}

func (configuration *TezbakeModuleConfiguration) loadBakersFromNodeConfiguration() {
	aliases := []string{"baker"} // baker is used by default

	nodeAppPath := configuration.GetNodeAppPath()
	if nodeAppPath == "" {
		slog.Error("node app path is not set")
		return
	}

	nodeAppJsonPath := filepath.Join(nodeAppPath, "app.json")
	if _, err := os.Stat(nodeAppJsonPath); os.IsNotExist(err) {
		nodeAppJsonPath = filepath.Join(nodeAppPath, "app.hjson")
	}

	nodeAppJsonBytes, err := os.ReadFile(nodeAppJsonPath)
	if err != nil {
		slog.Error("failed to read node app.json file", "error", err.Error())
		return
	}

	var nodeApp nodeAppJson
	err = hjson.Unmarshal(nodeAppJsonBytes, &nodeApp)
	if err != nil {
		slog.Error("failed to parse node app.json file", "error", err.Error())
		return
	}

	aliases = append(aliases, nodeApp.Configuration.AdditionalKeysAliases...)

	// r.Bakers
	publicKeyHashesPath := filepath.Join(nodeAppPath, "data", ".tezos-client", "public_key_hashs")
	publicKeyHashesBytes, err := os.ReadFile(publicKeyHashesPath)
	if err != nil {
		slog.Error("failed to read node public_key_hashs file", "error", err.Error())
		return
	}

	var publicKeys nodePublicKeys
	err = hjson.Unmarshal(publicKeyHashesBytes, &publicKeys)
	if err != nil {
		slog.Error("failed to parse node public_key_hashs file", "error", err.Error())
		return
	}

	for _, publicKey := range publicKeys {
		for _, alias := range aliases {
			if publicKey.Name == alias {
				configuration.Bakers = append(configuration.Bakers, publicKey.Hash)
			}
		}
	}

}

func (c *TezbakeModuleConfiguration) Hydrate() {
	if len(c.Bakers) == 0 {
		c.loadBakersFromNodeConfiguration()
	}

	validBakers := []string{}
	for _, baker := range c.Bakers {
		if _, err := tezos.ParseAddress(baker); err == nil {
			validBakers = append(validBakers, baker)
		} else {
			slog.Warn("Invalid baker address", "address", baker)
		}
	}
	c.Bakers = validBakers

	if c.RightsBlockWindow <= 0 {
		c.RightsBlockWindow = constants.DEFAULT_RIGHTS_BLOCK_WINDOW
	}

	if c.ArcBinaryPath == "" {
		exePath, err := os.Executable()
		if err == nil {
			exeDir := filepath.Dir(exePath)
			arcBinaryPath := filepath.Join(exeDir, "arc")
			if _, err := os.Stat(arcBinaryPath); err == nil {
				c.ArcBinaryPath = arcBinaryPath
			}
		}
	}
	if c.ArcBinaryPath == "" {
		cwd, err := os.Getwd()
		if err == nil {
			arcBinaryPath := filepath.Join(cwd, "arc")
			if _, err := os.Stat(arcBinaryPath); err == nil {
				c.ArcBinaryPath = arcBinaryPath
			}
		}
	}

	if len(c.LedgerWallets) == 0 {
		c.LedgerWallets = constants.DEFAULT_LEDGER_WALLETS
	}
}

func normalizeVersion(version string) string {
	if strings.HasPrefix(version, "v") {
		return version
	}
	return "v" + version
}

func (c *TezbakeModuleConfiguration) Validate() error {
	if _, err := url.Parse(c.SignerUrl); err != nil {
		return constants.ErrInvalidSignerUrl
	}

	if len(c.Bakers) == 0 {
		return constants.ErrNoValidBakers
	}

	arcBinaryPath, err := exec.LookPath(c.ArcBinaryPath)
	if err == nil {
		cmd := exec.Command(arcBinaryPath, "--version")
		output, err := cmd.Output()
		if err != nil {
			slog.Warn("Failed to get arc binary version", "error", err.Error())
		} else {
			version := normalizeVersion(strings.TrimSpace(string(output)))
			if !semver.IsValid(version) {
				slog.Warn("Invalid arc binary version", "version", version)
			}

			if semver.Compare(version, "v0.0.11") < 0 {
				slog.Warn("Arc binary version is too old, please update", "version", version)
			}
		}
	} else {
		slog.Warn("Failed to find arc binary, ledger monitoring disabled", "error", err.Error())
		// TODO: whether to monitor would be nicer through a flag in runtime configuration
		c.ArcBinaryPath = ""
	}

	return nil
}

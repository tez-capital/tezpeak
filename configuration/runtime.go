package configuration

import (
	"errors"
	"log/slog"
	"net"
	"net/url"
	"os"
	"path/filepath"

	"github.com/hjson/hjson-go/v4"
	"github.com/samber/lo"
	"github.com/tez-capital/tezpeak/constants"
	"github.com/tez-capital/tezpeak/constants/enums"
)

type versionedConfig interface {
	ToRuntime() *Runtime
}

type deserializedConfigVersion struct {
	Version int `json:"version,omitempty"`
}

type RuntimeReferenceNode struct {
	Address              string
	IsRightsProvider     bool
	IsBlockProvider      bool
	IsGovernanceProvider bool
}

type PeakMode string

const (
	PrivatePeakMode PeakMode = "private"
	PublicPeakMode  PeakMode = "public"
	AutoPeakMode    PeakMode = "auto"
)

type Providers struct {
	Services enums.ServiceStatusProviderKind
}

type Runtime struct {
	Id               string
	Listen           string
	Bakers           []string
	WorkingDirectory string
	TezbakeHome      string
	TezpayHome       string
	Node             string
	Signer           string
	ReferenceNodes   map[string]RuntimeReferenceNode
	BlockWindow      int64
	Providers        Providers
	Mode             PeakMode
}

func gerDefaultRuntime() *Runtime {
	return &Runtime{
		Id:               "",
		Listen:           constants.DEFAULT_LISTEN_ADDRESS,
		WorkingDirectory: "",
		TezbakeHome:      "",
		ReferenceNodes:   make(map[string]RuntimeReferenceNode),
		BlockWindow:      50,
		Mode:             AutoPeakMode,
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

func (r *Runtime) loadBakersFromNodeConfiguration() {
	aliases := []string{"baker"} // baker is used by default

	nodeDirectory := filepath.Join(r.WorkingDirectory, "node")
	nodeAppJsonPath := filepath.Join(nodeDirectory, "app.json")
	if _, err := os.Stat(nodeAppJsonPath); os.IsNotExist(err) {
		nodeAppJsonPath = filepath.Join(nodeDirectory, "app.hjson")
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
	publicKeyHashesPath := filepath.Join(nodeDirectory, "data", ".tezos-client", "public_key_hashs")
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
				r.Bakers = append(r.Bakers, publicKey.Hash)
			}
		}
	}
}

func (r *Runtime) Hydrate() *Runtime {
	if r.Listen == "" {
		r.Listen = constants.DEFAULT_LISTEN_ADDRESS
	}

	if r.Mode == AutoPeakMode {
		if host, _, err := net.SplitHostPort(r.Listen); err == nil {
			switch lo.Contains(constants.PRIVATE_NETWORK_HOSTS, host) {
			case true:
				r.Mode = PrivatePeakMode
			default:
				r.Mode = PublicPeakMode
			}
		} else {
			r.Mode = PublicPeakMode
		}
	}

	if r.WorkingDirectory == "" {
		r.WorkingDirectory, _ = os.Getwd()
	}

	if r.TezbakeHome == "" {
		envTezbakeHome := os.Getenv(constants.ENV_TEZBAKE_HOME)
		if envTezbakeHome != "" {
			r.TezbakeHome = envTezbakeHome
		} else {
			r.TezbakeHome, _ = os.Getwd()
		}
	}

	if r.Node == "" {
		r.Node = constants.DEFAULT_BAKER_NODE_URL
	}

	if r.Signer == "" {
		r.Signer = constants.DEFAULT_BAKER_SIGNER_URL
	}

	if r.BlockWindow == 0 {
		r.BlockWindow = 50
	}

	if len(r.ReferenceNodes) == 0 {
		r.ReferenceNodes = map[string]RuntimeReferenceNode{
			"Tezos Foundation": {
				Address:          constants.DEFAULT_REFERENCE_NODE_URL,
				IsRightsProvider: constants.DEFAULT_REFERENCE_NODE_IS_RIGHTS_PROVIDER,
				IsBlockProvider:  constants.DEFAULT_REFERENCE_NODE_IS_BLOCK_PROVIDER,
			},
			"tzkt": {
				Address:          constants.DEFAULT_REFERENCE_NODE_2_URL,
				IsRightsProvider: constants.DEFAULT_REFERENCE_NODE_2_IS_RIGHTS_PROVIDER,
				IsBlockProvider:  constants.DEFAULT_REFERENCE_NODE_2_IS_BLOCK_PROVIDER,
			},
		}
	}

	if r.Providers.Services == "" {
		r.Providers.Services = enums.TezbakeServiceStatusProvider
	}

	if len(r.Bakers) == 0 {
		r.loadBakersFromNodeConfiguration()
	}
	return r
}

func (r *Runtime) Validate() (*Runtime, error) {
	if r.Listen != "" {
		_, _, err := net.SplitHostPort(r.Listen)
		if err != nil {
			return nil, constants.ErrInvalidListenAddress
		}
	}

	if r.WorkingDirectory == "" {
		return nil, constants.ErrInvalidWorkingDirectory
	}

	if len(r.ReferenceNodes) == 0 {
		return nil, constants.ErrInvalidBlockWindow
	}

	if _, err := url.Parse(r.Signer); err != nil {
		return nil, constants.ErrInvalidSignerUrl
	}

	return r, nil
}

func Load() (*Runtime, error) {
	var err error
	configFilePath := os.Getenv(constants.ENV_TEZPEAK_CONFIG_FILE)
	if configFilePath == "" {
		configFilePath = "config.hjson"
	}

	configBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		slog.Debug("failed to read config file", "error", err.Error())

		// return config loaded from environment variables
		return gerDefaultRuntime().Hydrate().Validate()
	}

	var configVersion deserializedConfigVersion
	err = hjson.Unmarshal(configBytes, &configVersion)
	if err != nil {
		return nil, errors.Join(constants.ErrInvalidConfigVersion, err)
	}

	var configuration versionedConfig
	switch configVersion.Version {
	case 0:
		configuration, err = load_v0(configBytes)
	default:
		return nil, constants.ErrInvalidConfigVersion
	}

	if err != nil {
		return nil, errors.Join(constants.ErrInvalidConfig, err)
	}

	return configuration.ToRuntime().Hydrate().Validate()
}

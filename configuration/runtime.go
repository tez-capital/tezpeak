package configuration

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"os"
	"path/filepath"

	"github.com/hjson/hjson-go/v4"
	"github.com/tez-capital/tezpeak/constants"
)

type moduleConfigurationbase struct {
	Applications map[string]string `json:"applications,omitempty"`

	Mode PeakMode `json:"mode,omitempty"`
}

type PeakMode string

const (
	PrivatePeakMode PeakMode = "private"
	PublicPeakMode  PeakMode = "public"
	AutoPeakMode    PeakMode = "auto"
)

type versionedConfig interface {
	ToRuntime() *Runtime
}

type deserializedConfigVersion struct {
	Version int `json:"version,omitempty"`
}

type TezosNode struct {
	Address                string `json:"address"`
	IsRightsProvider       bool   `json:"is_rights_provider,omitempty"`
	IsBlockProvider        bool   `json:"is_block_provider,omitempty"`
	IsGovernanceProvider   bool   `json:"is_governance_provider,omitempty"`
	IsNetworkSInfoProvider bool   `json:"is_network_info_provider,omitempty"`
	IsEssential            bool   `json:"is_essential,omitempty"`
}

var (
	TF_RPC = TezosNode{
		Address:              "https://rpc.tzbeta.net/",
		IsGovernanceProvider: true,
	}
	TZKT_RPC = TezosNode{
		Address:              "https://rpc.tzkt.io/mainnet/",
		IsBlockProvider:      true,
		IsRightsProvider:     true,
		IsGovernanceProvider: true,
	}
	TZC_EU_RPC = TezosNode{
		Address:              "https://eu.rpc.tez.capital/",
		IsBlockProvider:      true,
		IsRightsProvider:     true,
		IsGovernanceProvider: true,
	}
	TZC_US_RPC = TezosNode{
		Address:              "https://us.rpc.tez.capital/",
		IsBlockProvider:      true,
		IsRightsProvider:     true,
		IsGovernanceProvider: true,
	}
	BAKER_NODE = TezosNode{
		Address:                "http://127.0.0.1:8732/",
		IsRightsProvider:       true,
		IsBlockProvider:        true,
		IsGovernanceProvider:   true,
		IsNetworkSInfoProvider: true,
		IsEssential:            true,
	}
)

type Runtime struct {
	Id     string
	Listen string
	Mode   PeakMode
	// path to the root where are the apps located e.g. /bake-buddy
	AppRoot string

	Modules map[string]json.RawMessage `json:"modules,omitempty"`

	Nodes map[string]TezosNode
}

func gerDefaultRuntime() *Runtime {
	return &Runtime{
		Id:      "",
		Listen:  constants.DEFAULT_LISTEN_ADDRESS,
		Mode:    AutoPeakMode,
		Modules: map[string]json.RawMessage{},
	}
}

func (v *Runtime) GetTezbakeModuleConfiguration() (bool, *TezbakeModuleConfiguration) {
	rawConfiguration, ok := v.Modules[constants.TEZBAKE_MODULE_ID]
	if !ok {
		return false, nil
	}

	configuration := getDefaultTezbakeModuleConfiguration()
	err := hjson.Unmarshal(rawConfiguration, configuration)
	if err != nil {
		slog.Error("failed to parse tezbake module configuration", "error", err.Error())
		return false, nil
	}

	for key, value := range configuration.Applications {
		if filepath.IsAbs(value) {
			continue // skip absolute paths
		}
		configuration.Applications[key] = filepath.Join(v.AppRoot, value)
	}

	if configuration.Mode == "" {
		configuration.Mode = v.Mode
	}
	configuration.Hydrate()

	if err := configuration.Validate(); err != nil {
		slog.Error("failed to validate tezbake module configuration", "error", err.Error())
		return false, nil
	}

	return true, configuration
}

func (v *Runtime) GetTezpayModuleConfiguration() (bool, *TezpayModuleConfiguration) {
	rawConfiguration, ok := v.Modules[constants.TEZPAY_MODULE_ID]
	if !ok {
		return false, nil
	}

	configuration := getDefaultTezpayModuleConfiguration()
	err := hjson.Unmarshal(rawConfiguration, configuration)
	if err != nil {
		slog.Error("failed to parse tezpay module configuration", "error", err.Error())
		return false, nil
	}

	for key, value := range configuration.Applications {
		if filepath.IsAbs(value) {
			continue // skip absolute paths
		}
		configuration.Applications[key] = filepath.Join(v.AppRoot, value)
	}

	if configuration.Mode == "" {
		configuration.Mode = v.Mode
	}
	configuration.Hydrate()

	if err := configuration.Validate(); err != nil {
		slog.Error("failed to validate tezpay module configuration", "error", err.Error())
		return false, nil
	}

	return true, configuration
}

func (r *Runtime) Validate() (*Runtime, error) {
	if r.Listen != "" {
		_, _, err := net.SplitHostPort(r.Listen)
		if err != nil {
			return nil, constants.ErrInvalidListenAddress
		}
	}

	if r.AppRoot == "" {
		return nil, constants.ErrInvalidWorkingDirectory
	}

	// TODO: optimize?
	if ok, module := r.GetTezbakeModuleConfiguration(); ok {
		if err := module.Validate(); err != nil {
			return nil, err
		}
	}

	if ok, module := r.GetTezpayModuleConfiguration(); ok {
		if err := module.Validate(); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *Runtime) Hydrate() *Runtime {
	if r.AppRoot == "" {
		r.AppRoot, _ = os.Getwd()
	}

	if len(r.Nodes) == 0 {
		r.Nodes = map[string]TezosNode{
			"baker":  BAKER_NODE,
			"TzC-EU": TZC_EU_RPC,
			"TzC-US": TZC_US_RPC,
			"TF":     TF_RPC,
			"TZKT":   TZKT_RPC,
		}
	}

	return r
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

package configuration

import (
	"encoding/json"

	"github.com/hjson/hjson-go/v4"
	"github.com/tez-capital/tezpeak/constants"
)

type v0 struct {
	Version int      `json:"version,omitempty"`
	Id      string   `json:"id,omitempty"`
	AppRoot string   `json:"app_root,omitempty"`
	Listen  string   `json:"listen,omitempty"`
	Mode    PeakMode `json:"mode,omitempty"`

	Modules map[string]json.RawMessage `json:"modules,omitempty"`

	Nodes map[string]TezosNode `json:"nodes,omitempty"`
}

func getDefault_v0() *v0 {
	return &v0{
		Version: 0,
		AppRoot: "",
		Listen:  constants.DEFAULT_LISTEN_ADDRESS,
		Mode:    AutoPeakMode,
		Modules: map[string]json.RawMessage{},
	}
}

func load_v0(configBytes []byte) (*v0, error) {
	configuration := getDefault_v0()

	err := hjson.Unmarshal(configBytes, &configuration)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}

func (v *v0) ToRuntime() *Runtime {
	result := &Runtime{
		Id:     v.Id,
		Listen: v.Listen,
		Mode:   v.Mode,

		AppRoot: v.AppRoot,

		Modules: v.Modules,
	}
	return result
}

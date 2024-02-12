package configuration

import (
	"github.com/hjson/hjson-go/v4"
	"github.com/tez-capital/tezpeak/constants"
)

type referenceNode struct {
	Address          string `json:"address"`
	IsRightsProvider *bool  `json:"is_rights_provider,omitempty"`
	IsBlockProvider  *bool  `json:"is_block_provider,omitempty"`
}

type v0 struct {
	Version          int                       `json:"version,omitempty"`
	Id               string                    `json:"id,omitempty"`
	Listen           string                    `json:"listen,omitempty"`
	Bakers           []string                  `json:"bakers,omitempty"`
	WorkingDirectory string                    `json:"working_directory,omitempty"`
	TezbakeHome      string                    `json:"tezbake_home,omitempty"`
	Node             string                    `json:"node,omitempty"`
	ReferenceNodes   *map[string]referenceNode `json:"reference_nodes,omitempty"`
	BlockWindow      int64                     `json:"block_window,omitempty"`
	Mode             PeakMode                  `json:"mode,omitempty"`
}

func getDefault_v0() *v0 {
	isRightsProvider := constants.DEFAULT_REFERENCE_NODE_IS_RIGHTS_PROVIDER
	isBlockProvider := constants.DEFAULT_REFERENCE_NODE_IS_BLOCK_PROVIDER

	isRightsProvider2 := constants.DEFAULT_REFERENCE_NODE_2_IS_RIGHTS_PROVIDER
	isBlockProvider2 := constants.DEFAULT_REFERENCE_NODE_2_IS_BLOCK_PROVIDER

	return &v0{
		Version:          0,
		Listen:           constants.DEFAULT_LISTEN_ADDRESS,
		Bakers:           []string{},
		WorkingDirectory: "",
		ReferenceNodes: &map[string]referenceNode{
			"Tezos Foundation": {
				Address:          constants.DEFAULT_REFERENCE_NODE_URL,
				IsRightsProvider: &isRightsProvider,
				IsBlockProvider:  &isBlockProvider,
			},
			"tzkt": {
				Address:          constants.DEFAULT_REFERENCE_NODE_2_URL,
				IsRightsProvider: &isRightsProvider2,
				IsBlockProvider:  &isBlockProvider2,
			},
		},
		BlockWindow: 50,
		Mode:        AutoPeakMode,
	}
}

func load_v0(configBytes []byte) (*v0, error) {
	configuration := getDefault_v0()

	err := hjson.Unmarshal(configBytes, &configuration)
	if err != nil {
		return nil, err
	}

	isProvider := true
	if configuration.ReferenceNodes != nil {
		for _, node := range *configuration.ReferenceNodes {
			if node.IsRightsProvider == nil {
				node.IsRightsProvider = &isProvider
			}
		}
	}
	return configuration, nil
}

func (v *v0) ToRuntime() *Runtime {
	result := &Runtime{
		Id:               v.Id,
		Listen:           v.Listen,
		Bakers:           v.Bakers,
		WorkingDirectory: v.WorkingDirectory,
		TezbakeHome:      v.TezbakeHome,
		Node:             v.Node,
		ReferenceNodes:   make(map[string]RuntimeReferenceNode),
		BlockWindow:      v.BlockWindow,
		Mode:             v.Mode,
	}

	if v.ReferenceNodes != nil {
		for name, node := range *v.ReferenceNodes {
			runtimeReferenceNode := RuntimeReferenceNode{
				Address:          node.Address,
				IsRightsProvider: true,
				IsBlockProvider:  true,
			}
			if node.IsRightsProvider != nil {
				runtimeReferenceNode.IsRightsProvider = *node.IsRightsProvider
			}
			if node.IsBlockProvider != nil {
				runtimeReferenceNode.IsBlockProvider = *node.IsBlockProvider
			}
			result.ReferenceNodes[name] = runtimeReferenceNode
		}
	}

	return result
}

package util

import (
	"encoding/json"
	"errors"

	"github.com/hjson/hjson-go/v4"
	"github.com/trilitech/tzgo/rpc"
)

func TryUnwrapRPCError(err error) error {
	if rpcError, ok := err.(rpc.RPCError); ok {
		body := rpcError.Body()
		if len(body) == 0 {
			return err
		}

		var message any
		err := json.Unmarshal(body, &message)
		if err != nil {
			return errors.New(string(body))
		}

		content, err := hjson.Marshal(message)
		if err != nil {
			return errors.New(string(body))
		}
		return errors.New(string(content))
	}
	return err
}

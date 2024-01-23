package providers

import (
	"log/slog"

	"blockwatch.cc/tzgo/rpc"
)


func attemptWithClients[T interface{}](clients []*rpc.Client, f func(client *rpc.Client) (T, error)) (T, error) {
	var err error
	var result T

	for _, client := range clients {
		slog.Debug("attempting with client", "client", client.BaseURL.Host)
		result, err = f(client)
		if err != nil {
			continue
		}
		return result, nil
	}
	return result, err
}

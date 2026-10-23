package common

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/tez-capital/tezpeak/configuration"
	"github.com/trilitech/tzgo/rpc"
)

func TestAttemptWithRpcClientsPriorityAndFallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chains/main/is_bootstrapped" {
			t.Errorf("unexpected RPC path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"bootstrapped":true,"sync_state":"synced"}`)
	}))
	defer server.Close()

	previousNodes := activeRpcNodes
	t.Cleanup(func() { activeRpcNodes = previousNodes })
	activeRpcNodes = make(map[string]*ActiveRpcNode)
	for _, priority := range []int{30, 10, 40, 20} {
		client, err := rpc.NewClient(server.URL, server.Client())
		if err != nil {
			t.Fatal(err)
		}
		activeRpcNodes[fmt.Sprint(priority)] = &ActiveRpcNode{
			TezosNode: configuration.TezosNode{Priority: priority},
			Client:    client,
		}
	}

	var attempts []int
	result, err := AttemptWithRpcClients(t.Context(), func(node *ActiveRpcNode) (int, error) {
		attempts = append(attempts, node.Priority)
		if node.Priority > 20 {
			return 0, errors.New("RPC failed")
		}
		return node.Priority, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if result != 20 {
		t.Errorf("result = %d, want 20", result)
	}
	if want := []int{40, 30, 20}; !slices.Equal(attempts, want) {
		t.Errorf("attempted priorities = %v, want %v", attempts, want)
	}
}

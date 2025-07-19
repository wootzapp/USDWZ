//go:build integration

package integration_test

import (
	"fmt"
	"sort"
	"sync"
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/example/usdwz/chain/app"
	"github.com/example/usdwz/chain/vm"
	escrowtypes "github.com/example/usdwz/chain/x/escrow/types"
	stabletypes "github.com/example/usdwz/chain/x/stablecoin/types"
	validatortypes "github.com/example/usdwz/chain/x/validator/types"
	yieldtypes "github.com/example/usdwz/chain/x/yield/types"
)

// testNode simulates a single chain node.
type testNode struct {
	id  string
	app *app.UsdWzApp
	ctx sdk.Context
}

func newTestNode(id string) *testNode {
	a := app.New(log.NewNopLogger())
	ctx := contextWithMemDB(a)
	return &testNode{id: id, app: a, ctx: ctx}
}

func contextWithMemDB(a *app.UsdWzApp) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	keys := []string{
		stabletypes.StoreKey,
		escrowtypes.StoreKey,
		validatortypes.StoreKey,
		yieldtypes.StoreKey,
	}
	for _, name := range keys {
		ms.MountStoreWithDB(a.KVStoreKey(name), storetypes.StoreTypeIAVL, db)
	}
	_ = ms.LoadLatestVersion()
	return sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger())
}

// network emulates a minimal multi-node environment using separate
// in-memory stores for each validator. This mirrors the approach in
// Cosmos SDK integration tests, where each node runs its own app
// instance with an isolated context.
type network struct {
	nodes map[string]*testNode
}

func newNetwork(n int) *network {
	ids := make([]string, n)
	for i := 0; i < n; i++ {
		ids[i] = fmt.Sprintf("v%d", i+1)
	}
	return newNetworkWithIDs(ids)
}

func newNetworkWithIDs(ids []string) *network {
	net := &network{
		nodes: make(map[string]*testNode),
	}
	for _, id := range ids {
		net.nodes[id] = newTestNode(id)
	}
	return net
}

func (net *network) submitVote(id string, approve bool) {
	node, ok := net.nodes[id]
	if !ok {
		return
	}
	node.app.ValidatorKeeper.SubmitVote(node.ctx, id, approve)
}

func (net *network) aggregated() map[string]bool {
	m := make(map[string]bool, len(net.nodes))
	for id, n := range net.nodes {
		if v, ok := n.app.ValidatorKeeper.Vote(n.ctx, id); ok {
			m[id] = v
		}
	}
	return m
}

func TestDistributedQuorum(t *testing.T) {
	script := []vm.Instruction{{vm.OP_SET, "A"}, {vm.OP_QUORUM, "2"}}
	sets := map[string][]string{"A": {"v1", "v2", "v3"}}
	cases := []struct {
		name  string
		votes []bool
		pass  bool
	}{
		{"quorum yes", []bool{true, true, false}, true},
		{"quorum not reached", []bool{true, false, false}, false},
		{"all yes", []bool{true, true, true}, true},
		{"all no", []bool{false, false, false}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			net := newNetworkWithIDs([]string{"v1", "v2", "v3"})
			var wg sync.WaitGroup
			for i, v := range tc.votes {
				wg.Add(1)
				id := fmt.Sprintf("v%d", i+1)
				go func(vid string, vote bool) {
					defer wg.Done()
					net.submitVote(vid, vote)
				}(id, v)
			}
			wg.Wait()
			ag := net.aggregated()
			for id, n := range net.nodes {
				res := n.app.VM.Validate(script, ag, sets)
				if res != tc.pass {
					t.Fatalf("node %s expect %v got %v", id, tc.pass, res)
				}
			}
		})
	}
}

// collectIDs returns the union of validator IDs present in sets and votes.
func collectIDs(sets map[string][]string, votes map[string]bool) []string {
	m := make(map[string]struct{})
	for _, vs := range sets {
		for _, id := range vs {
			m[id] = struct{}{}
		}
	}
	for id := range votes {
		m[id] = struct{}{}
	}
	ids := make([]string, 0, len(m))
	for id := range m {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func TestDistributedMultiStageQuorum(t *testing.T) {
	script := []vm.Instruction{
		{vm.OP_SET, "A"}, {vm.OP_QUORUM, "2"}, {vm.OP_THEN, ""},
		{vm.OP_SET, "B"}, {vm.OP_ALL, ""},
	}
	cases := []struct {
		name  string
		sets  map[string][]string
		votes map[string]bool
		pass  bool
	}{
		{
			"happy path",
			map[string][]string{"A": {"v1", "v2", "v3"}, "B": {"v4", "v5"}},
			map[string]bool{"v1": true, "v2": true, "v4": true, "v5": true},
			true,
		},
		{
			"fail second",
			map[string][]string{"A": {"v1", "v2", "v3"}, "B": {"v4", "v5"}},
			map[string]bool{"v1": true, "v2": true, "v5": true},
			false,
		},
		{
			"fail first",
			map[string][]string{"A": {"v1", "v2", "v3"}, "B": {"v4", "v5"}},
			map[string]bool{"v1": true, "v4": true, "v5": true},
			false,
		},
		{
			"different sets",
			map[string][]string{"A": {"x1", "x2"}, "B": {"y1", "y2"}},
			map[string]bool{"x1": true, "x2": true, "y1": true, "y2": true},
			true,
		},
		{
			"extra votes",
			map[string][]string{"A": {"v1", "v2"}, "B": {"v3"}},
			map[string]bool{"v1": true, "v2": true, "v3": true, "v4": true},
			true,
		},
		{
			"missing all second",
			map[string][]string{"A": {"v1", "v2"}, "B": {"v3"}},
			map[string]bool{"v1": true, "v2": true},
			false,
		},
		{
			"quorum not met",
			map[string][]string{"A": {"v1", "v2"}, "B": {"v3"}},
			map[string]bool{"v1": true, "v3": true},
			false,
		},
		{
			"duplicate validators",
			map[string][]string{"A": {"v1", "v1"}, "B": {"v2", "v2"}},
			map[string]bool{"v1": true, "v2": true},
			true,
		},
		{
			"empty sets",
			map[string][]string{"A": []string{}, "B": []string{}},
			map[string]bool{},
			false,
		},
		{
			"extra then ignored",
			map[string][]string{"A": {"v1"}, "B": {"v2"}},
			map[string]bool{"v1": true, "v2": true},
			false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ids := collectIDs(tc.sets, tc.votes)
			net := newNetworkWithIDs(ids)
			var wg sync.WaitGroup
			for id, vote := range tc.votes {
				wg.Add(1)
				go func(id string, vote bool) {
					defer wg.Done()
					net.submitVote(id, vote)
				}(id, vote)
			}
			wg.Wait()
			ag := net.aggregated()
			for id, n := range net.nodes {
				res := n.app.VM.Validate(script, ag, tc.sets)
				if res != tc.pass {
					t.Fatalf("node %s expect %v got %v", id, tc.pass, res)
				}
			}
		})
	}
}

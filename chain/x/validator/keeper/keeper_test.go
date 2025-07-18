package keeper

import (
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	stablecoinkeeper "github.com/example/usdwz/chain/x/stablecoin/keeper"
	"github.com/example/usdwz/chain/x/validator/types"
)

func TestNewKeeper(t *testing.T) {
	key := storetypes.NewKVStoreKey("validator")
	k := NewKeeper(nil, key)
	if k.storeKey != key {
		t.Fatal("store key mismatch")
	}
}

func TestModuleName(t *testing.T) {
	if ModuleName() != types.ModuleName {
		t.Fatalf("expected %s got %s", types.ModuleName, ModuleName())
	}
}

func TestVoting(t *testing.T) {
	key := storetypes.NewKVStoreKey("validator")
	k := NewKeeper(nil, key)
	ctx := stablecoinkeeper.TestContext(key)

	k.SubmitVote(ctx, "v1", true)
	v, ok := k.Vote(ctx, "v1")
	if !ok || !v {
		t.Fatal("vote not stored")
	}
	k.SubmitVote(ctx, "v2", false)
	yes, no := k.Tally(ctx)
	if yes != 1 || no != 1 {
		t.Fatalf("unexpected tally %d/%d", yes, no)
	}
}

func TestVoteEdgeCases(t *testing.T) {
	key := storetypes.NewKVStoreKey("validator")
	k := NewKeeper(nil, key)
	ctx := stablecoinkeeper.TestContext(key)

	cases := []struct {
		id  string
		val bool
	}{
		{"dup", true},
		{"dup", false},
		{"missing", false},
		{"unicode☃", true},
		{"long", true},
		{"upper", false},
		{"", true},
	}

	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			k.SubmitVote(ctx, tc.id, tc.val)
			v, _ := k.Vote(ctx, tc.id)
			if v != tc.val {
				t.Fatalf("vote mismatch for %s", tc.id)
			}
		})
	}
}

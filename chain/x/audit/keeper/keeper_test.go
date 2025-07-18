package keeper

import (
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	store "github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupContext(key *storetypes.KVStoreKey) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	_ = ms.LoadLatestVersion()
	return sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger())
}

func TestPublishAttestation(t *testing.T) {
	key := storetypes.NewKVStoreKey("audit")
	ctx := setupContext(key)
	k := NewKeeper(nil, key)
	if err := k.PublishAttestation(ctx, "hash1"); err != nil {
		t.Fatalf("publish failed: %v", err)
	}
	if !ctx.KVStore(key).Has([]byte("attestation:hash1")) {
		t.Fatal("attestation not stored")
	}
}

func TestAttestationQueries(t *testing.T) {
	key := storetypes.NewKVStoreKey("audit")
	ctx := setupContext(key)
	k := NewKeeper(nil, key)

	_ = k.PublishAttestation(ctx, "h1")
	att, ok := k.GetAttestation(ctx, "h1")
	if !ok || att.Hash != "h1" {
		t.Fatal("get attestation failed")
	}
	if _, ok := k.GetAttestation(ctx, "missing"); ok {
		t.Fatal("expect missing attestation")
	}

	_ = k.PublishAttestation(ctx, "h2")
	list := k.ListAttestations(ctx)
	if len(list) != 2 {
		t.Fatalf("expected 2 got %d", len(list))
	}
}

func TestAttestationScenarios(t *testing.T) {
	key := storetypes.NewKVStoreKey("audit")
	ctx := setupContext(key)
	k := NewKeeper(nil, key)

	cases := []struct {
		name string
		hash string
	}{
		{"simple", "a1"},
		{"empty", ""},
		{"numbers", "123"},
		{"symbols", "@@"},
		{"long", "0123456789abcdef"},
		{"unicode", "测试"},
		{"spaces", "hash with space"},
		{"repeat", "dup"},
		{"repeat-second", "dup"},
		{"final", "end"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_ = k.PublishAttestation(ctx, tc.hash)
			if _, ok := k.GetAttestation(ctx, tc.hash); !ok {
				t.Fatalf("attestation %s missing", tc.hash)
			}
		})
	}
	unique := 0
	seen := map[string]bool{}
	for _, c := range cases {
		if !seen[c.hash] {
			seen[c.hash] = true
			unique++
		}
	}
	if len(k.ListAttestations(ctx)) != unique {
		t.Fatal("list length mismatch")
	}
}

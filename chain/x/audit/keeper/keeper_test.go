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

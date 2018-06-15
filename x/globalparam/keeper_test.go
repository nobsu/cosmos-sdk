package globalparam

import (
	"testing"

	"github.com/stretchr/testify/assert"

	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

func defaultContext(key sdk.StoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, abci.Header{}, false, nil, log.NewNopLogger())
	return ctx
}

func TestKeeper(t *testing.T) {
	kvs := []struct {
		key   string
		param int64
	}{
		{"key1", 10},
		{"key2", 55},
		{"key3", 182},
		{"key4", 17582},
		{"key5", 2768554},
	}

	skey := sdk.NewKVStoreKey("test")
	ctx := defaultContext(skey)
	setter := NewKeeper(skey, wire.NewCodec()).Setter()

	for _, kv := range kvs {
		err := setter.Set(ctx, kv.key, kv.param)
		assert.Nil(t, err)
	}

	for _, kv := range kvs {
		var param int64
		err := setter.Get(ctx, kv.key, &param)
		assert.Nil(t, err)
		assert.Equal(t, kv.param, param)
	}

	cdc := wire.NewCodec()
	for _, kv := range kvs {
		var param int64
		bz := setter.GetBytes(ctx, kv.key)
		err := cdc.UnmarshalBinary(bz, &param)
		assert.Nil(t, err)
		assert.Equal(t, kv.param, param)
	}

	for _, kv := range kvs {
		var param bool
		err := setter.Get(ctx, kv.key, &param)
		assert.NotNil(t, err)
	}

	for _, kv := range kvs {
		err := setter.Set(ctx, kv.key, true)
		assert.NotNil(t, err)
	}
}

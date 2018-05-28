package oracle

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"

	"github.com/tendermint/go-crypto"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

func defaultContext(keys ...sdk.StoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	for _, key := range keys {
		cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	}
	cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, abci.Header{}, false, nil, nil)
	return ctx
}

type validator struct {
	address sdk.Address
	power   sdk.Rat
}

func (v validator) GetStatus() sdk.BondStatus {
	return sdk.Bonded
}

func (v validator) GetOwner() sdk.Address {
	return v.address
}

func (v validator) GetPubKey() crypto.PubKey {
	return nil
}

func (v validator) GetPower() sdk.Rat {
	return v.power
}

func (v validator) GetBondHeight() int64 {
	return 0
}

type validatorSet struct {
	validators []validator
}

func (vs *validatorSet) IterateValidators(ctx sdk.Context, fn func(index int64, validator sdk.Validator) bool) {
	for i, val := range vs.validators {
		if fn(int64(i), val) {
			break
		}
	}
}

func (vs *validatorSet) IterateValidatorsBonded(ctx sdk.Context, fn func(index int64, validator sdk.Validator) bool) {
	vs.IterateValidators(ctx, fn)
}

func (vs *validatorSet) Validator(ctx sdk.Context, addr sdk.Address) sdk.Validator {
	for _, val := range vs.validators {
		if bytes.Equal(val.address, addr) {
			return val
		}
	}
	return nil
}

func (vs *validatorSet) TotalPower(ctx sdk.Context) sdk.Rat {
	res := sdk.ZeroRat()
	for _, val := range vs.validators {
		res = res.Add(val.power)
	}
	return res
}

type seqOracle struct {
	Seq   int
	Nonce int
}

func (o seqOracle) Type() string {
	return "seq"
}

func (o seqOracle) ValidateBasic() sdk.Error {
	return nil
}

func makeCodec() *wire.Codec {
	var cdc = wire.NewCodec()

	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterConcrete(OracleMsg{}, "test/Oracle", nil)

	cdc.RegisterInterface((*Payload)(nil), nil)
	cdc.RegisterConcrete(seqOracle{}, "test/oracle/seqOracle", nil)

	return cdc
}

func seqHandler(ork Keeper, key sdk.StoreKey, codespace sdk.CodespaceType) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case OracleMsg:
			return ork.Handle(func(ctx sdk.Context, p Payload) sdk.Error {
				switch p := p.(type) {
				case seqOracle:
					return handleSeqOracle(ctx, key, p)
				default:
					return sdk.ErrUnknownRequest("")
				}
			}, ctx, msg, codespace)
		default:
			return sdk.ErrUnknownRequest("").Result()
		}
	}
}

func getSequence(ctx sdk.Context, key sdk.StoreKey) int {
	store := ctx.KVStore(key)
	seqbz := store.Get([]byte("seq"))

	var seq int
	if seqbz == nil {
		seq = 0
	} else {
		wire.NewCodec().MustUnmarshalBinary(seqbz, &seq)
	}

	return seq
}

func handleSeqOracle(ctx sdk.Context, key sdk.StoreKey, o seqOracle) sdk.Error {
	store := ctx.KVStore(key)

	seq := getSequence(ctx, key)
	if seq != o.Seq {
		return sdk.NewError(sdk.CodespaceUndefined, 1, "")
	}

	bz := wire.NewCodec().MustMarshalBinary(seq + 1)
	store.Set([]byte("seq"), bz)

	return nil
}

func TestOracle(t *testing.T) {
	cdc := makeCodec()

	addr1 := []byte("addr1")
	addr2 := []byte("addr2")
	addr3 := []byte("addr3")
	addr4 := []byte("addr4")
	valset := &validatorSet{[]validator{
		validator{addr1, sdk.NewRat(7)},
		validator{addr2, sdk.NewRat(7)},
		validator{addr3, sdk.NewRat(1)},
	}}

	okey := sdk.NewKVStoreKey("oracle")
	key := sdk.NewKVStoreKey("key")
	ctx := defaultContext(okey, key)

	ork := NewKeeper(okey, cdc, valset)
	h := seqHandler(ork, key, sdk.CodespaceUndefined)

	// Nonvalidator signed, transaction failed
	msg := OracleMsg{seqOracle{0, 0}, []byte("randomguy")}
	res := h(ctx, msg)
	assert.False(t, res.IsOK())
	assert.Equal(t, 0, getSequence(ctx, key))

	// Less than 2/3 signed, msg not processed
	msg.Signer = addr1
	res = h(ctx, msg)
	assert.True(t, res.IsOK())
	assert.Equal(t, 0, getSequence(ctx, key))

	// Double signed, transaction failed
	res = h(ctx, msg)
	assert.False(t, res.IsOK())
	assert.Equal(t, 0, getSequence(ctx, key))

	// More than 2/3 signed, msg processed
	msg.Signer = addr2
	res = h(ctx, msg)
	assert.True(t, res.IsOK())
	assert.Equal(t, 1, getSequence(ctx, key))

	// Already processed, transaction failed
	msg.Signer = addr3
	res = h(ctx, msg)
	assert.False(t, res.IsOK())
	assert.Equal(t, 1, getSequence(ctx, key))

	// Less than 2/3 signed, msg not processed
	msg = OracleMsg{seqOracle{100, 1}, addr1}
	res = h(ctx, msg)
	assert.True(t, res.IsOK())
	assert.Equal(t, 1, getSequence(ctx, key))

	// More than 2/3 signed but payload is invalid
	msg.Signer = addr2
	res = h(ctx, msg)
	assert.True(t, res.IsOK())
	assert.NotEqual(t, "", res.Log)
	assert.Equal(t, 1, getSequence(ctx, key))

	// Already processed, transaction failed
	msg.Signer = addr3
	res = h(ctx, msg)
	assert.False(t, res.IsOK())
	assert.Equal(t, 1, getSequence(ctx, key))

	// Should handle validator set change
	valset.validators = append(valset.validators, validator{addr4, sdk.NewRat(12)})

	// Less than 2/3 signed, msg not processed
	msg = OracleMsg{seqOracle{1, 2}, addr1}
	res = h(ctx, msg)
	assert.True(t, res.IsOK())
	assert.Equal(t, 1, getSequence(ctx, key))

	// Less than 2/3 signed, msg not processed
	msg.Signer = addr2
	res = h(ctx, msg)
	assert.True(t, res.IsOK())
	assert.Equal(t, 1, getSequence(ctx, key))

	// More than 2/3 signed, msg processed
	msg.Signer = addr4
	res = h(ctx, msg)
	assert.True(t, res.IsOK())
	assert.Equal(t, 2, getSequence(ctx, key))

	// Should handle validator set change while oracle process is happening
	msg = OracleMsg{seqOracle{2, 3}, addr4}

	// Less than 2/3 signed, msg not processed
	res = h(ctx, msg)
	assert.True(t, res.IsOK())
	assert.Equal(t, 2, getSequence(ctx, key))

	// Signed validator is kicked out
	valset.validators = valset.validators[:len(valset.validators)-1]

	// Less than 2/3 signed, msg not processed
	msg.Signer = addr1
	res = h(ctx, msg)
	assert.True(t, res.IsOK())
	assert.Equal(t, 2, getSequence(ctx, key))

	// More than 2/3 signed, msg processed
	msg.Signer = addr2
	res = h(ctx, msg)
	assert.True(t, res.IsOK())
	assert.Equal(t, 3, getSequence(ctx, key))
}

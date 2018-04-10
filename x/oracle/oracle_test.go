package oracle

import (
	"testing"

	//"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"

	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	//"github.com/cosmos/cosmos-sdk/x/stake"
)

func defaultContext(keys ...sdk.StoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	for _, key := range keys {
		cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	}
	cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, abci.Header{}, false, nil)
	return ctx
}

type seqOracle struct {
	seq int
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

	cdc.RegisterInterface((*Oracle)(nil), nil)
	cdc.RegisterConcrete(seqOracle{}, "test/oracle/seqOracle", nil)

	return cdc
}

func seqHandler(ork Keeper, key sdk.StoreKey) sdk.Handler {
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
			}, ctx, msg)
		default:
			return sdk.ErrUnknownRequest("").Result()
		}
	}
}

func handleSeqOracle(ctx sdk.Context, key sdk.StoreKey, o seqOracle) sdk.Error {
	store := ctx.KVStore(key)
	cdc := makeCodec()

	seqbz := store.Get([]byte("seq"))

	var seq int
	if seqbz == nil {
		seq = 0
	} else {
		if err := cdc.UnmarshalBinary(seqbz, &seq); err != nil {
			return sdk.NewError(1, "")
		}
	}

	if seq != o.seq {
		return sdk.NewError(1, "")
	}

	bz, _ := cdc.MarshalBinary(seq + 1)
	store.Set([]byte("seq"), bz)

	return nil
}

func TestOracle(t *testing.T) {
	cdc := makeCodec()

	addrs := []sdk.Address{[]byte("0"), []byte("1"), []byte("2")}

	akey := sdk.NewKVStoreKey("auth")
	skey := sdk.NewKVStoreKey("stake")
	okey := sdk.NewKVStoreKey("oracle")
	key := sdk.NewKVStoreKey("key")
	ctx := defaultContext(skey, okey, key)

	am := auth.NewAccountMapper(cdc, akey, &auth.BaseAccount{})

	ck := bank.NewCoinKeeper(am)
	ck.AddCoins(ctx, addrs[0], sdk.Coins{{"fermion", 7}})
	ck.AddCoins(ctx, addrs[1], sdk.Coins{{"fermion", 7}})
	ck.AddCoins(ctx, addrs[2], sdk.Coins{{"fermion", 1}})
	/*
		sk := stake.NewKeeper(ctx, cdc, skey, ck)
		c0 := stake.Candidate{
			Address: sdk.Address(addrs[0]),
			Assets:  sdk.NewRat(7),
		}
		c1 := stake.Candidate{
			Address: sdk.Address(addrs[1]),
			Assets:  sdk.NewRat(7),
		}
		c2 := stake.Candidate{
			Address: sdk.Address(addrs[2]),
			Assets:  sdk.NewRat(1),
		}
			sk.setCandidate(ctx, c1)
			sk.setCandidate(ctx, c2)
			sk.setCandidate(ctx, c3)
	*/

	//ork := NewKeeper(okey, cdc, sk)

	//h := seqHandler(ork, key)
}

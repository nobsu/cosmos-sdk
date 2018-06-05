# Oracle Module

`x/oracle` provides a way to receive external information(real world price, events from other chains, etc.) with validators' vote. Each validator make transaction which contains those informations, and Oracle aggregates them until the supermajority signed on it. After then, Oracle sends the information to the actual module that processes the information, and prune the votes from the state.

## Integration

See `x/oracle/oracle_test.go` for the code that using Oracle

To use Oracle in your module, first define a `payload`. It should implement `oracle.Payload` and contain nessesary information for your module. Including nonce is recommended.

```go
type MyPayload struct {
    Data  int
    Nonce int
}
```

When you write a payload, its `.Type()` should return same name with your module is registered on the router. It is because `oracle.Msg` inherits `.Type()` from its embedded payload and it should be handled on the user modules.

Then route every incoming `oracle.Msg` to `oracle.Keeper.Handler()` with the function that implements `oracle.Handler`.

```go
func NewHandler(keeper Keeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
        switch msg := msg.(type) {
        case oracle.Msg: 
            return keeper.oracle.Handle(ctx sdk.Context, p oracle.Payload) sdk.Error {
                switch p := p.(type) {
                case MyPayload:
                    return handleMyPayload(ctx, keeper, p)
                }
            }
        }
    }
}
```

In the previous example, the keeper has an `oracle.Keeper`. To store an `oracle.Keeper`, your `NewKeeper` has to receive an `oracle.KeeperGen`. `oracle.KeeperGen` is a function that takes configurations and returns an `oracle.Keeper`

```go
func NewKeeper(key sdk.StoreKey, cdc *wire.Codec, ogen oracle.KeeperGen) Keeper {
    return Keeper {
        cdc: cdc,
        key: key,
        // The oracle keeper will pass payload
        // when more than 2/3 signed on it
        // and will prune votes after 100 blocks from last sign
        ork: ogen(cdc, sdk.NewRat(2, 3), 100),
    }
}
```

Now the validators can send `oracle.Msg`s with `MyPayload` when they want to witness external events. 

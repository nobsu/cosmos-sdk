package types

type Oracle interface {
	Type() string
	ValidateBasic() sdk.Error
}

package shift

type Symbol interface {
	IsTerminal() bool
	IsEqual(other Symbol) bool
	ToSymbolString() string
}

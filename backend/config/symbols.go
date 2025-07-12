// supported symbols for the backend

package config

// SupportedSymbols contains a list of symbols that are supported by the backend.
var SupportedSymbols = []string{
	"6A", "6B", "6C", "6E", "6J",
	"6N", "6S", "CL", "E7", "ES", "GC",
	"HE", "HG", "HO", "LE", "M2K", "M6A", "M6E", "MCL", "MES", "MGC",
	"MNQ", "MYM", "NG", "NQ", "QI", "QM", "QO", "RB", "RTY", "SI", "UB",
	"YM", "ZB", "ZC", "ZF", "ZL", "ZN", "ZQ", "ZS", "ZT", "ZW",
}

// IsSupported checks if a given symbol is supported by the backend.
func IsSupported(symbol string) bool {
	for _, sym := range SupportedSymbols {
		if symbol == sym {
			return true
		}
	}
	return false
}

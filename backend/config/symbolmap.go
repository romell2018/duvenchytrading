package config

var SymbolMap = map[string]string{
	// Currency Futures
	"6E": "EUR/USD",
	"6A": "AUD/USD",
	"6B": "GBP/USD",
	"6C": "CAD/USD",
	"6J": "USD/JPY",
	"6N": "NZD/USD",
	"6S": "USD/CHF",

	// Commodities / Index Futures
	"CL":  "CL=F",
	"NG":  "NG=F",
	"GC":  "GC=F",
	"SI":  "SI=F",
	"ES":  "ES=F",
	"NQ":  "NQ=F",
	"YM":  "YM=F",
	"RTY": "RTY=F",

	// Others: Add more mappings when you confirm which Twelve Data symbol matches
}

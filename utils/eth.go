package utils

import (
	"math"
	"math/big"
)

func WeiToEtherFloatByDecimals(decimals int, wei *big.Int) *big.Float {
	if decimals == 0 {
		decimals = 18
	}
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(math.Pow10(decimals)))
}

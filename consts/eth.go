package consts

import (
	"github.com/ethereum/go-ethereum/common"
)

const (
	WBNBAddress             = "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"
	WETHAddress             = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	UniSwapV2RouterAddress  = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
	PancakeRouterAddress    = "0x10ED43C718714eb63d5aA57B78B54704E256024E"
	UniSwapV2FactoryAddress = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"
	PancakeFactoryAddress   = "0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73"
)

const (
	BscRpcAddr = "https://bsc-dataseed.binance.org/"
	EthRpcAddr = "https://main-light.eth.linkpool.io/"
)

const (
	ChainTypeEth = "eth"
	ChainTypeBsc = "bsc"
)

var (
	ZeroAddress = common.Address{}
	ChainIdMap  = map[string]int64{
		ChainTypeEth: 1,
		ChainTypeBsc: 56,
	}
	UniSwapWrapperTokenContractMap = map[string]common.Address{
		ChainTypeEth: common.HexToAddress(WETHAddress),
		ChainTypeBsc: common.HexToAddress(WBNBAddress),
	}

	UniSwapFactoryContractMap = map[string]common.Address{
		ChainTypeEth: common.HexToAddress(UniSwapV2FactoryAddress),
		ChainTypeBsc: common.HexToAddress(PancakeFactoryAddress),
	}

	UniSwapRouterContractMap = map[string]common.Address{
		ChainTypeEth: common.HexToAddress(UniSwapV2RouterAddress),
		ChainTypeBsc: common.HexToAddress(PancakeRouterAddress),
	}
)

package consts

const (
	WBNBAddress            = ""
	WETHAddress            = ""
	UniSwapV2RouterAddress = ""
	PancakeRouterAddress   = ""
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
	ChainIdMap = map[string]int64{
		ChainTypeEth: 1,
		ChainTypeBsc: 56,
	}
)

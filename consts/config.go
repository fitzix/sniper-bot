package consts

var DefaultConfig = `
gasLimit: 200000
gasPrice: 10
# eip1559 gas tip, bsc not support
gasTip: 10
buyingBnbOrEthAmount: 0.3
# only when the number of eth/bnb in the pool liquidity is greater than the minPoolLiquidityAdded will trigger. only for cake and uni,
minPoolLiquidityAdded: 3
# percent 12 mean 12%, 0 mean auto. When buying start up, it’s best to set it to 0
slippage: 12
privateKey: your private key
# sniper interval when contract not active. 1s = 1000ms
sniperInterval: 1000
targetContract: 0x31e7ddebc4b4c1a9ba91a761390445f887354b25
`

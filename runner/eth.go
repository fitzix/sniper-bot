package runner

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/fitzix/sniper-bot/consts"
	"github.com/fitzix/sniper-bot/contract/uniswap"
	"github.com/fitzix/sniper-bot/utils"
	"github.com/spf13/viper"
)

type ethRunner struct {
	privateKey   *ecdsa.PrivateKey
	fromAddress  common.Address
	uniAbi       abi.ABI
	ethClientMap map[string]*ethclient.Client
}

func NewEthRunner() *ethRunner {
	uniAbi, err := abi.JSON(strings.NewReader(uniswap.UniswapV2ABI))
	if err != nil {
		panic(err)
	}

	privateKey, err := crypto.HexToECDSA(viper.GetString("privateKey"))
	if err != nil {
		panic(err)
	}

	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &ethRunner{
		privateKey:   privateKey,
		fromAddress:  fromAddress,
		uniAbi:       uniAbi,
		ethClientMap: make(map[string]*ethclient.Client),
	}
}

func (e *ethRunner) SniperDxsale(chain string) {
	dxsaleContractAddress := common.HexToAddress(viper.GetString("targetContract"))
	value, _ := big.NewFloat(viper.GetFloat64("buyingBnbOrEthAmount") * params.Ether).Int(nil)
	estimateTransferGasData, err := e.uniAbi.Pack("transfer", dxsaleContractAddress, value)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	gas, err := e.getClient(chain).EstimateGas(ctx, ethereum.CallMsg{
		To:   &dxsaleContractAddress,
		Data: estimateTransferGasData,
	})

	if err != nil {
		if err.Error() != "execution reverted" {
			panic(err)
		}
		interval := viper.GetInt64("sniperInterval")
		log.Printf("contract not active, retry in %d ms", interval)
		time.AfterFunc(time.Duration(interval)*time.Millisecond, func() {
			e.SniperDxsale(chain)
		})
		return
	}

	if viper.GetUint64("gasLimit") < gas {
		log.Println("config gas limit less than estimate gas ", gas, "auto set to estimate gasLimit")
		viper.Set("gasLimit", gas)
	}

	log.Println("EstimateGas", gas, "ready to transfer")

	e.transfer(ctx, chain, dxsaleContractAddress, value)
}

func (e *ethRunner) transfer(ctx context.Context, chain string, toAddress common.Address, value *big.Int) {

	nonce, err := e.getClient(chain).PendingNonceAt(ctx, e.fromAddress)
	if err != nil {
		panic(err)
	}

	println(nonce)

	tx, err := types.SignNewTx(e.privateKey, types.LatestSignerForChainID(big.NewInt(consts.ChainIdMap[chain])), &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: big.NewInt(viper.GetInt64("gasPrice") * params.GWei),
		Gas:      viper.GetUint64("gasLimit"),
		To:       &toAddress,
		Value:    value,
	})

	if err != nil {
		panic(err)
	}

	err = e.getClient(chain).SendTransaction(ctx, tx)
	if err != nil {
		panic(err)
	}

	log.Printf("Transaction has been sent, tx hash: %s", tx.Hash().Hex())
}

func (e *ethRunner) getClient(chain string) *ethclient.Client {
	client, ok := e.ethClientMap[chain]
	if ok {
		return client
	}
	var rpcAddress string
	switch chain {
	case consts.ChainTypeBsc:
		rpcAddress = consts.BscRpcAddr
	case consts.ChainTypeEth:
		rpcAddress = consts.EthRpcAddr
	default:
		panic("not support chain")
	}

	client, err := ethclient.Dial(rpcAddress)
	if err != nil {
		panic(err)
	}
	e.ethClientMap[chain] = client
	return client
}

func (e *ethRunner) SniperUniCake(chain string) {
	factory, err := uniswap.NewUniswapV2(common.HexToAddress(getFactoryAddressByChain(chain)), e.getClient(chain))
	if err != nil {
		panic(err)
	}

	wrapperTokenAddress := getWrapperTokenByChain(chain)
	targetTokenAddress := common.HexToAddress(viper.GetString("targetContract"))
	pairAddress, err := factory.GetPair(nil, wrapperTokenAddress, targetTokenAddress)
	if err != nil {
		panic(err)
	}
	interval := time.Duration(viper.GetInt64("sniperInterval"))
	if pairAddress == consts.ZeroAddress {
		log.Printf("pair not create, retry in %d ms", interval)
		time.AfterFunc(interval*time.Millisecond, func() {
			e.SniperUniCake(chain)
		})
		return
	}

	ethToken, err := uniswap.NewUniswapV2(wrapperTokenAddress, e.getClient(chain))
	if err != nil {
		panic(err)
	}

	minPoolLiquidityAdded, _ := big.NewFloat(viper.GetFloat64("minPoolLiquidityAdded") * params.Ether).Int(nil)
	for {
		balance, err := ethToken.BalanceOf(nil, pairAddress)
		if err != nil {
			panic(err)
		}

		if balance.Cmp(minPoolLiquidityAdded) >= 1 {
			log.Printf("pool liquidity %s", utils.WeiToEtherFloatByDecimals(18, balance).String())
			break
		}

		time.Sleep(interval * time.Millisecond)
	}

	////
	path := []common.Address{wrapperTokenAddress, targetTokenAddress}
	amountIn, _ := big.NewFloat(viper.GetFloat64("buyingBnbOrEthAmount") * params.Ether).Int(nil)
	router, err := uniswap.NewUniswapV2(getRouterAddressByChain(chain), e.getClient(chain))
	if err != nil {
		panic(err)
	}
	amountOutMin := big.NewInt(0)
	if slippage := viper.GetInt64("slippage"); slippage != 0 && slippage < 100 {
		amounts, err := router.GetAmountsOut(nil, amountIn, path)
		if err != nil {
			panic(err)
		}
		amountOutMin = new(big.Int).Div(new(big.Int).Mul(amounts[0], big.NewInt(100-slippage)), big.NewInt(100))
	}

	auth, err := bind.NewKeyedTransactorWithChainID(e.privateKey, big.NewInt(consts.ChainIdMap[chain]))
	if err != nil {
		panic(err)
	}
	auth.Value = amountIn
	auth.GasLimit = viper.GetUint64("gasLimit")
	auth.GasPrice = big.NewInt(viper.GetInt64("gasPrice") * params.GWei)

	tx, err := router.SwapExactETHForTokens(auth, amountOutMin, path, e.fromAddress, big.NewInt(time.Now().Add(2*time.Minute).Unix()))
	if err != nil {
		panic(err)
	}

	log.Printf("Transaction has been sent, tx hash: %s", tx.Hash().Hex())
}

func getFactoryAddressByChain(chain string) string {
	switch chain {
	case consts.ChainTypeEth:
		return consts.UniSwapV2FactoryAddress
	case consts.ChainTypeBsc:
		return consts.PancakeFactoryAddress
	}
	return ""
}

func getRouterAddressByChain(chain string) common.Address {
	var address string
	switch chain {
	case consts.ChainTypeEth:
		address = consts.UniSwapV2RouterAddress
	case consts.ChainTypeBsc:
		address = consts.PancakeRouterAddress
	}
	return common.HexToAddress(address)
}

func getWrapperTokenByChain(chain string) common.Address {
	var address string
	switch chain {
	case consts.ChainTypeEth:
		address = consts.WETHAddress
	case consts.ChainTypeBsc:
		address = consts.WBNBAddress
	}

	return common.HexToAddress(address)
}

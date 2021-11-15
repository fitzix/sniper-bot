package runner

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/fitzix/sniper-bot/consts"
	"github.com/fitzix/sniper-bot/contract/uniswap"
	"github.com/spf13/viper"
	"log"
	"math/big"
	"strings"
	"time"
)

type ethRunner struct {
	privateKey   *ecdsa.PrivateKey
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

	return &ethRunner{
		privateKey:   privateKey,
		uniAbi:       uniAbi,
		ethClientMap: make(map[string]*ethclient.Client),
	}
}

func (e *ethRunner) SniperDxsale(chain string) {
	dxsaleContractAddress := common.HexToAddress(viper.GetString("dxsale.contract"))
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
		time.Sleep(time.Duration(interval) * time.Millisecond)
		e.SniperDxsale(chain)
	}

	if viper.GetUint64("gasLimit") < gas {
		log.Println("config gas limit less than estimate gas ", gas, "auto set to estimate gasLimit")
		viper.Set("gasLimit", gas)
	}

	log.Println("EstimateGas", gas, "ready to transfer")

	e.transfer(ctx, chain, dxsaleContractAddress, value)
}

func (e *ethRunner) transfer(ctx context.Context, chain string, toAddress common.Address, value *big.Int) {

	//
	publicKeyECDSA, ok := e.privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := e.getClient(chain).PendingNonceAt(ctx, fromAddress)
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

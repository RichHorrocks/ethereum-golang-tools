package ethTxn

import (
    "github.com/DaveAppleton/ether_go/ethIpc"
    "github.com/DaveAppleton/ether_go/ethKeys"

    "errors"
    "fmt"
    "math/big"
    "time"

    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/common/math"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/params"
    "golang.org/x/net/context"
)

func SendEthereum(sender *ethKeys.AccountKey, 
                  recipient common.Address, 
                  amountToSend int64) (interface{}, error) {

    var zero interface{}

    myEipc, err := ethIpc.NewEthIpc()
    if err != nil {
        return zero, err
    }
    defer myEipc.Close()

    ec, _ := myEipc.EthClient()

    nonce, err := ec.PendingNonceAt(context.TODO(), sender.PublicKey())
    gasPrice, err := ec.SuggestGasPrice(context.TODO())
    if err != nil {
        return zero, err
    }

    fmt.Println("Nonce : ", nonce)
    fmt.Println("GasPrice : ", gasPrice)
    s := types.NewEIP155Signer(params.TestnetChainConfig.ChainId)

    var amount big.Int
    amount.SetInt64(amountToSend)
    
    var gasLimit big.Int
    gasLimit.SetInt64(121000) // because it is a send - quite standard
    
    data := common.FromHex("0x")
    t := types.NewTransaction(nonce, 
                              recipient, 
                              &amount, 
                              &gasLimit, 
                              gasPrice, 
                              data)

    nt, err := types.SignTx(t, s, sender.GetKey())
    if err != nil {
        return zero, err
    }
    
    err = ec.SendTransaction(context.TODO(), nt)
    
    return nt.Hash(), err

    // rlpEncodedTx, err := rlp.EncodeToBytes(nt)
    // if err != nil {
    //  return zero, err
    // }
    // err = myEipc.Call(&ret, "eth_sendRawTransaction", common.ToHex(rlpEncodedTx))
    // return ret, err
}

type tx struct {
    From     string
    To       string
    Value    interface{}
    Gas      interface{}
    GasPrice interface{}
    Data     string
}

// Estimate the gas required for a contract to run
func estimateGas(sender *ethKeys.AccountKey, 
                 contract string) (big.Int, error) {
    
    var txStruct tx
    var zero big.Int

    zero.SetInt64(0)
    fmt.Println("Estimate Gas")

    myEipc, err := ethIpc.NewEthIpc()
    if err != nil {
        return zero, errors.New("error in IPC")
    }

    txStruct.Data = contract

    var gasLimitStr string
    err = myEipc.Call(&gasLimitStr, "eth_estimateGas", &txStruct)
    if err != nil {
        return zero, err
    }

    fmt.Println("Gastimate: ", gasLimitStr)    

    gasLimit := math.MustParseBig256(gasLimitStr)

    return *gasLimit, nil

}

// PostContract - post a contract to the blockchain.
func PostContract(sender *ethKeys.AccountKey, 
                  contract []byte) (interface{}, error) {

    var zero interface{}

    myEipc, err := ethIpc.NewEthIpc()
    if err != nil {
        return zero, err
    }

    ec, _ := myEipc.EthClient()

    nonce, err := ec.NonceAt(context.TODO(), sender.PublicKey(), nil)
    if err != nil {
        return zero, err
    }

    gasPrice, err := ec.SuggestGasPrice(context.TODO())
    if err != nil {
        return zero, err
    }

    var amountZero big.Int
    amountZero.SetInt64(00)

    var gasLimit big.Int
    gasLimit.SetInt64(90000000)

    cm := ethereum.CallMsg{
        From:     sender.PublicKey(),
        To:       nil,
        Gas:      &gasLimit,
        GasPrice: gasPrice,
        Value:    &amountZero,
        Data:     contract,
    }

    estGas, err := ec.EstimateGas(context.TODO(), cm)
    if err != nil {
        return zero, err
    }

    newContractTx := types.NewContractCreation(nonce, 
                                               &amountZero, 
                                               estGas, 
                                               gasPrice, 
                                               contract)
    nt, err := sender.Sign(newContractTx)

    err = ec.SendTransaction(context.TODO(), nt)
    
    return nt.Hash(), err
}

func WaitForTxnReceipt(txn interface{}) (interface{}, error) {
    var ret interface{}
    var zero interface{}

    myEipc, err := ethIpc.NewEthIpc()
    if err != nil {
        return zero, errors.New("error in IPC")
    }

    fmt.Println()
    count := 100
    err = errors.New("43")
    for err != nil {
        err = myEipc.Call(&ret, "eth_getTransactionReceipt", txn)

        if err == nil {
            break
        }
        fmt.Print(".")
        time.Sleep(500 * time.Millisecond)

        if count < 0 {
            return zero, errors.New("Timeout")

        }
        count--
    }
    fmt.Println(ret)
    
    return ret, nil
}   
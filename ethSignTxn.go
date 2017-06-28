package main

import (
    "bytes"
    "fmt"
    "math/big"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
)

func main() {

    chainId := big.NewInt(3) // ropsten

    senderPrivKey, _ := crypto.HexToECDSA("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
    recipientAddr := common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

    nonce := uint64(7)
    amount := big.NewInt(1000000000000000000) // 1 ether
    gasLimit := big.NewInt(100000)
    gasPrice := big.NewInt(20000000000) // 20 gwei

    tx := types.NewTransaction(nonce, recipientAddr, amount, gasLimit, gasPrice, nil)

    signer := types.NewEIP155Signer(chainId)
    signedTx, _ := types.SignTx(tx, signer, senderPrivKey)
    fmt.Println(signedTx)

    var buff bytes.Buffer
    signedTx.EncodeRLP(&buff)
    fmt.Printf("0x%x\n", buff.Bytes())
}

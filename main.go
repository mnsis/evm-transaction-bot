package main

import (
        "bufio"
        "context"
        "crypto/ecdsa"
        "errors"
        "fmt"
        "log"
        "math/big"
        "math/rand"
        "os"
        "os/exec"
        "strings"
        "sync"
        "time"

        "github.com/ethereum/go-ethereum/common"
        "github.com/ethereum/go-ethereum/core/types"
        "github.com/ethereum/go-ethereum/crypto"
        "github.com/ethereum/go-ethereum/ethclient"
)

type wallet struct {
        privateKey *ecdsa.PrivateKey
        address    common.Address
        nonce      uint64
}

func createWalletFromPrivateKey(privateKeyHex string) (*wallet, error) {
        privateKey, err := crypto.HexToECDSA(privateKeyHex)
        if err != nil {
                return nil, fmt.Errorf("error parsing private key: %w", err)
        }

        publicKey := privateKey.Public()
        publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
        if !ok {
                return nil, errors.New("error converting public key to ECDSA type")
        }

        address := crypto.PubkeyToAddress(*publicKeyECDSA)

        return &wallet{
                privateKey: privateKey,
                address:    address,
        }, nil
}

var nonceMutex sync.Mutex

func main() {
        rand.Seed(time.Now().UnixNano()) // Initialize random seed

        reader := bufio.NewReader(os.Stdin)

        readString := func(prompt string) string {
                fmt.Print(prompt)
                input, _ := reader.ReadString('\n')
                return strings.TrimSpace(input)
        }

        privateKeyHex := readString("Enter EVM Private Key: ")
        url := readString("Enter RPC URL: ")

        wallet, err := createWalletFromPrivateKey(privateKeyHex)
        if err != nil {
                fmt.Println("Error creating wallet:", err)
                return
        }
        fmt.Println(wallet.address)
        chainId := big.NewInt(1234)
        client, err := ethclient.Dial(url)
        if err != nil {
                panic(err)
        }

        balance, err := client.BalanceAt(context.Background(), wallet.address, nil)
        if err != nil {
                panic(err)
        }
        fmt.Println("balance: ", balance)

        nonce, err := client.PendingNonceAt(context.Background(), wallet.address)
        if err != nil {
                panic(err)
        }
        wallet.nonce = nonce

        // Send 25 initial transactions
        fmt.Println("Sending initial 25 transactions...")
        sendTransactionsSequentially(client, wallet, chainId, 25)

        // Start monitoring logs for further transactions
        go monitorStationLogs(client, wallet, chainId)

        select {} // Keep the main function running indefinitely
}

func monitorStationLogs(client *ethclient.Client, wallet *wallet, chainId *big.Int) {
        cmd := exec.Command("sudo", "journalctl", "-u", "stationd", "-f", "--no-hostname", "-o", "cat")
        stdout, err := cmd.StdoutPipe()
        if err != nil {
                log.Fatalf("Failed to get stdout pipe: %v", err)
        }

        if err := cmd.Start(); err != nil {
                log.Fatalf("Failed to start command: %v", err)
        }

        scanner := bufio.NewScanner(stdout)
        for scanner.Scan() {
                line := scanner.Text()
                if strings.Contains(line, "Generating New unverified pods") {
                        fmt.Println("New pod detected, sending 107 transactions...")
                        sendTransactionsSequentially(client, wallet, chainId, 107)
                }
        }

        if err := scanner.Err(); err != nil {
                log.Fatalf("Error reading standard output: %v", err)
        }

        if err := cmd.Wait(); err != nil {
                log.Fatalf("Command finished with error: %v", err)
        }
}

func sendTransactionsSequentially(client *ethclient.Client, wallet *wallet, chainId *big.Int, count int) {
        for i := 0; i < count; i++ {
                nonceMutex.Lock()
                currentNonce := wallet.nonce
                wallet.nonce++
                nonceMutex.Unlock()

                sendTx(client, wallet.address, wallet.address, big.NewInt(int64(rand.Intn(10)+1)), chainId, wallet.privateKey, currentNonce)
                time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000))) // Random delay to mimic human behavior
        }
}

func sendTx(client *ethclient.Client, from, to common.Address, amount *big.Int, chainId *big.Int, pk *ecdsa.PrivateKey, nonce uint64) {
        gasPrice, err := client.SuggestGasPrice(context.Background())
        if err != nil {
                log.Printf("Failed to suggest gas price: %v", err)
                return
        }

        gasLimit := uint64(22000)

        tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)
        signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), pk)
        if err != nil {
                log.Printf("Failed to sign transaction: %v", err)
                return
        }

        err = client.SendTransaction(context.Background(), signedTx)
        if err != nil {
                log.Printf("Failed to send transaction: %v", err)
                return
        }

        fmt.Printf("Tx %d: %s\n", nonce, tx.Hash())
}

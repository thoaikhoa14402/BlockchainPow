package main

import (
	"Lab01/blockchain"
	"Lab01/utils"
	"Lab01/wallet"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type BlockchainService struct {
	port    uint16
	gateway string
}

func NewBlockchainService(port uint16, gateway string) *BlockchainService {
	return &BlockchainService{port, gateway}
}

func (bs *BlockchainService) GetPort() uint16 {
	return bs.port
}

func (bs *BlockchainService) GetGateway() string {
	return bs.gateway
}

func (bs *BlockchainService) CreateTransaction() {
	if UserProfile.PublicKey == "" || UserProfile.BlockchainAddress == "" || UserProfile.PrivateKey == "" {
		fmt.Print("\n\nERROR MESSAGE: Please create your wallet first!")
		return
	}

	if UserProfile.Balance <= 0 {
		fmt.Print("\n\nERROR MESSAGE: Your account balance is insufficient!")
		return
	}

	input := bufio.NewScanner(os.Stdin)
	var recipientAddress string
	fmt.Print("Enter recipient blockchain address: ")
	input.Scan()
	recipientAddress = input.Text()

	var value float64
	fmt.Print("Enter amount of value: ")
	input.Scan()
	value, _ = strconv.ParseFloat(input.Text(), 64)

	var isSignedTransaction string
	var senderPrivateKey string
	fmt.Print("Sign your transaction: (y/yes, n/no): ")
	input.Scan()
	isSignedTransaction = input.Text()

	if isSignedTransaction == "yes" || isSignedTransaction == "y" {
		senderPrivateKey = UserProfile.PrivateKey
	} else {
		senderPrivateKey = ""
		fmt.Print("\n\nERROR MESSAGE: Your transaction has been canceled!")
		return
	}

	var t wallet.TransactionRequest = wallet.TransactionRequest{
		SenderPrivateKey:           senderPrivateKey,
		SenderBlockchainAddress:    UserProfile.BlockchainAddress,
		RecipientBlockchainAddress: recipientAddress,
		SenderPublicKey:            UserProfile.PublicKey,
		Value:                      fmt.Sprintf("%f", value),
	}

	if !t.Validate() {
		log.Print("\n\nERROR MESSAGE: Required Fields are Missing...Try again!")
		return
	}

	publicKey := utils.PublicKeyFromString(t.SenderPublicKey)
	privateKey := utils.PrivateKeyFromString(t.SenderPrivateKey, publicKey)
	value, err := strconv.ParseFloat(t.Value, 32)
	if err != nil {
		log.Println("ERROR: parse error")
		return
	}
	value32 := float32(value)

	transaction := wallet.NewTransaction(privateKey, publicKey,
		t.SenderBlockchainAddress, t.RecipientBlockchainAddress, value32)

	signature := transaction.GenerateSignature()
	signatureStr := signature.String() // convert to string

	bt := &blockchain.TransactionRequest{
		&t.SenderBlockchainAddress,
		&t.RecipientBlockchainAddress,
		&t.SenderPublicKey,
		&value32,
		&signatureStr,
	}

	m, _ := json.Marshal(bt)
	buf := bytes.NewBuffer(m)

	resp, _ := http.Post(bs.GetGateway()+"/transactions", "application/json", buf)

	if resp.StatusCode == 201 {
		fmt.Print("\n\nSUCCESS MESSAGE: Your transaction has been sent successfully...")
		return
	} else {
		fmt.Println("\n\nERROR MESSAGE: Your transaction has failed, please try again...")
		return
	}
}

type Transaction struct {
	SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
	Value                      float32 `json:"value"`
}
type Block struct {
	Timestamp      int64
	Nonce          int
	Transactions   []*Transaction
	PrevBlockHash  []byte
	MerkleRootHash []byte
	Hash           []byte
}

type Balance struct {
	Balance float32 `json:"balance"`
}

func (b *Block) UnmarshalJSON(data []byte) error {
	v := &struct {
		Timestamp      *int64          `json:"timestamp"`
		Nonce          *int            `json:"nonce"`
		Transactions   *[]*Transaction `json:"transactions"`
		PrevBlockHash  *[]byte         `json:"previous_hash"`
		MerkleRootHash *[]byte         `json:"merkle_root_hash"`
		Hash           *[]byte         `json:"hash"`
	}{
		Timestamp:      &b.Timestamp,
		Nonce:          &b.Nonce,
		Transactions:   &b.Transactions,
		PrevBlockHash:  &b.PrevBlockHash,
		MerkleRootHash: &b.MerkleRootHash,
		Hash:           &b.Hash,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

type IUserProfile struct {
	PrivateKey        string
	PublicKey         string
	BlockchainAddress string
	Balance           float32
}

var UserProfile = IUserProfile{
	PrivateKey:        "",
	PublicKey:         "",
	BlockchainAddress: "",
	Balance:           0,
}

// helper function
func printResponseBody(response *http.Response) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}

func PrintBlockchain(blocks []Block) {
	fmt.Print("\n\n")
	for i, block := range blocks {
		fmt.Printf("\n%s Block %d %s\n", strings.Repeat("=", 25), i,
			strings.Repeat("=", 25))
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Previous Block Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Merkle Root Hash: %x\n", block.MerkleRootHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println("List of all transactions:")
		for _, transaction := range block.Transactions {
			// Show list transactions ìnformation
			fmt.Printf("%s\n", strings.Repeat("-", 40))
			fmt.Println("Sender:", transaction.SenderBlockchainAddress)
			fmt.Println("Recipient:", transaction.RecipientBlockchainAddress)
			fmt.Println("Value:", transaction.Value)
		}
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func (bs *BlockchainService) PrintUserProfile() {
	bs.CheckAssetBalance(true)
	profileInformation := [][]string{
		[]string{"Your private key", UserProfile.PrivateKey},
		[]string{"Your public key", UserProfile.PublicKey},
		[]string{"Your blockchain address", UserProfile.BlockchainAddress},
		[]string{"Your balance", fmt.Sprintf("%f", UserProfile.Balance)},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Your wallet details"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgHiRedColor, tablewriter.BgBlackColor}, tablewriter.Colors{tablewriter.FgHiRedColor, tablewriter.Bold, tablewriter.BgBlackColor})
	table.SetColumnColor(tablewriter.Colors{tablewriter.FgHiRedColor}, tablewriter.Colors{tablewriter.FgHiRedColor})
	for _, content := range profileInformation {
		table.Append(content)
	}
	table.Render()
}

func (bs *BlockchainService) PrintMenuOptions() {
	menuOptions := [][]string{
		[]string{"1", "Create wallet"},
		[]string{"2", "Create transaction"},
		[]string{"3", "Scan blockchain information"},
		[]string{"4", "Scan transactions in Mempool (transactions that have not been included in a new block)"},
		[]string{"5", "Scan your transaction history"},
		[]string{"6", "Scan transaction history based on specific address"},
		[]string{"7", "Check asset balance based on specific address"},
		[]string{"8", "Request for a new block from the blockchain gateway (used for testing)"},
		[]string{"9", "Request all nodes to start mining automatically"},
		[]string{"10", "Reload an application"},
		[]string{"0", "Exit"},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetColWidth(150)
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_CENTER, // Cột 1 căn giữa
		tablewriter.ALIGN_LEFT,   // Cột 2 căn trái
	})
	table.SetHeader([]string{"Option", "Description"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.BgBlackColor}, tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.BgBlackColor})
	table.SetColumnColor(tablewriter.Colors{tablewriter.FgCyanColor}, tablewriter.Colors{tablewriter.FgCyanColor})
	table.SetFooter([]string{"COPYRIGHT", "Blockchain service provider by Group 1: Distributed-Decentralized-Security "})
	table.SetFooterColor(tablewriter.Colors{tablewriter.FgCyanColor}, tablewriter.Colors{tablewriter.FgCyanColor})
	//table.SetAlignment(tablewriter.ALIGN_LEFT)
	for _, content := range menuOptions {
		table.Append(content)
	}
	table.Render()
}

func (bs *BlockchainService) CreateWallet() {
	myWallet := wallet.NewWallet()
	if myWallet.PublicKey() == nil || myWallet.PrivateKey() == nil || myWallet.BlockchainAddress() == "" {
		fmt.Print("\n\nERROR MESSAGE: Create wallet failed...Please try again!")
		return
	}
	// Create a wallet for a new user
	UserProfile.PrivateKey = myWallet.PrivateKeyStr()
	UserProfile.PublicKey = myWallet.PublicKeyStr()
	UserProfile.BlockchainAddress = myWallet.BlockchainAddress()
	UserProfile.Balance = 0

	// System send 10 coins as a default for new user
	systemWallet := wallet.NewWallet()
	transaction := wallet.NewTransaction(systemWallet.PrivateKey(), systemWallet.PublicKey(),
		"BLOCKCHAIN_SERVICE_PROVIDER", UserProfile.BlockchainAddress, 10.0)

	signature := transaction.GenerateSignature()
	signatureStr := signature.String()
	// Make transaction request
	var systemWalletAddress = "BLOCKCHAIN_SERVICE_PROVIDER"
	var recipientBlockchainAddress = UserProfile.BlockchainAddress
	var systemPublicKey = systemWallet.PublicKeyStr()
	var amount = float32(10.0)

	bt := &blockchain.TransactionRequest{
		&systemWalletAddress,
		&recipientBlockchainAddress,
		&systemPublicKey,
		&amount,
		&signatureStr,
	}

	m, _ := json.Marshal(bt)
	buf := bytes.NewBuffer(m)
	resp, _ := http.Post(bs.GetGateway()+"/transactions", "application/json", buf)
	if resp.StatusCode == 201 {
		fmt.Print("\n\nSUCCESS MESSAGE: Your wallet has been created. You can now begin using our services!")
	}
}

func (bs *BlockchainService) ScanBlockchain() {
	endpoint := fmt.Sprintf("%s/chain", bs.GetGateway())
	client := &http.Client{}
	newRequest, _ := http.NewRequest("GET", endpoint, nil)
	response, err := client.Do(newRequest)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}

	if response.StatusCode == 200 {
		var data, _ = ioutil.ReadAll(response.Body)
		var chainResp struct {
			Chain []Block `json:"chain"`
		}

		if err := json.Unmarshal(data, &chainResp); err != nil {
			fmt.Println("Error decoding response:", err)
			return
		}
		if err != nil {
			fmt.Printf("\n\nERROR MESSAGE: %v", err)
			return
		}
		PrintBlockchain(chainResp.Chain)
		fmt.Println("THE TOTAL NUMBER OF BLOCKS: ", len(chainResp.Chain))
	}

}

func (bs *BlockchainService) ScanTransactionInMemPool() {
	endpoint := fmt.Sprintf("%s/transactions", bs.GetGateway())
	client := &http.Client{}
	newRequest, _ := http.NewRequest("GET", endpoint, nil)
	response, err := client.Do(newRequest)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}

	if response.StatusCode == 200 {
		var data, _ = ioutil.ReadAll(response.Body)
		var transactionList struct {
			Transactions []Transaction `json:"transactions"`
			Length       int           `json:"length"`
		}
		if err := json.Unmarshal(data, &transactionList); err != nil {
			fmt.Println("Error decoding response:", err)
			return
		}
		if err != nil {
			fmt.Printf("\n\nERROR MESSAGE: %v", err)
			return
		}

		if transactionList.Length == 0 {
			fmt.Printf("\nSTATUS MESSAGE: Mempool is empty now!")
			return
		}

		fmt.Println("\nLIST TRANSACTIONS IN MEMPOOL: ")
		for i, transaction := range transactionList.Transactions {
			fmt.Printf("\n%s Transaction %d %s\n", strings.Repeat("=", 25), i+1,
				strings.Repeat("=", 25))
			//fmt.Printf("%s\n", strings.Repeat("-", 40))
			fmt.Println("Sender:", transaction.SenderBlockchainAddress)
			fmt.Println("Recipient:", transaction.RecipientBlockchainAddress)
			fmt.Println("Value:", transaction.Value)
		}
	}
}

func (bs *BlockchainService) BlockMining() {
	endpoint := fmt.Sprintf("%s/mine", bs.GetGateway())
	client := &http.Client{}
	newRequest, _ := http.NewRequest("GET", endpoint, nil)
	response, err := client.Do(newRequest)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}
	if response.StatusCode == http.StatusOK {
		fmt.Printf("\n\nSUCCESS MESSAGE: Blockchain Node is mining a new block")
	}
}

func (bs *BlockchainService) BlockAutoMining() {
	endpoint := fmt.Sprintf("%s/mine/auto", bs.GetGateway())
	client := &http.Client{}
	newRequest, _ := http.NewRequest("GET", endpoint, nil)
	response, err := client.Do(newRequest)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}
	if response.StatusCode == http.StatusOK {
		fmt.Printf("\n\nSUCCESS MESSAGE: Blockchain Node is mining a new block")
	}
}

func (bs *BlockchainService) CheckAssetBalance(mySelf bool) {
	var blockchainAddress string
	if mySelf == true { // check asset balance of client user
		blockchainAddress = UserProfile.BlockchainAddress
	} else {
		fmt.Printf("Enter your blockchain address: ")
		fmt.Scanf("%s", &blockchainAddress)
	}

	endpoint := fmt.Sprintf("%s/balance", bs.GetGateway())

	client := &http.Client{}
	newRequest, _ := http.NewRequest("GET", endpoint, nil)
	q := newRequest.URL.Query()
	q.Add("blockchain_address", blockchainAddress)
	newRequest.URL.RawQuery = q.Encode()

	response, err := client.Do(newRequest)
	if err != nil {
		fmt.Print("\n\nERROR MESSSAGE: Unexpeced error when scanning asset balance!")
	}

	if response.StatusCode == 200 {
		decoder := json.NewDecoder(response.Body)
		var accountBalance Balance
		err := decoder.Decode(&accountBalance)
		if err != nil {
			fmt.Print("ERROR MESSAGE: %v", err)
		}

		if mySelf == true {
			log.Println("account balance response: ", accountBalance.Balance)
			UserProfile.Balance = accountBalance.Balance
		} else {
			fmt.Printf("\n\nThis account balance: %f", accountBalance.Balance)
		}

	} else {
		fmt.Print("\n\nERROR MESSSAGE: Address Not Found!")
	}
}

func (bs *BlockchainService) RunCli() {
	for {
		fmt.Print("\n")
		bs.PrintUserProfile()
		bs.PrintMenuOptions()
		fmt.Printf("Enter your choice: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print("\n\nError reading input:", err)
			continue
		}
		// Xóa kí tự xuống dòng từ chuỗi đầu vào
		input = strings.TrimSpace(input)
		// Chuyển đổi chuỗi thành số nguyên
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Print("\n\nERROR MESSAGE: Invalid input. Please enter a valid option!")
			continue
		}
		switch choice {
		case 1:
			bs.CreateWallet()
		case 2:
			bs.CreateTransaction()
		case 3:
			bs.ScanBlockchain()
		case 4:
			bs.ScanTransactionInMemPool()
		case 5:
			// doOperation()
		case 6:
			// doOperation()
		case 7:
			bs.CheckAssetBalance(false)
		case 8:
			bs.BlockMining()
		case 9:
			bs.BlockAutoMining()
		case 10:
			fmt.Print("\n\nSTATUS MESSAGE: Reloaded successfully...")
		case 0:
			fmt.Println("Exiting...")
			os.Exit(0)
			return
		default:
			if choice != 10 {
				fmt.Print("\n\nERROR MESSAGE: Invalid choice. Please enter a valid option!")
			}
			time.Sleep(2 * time.Second)
		}
	}
}

func init() {
	log.SetPrefix("Wallet Server: ")
}

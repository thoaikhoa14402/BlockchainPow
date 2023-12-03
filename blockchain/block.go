package blockchain

import (
	"Lab01/utils"
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	SYSTEM_SENDER     = "MY_BLOCKCHAIN"
	MINING_REWARD     = 1
	MINING_TIMER      = 20

	BLOCKCHAIN_SERVER_PORT_RANGE_START = 5000
	BLOCKCHAIN_SERVER_PORT_RANGE_END   = 5003
	PEER_NODE_IP_RANGE_START           = 0
	PEER_NODE_IP_RANGE_END             = 0
	BLOCKCHAIN_PEER_NODE_SYNC_TIMER    = 10
)

// ======= Block =======

type Block struct {
	Timestamp      int64
	Nonce          int
	Transactions   []*Transaction
	PrevBlockHash  []byte
	MerkleRootHash []byte
	Hash           []byte
}

func (b *Block) SetHash() {
	// 1. Declare headers buffer for writing data to
	var headersBuffer bytes.Buffer
	// 2. Write previous block hash to the headers buffer
	headersBuffer.Write(b.PrevBlockHash)
	// 3. Write merkle root hash to the headers buffer
	headersBuffer.Write(b.MerkleRootHash)
	//// Loop through the Transactions array to write each transaction data to the headers buffer
	//for _, t := range b.Transactions {
	//	headersBuffer.Write(t.Data)
	//}
	// 4. Write timestamp to headers buffer
	headersBuffer.WriteString(strconv.FormatInt(b.Timestamp, 10))
	// 5. Write nonce to headers buffer
	headersBuffer.WriteString(strconv.Itoa(b.Nonce))
	// 6. Calculate sha256 of bytes slice of headersBuffer
	hash := sha256.Sum256(headersBuffer.Bytes())
	// 7. Set result to this 'Hash' field of this block
	b.Hash = hash[:]
}

func NewBlock(nonce int, transactions []*Transaction, previousBlockHash []byte, merkleRootHash []byte) *Block {
	block := &Block{time.Now().Unix(), nonce, []*Transaction(transactions), previousBlockHash, merkleRootHash, []byte{}}
	block.SetHash()
	return block
}

func NewGenesisBlock() *Block {
	transactionList := []*Transaction{
		NewTransaction("MY_BLOCKCHAIN", "GENESIS", 0),
	}
	genesisBlock := NewBlock(0, transactionList, []byte{}, []byte{})
	genesisBlock.SetHash()
	return genesisBlock
}

func (b *Block) GetPreviousHash() []byte {
	return b.PrevBlockHash
}

func (b *Block) GetMerkleRootHash() []byte {
	return b.MerkleRootHash
}

func (b *Block) GetHash() []byte {
	return b.Hash
}

func (b *Block) GetNonce() int {
	return b.Nonce
}

func (b *Block) GetTransactions() []*Transaction {
	return b.Transactions
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp      int64          `json:"timestamp"`
		Nonce          int            `json:"nonce"`
		Transactions   []*Transaction `json:"transactions"`
		PrevBlockHash  []byte         `json:"previous_hash"`
		MerkleRootHash []byte         `json:"merkle_root_hash"`
		Hash           []byte         `json:"hash"`
	}{
		Timestamp:      b.Timestamp,
		Nonce:          b.Nonce,
		Transactions:   b.Transactions,
		PrevBlockHash:  b.PrevBlockHash,
		MerkleRootHash: b.MerkleRootHash,
		Hash:           b.Hash,
	})
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

// ====== Transaction =======

type Transaction struct {
	Data []byte
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	type TransactionInfo struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}
	var transInfo TransactionInfo
	err := json.Unmarshal(t.Data, &transInfo)
	if err != nil {
		fmt.Println("Unexpected error in MarshalJSON of Transaction")
	}
	return json.Marshal(transInfo)
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	t.Data = data
	return nil
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	// convert this transaction struct to byte array represented for JSON data
	data, _ := json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    sender,
		Recipient: recipient,
		Value:     value,
	})
	return &Transaction{Data: data}
}

// ====== Merkle Tree =======

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

type MerkleTree struct {
	Root *MerkleNode
}

func CreateMerkleNode(new_data []byte) *MerkleNode {
	return &MerkleNode{nil, nil, new_data}
}

func CreateMerkleTree(transactionList []*Transaction) *MerkleNode {
	queue_MerkleNode := make([]*MerkleNode, 0)
	queue_height := make([]int, 0)

	// Tạo queue chứa các node lá của Merkle Tree
	for _, tran := range transactionList {
		hash_val := sha256.Sum256(tran.Data)
		leaf_node := CreateMerkleNode(hash_val[:])

		queue_MerkleNode = append(queue_MerkleNode, leaf_node)
		queue_height = append(queue_height, 0)
	}

	// Thực hiện cho tới khi queue còn 1 phần tử => là Root của cây
	for len(queue_MerkleNode) != 1 {
		val1 := queue_MerkleNode[0]
		h1 := queue_height[0]
		val2 := queue_MerkleNode[1]
		h2 := queue_height[1]

		if h1 != h2 {
			// Xoá node sử dụng ra khỏi queue
			queue_MerkleNode = queue_MerkleNode[1:]
			queue_height = queue_height[1:]

			// Tính giá trị hash của node cha
			hash_val := sha256.Sum256(bytes.Join([][]byte{val1.Data, val1.Data}, []byte{}))

			// Tạo node cha
			nonleaf_node := CreateMerkleNode(hash_val[:])
			nonleaf_node.Left = val1
			nonleaf_node.Right = val1

			// Đưa node mới vào queue
			queue_MerkleNode = append(queue_MerkleNode, nonleaf_node)
			queue_height = append(queue_height, h1+1)
		} else {
			queue_MerkleNode = queue_MerkleNode[2:]
			queue_height = queue_height[2:]

			hash_val := sha256.Sum256(bytes.Join([][]byte{val1.Data, val2.Data}, []byte{}))

			nonleaf_node := CreateMerkleNode(hash_val[:])
			nonleaf_node.Left = val1
			nonleaf_node.Right = val2

			queue_MerkleNode = append(queue_MerkleNode, nonleaf_node)
			queue_height = append(queue_height, h1+1)
		}
	}

	return queue_MerkleNode[0]
}

// ====== Blockchain =======

type Blockchain struct {
	blocks            []*Block
	transactionPool   []*Transaction
	blockchainAddress string
	port              uint16
	peerNodes         []string
	mux               sync.Mutex
	muxPeerNodes      sync.Mutex
}

func NewBlockchain(blockchainAddress string, port uint16) *Blockchain {
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.port = port
	// add a first genesis block
	b := NewGenesisBlock()
	bc.blocks = append(bc.blocks, b)
	return bc
}

func (bc *Blockchain) Run() {
	bc.StartSyncPeerNodes()
	bc.HandleConflicts()
}

func (bc *Blockchain) GetChain() []*Block {
	return bc.blocks
}

func (bc *Blockchain) SetPeerNodes() {
	bc.peerNodes = utils.FindPeerNodes("127.0.0.1", bc.port, PEER_NODE_IP_RANGE_START, PEER_NODE_IP_RANGE_END,
		BLOCKCHAIN_SERVER_PORT_RANGE_START, BLOCKCHAIN_SERVER_PORT_RANGE_END)
	log.Printf("%v", bc.peerNodes)
}

func (bc *Blockchain) SyncPeerNodes() {
	bc.muxPeerNodes.Lock()
	defer bc.muxPeerNodes.Unlock()
	bc.SetPeerNodes()
}

func (bc *Blockchain) StartSyncPeerNodes() {
	bc.SyncPeerNodes()
	_ = time.AfterFunc(time.Second*BLOCKCHAIN_PEER_NODE_SYNC_TIMER, bc.SyncPeerNodes)
}

func (bc *Blockchain) GetTransactionPool() []*Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) ClearTransactionPool() {
	bc.transactionPool = bc.transactionPool[:0]
}

func (bc *Blockchain) AddBlock(nonce int, previousHash []byte, merkleRootHash []byte) *Block {
	b := NewBlock(nonce, bc.transactionPool, previousHash, merkleRootHash)
	bc.blocks = append(bc.blocks, b)
	bc.transactionPool = []*Transaction{}
	// When creating a new block, transactions in the pool will be empty, and the other nodes will also be emptied
	for _, n := range bc.peerNodes {
		endpoint := fmt.Sprintf("http://%s/transactions", n)
		client := &http.Client{}
		req, _ := http.NewRequest("DELETE", endpoint, nil)
		resp, _ := client.Do(req)
		log.Printf("%v", resp)
	}
	return b
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		var transactionInfo struct {
			Sender    string  `json:"sender_blockchain_address"`
			Recipient string  `json:"recipient_blockchain_address"`
			Value     float32 `json:"value"`
		}
		err := json.Unmarshal(t.Data, &transactionInfo)
		if err != nil {
			fmt.Println("Unexpected error in copy transaction pool")
			return nil
		}
		transactions = append(transactions, NewTransaction(transactionInfo.Sender, transactionInfo.Recipient, transactionInfo.Value))
	}
	return transactions
}

func (bc *Blockchain) GetLastBlock() *Block {
	return bc.blocks[len(bc.blocks)-1]
}

func (bc *Blockchain) ValidProof(nonce int, previousBlockHash []byte, merkleRootHash []byte, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := &Block{0, nonce, nil, previousBlockHash, merkleRootHash, []byte{}}
	guessBlock.SetHash()
	guessHashStr := fmt.Sprintf("%x", guessBlock.GetHash())
	return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) CreateTransaction(sender string, recipient string, value float32,
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	isTransactionAdded := bc.AddTransaction(sender, recipient, value, senderPublicKey, s)
	// Add a new transaction to transactions pool and then synchronize to other nodes
	if isTransactionAdded {
		for _, n := range bc.peerNodes {
			publicKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(), senderPublicKey.Y.Bytes())
			signatureStr := s.String()
			bt := &TransactionRequest{&sender, &recipient,
				&publicKeyStr, &value, &signatureStr}

			m, _ := json.Marshal(bt)
			buf := bytes.NewBuffer(m)

			endpoint := fmt.Sprintf("http://%s/transactions", n)
			client := &http.Client{}
			req, _ := http.NewRequest("PUT", endpoint, buf)
			resp, _ := client.Do(req)
			log.Printf("%v", resp)
		}
	}
	return isTransactionAdded
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32,
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, recipient, value)
	if sender == SYSTEM_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}
	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("ERROR: Verify Transaction")
	}
	return false
}

func (bc *Blockchain) VerifyTransactionSignature(
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	h := sha256.Sum256(t.Data)
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.GetLastBlock().GetHash()
	nonce := 0
	merkleTreeRoot := CreateMerkleTree(transactions)
	for !bc.ValidProof(nonce, previousHash, merkleTreeRoot.Data, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) CalculateTotalBalance(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.blocks {
		for _, t := range b.Transactions {
			type TransactionInfo struct {
				Sender    string  `json:"sender_blockchain_address"`
				Recipient string  `json:"recipient_blockchain_address"`
				Value     float32 `json:"value"`
			}
			var transInfo TransactionInfo
			err := json.Unmarshal(t.Data, &transInfo)
			if err != nil {
				log.Println(err)
			}
			if transInfo.Recipient == blockchainAddress {
				totalAmount += transInfo.Value
			}
			if transInfo.Sender == blockchainAddress {
				totalAmount -= transInfo.Value
			}
		}
	}
	return totalAmount
}

func (bc *Blockchain) Print() {
	for i, block := range bc.blocks {
		fmt.Printf("\n%s Block %d %s\n", strings.Repeat("=", 25), i,
			strings.Repeat("=", 25))
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Previous Block Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println("List of all transactions:")
		for _, transaction := range block.Transactions {
			var transactionInfo struct {
				Sender    string  `json:"sender_blockchain_address"`
				Recipient string  `json:"recipient_blockchain_address"`
				Value     float32 `json:"value_blockchain_address"`
			}
			err := json.Unmarshal(transaction.Data, &transactionInfo)
			if err != nil {
				fmt.Println("Error decoding transaction data:", err)
				return
			}
			// Show list transactions ìnformation
			fmt.Printf("%s\n", strings.Repeat("-", 40))
			fmt.Println("Sender:", transactionInfo.Sender)
			fmt.Println("Recipient:", transactionInfo.Recipient)
			fmt.Println("Value:", transactionInfo.Value)
		}
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func (bc *Blockchain) Mining() bool {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	bc.AddTransaction(SYSTEM_SENDER, bc.blockchainAddress, MINING_REWARD, nil, nil) // coinbase transaction does not need to verify signature
	nonce := bc.ProofOfWork()
	previousHash := bc.GetLastBlock().GetHash()
	transactions := bc.CopyTransactionPool()
	merkleTreeRoot := CreateMerkleTree(transactions)
	bc.AddBlock(nonce, previousHash, merkleTreeRoot.Data)

	for _, n := range bc.peerNodes {
		endpoint := fmt.Sprintf("http://%s/consensus", n)
		client := &http.Client{}
		req, _ := http.NewRequest("PUT", endpoint, nil)
		resp, _ := client.Do(req)
		log.Printf("%v", resp)
	}

	return true
}

func (bc *Blockchain) StartMining() {
	bc.Mining()
	_ = time.AfterFunc(time.Second*MINING_TIMER, bc.StartMining)
}

func (bc *Blockchain) ValidChain(chain []*Block) bool {
	preBlock := chain[0]
	currentIndex := 1
	for currentIndex < len(chain) {
		b := chain[currentIndex]
		if !bytes.Equal(b.GetPreviousHash(), preBlock.GetHash()) {
			return false
		}
		merkleRoot := CreateMerkleTree(b.GetTransactions())
		if !bc.ValidProof(b.GetNonce(), b.GetPreviousHash(), merkleRoot.Data, MINING_DIFFICULTY) {
			return false
		}
		preBlock = b
		currentIndex += 1
	}
	return true
}

func (bc *Blockchain) HandleConflicts() bool {
	var longestChain []*Block = nil
	maxLength := len(bc.blocks)

	for _, n := range bc.peerNodes {
		endpoint := fmt.Sprintf("http://%s/chain", n)
		resp, _ := http.Get(endpoint)
		var data, _ = ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			var chainResp struct {
				Chain []*Block `json:"chain"`
			}

			if err := json.Unmarshal(data, &chainResp); err != nil {
				fmt.Println("Error decoding response:", err)
				continue // Skip to the next iteration
			}

			if len(chainResp.Chain) > maxLength && bc.ValidChain(chainResp.Chain) {
				maxLength = len(chainResp.Chain)
				longestChain = chainResp.Chain
			}
		}
	}

	if longestChain != nil {
		bc.blocks = longestChain
		log.Printf("Conflicts have been replaced")
		return true
	}
	log.Printf("Conflicts have been not replaced")
	return false
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chain"`
	}{
		Blocks: bc.blocks,
	})
}

// ------------------------- Transaction Request From Client -------------------------

type TransactionRequest struct {
	SenderBlockchainAddress    *string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string  `json:"recipient_blockchain_address"`
	SenderPublicKey            *string  `json:"sender_public_key"`
	Value                      *float32 `json:"value"`
	Signature                  *string  `json:"signature"`
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderBlockchainAddress == nil || tr.RecipientBlockchainAddress == nil ||
		tr.SenderPublicKey == nil || tr.Value == nil || tr.Signature == nil {
		return false
	}
	return true
}

// ------------------------- Balance Response For Client -------------------------

type BalanceResponse struct {
	Balance float32 `json:"balance"`
}

func (br *BalanceResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Balance float32 `json:"balance"`
	}{
		Balance: br.Balance,
	})
}

package main

//
//import (
//	"Lab01/blockchain"
//	"Lab01/utils"
//	"Lab01/wallet"
//	"bytes"
//	"encoding/json"
//	"flag"
//	"fmt"
//	"github.com/olekukonko/tablewriter"
//	"io"
//	"log"
//	"net/http"
//	"os"
//	"strconv"
//)
//
//type WalletServerCli struct {
//	port    uint16
//	gateway string
//}
//
//func NewWalletServerCli(port uint16, gateway string) *WalletServerCli {
//	return &WalletServerCli{port, gateway}
//}
//
//func (ws *WalletServer) PortCli() uint16 {
//	return ws.port
//}
//
//func (ws *WalletServer) GatewayCli() string {
//	return ws.gateway
//}
//
////func (ws *WalletServer) IndexCli(w http.ResponseWriter, req *http.Request) {
////	switch req.Method {
////	case http.MethodGet:
////		t, _ := template.ParseFiles(path.Join(tempDir, "index.html"))
////		t.Execute(w, "")
////	default:
////		log.Printf("ERROR: Invalid HTTP Method")
////	}
////}
//
//func (ws *WalletServer) WalletCli(w http.ResponseWriter, req *http.Request) {
//	switch req.Method {
//	case http.MethodPost:
//		w.Header().Add("Content-Type", "application/json")
//		myWallet := wallet.NewWallet()
//		m, _ := myWallet.MarshalJSON()
//		io.WriteString(w, string(m[:]))
//	default:
//		w.WriteHeader(http.StatusBadRequest)
//		log.Println("ERROR: Invalid HTTP Method")
//	}
//}
//
//func (ws *WalletServer) CreateTransactionCli(w http.ResponseWriter, req *http.Request) {
//	switch req.Method {
//	case http.MethodPost:
//		decoder := json.NewDecoder(req.Body)
//		var t wallet.TransactionRequest
//		err := decoder.Decode(&t)
//		if err != nil {
//			log.Printf("ERROR: %v", err)
//			io.WriteString(w, string(utils.JsonStatus("fail")))
//			return
//		}
//		if !t.Validate() {
//			log.Println("ERROR: Missing field(s)")
//			io.WriteString(w, string(utils.JsonStatus("fail")))
//			return
//		}
//
//		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
//		privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
//		value, err := strconv.ParseFloat(*t.Value, 32)
//		if err != nil {
//			log.Println("ERROR: parse error")
//			io.WriteString(w, string(utils.JsonStatus("fail")))
//			return
//		}
//
//		value32 := float32(value)
//		w.Header().Add("Content-Type", "application/json")
//
//		transaction := wallet.NewTransaction(privateKey, publicKey,
//			*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, value32)
//
//		signature := transaction.GenerateSignature()
//		signatureStr := signature.String() // convert to string
//
//		bt := &blockchain.TransactionRequest{
//			t.SenderBlockchainAddress,
//			t.RecipientBlockchainAddress,
//			t.SenderPublicKey,
//			&value32,
//			&signatureStr,
//		}
//
//		m, _ := json.Marshal(bt)
//		buf := bytes.NewBuffer(m)
//
//		resp, _ := http.Post(ws.Gateway()+"/transactions", "application/json", buf)
//
//		if resp.StatusCode == 201 {
//			io.WriteString(w, string(utils.JsonStatus("success")))
//			return
//		}
//
//		io.WriteString(w, string(utils.JsonStatus("fail")))
//
//	default:
//		w.WriteHeader(http.StatusBadRequest)
//		log.Println("ERROR: Invalid HTTP Method")
//	}
//}
//
//func (ws *WalletServer) WalletAmountCli(w http.ResponseWriter, req *http.Request) {
//	switch req.Method {
//	case http.MethodGet:
//		blockchainAddress := req.URL.Query().Get("blockchain_address")
//		endpoint := fmt.Sprintf("%s/amount", ws.Gateway())
//
//		client := &http.Client{}
//		bcsReq, _ := http.NewRequest("GET", endpoint, nil)
//		q := bcsReq.URL.Query()
//		q.Add("blockchain_address", blockchainAddress)
//		bcsReq.URL.RawQuery = q.Encode()
//
//		bcsResp, err := client.Do(bcsReq)
//		if err != nil {
//			log.Printf("ERROR: %v", err)
//			io.WriteString(w, string(utils.JsonStatus("fail")))
//			return
//		}
//
//		w.Header().Add("Content-Type", "application/json")
//		if bcsResp.StatusCode == 200 {
//			decoder := json.NewDecoder(bcsResp.Body)
//			var bar blockchain.AmountResponse
//			err := decoder.Decode(&bar)
//			if err != nil {
//				log.Printf("ERROR: %v", err)
//				io.WriteString(w, string(utils.JsonStatus("fail")))
//				return
//			}
//
//			m, _ := json.Marshal(struct {
//				Message string  `json:"message"`
//				Amount  float32 `json:"amount"`
//			}{
//				Message: "success",
//				Amount:  bar.Amount,
//			})
//
//			io.WriteString(w, string(m[:]))
//		} else {
//			io.WriteString(w, string(utils.JsonStatus("fail")))
//		}
//	default:
//		w.WriteHeader(http.StatusBadRequest)
//		log.Println("ERROR: Invalid HTTP Method")
//	}
//}
//func printUserProfile() {
//	privateKey := "53d29a7d389680a3a799133e642e8ee8b20c878e7117d7918d5c652326883e46"
//	publicKey := "c46fe426f439428d8c0c46bff2266f80970dadc382347437ababeeafff2ccdb2f5bacd8d0caf2e2598218361504aa55c9ccde1bacadf4b8773735d50c3a962b5"
//	address := "12s2Kof7dwC4qhNwbQSRYoqcPr8o7N4mLD"
//	balance := 14.67
//
//	profileInformation := [][]string{
//		[]string{"Your private key", privateKey},
//		[]string{"Your public key", publicKey},
//		[]string{"Your blockchain address", address},
//		[]string{"Your balance", fmt.Sprintf("%f", balance)},
//	}
//	table := tablewriter.NewWriter(os.Stdout)
//	table.SetHeader([]string{"Field", "Your wallet details"})
//	table.SetAlignment(tablewriter.ALIGN_LEFT)
//	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgHiRedColor, tablewriter.BgBlackColor}, tablewriter.Colors{tablewriter.FgHiRedColor, tablewriter.Bold, tablewriter.BgBlackColor})
//	table.SetColumnColor(tablewriter.Colors{tablewriter.FgHiRedColor}, tablewriter.Colors{tablewriter.FgHiRedColor})
//	for _, content := range profileInformation {
//		table.Append(content)
//	}
//	table.Render()
//}
//
//func printMenuOptions() {
//	menuOptions := [][]string{
//		[]string{"1", "Create wallet"},
//		[]string{"2", "Create transaction"},
//		[]string{"3", "Look up your transaction history"},
//		[]string{"4", "Look up blockchain information"},
//		[]string{"5", "Check the current block number on the blockchain"},
//		[]string{"6", "List total transactions on blockchain"},
//		[]string{"7", "Exit"},
//	}
//	table := tablewriter.NewWriter(os.Stdout)
//	table.SetColWidth(150)
//	table.SetHeader([]string{"Option", "Description"})
//	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.BgBlackColor}, tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.BgBlackColor})
//	table.SetColumnColor(tablewriter.Colors{tablewriter.FgCyanColor}, tablewriter.Colors{tablewriter.FgCyanColor})
//	table.SetFooter([]string{"COPYRIGHT", "Blockchain service provider by Group 1"})
//	table.SetFooterColor(tablewriter.Colors{tablewriter.FgCyanColor}, tablewriter.Colors{tablewriter.FgCyanColor})
//	table.SetAlignment(tablewriter.ALIGN_CENTER)
//	for _, content := range menuOptions {
//		table.Append(content)
//	}
//	table.Render()
//}
//
//func (ws *WalletServerCli) RunCli() {
//	for {
//		printUserProfile()
//		printMenuOptions()
//		fmt.Printf("Enter your choice: ")
//		var choice int
//		fmt.Scanln(&choice)
//
//		switch choice {
//		case 1:
//			// doOperation()
//		case 2:
//			// doOperation()
//		case 3:
//			// doOperation()
//		case 4:
//			// doOperation()
//		case 5:
//			// doOperation()
//		case 6:
//			// doOperation()
//		case 7:
//			fmt.Println("Exiting...")
//			os.Exit(0)
//			return
//		}
//	}
//}
//
//func init() {
//	log.SetPrefix("Wallet Server: ")
//}
//
//func main() {
//	port := flag.Uint("port", 8080, "TCP Port Number for Wallet Server")
//	gateway := flag.String("gateway", "http://127.0.0.1:5000", "Blockchain Gateway")
//	flag.Parse()
//	app := NewWalletServerCli(uint16(*port), *gateway)
//	app.RunCli()
//}

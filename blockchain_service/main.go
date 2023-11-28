package main

import "flag"

func main() {
	port := flag.Uint("port", 8080, "TCP Port Number for Wallet Server")
	gateway := flag.String("gateway", "http://127.0.0.1:5000", "Blockchain Gateway IP")
	flag.Parse()
	app := NewBlockchainService(uint16(*port), *gateway)
	app.RunCli()
}

package main

import (
	"Lab01/utils"
	"fmt"
)

func main() {
	fmt.Println(utils.FindPeerNodes("127.0.0.1", 5000, 0, 2, 5000, 5003))
}

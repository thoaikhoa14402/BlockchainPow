package utils

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

var IP_PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

func IsFoundHost(host string, port uint16) bool {
	target := fmt.Sprintf("%s:%d", host, port)
	_, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		fmt.Printf("%s %v\n", target, err)
		return false
	}
	return true
}

func FindPeerNodes(myHost string, myPort uint16, startIP uint8, endIP uint8, startPort uint16, endPort uint16) []string {
	address := fmt.Sprintf("%s:%d", myHost, myPort)

	m := IP_PATTERN.FindStringSubmatch(myHost)
	if m == nil {
		return nil
	}
	prefixHost := m[1]
	lastIP, _ := strconv.Atoi(m[len(m)-1])
	peerNodes := make([]string, 0)

	for port := startPort; port <= endPort; port += 1 {
		for ip := startIP; ip <= endIP; ip += 1 {
			guessHost := fmt.Sprintf("%s%d", prefixHost, lastIP+int(ip))
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
			if guessTarget != address && IsFoundHost(guessHost, port) {
				peerNodes = append(peerNodes, guessTarget)
			}
		}
	}
	return peerNodes
}

func GetHost() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "127.0.0.1"
	}
	address, err := net.LookupHost(hostname)
	if err != nil {
		return "127.0.0.1"
	}
	return address[0]

}

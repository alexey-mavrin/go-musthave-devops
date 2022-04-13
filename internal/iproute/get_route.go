// Package iproute provides a way to find the IP address to be used
// for reaching the specific target ip.
// There is also a package https://pkg.go.dev/github.com/google/gopacket/routing
// but unfortunately it gives me an error on a machine with Docker networks
// configured
package iproute

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os/exec"
	"strings"
)

// GetSrcIPURL parses incoming url to find the server name, then uses
// GetSrcIP to get src IP for this request
func GetSrcIPURL(dstURL string) (string, error) {
	u, err := url.Parse(dstURL)
	if err != nil {
		return "", err
	}

	hostPort := u.Host
	server := strings.Split(hostPort, ":")
	return GetSrcIP(server[0])
}

// GetSrcIP returns the src IP address to reach the specific server name
// If server name resolves to multiple IPs, this func returns the result
// only for the first one.
func GetSrcIP(dstAddr string) (string, error) {
	ips, err := net.LookupIP(dstAddr)
	if err != nil {
		return "", err
	}

	if len(ips) == 0 {
		return "", errors.New("IP address list is empty for " + dstAddr)
	}

	return GetSrcIPToIP(ips[0].String())
}

// GetSrcIPToIP returns the src IP address to reach the specific server IP
func GetSrcIPToIP(dstIP string) (string, error) {
	cmdLine := "ip route get " + dstIP
	ipCmdArgs := strings.Fields(cmdLine)
	cmd := exec.Command(ipCmdArgs[0], ipCmdArgs[1:]...)
	var srcIP string

	var o, e bytes.Buffer
	cmd.Stdout = &o
	cmd.Stderr = &e

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s: %s (%w)", cmdLine, e.String(), err)
	}

	outFields := strings.Fields(o.String())
	for i, s := range outFields {
		if s != "src" {
			continue
		}

		if i >= len(outFields)-1 {
			return "", errors.New(cmdLine + ": " +
				o.String() +
				": parsing error")
		}

		srcIP = outFields[i+1]
		break
	}
	return srcIP, nil
}

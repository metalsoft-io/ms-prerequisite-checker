package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func (app *application) testHTTPConnection(ctx context.Context, host string, port int) int {
	slog.Debug(fmt.Sprintf("Testing HTTP connection to %s:%d", host, port))

	response, err := http.Get("http://" + net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for HTTP connection to %s:%d - %s", host, port, err.Error()))
		return 1
	}
	defer response.Body.Close()

	slog.Debug(fmt.Sprintf("HTTP GET %s:%d returned status %s", host, port, response.Status))

	return 0
}

func (app *application) testHTTPSConnection(ctx context.Context, host string, port int) int {
	slog.Debug(fmt.Sprintf("Testing HTTPS connection to %s:%d", host, port))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	response, err := client.Get("https://" + net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for HTTPS connection to %s:%d - %s", host, port, err.Error()))
		return 1
	}
	defer response.Body.Close()

	slog.Debug(fmt.Sprintf("HTTPS GET %s:%d returned status %s", host, port, response.Status))

	return 0
}

func (app *application) testLink(ctx context.Context, url string) int {
	slog.Debug(fmt.Sprintf("Testing connection to %s", url))

	response, err := http.Get(url)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for link %s - %s", url, err.Error()))
		return 1
	}
	defer response.Body.Close()

	slog.Debug(fmt.Sprintf("GET %s returned status %s", url, response.Status))

	return 0
}

func (app *application) testTCPConnection(ctx context.Context, host string, port int) int {
	slog.Debug(fmt.Sprintf("Testing TCP connection to %s:%d", host, port))

	conn, err := net.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for TCP connection to %s:%d - %s", host, port, err.Error()))
		return 1
	}
	defer conn.Close()

	dataIn := []byte("PING")
	bytesWritten, err := conn.Write(dataIn)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for TCP connection to %s:%d - %s", host, port, err.Error()))
		return 1
	}

	slog.Debug(fmt.Sprintf("Wrote %d bytes to TCP %s:%d - %s", bytesWritten, host, port, string(dataIn)))

	dataOut := make([]byte, 1024)
	bytesRead, err := conn.Read(dataOut)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for TCP connection to %s:%d - %s", host, port, err.Error()))
		return 1
	}

	slog.Debug(fmt.Sprintf("Read %d bytes from TCP %s:%d - %s", bytesRead, host, port, string(dataOut[:bytesRead])))

	return 0
}

func (app *application) testUDPConnection(ctx context.Context, host string, port int) int {
	slog.Debug(fmt.Sprintf("Testing UDP connection to %s:%d", host, port))

	conn, err := net.Dial("udp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for UDP connection to %s:%d - %s", host, port, err.Error()))
		return 1
	}
	defer conn.Close()

	dataIn := []byte("PING")
	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for UDP connection to %s:%d - %s", host, port, err.Error()))
		return 1
	}
	bytesWritten, err := conn.Write(dataIn)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for UDP connection to %s:%d - %s", host, port, err.Error()))
		return 1
	}

	slog.Debug(fmt.Sprintf("Wrote %d bytes to UDP %s:%d - %s", bytesWritten, host, port, string(dataIn)))

	dataOut := make([]byte, 1024)
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for UDP connection to %s:%d - %s", host, port, err.Error()))
		return 1
	}
	bytesRead, err := conn.Read(dataOut)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for UDP connection to %s:%d - %s", host, port, err.Error()))
		return 1
	}

	slog.Debug(fmt.Sprintf("Read %d bytes from UDP %s:%d - %s", bytesRead, host, port, string(dataOut[:bytesRead])))

	return 0
}

func (app *application) testSSHConnection(ctx context.Context, host string, port int, username string, password string) int {
	slog.Debug(fmt.Sprintf("Testing SSH connection to %s:%d", host, port))

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)), config)
	if err != nil {
		slog.Error(fmt.Sprintf("Unable to connect to SSH server %s:%d - %s", host, port, err.Error()))
		return 1
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		slog.Error(fmt.Sprintf("Unable to create session with SSH server %s:%d - %s", host, port, err.Error()))
		return 1
	}
	defer session.Close()

	slog.Debug(fmt.Sprintf("Connected to SSH server %s:%d", host, port))

	return 0
}

func (app *application) testICMPConnection(ctx context.Context, host string) int {
	slog.Debug(fmt.Sprintf("Testing ICMP connection to %s", host))

	// Listen for ICMP packets
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for ICMP connection to %s - %s", host, err.Error()))
		return 1
	}
	defer conn.Close()

	// Create an ICMP echo request
	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte(""),
		},
	}

	// Encode the request
	b, err := message.Marshal(nil)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for ICMP connection to %s - %s", host, err.Error()))
		return 1
	}

	// Send the request
	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for ICMP connection to %s - %s", host, err.Error()))
		return 1
	}
	if _, err := conn.WriteTo(b, &net.IPAddr{IP: net.ParseIP(host)}); err != nil {
		slog.Error(fmt.Sprintf("Failed test for ICMP connection to %s - %s", host, err.Error()))
		return 1
	}

	// Wait for a reply
	reply := make([]byte, 1500)
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for ICMP connection to %s - %s", host, err.Error()))
		return 1
	}
	n, peer, err := conn.ReadFrom(reply)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for ICMP connection to %s - %s", host, err.Error()))
		return 1
	}

	// Parse the reply
	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), reply[:n])
	if err != nil {
		slog.Error(fmt.Sprintf("Failed test for ICMP connection to %s - %s", host, err.Error()))
		return 1
	}

	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		slog.Debug(fmt.Sprintf("Got ICMP echo reply from %v", peer))
	default:
		slog.Debug(fmt.Sprintf("Got ICMP %+v reply from %v - expected echo", rm, peer))
		return 1
	}

	return 0
}

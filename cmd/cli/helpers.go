package main

import (
	"context"
	"crypto/tls"
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
	app.logger.Debug().Msgf("Testing HTTP connection to %s:%d", host, port)

	response, err := http.Get("http://" + net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		app.logger.Error().Err(err).Msgf("Failed test for HTTP connection to %s:%d", host, port)
		return 1
	}
	defer response.Body.Close()

	app.logger.Debug().Msgf("HTTP GET %s:%d returned status %s", host, port, response.Status)

	return 0
}

func (app *application) testHTTPSConnection(ctx context.Context, host string, port int) int {
	app.logger.Debug().Msgf("Testing HTTPS connection to %s:%d", host, port)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	response, err := client.Get("https://" + net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		app.logger.Error().Err(err).Msgf("Failed test for HTTPS connection to %s:%d", host, port)
		return 1
	}
	defer response.Body.Close()

	app.logger.Debug().Msgf("HTTPS GET %s:%d returned status %s", host, port, response.Status)

	return 0
}

func (app *application) testLink(ctx context.Context, url string) int {
	app.logger.Debug().Msgf("Testing connection to %s", url)

	response, err := http.Get(url)
	if err != nil {
		app.logger.Error().Err(err).Msgf("Failed test for link %s", url)
		return 1
	}
	defer response.Body.Close()

	app.logger.Debug().Msgf("GET %s returned status %s", url, response.Status)

	return 0
}

func (app *application) testTCPConnection(ctx context.Context, host string, port int) int {
	app.logger.Debug().Msgf("Testing TCP connection to %s:%d", host, port)

	conn, err := net.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		app.logger.Error().Err(err).Msgf("Failed test for TCP connection to %s:%d", host, port)
		return 1
	}
	defer conn.Close()

	dataIn := []byte("PING")
	bytesWritten, err := conn.Write(dataIn)
	if err != nil {
		app.logger.Error().Err(err).Msgf("Failed test for TCP connection to %s:%d", host, port)
		return 1
	}

	app.logger.Debug().Msgf("Wrote %d bytes to TCP %s:%d - %s", bytesWritten, host, port, string(dataIn))

	dataOut := make([]byte, 1024)
	bytesRead, err := conn.Read(dataOut)
	if err != nil {
		app.logger.Error().Err(err).Msgf("Failed test for TCP connection to %s:%d", host, port)
		return 1
	}

	app.logger.Debug().Msgf("Read %d bytes from TCP %s:%d - %s", bytesRead, host, port, string(dataOut[:bytesRead]))

	return 0
}

func (app *application) testUDPConnection(ctx context.Context, host string, port int) int {
	app.logger.Debug().Msgf("Testing UDP connection to %s:%d", host, port)

	conn, err := net.Dial("udp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		app.logger.Error().Err(err).Msgf("Failed test for UDP connection to %s:%d", host, port)
		return 1
	}
	defer conn.Close()

	dataIn := []byte("PING")
	bytesWritten, err := conn.Write(dataIn)
	if err != nil {
		app.logger.Error().Err(err).Msgf("Failed test for UDP connection to %s:%d", host, port)
		return 1
	}

	app.logger.Debug().Msgf("Wrote %d bytes to UDP %s:%d - %s", bytesWritten, host, port, string(dataIn))

	dataOut := make([]byte, 1024)
	bytesRead, err := conn.Read(dataOut)
	if err != nil {
		app.logger.Error().Err(err).Msgf("Failed test for UDP connection to %s:%d", host, port)
		return 1
	}

	app.logger.Debug().Msgf("Read %d bytes from UDP %s:%d - %s", bytesRead, host, port, string(dataOut[:bytesRead]))

	return 0
}

func (app *application) testSSHConnection(ctx context.Context, host string, port int, username string, password string) int {
	app.logger.Debug().Msgf("Testing SSH connection to %s:%d", host, port)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)), config)
	if err != nil {
		app.logger.Error().Err(err).Msgf("Unable to connect to SSH server %s:%d", host, port)
		return 1
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		app.logger.Error().Err(err).Msgf("Unable to create session with SSH server %s:%d", host, port)
		return 1
	}
	defer session.Close()

	app.logger.Debug().Msgf("Connected to SSH server %s:%d", host, port)

	return 0
}

func (app *application) testICMPConnection(ctx context.Context, host string) int {
	app.logger.Debug().Msgf("Testing ICMP connection to %s", host)

	// Listen for ICMP packets
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		app.logger.Error().Err(err).Msgf("Failed test for ICMP connection to %s", host)
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
		app.logger.Error().Err(err).Msgf("Failed test for ICMP connection to %s", host)
		return 1
	}

	// Send the request
	if _, err := conn.WriteTo(b, &net.IPAddr{IP: net.ParseIP(host)}); err != nil {
		app.logger.Error().Err(err).Msgf("Failed test for ICMP connection to %s", host)
		return 1
	}

	// Wait for a reply
	reply := make([]byte, 1500)
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		app.logger.Error().Err(err).Msgf("Failed test for ICMP connection to %s", host)
		return 1
	}
	n, peer, err := conn.ReadFrom(reply)
	if err != nil {
		app.logger.Error().Err(err).Msgf("Failed test for ICMP connection to %s", host)
		return 1
	}

	// Parse the reply
	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), reply[:n])
	if err != nil {
		app.logger.Error().Err(err).Msgf("Failed test for ICMP connection to %s", host)
		return 1
	}

	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		app.logger.Debug().Msgf("Got ICMP echo reply from %v", peer)
	default:
		app.logger.Debug().Msgf("Got ICMP %+v reply from %v - expected echo", rm, peer)
		return 1
	}

	return 0
}

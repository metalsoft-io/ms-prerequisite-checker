package main

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"strconv"

	"golang.org/x/crypto/ssh"
)

func (app *application) testHTTPConnection(ctx context.Context, host string, port int) int {
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

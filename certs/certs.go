package certs

import (
	"embed"
)

//go:embed "certs"
var certsFS embed.FS

func GetCert(name string) ([]byte, error) {
	cert, err := certsFS.ReadFile("certs/" + name)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

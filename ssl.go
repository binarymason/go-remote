package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// initSSL creates a SSL certificate and key using system's openssl.
// TODO: don't blow away PKI if it already exists and is valid.
func initSSL(certPath, keyPath string) ([]byte, error) {
	if out, err := initRndFile(); err != nil {
		return out, err
	}

	fqdn, err := getFQDN()

	if err != nil {
		return []byte{}, err
	}

	mkdirP(certPath)
	mkdirP(keyPath)

	command := "openssl"
	args := []string{
		"req",
		"-new",
		"-newkey",
		"rsa:4096",
		"-days",
		"3650",
		"-nodes",
		"-x509",
		"-subj",
		fmt.Sprintf("/C=US/ST=Somewhere/L=Unknown/O=Idk/CN=%s", fqdn),
		"-keyout",
		keyPath,
		"-out",
		certPath,
	}

	return runCommand(command + " " + strings.Join(args, " "))
}

func initRndFile() ([]byte, error) {
	return runCommand(`openssl rand -out "$HOME/.rnd" -hex 256`)
}

func getFQDN() (fqdn string, err error) {
	out, err := runCommand("hostname --fqdn")

	if err != nil {
		return
	}

	fqdn = string(out)
	fqdn = fqdn[:len(fqdn)-1] // removing EOL

	return
}

func mkdirP(p string) {
	absPath, _ := filepath.Abs(p)
	dir := filepath.Dir(absPath)
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		os.MkdirAll(dir, 0700)

	}

}

package meta

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

var githubPublicKeyFingerprints = []string{
	"SHA256:p2QAMXNIC1TJYWeIOttrVc98/R1BUFWu3/LiyKgUfQM", // ECDSA
	"SHA256:nThbg6kXUpJWGl7E1IGOCspRomTxdCARLviKw6E5SY8", // RSA
	"SHA256:+DiY3wvvV6TuJJhbpZisF/zLDA0zPMSvHdkr4UvCOqU", // ED25519
}

func acceptGithubHostKeys(hostname string, _ net.Addr, key ssh.PublicKey) error {
	hostname, _, err := net.SplitHostPort(hostname)
	if err != nil {
		return err
	}

	if hostname != "github.com" {
		return fmt.Errorf("only connections to github.com are allowed: %q", hostname)
	}

	fingerprint := ssh.FingerprintSHA256(key)

	for _, candidate := range githubPublicKeyFingerprints {
		if candidate == fingerprint {
			return nil
		}
	}

	return fmt.Errorf("fingerprint %q did not match any valid github.com fingerprint", fingerprint)
}

func loadDeployKey() func() ([]ssh.Signer, error) {
	rawDeployKey := os.Getenv("DOCKERIZED_DEPLOY_KEY")
	privateKey, err := ssh.ParsePrivateKey([]byte(rawDeployKey))
	signers := []ssh.Signer{privateKey}

	if err != nil {
		err = fmt.Errorf("could not load private key from environment: %w", err)
	}

	return func() ([]ssh.Signer, error) {
		return signers, err
	}
}

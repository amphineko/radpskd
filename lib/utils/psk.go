package utils

import (
	"github.com/amphineko/radpskd/lib/log"
	"github.com/tyler-smith/go-bip39"
)

const (
	DefaultEntropySize = 128
	DefaultPskLength   = 63
)

func GeneratePsk() (string, error) {
	for {
		entropy, err := bip39.NewEntropy(DefaultEntropySize)
		if err != nil {
			log.Error.Printf("[utils.GeneratePsk] failed to generate entropy: %s", err)
			return "", err
		}

		mnemonic, err := bip39.NewMnemonic(entropy)
		if err != nil {
			log.Error.Printf("[utils.GeneratePsk] failed to generate mnemonic: %s", err)
			return "", err
		}

		if len(mnemonic) < DefaultPskLength {
			return mnemonic, nil
		}
	}
}

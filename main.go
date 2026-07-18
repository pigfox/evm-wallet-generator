// Command evm-wallet-generator creates BIP-39 mnemonic backed EVM wallets.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/crypto"
	bip39 "github.com/tyler-smith/go-bip39"
)

// Wallet is a single derived account.
type Wallet struct {
	Path       string `json:"path"`
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

// Result is the full generator output.
type Result struct {
	Mnemonic string   `json:"mnemonic"`
	Wallets  []Wallet `json:"wallets"`
}

func main() {
	var (
		mnemonic = flag.String("mnemonic", "", "existing BIP-39 mnemonic to derive from (default: generate a new one)")
		passwd   = flag.String("passphrase", "", "optional BIP-39 passphrase (25th word)")
		bits     = flag.Int("bits", 128, "mnemonic entropy in bits: 128 (12 words) or 256 (24 words)")
		count    = flag.Int("count", 1, "number of accounts to derive")
		account  = flag.Int("account", 0, "BIP-44 account index")
		asJSON   = flag.Bool("json", false, "print output as JSON")
	)
	flag.Parse()

	res, err := generate(*mnemonic, *passwd, *bits, *count, *account)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if *asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(res); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		return
	}

	fmt.Printf("Mnemonic: %s\n", res.Mnemonic)
	for _, w := range res.Wallets {
		fmt.Printf("\nPath:        %s\nAddress:     %s\nPrivate key: %s\nPublic key:  %s\n",
			w.Path, w.Address, w.PrivateKey, w.PublicKey)
	}
}

func generate(mnemonic, passphrase string, bits, count, account int) (*Result, error) {
	if count < 1 {
		return nil, fmt.Errorf("count must be at least 1")
	}

	if mnemonic == "" {
		if bits != 128 && bits != 256 {
			return nil, fmt.Errorf("bits must be 128 or 256, got %d", bits)
		}
		entropy, err := bip39.NewEntropy(bits)
		if err != nil {
			return nil, fmt.Errorf("generating entropy: %w", err)
		}
		mnemonic, err = bip39.NewMnemonic(entropy)
		if err != nil {
			return nil, fmt.Errorf("generating mnemonic: %w", err)
		}
	} else if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}

	seed := bip39.NewSeed(mnemonic, passphrase)
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("deriving master key: %w", err)
	}

	// m/44'/60'/account'/0
	branch := master
	for _, i := range []uint32{
		hdkeychain.HardenedKeyStart + 44,
		hdkeychain.HardenedKeyStart + 60,
		hdkeychain.HardenedKeyStart + uint32(account),
		0,
	} {
		if branch, err = branch.Derive(i); err != nil {
			return nil, fmt.Errorf("deriving path: %w", err)
		}
	}

	res := &Result{Mnemonic: mnemonic}
	for i := 0; i < count; i++ {
		child, err := branch.Derive(uint32(i))
		if err != nil {
			return nil, fmt.Errorf("deriving account %d: %w", i, err)
		}
		btcKey, err := child.ECPrivKey()
		if err != nil {
			return nil, fmt.Errorf("extracting private key %d: %w", i, err)
		}
		key := btcKey.ToECDSA()

		res.Wallets = append(res.Wallets, Wallet{
			Path:       fmt.Sprintf("m/44'/60'/%d'/0/%d", account, i),
			Address:    crypto.PubkeyToAddress(key.PublicKey).Hex(),
			PrivateKey: fmt.Sprintf("0x%x", crypto.FromECDSA(key)),
			PublicKey:  fmt.Sprintf("0x%x", crypto.FromECDSAPub(&key.PublicKey)),
		})
	}
	return res, nil
}

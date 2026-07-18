# evm-wallet-generator

Simple EVM wallet generator in Go. Creates a BIP-39 mnemonic and derives Ethereum-compatible
accounts from it along the standard BIP-44 path `m/44'/60'/<account>'/0/<index>`.

Generated wallets are compatible with MetaMask, Ledger, Trezor, and any other wallet that
follows the same derivation path.

## Install

```bash
go install github.com/pigfox/evm-wallet-generator@latest
```

Or build from source:

```bash
git clone https://github.com/pigfox/evm-wallet-generator.git
cd evm-wallet-generator
go build -o evm-wallet-generator .
```

## Usage

Generate a new wallet with a fresh 12-word mnemonic:

```bash
./evm-wallet-generator
```

```
Mnemonic: pledge cushion tumble hover carbon useful ...

Path:        m/44'/60'/0'/0/0
Address:     0x2f0b23f3f0f2c8b0c9e0f7c1f2a3b4c5d6e7f809
Private key: 0x...
Public key:  0x...
```

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `-mnemonic` | *(generate new)* | Derive from an existing BIP-39 mnemonic instead of creating one |
| `-passphrase` | `""` | Optional BIP-39 passphrase (the "25th word") |
| `-bits` | `128` | Entropy for a new mnemonic: `128` (12 words) or `256` (24 words) |
| `-count` | `1` | Number of accounts to derive |
| `-account` | `0` | BIP-44 account index |
| `-json` | `false` | Print output as JSON |

### Examples

24-word mnemonic with five derived accounts:

```bash
./evm-wallet-generator -bits 256 -count 5
```

Recover addresses from an existing mnemonic, as JSON:

```bash
./evm-wallet-generator -mnemonic "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about" -json
```

```json
{
  "mnemonic": "abandon abandon ... about",
  "wallets": [
    {
      "path": "m/44'/60'/0'/0/0",
      "address": "0x9858EfFD232B4033E47d90003D41EC34EcaEda94",
      "private_key": "0x1ab42cc412b618bdea3a599e3c9bae199ebf030895b039e9db1e30dafb12b727",
      "public_key": "0x0437b0bb..."
    }
  ]
}
```

## Security

- Private keys and mnemonics are printed to **stdout**. Never run this on a shared or untrusted
  machine, and be aware that piping to a file or leaving it in scrollback exposes the secret.
- The mnemonic in the example above is a well-known public test vector. Never send funds to
  addresses derived from it.
- Entropy comes from `crypto/rand` via [`go-bip39`](https://github.com/tyler-smith/go-bip39).
- This tool is provided as-is; review the code before using it to secure real funds.

## Dependencies

- [`github.com/tyler-smith/go-bip39`](https://github.com/tyler-smith/go-bip39) — mnemonic generation and seed derivation
- [`github.com/btcsuite/btcd`](https://github.com/btcsuite/btcd) — BIP-32 hierarchical key derivation
- [`github.com/ethereum/go-ethereum`](https://github.com/ethereum/go-ethereum) — secp256k1 keys and EIP-55 addresses

## License

MIT — see [LICENSE](LICENSE).

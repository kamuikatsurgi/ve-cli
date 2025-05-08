# ve-cli

ve-cli is a command‑line tool for fetching and decoding Vote Extensions (VEs) from Heimdall‑v2 blocks. It supports querying a single block or a range of blocks. It prints human‑readable details for each vote, including side‑transactions and checkpoint messages.

## Usage

```bash
Use the 'block' subcommand for a single block, or the 'blocks' subcommand to process a range of blocks.

Usage:
  ve-cli [command]

Available Commands:
  block       Process a single block height
  blocks      Process a range of block heights
  help        Help about any command

Flags:
  -c, --comet-endpoint string      CometBFT Endpoint (default "http://localhost:26657")
  -e, --heimdall-endpoint string   Heimdall Endpoint (default "http://localhost:1317")
  -h, --help                       help for ve-cli

Use "ve-cli [command] --help" for more information about a command.
```

## Examples

```bash
bin/ve-cli block 10 --comet-endpoint http://localhost:26657 --heimdall-endpoint http://localhost:1317

bin/ve-cli blocks 10 20 --comet-endpoint http://localhost:26657 --heimdall-endpoint http://localhost:1317
```

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
  -e, --endpoint string   Heimdall-v2 RPC URL (default "http://localhost:26657")
  -h, --help              help for ve-cli

Use "ve-cli [command] --help" for more information about a command.
```

## Examples

```bash
ve-cli block 10 --endpoint http://localhost:26657

ve-cli blocks 10 20 --endpoint http://localhost:26657
```

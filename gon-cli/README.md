# gon-cli

The GON cli is a CLI for Game of NFTs. It can create classes/collections/denoms, mint NFTs, transfer NFTs, and query NFTs from a single CLI.

## Feature matrix

| Chain    | Create Class | Mint NFTs | Transfer NFTs (over IBC) | Query NFTs |
|----------|--------------|-----------|--------------------------|------------|
| IRIS     | ✅            | ✅         | ✅                        | ✅          |
| Stargaze | ❌            | ❌         | ✅                        | ✅          |
| Juno     | ❌            | ❌         | ✅                        | ✅          |
| Uptick   | ❌            | ❌         | ❌                        | ❌          |
| OmniFlix | ❌            | ❌         | ✅                        | ✅          |

## Usage

Simply run `gon` and follow the instructions:

![gon.gif](./gon.gif)

### Key management
To manage keys directly, you can use the familiar `gon keys [command]` commands like `gon keys add --recover`, `gon keys list`, etc.

## Installation

### Download binary

[Download](https://github.com/gjermundgaraba/gon-tools/releases/latest/download/gon) the latest release and use it directly.

### Install from source

```bash
$ go install
```
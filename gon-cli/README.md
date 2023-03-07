# gon-cli

The GON cli is a CLI for Game of NFTs. It can create classes/collections/denoms, mint NFTs, transfer NFTs, and query NFTs from a single CLI.

## Feature matrix

| Chain    | Create Class | Mint NFTs | Transfer NFTs (over IBC) | Query NFTs |
|----------|--------------|-----------|--------------------------|------------|
| IRIS     | ✅            | ✅         | ✅                        | ✅          |
| Stargaze | ❌            | ❌         | ✅                        | ✅          |
| Juno     | ❌            | ❌         | ✅                        | ✅          |
| Uptick   | ❌            | ❌         | ✅                        | ✅          |
| OmniFlix | ❌            | ❌         | ✅                        | ✅          |

## Usage

Simply run `gon` and follow the instructions:

![ezgif-5-9d3d5d8072.gif](..%2F..%2F..%2FDownloads%2Fezgif-5-9d3d5d8072.gif)
## Installation

### Download binary

Download the latest release and use it directly.

### Install from source

```bash
$ go install
```
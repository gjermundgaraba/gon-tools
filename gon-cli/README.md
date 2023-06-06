# gon-cli

This repo is getting archived:

The development of the gon-cli has been moved to https://github.com/gjermundgaraba/nft-cli in the form of a less
testnet specific version that focuses on interchain NFT use cases (still transfering and self-relaying NFT transfers).

The GON cli is a CLI for Game of NFTs. It can create classes/collections/denoms, mint NFTs, transfer NFTs, and query NFTs from a single CLI.

## Feature matrix

| Chain    | Create Class | Mint NFTs | Transfer NFTs (over IBC) | Query NFTs |
|----------|--------------|-----------|--------------------------|------------|
| IRIS     | ✅            | ✅         | ✅                        | ✅          |
| Stargaze | ❌            | ❌         | ✅                        | ✅          |
| Juno     | ❌            | ❌         | ✅                        | ✅          |
| Uptick   | ❌            | ❌         | ✅                        | ✅          |
| OmniFlix | ❌            | ❌         | ✅                        | ✅          |

Extra helper tools and features:
- ✅ Self relaying ([docs](./self-relay.md))
- ✅ IBC Transaction lookup
- ✅ Key management
- ✅ Multiline editor for NFT metadata
- ✅ List GoN IBC Connections
- ✅ Calculate Class Hash
- ✅ Generate Class IBC trace (by choosing each hop)
- ✅ Generate Class IBC trace (by writing the flow, ie. i1j2u or i --(1)--> j --(2)--> u)
- ✅ Generate relay commands so that you can relay yourself with `rly` or `hermes` on the paths you need only
- ✅ Validate GoN evidence sheet
- ✅ Run the GoN race (if you have an NFT set up with the correct data)
- ✅ Test your NFT knowledge with the GoN quiz!

## Usage

Simply run `gon` and follow the instructions:

![gon.gif](./gon.gif)

### Editor
The editor, which is used for multiline fields, might default to `vim` in many peoples environment. You can override this by setting the `EDITOR` environment variable (to for instance nano), or simply set it when calling the cli like this:
```bash
$ EDITOR=nano gon
```

### Key management
To manage keys directly, you can use the familiar `gon keys [command]` commands like `gon keys add --recover`, `gon keys list`, etc.

## Installation

### Download binary

[Download](https://github.com/gjermundgaraba/gon-tools/releases/latest/download/gon) the latest release and use it directly.

### Install from source

```bash
$ go install
```

# CLI env

The CLI env tool is simply a bash script that sets a bunch of convenient environment variables needed for Game of NFTs cli usage.

## Setup

### Change the wallet name
The tx environment variables assumes you have the same wallet (key) name on every chain.
The first thing you need to do is set the wallet name in the `env-variables.sh` file:

```bash
$WALLET_NAME="the-name-of-your-wallet"
```

### Set the environment variables
It needs to be sourced with the following command:

```bash
$ source env-variables.sh
```

## Use

Read through the list of environment variables, but some good ones include:
- $`<CHAIN>`\_TX\_DEFAULT: All the tx flags you need, including gas, --from and --node
- $`<CHAIN>`Q\_DEFAULT: All the query flags you need, including --node
- $`<CHAIN>`\_TO\_`<ANOTHER_CHAIN>`_PORT: The nft port to use for cross-chain transfers
- $`<CHAIN>`\_TO\_`<ANOTHER_CHAIN>`_CHANNEL: The channel to use for cross-chain transfers

Where `<CHAIN>` is the name of the chain you want to use. For example, `IRIS` or `STARGAZE`.

This would let you do a cross-chain NFT transfer from IRISnet to Stargaze like this:

```bash
$ iris tx nft-transfer transfer $IRIS_TO_STARGAZE_PORT_1 $IRIS_TO_STARGAZE_CHANNEL_1 <stargaze_to_address> <nft-class-id> <nft-id> $IRIS_TX_DEFAULT
```
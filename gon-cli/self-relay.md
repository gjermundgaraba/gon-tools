# Self relay

There is a command in the GoN CLI that lets you relay your own transactions.

For this to work, however, the go-relayer must be configured correctly.

## Step 1: Install the go relayer

Install the go relayer binary from here: https://github.com/cosmos/relayer

## Step 2: Add your GoN keys to the relayer

You can use the same keys that you use for the GoN CLI, or you can create new keys for the relayer (in which case you need to also fund those accounts using the faucet).

(I am using `key-1` as the key name here because that is the same key name used in the GoN relayer config instructions)

By restoring from mnemonic:
```bash
$ rly keys restore elgafar-1 key-1 "mnemonic goes here etc"
$ rly keys restore gon-flixnet-1 key-1 "mnemonic goes here etc"
$ rly keys restore uni-6 key-1 "mnemonic goes here etc"
$ rly keys restore gon-irishub-1 key-1 "mnemonic goes here etc"
# NOTICE THE --coin-type 60 flag in the next command!! 
$ rly keys restore uptick_7000-2 key-1 "mnemonic goes here etc" --coin-type 60
```

Creating new keys:
```bash
$ rly keys add elgafar-1 key-1
$ rly keys add gon-flixnet-1 key-1
$ rly keys add uni-6 key-1
$ rly keys add gon-irishub-1 key-1
# NOTICE THE --coin-type 60 flag in the next command!! 
$ rly keys add uptick_7000-2 key-1 --coin-type 60
```

## Step 3: Create the configuration file

Set up your ~/.relayer/config/config.yaml according to this: https://github.com/game-of-nfts/gon-testnets/blob/main/doc/relayer-config.md#go-relayer

If you changed the key name in step 2, you need to change the key name in the config file as well.

## Step 4: Use the GoN CLI to relay your own transactions

You just need to know source chain, destination chain and the tx hash of the transaction you want to relay.

```bash
$ gon
? What would you like to do? Self Relay IBC message
This command requires the go relayer to have been set up according to the documentation see self-relay.md
? Source chain of transactions that needs relaying Stargaze GoN Testnet
? Destination chain of transactions that needs relaying Juno GoN Testnet
? Transaction hash to relay 690512AE975B86811C76C36F4416DD55CB5939826F79B529635153D27184302B
``` 
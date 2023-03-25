# Self relay

There is a command, and a flag, in the GoN CLI that lets you relay your own transactions.

## Option 1: Transfer with the `--self-relay` flag
Simply start the gon cli using the `--self-relay` flag, and then use the `transfer` command as usual.

```bash
$ gon --self-relay
? What would you like to do? Transfer NFT (Over IBC)
etc...
```

## Option 2: Relay an existing, pending, IBC message

You just need to know source chain, destination chain and the tx hash of the transaction you want to relay.

```bash
$ gon
? What would you like to do? Self Relay IBC message
This command requires the go relayer to have been set up according to the documentation see self-relay.md
? Source chain of transactions that needs relaying Stargaze GoN Testnet
? Destination chain of transactions that needs relaying Juno GoN Testnet
? Transaction hash to relay 690512AE975B86811C76C36F4416DD55CB5939826F79B529635153D27184302B
```
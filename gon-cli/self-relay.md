# Self relay

There is a command, and a flag, in the GoN CLI that lets you relay your own transactions.

For this to work, however, there are a few things you need to do first. In particular, you need to configure the go-relayer.
This is because the GoN CLI reuses a lot of the code from the go-relayer, and that includes configuration and key management - for now.

## Step 1: Install the go relayer

(This step is only necessary because of the key management in step 2).

Install the go relayer binary from here: https://github.com/cosmos/relayer

## Step 2: Create the configuration file

Set up your ~/.relayer/config/config.yaml according to this: https://github.com/game-of-nfts/gon-testnets/blob/main/doc/relayer-config.md#go-relayer.
Or, use the config I have created for you that you can find at the bottom of this document.

Take note of the name you use for the `key:` value. That will be required in step 3.

## Step 3: Add your GoN keys to the relayer

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

## Step 4: Use the GoN CLI to relay your own transactions

### Option 1: Transfer with the `--self-relay` flag
Simply start the gon cli using the `--self-relay` flag, and then use the `transfer` command as usual.

```bash
$ gon --self-relay
? What would you like to do? Transfer NFT (Over IBC)
etc...
```

### Option 2: Relay an existing, pending, IBC message

You just need to know source chain, destination chain and the tx hash of the transaction you want to relay.

```bash
$ gon
? What would you like to do? Self Relay IBC message
This command requires the go relayer to have been set up according to the documentation see self-relay.md
? Source chain of transactions that needs relaying Stargaze GoN Testnet
? Destination chain of transactions that needs relaying Juno GoN Testnet
? Transaction hash to relay 690512AE975B86811C76C36F4416DD55CB5939826F79B529635153D27184302B
``` 

## My go relayer config file

Provided as an example, and it works for me, so :shrug:

This config also includes the now "defunct" channels between stargaze and juno, so if you need to relay between those, this config might be helpful.

```yaml
global:
    api-listen-addr: :5183
    timeout: 10s
    memo: gon
    light-cache-size: 20
chains:
    elgafar-1:
        type: cosmos
        value:
            key: key-1
            chain-id: elgafar-1
            rpc-addr: https://rpc.elgafar-1.stargaze-apis.com:443
            account-prefix: stars
            keyring-backend: test
            gas-adjustment: 1.2
            gas-prices: 0.01ustars
            min-gas-amount: 0
            debug: true
            timeout: 10s
            output-format: json
            sign-mode: direct
            extra-codecs: []
    gon-flixnet-1:
        type: cosmos
        value:
            key: key-1
            chain-id: gon-flixnet-1
            rpc-addr: http://65.21.93.56:26657
            account-prefix: omniflix
            keyring-backend: test
            gas-adjustment: 1.2
            gas-prices: 0.01uflix
            min-gas-amount: 0
            debug: true
            timeout: 10s
            output-format: json
            sign-mode: direct
            extra-codecs: []
    gon-irishub-1:
        type: cosmos
        value:
            key: key-1
            chain-id: gon-irishub-1
            rpc-addr: http://34.80.93.133:26657
            account-prefix: iaa
            keyring-backend: test
            gas-adjustment: 1.2
            gas-prices: 0.01uiris
            min-gas-amount: 0
            debug: true
            timeout: 10s
            output-format: json
            sign-mode: direct
            extra-codecs: []
    uni-6:
        type: cosmos
        value:
            key: key-1
            chain-id: uni-6
            rpc-addr: https://rpc.uni.junonetwork.io:443
            account-prefix: juno
            keyring-backend: test
            gas-adjustment: 1.2
            gas-prices: 0.01ujunox
            min-gas-amount: 0
            debug: true
            timeout: 10s
            output-format: json
            sign-mode: direct
            extra-codecs: []
    uptick_7000-2:
        type: cosmos
        value:
            key: key-1
            chain-id: uptick_7000-2
            rpc-addr: http://52.220.252.160:26657
            account-prefix: uptick
            keyring-backend: test
            gas-adjustment: 1.2
            gas-prices: 0.01auptick
            min-gas-amount: 0
            debug: true
            timeout: 10s
            output-format: json
            sign-mode: direct
            extra-codecs:
                - ethermint
paths:
    gon-flixnet-1_channel-44-elgafar-1_channel-209:
        src:
            chain-id: gon-flixnet-1
            client-id: 07-tendermint-64
            connection-id: connection-59
        dst:
            chain-id: elgafar-1
            client-id: 07-tendermint-186
            connection-id: connection-175
        src-channel-filter:
            rule: ""
            channel-list: []
    gon-flixnet-1_channel-45-elgafar-1_channel-210:
        src:
            chain-id: gon-flixnet-1
            client-id: 07-tendermint-65
            connection-id: connection-60
        dst:
            chain-id: elgafar-1
            client-id: 07-tendermint-187
            connection-id: connection-176
        src-channel-filter:
            rule: ""
            channel-list: []
    gon-flixnet-1_channel-46-uni-6_channel-91:
        src:
            chain-id: gon-flixnet-1
            client-id: 07-tendermint-66
            connection-id: connection-61
        dst:
            chain-id: uni-6
            client-id: 07-tendermint-86
            connection-id: connection-94
        src-channel-filter:
            rule: ""
            channel-list: []
    gon-flixnet-1_channel-47-uni-6_channel-92:
        src:
            chain-id: gon-flixnet-1
            client-id: 07-tendermint-68
            connection-id: connection-62
        dst:
            chain-id: uni-6
            client-id: 07-tendermint-87
            connection-id: connection-95
        src-channel-filter:
            rule: ""
            channel-list: []
    gon-irishub-1_channel-0-gon-flixnet-1_channel-24:
        src:
            chain-id: gon-irishub-1
            client-id: 07-tendermint-0
            connection-id: connection-0
        dst:
            chain-id: gon-flixnet-1
            client-id: 07-tendermint-40
            connection-id: connection-37
        src-channel-filter:
            rule: ""
            channel-list: []
    gon-irishub-1_channel-1-gon-flixnet-1_channel-25:
        src:
            chain-id: gon-irishub-1
            client-id: 07-tendermint-1
            connection-id: connection-1
        dst:
            chain-id: gon-flixnet-1
            client-id: 07-tendermint-41
            connection-id: connection-38
        src-channel-filter:
            rule: ""
            channel-list: []
    gon-irishub-1_channel-17-uptick_7000-2_channel-3:
        src:
            chain-id: gon-irishub-1
            client-id: 07-tendermint-17
            connection-id: connection-17
        dst:
            chain-id: uptick_7000-2
            client-id: 07-tendermint-26
            connection-id: connection-25
        src-channel-filter:
            rule: ""
            channel-list: []
    gon-irishub-1_channel-19-uptick_7000-2_channel-4:
        src:
            chain-id: gon-irishub-1
            client-id: 07-tendermint-18
            connection-id: connection-19
        dst:
            chain-id: uptick_7000-2
            client-id: 07-tendermint-27
            connection-id: connection-26
        src-channel-filter:
            rule: ""
            channel-list: []
    gon-irishub-1_channel-22-elgafar-1_channel-207:
        src:
            chain-id: gon-irishub-1
            client-id: 07-tendermint-20
            connection-id: connection-21
        dst:
            chain-id: elgafar-1
            client-id: 07-tendermint-184
            connection-id: connection-173
        src-channel-filter:
            rule: ""
            channel-list: []
    gon-irishub-1_channel-23-elgafar-1_channel-208:
        src:
            chain-id: gon-irishub-1
            client-id: 07-tendermint-21
            connection-id: connection-22
        dst:
            chain-id: elgafar-1
            client-id: 07-tendermint-185
            connection-id: connection-174
        src-channel-filter:
            rule: ""
            channel-list: []
    gon-irishub-1_channel-24-uni-6_channel-89:
        src:
            chain-id: gon-irishub-1
            client-id: 07-tendermint-22
            connection-id: connection-23
        dst:
            chain-id: uni-6
            client-id: 07-tendermint-84
            connection-id: connection-92
        src-channel-filter:
            rule: ""
            channel-list: []
    gon-irishub-1_channel-25-uni-6_channel-90:
        src:
            chain-id: gon-irishub-1
            client-id: 07-tendermint-23
            connection-id: connection-24
        dst:
            chain-id: uni-6
            client-id: 07-tendermint-85
            connection-id: connection-93
        src-channel-filter:
            rule: ""
            channel-list: []
    uni-6_channel-93-elgafar-1_channel-211:
        src:
            chain-id: uni-6
            client-id: 07-tendermint-88
            connection-id: connection-96
        dst:
            chain-id: elgafar-1
            client-id: 07-tendermint-188
            connection-id: connection-177
        src-channel-filter:
            rule: ""
            channel-list: []
    uni-6_channel-94-elgafar-1_channel-213:
        src:
            chain-id: uni-6
            client-id: 07-tendermint-89
            connection-id: connection-97
        dst:
            chain-id: elgafar-1
            client-id: 07-tendermint-189
            connection-id: connection-179
        src-channel-filter:
            rule: ""
            channel-list: []
    uni-6_channel-120-elgafar-1_channel-230:
        src:
            chain-id: uni-6
            client-id: 07-tendermint-115
            connection-id: connection-124
        dst:
            chain-id: elgafar-1
            client-id: 07-tendermint-211
            connection-id: connection-200
        src-channel-filter:
            rule: ""
            channel-list: []
    uni-6_channel-122-elgafar-1_channel-234:
        src:
            chain-id: uni-6
            client-id: 07-tendermint-117
            connection-id: connection-126
        dst:
            chain-id: elgafar-1
            client-id: 07-tendermint-213
            connection-id: connection-204
        src-channel-filter:
            rule: ""
            channel-list: []
    uptick_7000-2_channel-5-gon-flixnet-1_channel-41:
        src:
            chain-id: uptick_7000-2
            client-id: 07-tendermint-28
            connection-id: connection-27
        dst:
            chain-id: gon-flixnet-1
            client-id: 07-tendermint-60
            connection-id: connection-56
        src-channel-filter:
            rule: ""
            channel-list: []
    uptick_7000-2_channel-6-elgafar-1_channel-203:
        src:
            chain-id: uptick_7000-2
            client-id: 07-tendermint-29
            connection-id: connection-28
        dst:
            chain-id: elgafar-1
            client-id: 07-tendermint-179
            connection-id: connection-169
        src-channel-filter:
            rule: ""
            channel-list: []
    uptick_7000-2_channel-7-uni-6_channel-86:
        src:
            chain-id: uptick_7000-2
            client-id: 07-tendermint-30
            connection-id: connection-29
        dst:
            chain-id: uni-6
            client-id: 07-tendermint-81
            connection-id: connection-89
        src-channel-filter:
            rule: ""
            channel-list: []
    uptick_7000-2_channel-9-gon-flixnet-1_channel-42:
        src:
            chain-id: uptick_7000-2
            client-id: 07-tendermint-31
            connection-id: connection-30
        dst:
            chain-id: gon-flixnet-1
            client-id: 07-tendermint-61
            connection-id: connection-57
        src-channel-filter:
            rule: ""
            channel-list: []
    uptick_7000-2_channel-12-elgafar-1_channel-206:
        src:
            chain-id: uptick_7000-2
            client-id: 07-tendermint-35
            connection-id: connection-34
        dst:
            chain-id: elgafar-1
            client-id: 07-tendermint-183
            connection-id: connection-172
        src-channel-filter:
            rule: ""
            channel-list: []
    uptick_7000-2_channel-13-uni-6_channel-88:
        src:
            chain-id: uptick_7000-2
            client-id: 07-tendermint-36
            connection-id: connection-35
        dst:
            chain-id: uni-6
            client-id: 07-tendermint-83
            connection-id: connection-91
        src-channel-filter:
            rule: ""
            channel-list: []

```

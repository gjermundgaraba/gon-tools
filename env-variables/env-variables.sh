#!/bin/bash

WALLET_NAME="CHANGE_ME"

IRIS_CHAIN_ID="gon-irishub-1"
IRIS_DENOM="uiris"
IRIS_RPC="http://34.80.93.133:26657"
IRIS_GRPC="http://34.80.93.133:9090"
IRIS_WS="ws://34.80.93.133/26657/websocket"
IRIS_TX_DEFAULT=(--gas auto --gas-adjustment 1.5 --gas-prices 0.025$IRIS_DENOM --from $WALLET_NAME --node $IRIS_RPC --chain-id $IRIS_CHAIN_ID)
IRIS_Q_DEFAULT=(--node $IRIS_RPC)
IRIS_PORT="nft-transfer"

STARGAZE_CHAIN_ID="elgafar-1"
STARGAZE_DENOM="ustars"
STARGAZE_RPC="https://rpc.elgafar-1.stargaze-apis.com:443"
STARGAZE_GRPC="http://grpc-1.elgafar-1.stargaze-apis.com:26660"
STARGAZE_WS="ws://rpc.elgafar-1.stargaze-apis.com:443/websocket"
STARGAZE_TX_DEFAULT=(--gas auto --gas-adjustment 1.5 --gas-prices 0.025$STARGAZE_DENOM --from $WALLET_NAME --node $STARGAZE_RPC --chain-id $STARGAZE_CHAIN_ID)
STARGAZE_Q_DEFAULT=(--node $STARGAZE_RPC)

JUNO_CHAIN_ID="uni-6"
JUNO_DENOM="ujunox"
JUNO_RPC="https://rpc.uni.juno.deuslabs.fi:443"
JUNO_GRPC="http://juno-testnet-grpc.polkachu.com:12690"
JUNO_WS="wss://rpc.uni.junonetwork.io/websocket"
JUNO_TX_DEFAULT=(--gas auto --gas-adjustment 1.5 --gas-prices 0.025$JUNO_DENOM --from $WALLET_NAME --node $JUNO_RPC --chain-id $JUNO_CHAIN_ID)
JUNO_Q_DEFAULT=(--node $JUNO_RPC)

UPTICK_CHAIN_ID="uptick_7000-2"
UPTICK_DENOM="auptick"
UPTICK_RPC="http://52.220.252.160:26657"
UPTICK_GRPC="http://52.220.252.160:9090"
UPTICK_WS="ws://52.220.252.160:26657/websocket"
UPTICK_TX_DEFAULT=(--gas auto --gas-adjustment 1.5 --gas-prices 0.025$UPTICK_DENOM --from $WALLET_NAME --node $UPTICK_RPC --chain-id $UPTICK_CHAIN_ID)
UPTICK_Q_DEFAULT=(--node $UPTICK_RPC)

OMNIFLIX_CHAIN_ID="gon-flixnet-1"
OMNIFLIX_DENOM="uflix"
OMNIFLIX_RPC="http://65.21.93.56:26657"
OMNIFLIX_GRPC="http://65.21.93.56:9090"
OMNIFLIX_WS="ws://65.21.93.56:26657/websocket"
OMNIFLIX_TX_DEFAULT=(--gas auto --gas-adjustment 1.5 --gas-prices 0.025$OMNIFLIX_DENOM --from $WALLET_NAME --node $OMNIFLIX_RPC --chain-id $OMNIFLIX_CHAIN_ID)
OMNIFLIX_Q_DEFAULT=(--node $OMNIFLIX_RPC)

IRIS_TO_STARGAZE_PORT_1="nft-transfer"
IRIS_TO_STARGAZE_CHANNEL_1="channel-22"
STARGAZE_TO_IRIS_PORT_1="wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh"
STARGAZE_TO_IRIS_CHANNEL_1="channel-207"

IRIS_TO_STARGAZE_PORT_2="nft-transfer"
IRIS_TO_STARGAZE_CHANNEL_2="channel-23"
STARGAZE_TO_IRIS_PORT_2="wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh"
STARGAZE_TO_IRIS_CHANNEL_2="channel-208"

IRIS_TO_JUNO_PORT_1="nft-transfer"
IRIS_TO_JUNO_CHANNEL_1="channel-24"
JUNO_TO_IRIS_PORT_1="wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a"
JUNO_TO_IRIS_CHANNEL_1="channel-89"

IRIS_TO_JUNO_PORT_2="nft-transfer"
IRIS_TO_JUNO_CHANNEL_2="channel-25"
JUNO_TO_IRIS_PORT_2="wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a"
JUNO_TO_IRIS_CHANNEL_2="channel-90"

IRIS_TO_UPTICK_PORT_1="nft-transfer"
IRIS_TO_UPTICK_CHANNEL_1="channel-17"
UPTICK_TO_IRIS_PORT_1="nft-transfer"
UPTICK_TO_IRIS_CHANNEL_1="channel-3"

IRIS_TO_UPTICK_PORT_2="nft-transfer"
IRIS_TO_UPTICK_CHANNEL_2="channel-19"
UPTICK_TO_IRIS_PORT_2="nft-transfer"
UPTICK_TO_IRIS_CHANNEL_2="channel-4"

IRIS_TO_OMNIFLIX_PORT_1="nft-transfer"
IRIS_TO_OMNIFLIX_CHANNEL_1="channel-0"
OMNIFLIX_TO_IRIS_PORT_1="nft-transfer"
OMNIFLIX_TO_IRIS_CHANNEL_1="channel-24"

IRIS_TO_OMNIFLIX_PORT_2="nft-transfer"
IRIS_TO_OMNIFLIX_CHANNEL_2="channel-1"
OMNIFLIX_TO_IRIS_PORT_2="nft-transfer"
OMNIFLIX_TO_IRIS_CHANNEL_2="channel-25"

STARGAZE_TO_JUNO_PORT_1="wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh"
STARGAZE_TO_JUNO_CHANNEL_1="channel-230"
JUNO_TO_STARGAZE_PORT_1="wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a"
JUNO_TO_STARGAZE_CHANNEL_1="channel-120"

STARGAZE_TO_JUNO_PORT_2="wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh"
STARGAZE_TO_JUNO_CHANNEL_2="channel-234"
JUNO_TO_STARGAZE_PORT_2="wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a"
JUNO_TO_STARGAZE_CHANNEL_2="channel-122"

STARGAZE_TO_UPTICK_PORT_1="wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh"
STARGAZE_TO_UPTICK_CHANNEL_1="channel-203"
UPTICK_TO_STARGAZE_PORT_1="nft-transfer"
UPTICK_TO_STARGAZE_CHANNEL_1="channel-6"

STARGAZE_TO_UPTICK_PORT_2="wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh"
STARGAZE_TO_UPTICK_CHANNEL_2="channel-206"
UPTICK_TO_STARGAZE_PORT_2="nft-transfer"
UPTICK_TO_STARGAZE_CHANNEL_2="channel-12"

STARGAZE_TO_OMNIFLIX_PORT_1="wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh"
STARGAZE_TO_OMNIFLIX_CHANNEL_1="channel-209"
OMNIFLIX_TO_STARGAZE_PORT_1="nft-transfer"
OMNIFLIX_TO_STARGAZE_CHANNEL_1="channel-44"

STARGAZE_TO_OMNIFLIX_PORT_2="wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh"
STARGAZE_TO_OMNIFLIX_CHANNEL_2="channel-210"
OMNIFLIX_TO_STARGAZE_PORT_2="nft-transfer"
OMNIFLIX_TO_STARGAZE_CHANNEL_2="channel-45"

JUNO_TO_UPTICK_PORT_1="wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a"
JUNO_TO_UPTICK_CHANNEL_1="channel-86"
UPTICK_TO_JUNO_PORT_1="nft-transfer"
UPTICK_TO_JUNO_CHANNEL_1="channel-7"

JUNO_TO_UPTICK_PORT_2="wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a"
JUNO_TO_UPTICK_CHANNEL_2="channel-88"
UPTICK_TO_JUNO_PORT_2="nft-transfer"
UPTICK_TO_JUNO_CHANNEL_2="channel-13"

JUNO_TO_OMNIFLIX_PORT_1="wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a"
JUNO_TO_OMNIFLIX_CHANNEL_1="channel-91"
OMNIFLIX_TO_JUNO_PORT_1="nft-transfer"
OMNIFLIX_TO_JUNO_CHANNEL_1="channel-46"

JUNO_TO_OMNIFLIX_PORT_2="wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a"
JUNO_TO_OMNIFLIX_CHANNEL_2="channel-92"
OMNIFLIX_TO_JUNO_PORT_2="nft-transfer"
OMNIFLIX_TO_JUNO_CHANNEL_2="channel-47"

UPTICK_TO_OMNIFLIX_PORT_1="nft-transfer"
UPTICK_TO_OMNIFLIX_CHANNEL_1="channel-5"
OMNIFLIX_TO_UPTICK_PORT_1="nft-transfer"
OMNIFLIX_TO_UPTICK_CHANNEL_1="channel-41"

UPTICK_TO_OMNIFLIX_PORT_2="nft-transfer"
UPTICK_TO_OMNIFLIX_CHANNEL_2="channel-9"
OMNIFLIX_TO_UPTICK_PORT_2="nft-transfer"
OMNIFLIX_TO_UPTICK_CHANNEL_2="channel-42"
package cmd

import (
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gjermundgaraba/gon/chains"
	"github.com/spf13/cobra"
	lensclient "github.com/strangelove-ventures/lens/client"
	"github.com/strangelove-ventures/lens/client/codecs/ethermint"
	"golang.org/x/term"
	"io"
	"os"
	"strings"
)

const (
	keyOptionCreate OptionString = "Create new key"
	keyOptionList   OptionString = "List keys"
	keyOptionDelete              = "Delete key"

	ethermintKeyNameSuffix        = "-ethermint"
	defaultCoinType        uint32 = sdk.CoinType
	ethermintCoinType      uint32 = 60
)

func chooseOrCreateKey(cmd *cobra.Command, chain chains.Chain) string {
	if from, _ := cmd.Flags().GetString(flags.FlagFrom); from != "" {
		return from
	}

	kr := getKeyring(cmd)
	records, err := kr.List()
	if err != nil {
		panic(err)
	}

	var keyName string
	if len(records) == 0 {
		fmt.Println("No keys found, creating new key...")
		keyName = createKeyInteractive(kr)
	} else {
		keyName = chooseKeyName(records)
	}

	if chain.KeyAlgo() == chains.KeyAlgoEthSecp256k1 {
		keyName = getEthermintKeyName(keyName)
	}

	return keyName
}

func manageKeys(cmd *cobra.Command) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		panic(err)
	}

	action := chooseOne("What would you like to do, key-wise?", []OptionString{keyOptionCreate, keyOptionList, keyOptionDelete})

	switch action {
	case keyOptionCreate:
		_ = createKeyInteractive(clientCtx.Keyring)
	case keyOptionList:
		listKeys(clientCtx.Keyring)
	case keyOptionDelete:
		deleteKey(clientCtx.Keyring)
	}
}

func createKeyInteractive(kr keyring.Keyring) string {
	keyName := askForString("Key name", survey.WithValidator(survey.Required))
	restore := askForConfirmation("Recover from mnemonic?", true)

	var mnemonic string
	var err error
	if restore {
		mnemonic, err = readMnemonic(os.Stdin, os.Stderr)
		if err != nil {
			panic(err)
		}
	} else {
		mnemonic, err = lensclient.CreateMnemonic()
		if err != nil {
			panic(err)
		}
	}

	record := createKey(kr, keyName, mnemonic)

	fmt.Printf("Key %q created successfully for default algo and ethermint (like evmos/uptick)\n", keyName)
	fmt.Println()
	printAddressesForKey(kr, record)

	return keyName
}

func createKey(kr keyring.Keyring, keyName, mnemonic string) *keyring.Record {
	defaultAlgo := keyring.SignatureAlgo(hd.Secp256k1)
	record, err := kr.NewAccount(keyName, mnemonic, "", hd.CreateHDPath(defaultCoinType, 0, 0).String(), defaultAlgo)
	if err != nil {
		panic(err)
	}

	ethermintAlgo := keyring.SignatureAlgo(ethermint.EthSecp256k1)
	keyNameEthermint := getEthermintKeyName(keyName)
	_, err = kr.NewAccount(keyNameEthermint, mnemonic, "", hd.CreateHDPath(ethermintCoinType, 0, 0).String(), ethermintAlgo)
	if err != nil {
		panic(err)
	}

	return record
}

func getDefaultKeyName(keyName string) string {
	if !strings.HasSuffix(keyName, ethermintKeyNameSuffix) {
		return keyName
	}

	return strings.TrimSuffix(keyName, ethermintKeyNameSuffix)
}

func getEthermintKeyName(keyName string) string {
	if strings.HasSuffix(keyName, ethermintKeyNameSuffix) {
		return keyName
	}

	return fmt.Sprintf("%s%s", keyName, ethermintKeyNameSuffix)
}

// readMnemonic reads a password in terminal mode if stdin is a terminal,
// otherwise it returns all of stdin with the trailing newline removed.
func readMnemonic(stdin io.Reader, stderr io.Writer) (string, error) {
	type fder interface {
		Fd() uintptr
	}

	if f, ok := stdin.(fder); ok {
		fmt.Fprint(stderr, "Enter mnemonic 🔑: ")
		mnemonic, err := term.ReadPassword(int(f.Fd()))
		fmt.Fprintln(stderr)
		return string(mnemonic), err
	}

	in, err := io.ReadAll(stdin)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSuffix(in, []byte("\n"))), nil
}

func listKeys(kr keyring.Keyring) {
	records, err := kr.List()
	if err != nil {
		panic(err)
	}

	if len(records) == 0 {
		fmt.Println("No keys found")
		return
	}

	for i, record := range records {
		if strings.HasSuffix(record.Name, ethermintKeyNameSuffix) {
			continue
		}
		fmt.Println(record.Name)
		printAddressesForKey(kr, record)

		if i < len(records)-1 {
			fmt.Println()
		}
	}
}

func printAddressesForKey(kr keyring.Keyring, record *keyring.Record) {
	if strings.HasSuffix(record.Name, ethermintKeyNameSuffix) {
		panic("ethermint key should not be printed")
	}

	for _, chain := range chains.Chains {
		accAddr, err := record.GetAddress()
		if err != nil {
			panic(err)
		}
		algoType := record.PubKey.GetTypeUrl()
		if chain.KeyAlgo() == chains.KeyAlgoEthSecp256k1 {
			// get the ethermint key
			ethermintKeyName := getEthermintKeyName(record.Name)
			ethermintRecord, err := kr.Key(ethermintKeyName)
			if err != nil {
				panic(err)
			}

			accAddr, err = ethermintRecord.GetAddress()
			if err != nil {
				panic(err)
			}

			algoType = ethermintRecord.PubKey.GetTypeUrl()
		}

		address, err := sdk.Bech32ifyAddressBytes(chain.Bech32Prefix(), accAddr.Bytes())
		if err != nil {
			panic(err)
		}
		fmt.Printf("Address for %s: %s (type: %s)\n", chain.Name(), address, algoType)
	}
}

func deleteKey(kr keyring.Keyring) {
	records, err := kr.List()
	if err != nil {
		panic(err)
	}
	if len(records) == 0 {
		fmt.Println("No keys found")
		return
	}

	keyName := chooseKeyName(records)
	ethermintKeyName := getEthermintKeyName(keyName)

	if err := kr.Delete(keyName); err != nil {
		panic(err)
	}
	if err := kr.Delete(ethermintKeyName); err != nil {
		panic(err)
	}

	fmt.Println("Key deleted successfully")
}

func chooseKeyName(records []*keyring.Record) string {
	// filter out ethermint keys
	var walletOptions []OptionString
	for _, record := range records {
		if strings.HasSuffix(record.Name, ethermintKeyNameSuffix) {
			continue
		}

		walletOptions = append(walletOptions, OptionString(record.Name))
	}

	return string(chooseOne("Choose wallet", walletOptions))
}

func getKeyring(cmd *cobra.Command) keyring.Keyring {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		panic(err)
	}

	return getKeyringFromCodec(clientCtx.Codec)
}

func getKeyringFromCodec(cdc codec.Codec) keyring.Keyring {
	kr, err := keyring.New("gon", "os", ".", os.Stdin, cdc, lensclient.LensKeyringAlgoOptions())
	if err != nil {
		panic(err)
	}

	return kr
}

func getCorrectedKeyName(originalKeyName string, chain chains.Chain) string {
	correctedKeyName := originalKeyName
	switch chain.KeyAlgo() {
	case chains.KeyAlgoEthSecp256k1:
		correctedKeyName = getEthermintKeyName(originalKeyName)
	case chains.KeyAlgoSecp256k1:
		correctedKeyName = getDefaultKeyName(originalKeyName)
	}

	return correctedKeyName
}

func getAddressForChain(cmd *cobra.Command, chain chains.Chain, keyName string) string {
	kr := getKeyring(cmd)
	keyName = getCorrectedKeyName(keyName, chain)

	record, err := kr.Key(keyName)
	if err != nil {
		panic(err)
	}

	accAddr, err := record.GetAddress()
	if err != nil {
		panic(err)
	}

	address, err := sdk.Bech32ifyAddressBytes(chain.Bech32Prefix(), accAddr.Bytes())
	if err != nil {
		panic(err)
	}

	return address
}

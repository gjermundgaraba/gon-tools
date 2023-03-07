package main

import (
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/gjermundgaraba/gon/cmd"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// TODO: Is this really needed?
	// check if .env file exists, if so, load it
	if _, err := os.Stat(".env"); err == nil || !os.IsNotExist(err) {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file, %v", err)
		}
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	appHomeDir := filepath.Join(userHomeDir, ".gon-cli")

	rootCmd := cmd.NewRootCmd(appHomeDir)
	if err := svrcmd.Execute(rootCmd, "", appHomeDir); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}

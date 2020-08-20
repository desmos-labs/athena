package main

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/djuno/cmd/djuno/flags"
	"github.com/desmos-labs/djuno/notifications"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/parse"
	"github.com/spf13/cobra"
)

// GetDesmosParseCmd returns the command that should be run when we want to start parsing a chain state
func GetDesmosParseCmd(cdc *codec.Codec, builder db.Builder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse [config-file] [firebase-private-key]",
		Short: "Start parsing the Desmos blockchain using the provided config file and Firebase private key",
		Long: `
Starts a series of worker that read the blockchain state, parse it and store the data on the database provided
by the configuration file located inside the given [config-file] path.

The second argument is used to tell where the file containing the Firebase private key is located. This file 
will be used to send push notifications when parsing the messages that might require them (e.g. new post, comment, etc).
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return parseCmdHandler(cdc, builder, args)
		},
	}

	cmd.Flags().Bool(flags.FlagEnableNotifications, true, "Enabled the sending of push notifications")

	cmd.Flags().String(flags.FlagDBHost, "", "Overwrite the database host written inside the configuration")
	cmd.Flags().Uint64(flags.FlagDBPort, 0, "Overwrite the database port written inside the configuration")
	cmd.Flags().String(flags.FlagDBUser, "", "Overwrite the database user written inside the configuration")
	cmd.Flags().String(flags.FlagDBPassword, "", "Overwrite the database password written inside the configuration")
	cmd.Flags().String(flags.FlagDBName, "", "Overwrite the database name written inside the configuration")

	return parse.SetupFlags(cmd)
}

// parseCmdHandler represents the function that should be called when the parse command is executed
func parseCmdHandler(codec *codec.Codec, dbBuilder db.Builder, args []string) error {
	// Setup Firebase
	err := notifications.SetupFirebase(args[1])
	if err != nil {
		return err
	}

	return parse.ParseCmdHandler(codec, dbBuilder, args[0], nil)
}

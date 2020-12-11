package main

import (
	"github.com/desmos-labs/djuno/cmd/djuno/flags"
	"github.com/desmos-labs/djuno/notifications"
	junocmd "github.com/desmos-labs/juno/cmd"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/spf13/cobra"
)

// ParseCmd returns the command that should be run when we want to start parsing a chain state
func ParseCmd(cdcBuilder config.CodecBuilder, setupCfg config.SdkConfigSetup, buildDb db.Builder) *cobra.Command {
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
			cdc, cp, database, modules, err := junocmd.SetupParsing(args, cdcBuilder, setupCfg, buildDb)
			if err != nil {
				return err
			}

			// Setup Firebase
			err = notifications.SetupFirebase(args[1])
			if err != nil {
				return err
			}

			return junocmd.StartParsing(cdc, cp, database, modules)
		},
	}

	cmd.Flags().Bool(flags.FlagEnableNotifications, true, "Enabled the sending of push notifications")

	return junocmd.SetupFlags(cmd)
}

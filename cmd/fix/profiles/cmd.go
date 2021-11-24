package profiles

import (
	"encoding/hex"
	"fmt"

	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"
	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/forbole/juno/v2/node/remote"
	"github.com/forbole/juno/v2/types/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/desmos-labs/djuno/v2/database"
	"github.com/desmos-labs/djuno/v2/utils"
	"github.com/desmos-labs/djuno/v2/x/profiles"
)

// NewProfilesCmd returns the Cobra command that allows to fix all the things related to the x/profiles module
func NewProfilesCmd(parseCfg *parse.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profiles",
		Short: "Fix things related to the x/profiles module",
	}

	cmd.AddCommand(
		chainLinksCmd(parseCfg),
	)

	return cmd
}

// chainLinksCmd returns a Cobra command that allows to fix the chain links for all the profiles
func chainLinksCmd(parseConfig *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "chain-links",
		Short: "Fix the chain links stored by re-parsing them",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug().Msg("fixing chain links")

			parseCtx, err := parse.GetParsingContext(parseConfig)
			if err != nil {
				return err
			}

			remoteCfg, ok := config.Cfg.Node.Details.(*remote.Details)
			if !ok {
				panic(fmt.Errorf("cannot run DJuno on local node"))
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Get the profiles
			addresses, err := db.GetProfilesAddresses()
			if err != nil {
				return err
			}

			grpcConnection := remote.MustCreateGrpcConnection(remoteCfg.GRPC)
			profilesClient := profilestypes.NewQueryClient(grpcConnection)

			profilesModule := profiles.NewModule(profilesClient, parseCtx.EncodingConfig.Marshaler, db)

			for _, address := range addresses {
				log.Debug().Str("address", address).Msg("deleting chain links")

				err = db.DeleteProfileChainLinks(address)
				if err != nil {
					return err
				}

				log.Debug().Str("address", address).Msg("querying transactions")

				// Chain links

				err = restoreChainLinks(address, parseCtx, profilesModule)
				if err != nil {
					return err
				}

				err = restoreChainUnlinks(address, parseCtx, profilesModule)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func restoreChainLinks(address string, parseCtx *parse.Context, profilesModule *profiles.Module) error {
	query := fmt.Sprintf("link_chain_account.chain_link_account_owner='%s'", address)
	txs, err := utils.QueryTxs(parseCtx.Node, query)
	if err != nil {
		return err
	}

	for _, tx := range txs {
		transaction, err := parseCtx.Node.Tx(hex.EncodeToString(tx.Tx.Hash()))
		if err != nil {
			return err
		}

		// Handle the MsgChainLink messages
		for index, msg := range transaction.GetMsgs() {
			if _, ok := msg.(*profilestypes.MsgLinkChainAccount); !ok {
				continue
			}

			log.Debug().Str("address", address).Msg("handling MsgLinkChainAccount message")

			err = profilesModule.HandleMsg(index, msg, transaction)
			if err != nil {
				return fmt.Errorf("error while handling MsgLinkChainAccount: %s", err)
			}
		}
	}

	return nil
}

func restoreChainUnlinks(address string, parseCtx *parse.Context, profilesModule *profiles.Module) error {
	query := fmt.Sprintf("unlink_chain_account.chain_link_account_owner='%s'", address)
	txs, err := utils.QueryTxs(parseCtx.Node, query)
	if err != nil {
		return err
	}

	for _, tx := range txs {
		transaction, err := parseCtx.Node.Tx(hex.EncodeToString(tx.Tx.Hash()))
		if err != nil {
			return err
		}

		// Handle the MsgChainLink messages
		for index, msg := range transaction.GetMsgs() {
			if _, ok := msg.(*profilestypes.MsgUnlinkChainAccount); !ok {
				continue
			}

			log.Debug().Str("address", address).Msg("handling MsgLinkChainAccount message")

			err = profilesModule.HandleMsg(index, msg, transaction)
			if err != nil {
				return fmt.Errorf("error while handling MsgLinkChainAccount: %s", err)
			}
		}
	}

	return nil
}

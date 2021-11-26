package profiles

import (
	"encoding/hex"
	"fmt"
	"sort"

	coretypes "github.com/tendermint/tendermint/rpc/core/types"

	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"
	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/forbole/juno/v2/node/remote"
	"github.com/forbole/juno/v2/types/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/utils"
	"github.com/desmos-labs/djuno/x/profiles"
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

				// Collect all the transactions
				var txs []*coretypes.ResultTx

				// Get all the MsgLinkChain txs
				query := fmt.Sprintf("link_chain_account.chain_link_account_owner='%s'", address)
				linkChainTxs, err := utils.QueryTxs(parseCtx.Node, query)
				if err != nil {
					return err
				}
				txs = append(txs, linkChainTxs...)

				// Get all the MsgUnlinkChain txs
				query = fmt.Sprintf("unlink_chain_account.chain_link_account_owner='%s'", address)
				unlinkChainTxs, err := utils.QueryTxs(parseCtx.Node, query)
				if err != nil {
					return err
				}
				txs = append(txs, unlinkChainTxs...)

				// Sort the txs based on their ascending height
				sort.Slice(txs, func(i, j int) bool {
					return txs[i].Height < txs[j].Height
				})

				// Parse all the transactions' messages
				for _, tx := range txs {
					log.Debug().Int64("height", tx.Height).Msg("parsing transaction")

					transaction, err := parseCtx.Node.Tx(hex.EncodeToString(tx.Tx.Hash()))
					if err != nil {
						return err
					}

					// Handle only the MsgLinkChain and MsgUnlinkChain instances
					for index, msg := range transaction.GetMsgs() {
						_, isMsgLinkChain := msg.(*profilestypes.MsgLinkChainAccount)
						_, isMsgUnlinkChain := msg.(*profilestypes.MsgUnlinkChainAccount)

						if !isMsgLinkChain && !isMsgUnlinkChain {
							continue
						}

						err = profilesModule.HandleMsg(index, msg, transaction)
						if err != nil {
							return fmt.Errorf("error while handling MsgLinkChainAccount: %s", err)
						}
					}
				}
			}

			return nil
		},
	}
}

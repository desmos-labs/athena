package profiles

import (
	"encoding/hex"
	"fmt"
	"sort"

	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/node/remote"
	"github.com/forbole/juno/v3/types/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"

	profilestypes "github.com/desmos-labs/desmos/v4/x/profiles/types"

	"github.com/desmos-labs/djuno/v2/database"
	"github.com/desmos-labs/djuno/v2/utils"
	"github.com/desmos-labs/djuno/v2/x/profiles"
)

// chainLinksCmd returns a Cobra command that allows to fix the chain links for all the profiles
func chainLinksCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "chain-links",
		Short: "Fetch the chain links stored on chain and save them",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug().Msg("parsing chain links")

			parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parseConfig)
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
			profilesModule := profiles.NewModule(parseCtx.Node, grpcConnection, parseCtx.EncodingConfig.Marshaler, db)

			for _, address := range addresses {
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

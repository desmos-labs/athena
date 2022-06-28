package profiles

import (
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewPostsCmd returns the Cobra command that allows to parse things related to the x/posts module
func NewPostsCmd(parseCfg *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "posts",
		Short: "Parse things related to the x/posts module",
	}

	cmd.AddCommand(
		postsCmd(parseCfg),
	)

	return cmd
}

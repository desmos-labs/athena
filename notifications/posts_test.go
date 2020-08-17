package notifications_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	poststypes "github.com/desmos-labs/desmos/x/posts/types"
	"github.com/desmos-labs/djuno/notifications"
	"github.com/stretchr/testify/require"
)

func TestGetPostMentions(t *testing.T) {
	first, _ := sdk.AccAddressFromBech32("desmos1p7c8h59nrc8e5hxvgvu2g7tpp0xwn4mzevzgg7")
	second, _ := sdk.AccAddressFromBech32("desmos1p7ad878nealg249qkkdl9ldxrllst23lklngcx")

	message := `Hello @desmos1p7c8h59nrc8e5hxvgvu2g7tpp0xwn4mzevzgg7! Who is it going 
@desmos1p7ad878nealg249qkkdl9ldxrllst23lklngcx?`
	result, err := notifications.GetPostMentions(poststypes.Post{Message: message})
	require.NoError(t, err)

	expected := []sdk.AccAddress{first, second}
	require.Len(t, result, len(expected))
	for index, address := range result {
		require.True(t, address.Equals(expected[index]))
	}
}

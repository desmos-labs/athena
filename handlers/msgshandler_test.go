package handlers_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/djuno/db"
	"github.com/desmos-labs/djuno/handlers"
	junoDb "github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/db/postgresql"
	"github.com/desmos-labs/juno/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_MsgHandler(t *testing.T) {
	var testOwner, _ = sdk.AccAddressFromBech32("cosmos1cjf97gpzwmaf30pzvaargfgr884mpp5ak8f7ns")
	var timeZone, _ = time.LoadLocation("UTC")
	var date = time.Date(2020, 1, 1, 12, 0, 0, 0, timeZone)

	var msgCreatePost = posts.NewMsgCreatePost(
		"My new post",
		posts.PostID(1),
		false,
		"",
		map[string]string{},
		testOwner,
		date,
		nil,
		nil,
	)

	//var editDate = time.Date(2010, 1, 1, 15, 0, 0, 0, timeZone)
	//var msgEditPost = posts.NewMsgEditPost(posts.PostID(94), "Edited post message", testOwner, editDate)

	event := sdk.StringEvent{Type: "post_created", Attributes: []sdk.Attribute{{Key: "post_id", Value: "1"}}}

	logs := sdk.ABCIMessageLogs{
		sdk.ABCIMessageLog{MsgIndex: uint16(1), Log: "log", Events: sdk.StringEvents{event}},
	}

	var tx = types.Tx{
		TxResponse: sdk.TxResponse{Logs: logs},
		Messages:   []sdk.Msg{},
		Fee:        auth.StdFee{},
		Signatures: nil,
		Memo:       "",
	}

	var database junoDb.Database = db.DesmosDb{
		Database: &postgresql.Database{},
	}

	err := handlers.MsgHandler(tx, 0, msgCreatePost, database)

	assert.Nil(t, err)

}

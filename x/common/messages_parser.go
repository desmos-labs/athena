package common

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
	poststypes "github.com/desmos-labs/desmos/x/staging/posts/types"
	reportstypes "github.com/desmos-labs/desmos/x/staging/reports/types"
	"github.com/desmos-labs/juno/modules/messages"
)

var MessagesParser = messages.JoinMessageParsers(
	messages.CosmosMessageAddressesParser,
	desmosMessagesParser,
)

var desmosMessagesParser = messages.JoinMessageParsers(
	postsMessagesParser,
	profilesMessagesParser,
	reportsMessagesParser,
)

func postsMessagesParser(_ codec.Marshaler, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {
	case *poststypes.MsgCreatePost:
		return []string{msg.Creator}, nil

	case *poststypes.MsgEditPost:
		return []string{msg.Editor}, nil

	case *poststypes.MsgAddPostReaction:
		return []string{msg.User}, nil

	case *poststypes.MsgRemovePostReaction:
		return []string{msg.User}, nil

	case *poststypes.MsgRegisterReaction:
		return []string{msg.Creator}, nil

	case *poststypes.MsgAnswerPoll:
		return []string{msg.Answerer}, nil
	}

	return nil, messages.MessageNotSupported(cosmosMsg)
}

func profilesMessagesParser(_ codec.Marshaler, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {
	case *profilestypes.MsgSaveProfile:
		return []string{msg.Creator}, nil

	case *profilestypes.MsgDeleteProfile:
		return []string{msg.Creator}, nil

	case *profilestypes.MsgRequestDTagTransfer:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgCancelDTagTransfer:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgAcceptDTagTransfer:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgRefuseDTagTransfer:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgCreateRelationship:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgDeleteRelationship:
		return []string{msg.User, msg.Counterparty}, nil

	case *profilestypes.MsgBlockUser:
		return []string{msg.Blocker, msg.Blocked}, nil

	case *profilestypes.MsgUnblockUser:
		return []string{msg.Blocker, msg.Blocked}, nil
	}

	return nil, messages.MessageNotSupported(cosmosMsg)
}

func reportsMessagesParser(_ codec.Marshaler, cosmosMsg sdk.Msg) ([]string, error) {
	// nolint:singleCaseSwitch
	switch msg := cosmosMsg.(type) {
	case *reportstypes.MsgReportPost:
		return []string{msg.User}, nil
	}

	return nil, messages.MessageNotSupported(cosmosMsg)
}

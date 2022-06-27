package types

import reactionstypes "github.com/desmos-labs/desmos/v4/x/reactions/types"

type Reaction struct {
	reactionstypes.Reaction
	Height int64
}

func NewReaction(reaction reactionstypes.Reaction, height int64) Reaction {
	return Reaction{
		Reaction: reaction,
		Height:   height,
	}
}

type RegisteredReaction struct {
	reactionstypes.RegisteredReaction
	Height int64
}

func NewRegisteredReaction(reaction reactionstypes.RegisteredReaction, height int64) RegisteredReaction {
	return RegisteredReaction{
		RegisteredReaction: reaction,
		Height:             height,
	}
}

type ReactionParams struct {
	reactionstypes.SubspaceReactionsParams
	Height int64
}

func NewReactionParams(params reactionstypes.SubspaceReactionsParams, height int64) ReactionParams {
	return ReactionParams{
		SubspaceReactionsParams: params,
		Height:                  height,
	}
}

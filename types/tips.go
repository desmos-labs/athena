package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type Tip struct {
	SubspaceID uint64
	Sender     string
	Target     Target
	Amount     sdk.Coins
	Height     int64
}

func NewTip(subspaceID uint64, sender string, target Target, amount sdk.Coins, height int64) Tip {
	return Tip{
		SubspaceID: subspaceID,
		Sender:     sender,
		Target:     target,
		Amount:     amount,
		Height:     height,
	}
}

type Target interface {
	isTarget()
}

var (
	_ Target = UserTarget{}
)

type UserTarget struct {
	Address string
}

func NewUserTarget(address string) UserTarget {
	return UserTarget{
		Address: address,
	}
}

func (u UserTarget) isTarget() {}

var (
	_ Target = PostTarget{}
)

type PostTarget struct {
	PostID uint64
}

func NewPostTarget(postID uint64) PostTarget {
	return PostTarget{
		PostID: postID,
	}
}

func (p PostTarget) isTarget() {}

package types

import profilestypes "github.com/desmos-labs/desmos/x/profiles/types"

type Profile struct {
	*profilestypes.Profile
	Height int64
}

func NewProfile(profile *profilestypes.Profile, height int64) Profile {
	return Profile{
		Profile: profile,
		Height:  height,
	}
}

// -------------------------------------------------------------------------------------------------------------------

type DTagTransferRequest struct {
	profilestypes.DTagTransferRequest
	Height int64
}

func NewDTagTransferRequest(request profilestypes.DTagTransferRequest, height int64) DTagTransferRequest {
	return DTagTransferRequest{
		DTagTransferRequest: request,
		Height:              height,
	}
}

type DTagTransferRequestAcceptance struct {
	DTagTransferRequest
	NewDTag string
}

func NewDTagTransferRequestAcceptance(request DTagTransferRequest, newDTag string) DTagTransferRequestAcceptance {
	return DTagTransferRequestAcceptance{
		DTagTransferRequest: request,
		NewDTag:             newDTag,
	}
}

// -------------------------------------------------------------------------------------------------------------------

type Relationship struct {
	profilestypes.Relationship
	Height int64
}

func NewRelationship(relationship profilestypes.Relationship, height int64) Relationship {
	return Relationship{
		Relationship: relationship,
		Height:       height,
	}
}

// -------------------------------------------------------------------------------------------------------------------

type Blockage struct {
	profilestypes.UserBlock
	Height int64
}

func NewBlockage(blockage profilestypes.UserBlock, height int64) Blockage {
	return Blockage{
		UserBlock: blockage,
		Height:    height,
	}
}

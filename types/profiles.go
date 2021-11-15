package types

import (
	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"
)

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

// -------------------------------------------------------------------------------------------------------------------

type ChainLink struct {
	profilestypes.ChainLink
	Height int64
}

func NewChainLink(link profilestypes.ChainLink, height int64) ChainLink {
	return ChainLink{
		ChainLink: link,
		Height:    height,
	}
}

// -------------------------------------------------------------------------------------------------------------------

type ApplicationLink struct {
	profilestypes.ApplicationLink
	Height int64
}

func NewApplicationLink(link profilestypes.ApplicationLink, height int64) ApplicationLink {
	return ApplicationLink{
		ApplicationLink: link,
		Height:          height,
	}
}

// -------------------------------------------------------------------------------------------------------------------

type ProfilesParams struct {
	profilestypes.Params
	Height int64
}

func NewProfilesParams(params profilestypes.Params, height int64) ProfilesParams {
	return ProfilesParams{
		Params: params,
		Height: height,
	}
}

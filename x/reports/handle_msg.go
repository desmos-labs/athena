package reports

import (
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/gogo/protobuf/proto"

	"github.com/desmos-labs/djuno/v2/x/filters"

	reportstypes "github.com/desmos-labs/desmos/v4/x/reports/types"

	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"
)

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ *authz.MsgExec, _ int, executedMsg sdk.Msg, tx *juno.Tx) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 || !filters.ShouldMsgBeParsed(msg) {
		return nil
	}

	switch desmosMsg := msg.(type) {
	case *reportstypes.MsgCreateReport:
		return m.handleMsgCreateReport(tx, index, desmosMsg)

	case *reportstypes.MsgDeleteReport:
		return m.handleMsgDeleteReport(tx, desmosMsg)

	case *reportstypes.MsgAddReason:
		return m.handleMsgAddReason(tx, index, desmosMsg)

	case *reportstypes.MsgSupportStandardReason:
		return m.handleMsgSupportStandardReason(tx, index, desmosMsg)

	case *reportstypes.MsgRemoveReason:
		return m.handleMsgRemoveReason(tx, desmosMsg)
	}

	log.Debug().Str("module", "reports").Str("message", proto.MessageName(msg)).
		Int64("height", tx.Height).Msg("handled message")

	return nil
}

// handleMsgCreateReport handles a MsgCreateReport
func (m *Module) handleMsgCreateReport(tx *juno.Tx, index int, msg *reportstypes.MsgCreateReport) error {
	event, err := tx.FindEventByType(index, reportstypes.EventTypeCreateReport)
	if err != nil {
		return err
	}
	reportIDStr, err := tx.FindAttributeByKey(event, reportstypes.AttributeKeyReportID)
	if err != nil {
		return err
	}
	reportID, err := reportstypes.ParseReportID(reportIDStr)
	if err != nil {
		return err
	}

	return m.updateReport(tx.Height, msg.SubspaceID, reportID)
}

// handleMsgDeleteReport handles a MsgDeleteReport
func (m *Module) handleMsgDeleteReport(tx *juno.Tx, msg *reportstypes.MsgDeleteReport) error {
	return m.db.DeleteReport(tx.Height, msg.SubspaceID, msg.ReportID)
}

// handleMsgAddReason handles a MsgAddReason
func (m *Module) handleMsgAddReason(tx *juno.Tx, index int, msg *reportstypes.MsgAddReason) error {
	event, err := tx.FindEventByType(index, reportstypes.EventTypeAddReason)
	if err != nil {
		return err
	}
	reasonIDStr, err := tx.FindAttributeByKey(event, reportstypes.AttributeKeyReasonID)
	if err != nil {
		return err
	}
	reasonID, err := reportstypes.ParseReasonID(reasonIDStr)
	if err != nil {
		return err
	}

	return m.updateReason(tx.Height, msg.SubspaceID, reasonID)
}

// handleMsgSupportStandardReason handles a MsgSupportStandardReason
func (m *Module) handleMsgSupportStandardReason(tx *juno.Tx, index int, msg *reportstypes.MsgSupportStandardReason) error {
	event, err := tx.FindEventByType(index, reportstypes.EventTypeSupportStandardReason)
	if err != nil {
		return err
	}
	reasonIDStr, err := tx.FindAttributeByKey(event, reportstypes.AttributeKeyReasonID)
	if err != nil {
		return err
	}
	reasonID, err := reportstypes.ParseReasonID(reasonIDStr)
	if err != nil {
		return err
	}

	return m.updateReason(tx.Height, msg.SubspaceID, reasonID)
}

// handleMsgRemoveReason handles a MsgRemoveReason
func (m *Module) handleMsgRemoveReason(tx *juno.Tx, msg *reportstypes.MsgRemoveReason) error {
	return m.db.DeleteReason(tx.Height, msg.SubspaceID, msg.ReasonID)
}

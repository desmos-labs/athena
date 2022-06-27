package reports

import (
	"time"

	reportstypes "github.com/desmos-labs/desmos/v4/x/reports/types"
	"github.com/gogo/protobuf/proto"

	"github.com/desmos-labs/djuno/v2/types"

	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
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

	creationDateStr, err := tx.FindAttributeByKey(event, reportstypes.AttributeKeyCreationTime)
	if err != nil {
		return err
	}
	creationDate, err := time.Parse(time.RFC3339, creationDateStr)
	if err != nil {
		return err
	}

	report := reportstypes.NewReport(
		msg.SubspaceID,
		reportID,
		msg.ReasonsIDs,
		msg.Message,
		msg.Target.GetCachedValue().(reportstypes.ReportTarget),
		msg.Reporter,
		creationDate,
	)
	return m.db.SaveReport(types.NewReport(report, tx.Height))
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

	reason := reportstypes.NewReason(msg.SubspaceID, reasonID, msg.Title, msg.Description)
	return m.db.SaveReason(types.NewReason(reason, tx.Height))
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

	stdReason, err := m.getStandardReason(tx.Height, msg.StandardReasonID)
	if err != nil {
		return err
	}

	reason := reportstypes.NewReason(msg.SubspaceID, reasonID, stdReason.Title, stdReason.Description)
	return m.db.SaveReason(types.NewReason(reason, tx.Height))
}

// handleMsgRemoveReason handles a MsgRemoveReason
func (m *Module) handleMsgRemoveReason(tx *juno.Tx, msg *reportstypes.MsgRemoveReason) error {
	return m.db.DeleteReason(tx.Height, msg.SubspaceID, msg.ReasonID)
}

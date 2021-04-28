package reports

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	reportstypes "github.com/desmos-labs/desmos/x/staging/reports/types"
	desmosdb "github.com/desmos-labs/djuno/database"
	juno "github.com/desmos-labs/juno/types"
)

// HandleMsg handles a message properly
func HandleMsg(tx *juno.Tx, msg sdk.Msg, db *desmosdb.DesmosDb) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	desmosMsg, ok := msg.(*reportstypes.MsgReportPost)
	if ok {
		return handleMsgReport(desmosMsg, db)
	}

	return nil
}

// handleMsgReport allows to handle a MsgReportPost properly
func handleMsgReport(msg *reportstypes.MsgReportPost, db *desmosdb.DesmosDb) error {
	return db.SaveReport(reportstypes.NewReport(
		msg.PostId,
		msg.ReportType,
		msg.Message,
		msg.User,
	))
}

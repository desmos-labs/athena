package reports

import (
	abci "github.com/cometbft/cometbft/abci/types"
	reportstypes "github.com/desmos-labs/desmos/v6/x/reports/types"
	juno "github.com/forbole/juno/v5/types"

	"github.com/desmos-labs/athena/utils"
)

// HandleTx implements modules.TransactionModule
func (m *Module) HandleTx(tx *juno.Tx) error {
	return utils.ParseTxEvents(tx, map[string]func(tx *juno.Tx, event abci.Event) error{
		reportstypes.EventTypeCreateReport:          m.parseCreateReportEvent,
		reportstypes.EventTypeDeleteReport:          m.parseDeleteReportEvent,
		reportstypes.EventTypeAddReason:             m.parseAddReasonEvent,
		reportstypes.EventTypeSupportStandardReason: m.parseSupportStandardReasonEvent,
	})
}

// -------------------------------------------------------------------------------------------------------------------

// parseCreateReportEvent handles a create report event
func (m *Module) parseCreateReportEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := utils.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	reportID, err := GetReportIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.updateReport(tx.Height, subspaceID, reportID)
}

// parseDeleteReportEvent handles a delete report event
func (m *Module) parseDeleteReportEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := utils.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	reportID, err := GetReportIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.db.DeleteReport(tx.Height, subspaceID, reportID)
}

// parseAddReasonEvent handles an add reason event
func (m *Module) parseAddReasonEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := utils.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	reasonID, err := GetReasonIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.updateReason(tx.Height, subspaceID, reasonID)
}

// parseSupportStandardReasonEvent handles a support standard reason event
func (m *Module) parseSupportStandardReasonEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := utils.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	reasonID, err := GetReasonIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.updateReason(tx.Height, subspaceID, reasonID)
}

// parseRemoveReasonEvent handles a remove reason event
func (m *Module) parseRemoveReasonEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := utils.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	reasonID, err := GetReasonIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.db.DeleteReason(tx.Height, subspaceID, reasonID)
}

package types

// RelationshipRow represents a database row containing the data of a relationship between two users
type RelationshipRow struct {
	Sender   string `db:"sender_address"`
	Receiver string `db:"receiver_address"`
	Subspace uint64 `db:"subspace"`
	Height   int64  `db:"height"`
}

func NewRelationshipRow(sender string, receiver string, subspace uint64, height int64) RelationshipRow {
	return RelationshipRow{
		Sender:   sender,
		Receiver: receiver,
		Subspace: subspace,
		Height:   height,
	}
}

func (row RelationshipRow) Equal(other RelationshipRow) bool {
	return row.Sender == other.Sender &&
		row.Receiver == other.Receiver &&
		row.Subspace == other.Subspace &&
		row.Height == other.Height
}

// ________________________________________________

// BlockageRow represents a single database row containing the data of a user blockage
type BlockageRow struct {
	Blocker  string `db:"blocker_address"`
	Blocked  string `db:"blocked_user_address"`
	Reason   string `db:"reason"`
	Subspace uint64 `db:"subspace"`
	Height   int64  `db:"height"`
}

func NewBlockageRow(blocker string, blocked string, reason string, subspace uint64, height int64) BlockageRow {
	return BlockageRow{
		Blocker:  blocker,
		Blocked:  blocked,
		Reason:   reason,
		Subspace: subspace,
		Height:   height,
	}
}

func (row BlockageRow) Equal(other BlockageRow) bool {
	return row.Blocker == other.Blocker &&
		row.Blocked == other.Blocked &&
		row.Reason == other.Reason &&
		row.Subspace == other.Subspace &&
		row.Height == other.Height
}

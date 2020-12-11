package types

// RelationshipRow represents a database row containing the data of a relationship between two users
type RelationshipRow struct {
	Sender   string `db:"sender_address"`
	Receiver string `db:"receiver_address"`
	Subspace string `db:"subspace"`
}

func (row RelationshipRow) Equal(other RelationshipRow) bool {
	return row.Sender == other.Sender &&
		row.Receiver == other.Receiver &&
		row.Subspace == other.Subspace
}

// ________________________________________________

// BlockageRow represents a single database row containing the data of a user blockage
type BlockageRow struct {
	Blocker  string `db:"blocker_address"`
	Blocked  string `db:"blocked_user_address"`
	Reason   string `db:"reason"`
	Subspace string `db:"subspace"`
}

func (row BlockageRow) Equal(other BlockageRow) bool {
	return row.Blocker == other.Blocker &&
		row.Blocked == other.Blocked &&
		row.Reason == other.Reason &&
		row.Subspace == other.Subspace
}

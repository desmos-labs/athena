package types

// RegisteredReactionRow represents a single PostgreSQL row containing the data of a post reaction
type RegisteredReactionRow struct {
	OwnerAddress string `db:"owner_address"`
	ShortCode    string `db:"short_code"`
	Value        string `db:"value"`
	Subspace     string `db:"subspace"`
}

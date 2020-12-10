package database

// SaveRelationship allows to save a relationship between the sender and receiver on the given subspace
func (db DesmosDb) SaveRelationship(sender, receiver string, subspace string) error {
	if err := db.SaveUserIfNotExisting(sender); err != nil {
		return err
	}

	if err := db.SaveUserIfNotExisting(receiver); err != nil {
		return err
	}

	stmt := `INSERT INTO relationship (sender, receiver, subspace) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(stmt, sender, receiver, subspace)
	return err
}

// DeleteRelationship allows to delete the relationship between the given sender and receiver on the specified subspace
func (db DesmosDb) DeleteRelationship(sender, counterparty string, subspace string) error {
	stmt := `DELETE FROM relationship WHERE sender = $1 AND receiver = $2 AND subspace = $3`
	_, err := db.Sql.Exec(stmt, sender, counterparty, subspace)
	return err
}

// SaveBlockage allows to save a user blockage
func (db DesmosDb) SaveBlockage(blocker, blocked string, reason, subspace string) error {
	if err := db.SaveUserIfNotExisting(blocker); err != nil {
		return err
	}

	if err := db.SaveUserIfNotExisting(blocked); err != nil {
		return err
	}

	stmt := `INSERT INTO user_block(blocker, blocked_user, reason, subspace) 
			 VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(stmt, blocker, blocked, reason, subspace)
	return err
}

// RemoveBlockage allow to remove a previously saved user blockage
func (db DesmosDb) RemoveBlockage(blocker, blocked string, subspace string) error {
	stmt := `DELETE FROM user_block WHERE blocker = $1 AND blocked_user = $2 AND subspace = $3`
	_, err := db.Sql.Exec(stmt, blocker, blocked, subspace)
	return err
}

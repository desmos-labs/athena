package database

// SaveRelationship allows to save a relationship between the sender and receiver on the given subspace
func (db DesmosDb) SaveRelationship(sender, receiver string, subspace string) error {
	err := db.SaveUserIfNotExisting(sender)
	if err != nil {
		return err
	}

	err = db.SaveUserIfNotExisting(receiver)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO relationship (sender_address, receiver_address, subspace) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err = db.Sql.Exec(stmt, sender, receiver, subspace)
	return err
}

// DeleteRelationship allows to delete the relationship between the given sender and receiver on the specified subspace
func (db DesmosDb) DeleteRelationship(sender, counterparty string, subspace string) error {
	stmt := `DELETE FROM relationship WHERE sender_address = $1 AND receiver_address = $2 AND subspace = $3`
	_, err := db.Sql.Exec(stmt, sender, counterparty, subspace)
	return err
}

// SaveBlockage allows to save a user blockage
func (db DesmosDb) SaveBlockage(blocker, blocked string, reason, subspace string) error {
	err := db.SaveUserIfNotExisting(blocker)
	if err != nil {
		return err
	}

	err = db.SaveUserIfNotExisting(blocked)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO user_block(blocker_address, blocked_user_address, reason, subspace) 
VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err = db.Sql.Exec(stmt, blocker, blocked, reason, subspace)
	return err
}

// RemoveBlockage allow to remove a previously saved user blockage
func (db DesmosDb) RemoveBlockage(blocker, blocked string, subspace string) error {
	stmt := `DELETE FROM user_block WHERE blocker_address = $1 AND blocked_user_address = $2 AND subspace = $3`
	_, err := db.Sql.Exec(stmt, blocker, blocked, subspace)
	return err
}

package queries

const (
	InsertOnboarding GoQuery = `
		INSERT INTO onboardings
		(name, email, password, phone_number, address, register_as, status, additional_info)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?)
	`

	UpdateOnboarding GoQuery = `
		UPDATE onboardings
		SET
			email = ?,
			password = ?,
			phone_number = ?,
			name = ?,
			address = ?,
			register_as = ?,
			status = ?
		WHERE
			id = ?
	`
)

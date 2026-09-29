package queries

const (
	InsertUser GoQuery = `
		INSERT INTO users (onboarding_id, user_name, status, additional_info)
		VALUES
			(?, ?, ?, ?)
	`

	UpdateUser GoQuery = `
		UPDATE users
		SET
			onboarding_id = ?,
			user_name = ?,
			status = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			id = ?
	`

	SelectUserByID GoQuery = `
		SELECT id, onboarding_id, user_name, status, additional_info, created_at, updated_at, deleted_at
		FROM users
		WHERE id = ?
		AND deleted_at IS NULL
	`
)

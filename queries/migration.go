package queries

const (
	MigrateOnboarding TableMigration = `
		CREATE TABLE IF NOT EXISTS onboardings (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			phone_number VARCHAR(64) NOT NULL UNIQUE,
			address TEXT NOT NULL,
			register_as VARCHAR(32) NOT NULL,
			status VARCHAR(32) NOT NULL,
			additional_info JSON NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL,
			INDEX idx_onboardings_register_as (register_as),
			INDEX idx_onboardings_status (status)
		)
	`

	MigrateUser TableMigration = `
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			onboarding_id BIGINT NOT NULL UNIQUE,
			user_name VARCHAR(255) NOT NULL UNIQUE,
			status VARCHAR(32) NOT NULL,
			additional_info JSON NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL,
			CONSTRAINT fk_users_onboarding
				FOREIGN KEY (onboarding_id) REFERENCES onboardings(id) ON UPDATE CASCADE,
			INDEX idx_users_status (status)
		)
	`
)

var MigrateTables = []string{
	MigrateOnboarding.ToString(),
	MigrateUser.ToString(),
}

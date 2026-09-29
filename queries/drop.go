package queries

const (
	DropUserTable       DropTable = `DROP TABLE IF EXISTS users`
	DropOnboardingTable DropTable = `DROP TABLE IF EXISTS onboardings`
)

var DropTables = []string{
	DropUserTable.ToString(),
	DropOnboardingTable.ToString(),
}

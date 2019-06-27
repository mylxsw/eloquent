package event

type MigrationsStartedEvent struct{}
type MigrationsEndedEvent struct{}

type MigrationStartedEvent struct {
	SQL string
}
type MigrationEndedEvent struct {
	SQL string
}

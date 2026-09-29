package apps

import (
	"context"
	"log"

	db_client "github.com/RandySteven/go-cook/db"
	temporal_client "github.com/RandySteven/go-cook/temporal"
	"github.com/RandySteven/onboard-be/configs"
	rest_handler "github.com/RandySteven/onboard-be/handlers"
	"github.com/RandySteven/onboard-be/queries"
	"github.com/RandySteven/onboard-be/repositories"
	"github.com/RandySteven/onboard-be/usecases"
)

type App struct {
	MySQL    db_client.DBClient
	Temporal temporal_client.Temporal
}

func NewDBApp(config *configs.Config) (*App, error) {
	mysqlClient, err := NewMySQLClient(prepareDBConfig(config))
	if err != nil {
		return nil, err
	}
	return &App{MySQL: mysqlClient}, nil
}

func NewApp(config *configs.Config) (*App, error) {
	app, err := NewDBApp(config)
	if err != nil {
		return nil, err
	}
	temporalClient, err := temporal_client.NewTemporalClient(prepareTemporalConfig(config))
	if err != nil {
		app.MySQL.Close()
		return nil, err
	}
	app.Temporal = temporalClient
	return app, nil
}

func (a *App) PrepareHttpHandler(ctx context.Context) *rest_handler.Handlers {
	repos := repositories.NewRepositories(a.MySQL.Client())
	usecaseSet := usecases.NewUsecases(repos, a.Temporal)
	return rest_handler.NewHandlers(usecaseSet)
}

func (a *App) ExecuteMigration(ctx context.Context) error {
	defer a.MySQL.Close()
	migrationWorker := db_client.MigrationWorker{
		DBClient: a.MySQL,
		Tables:   queries.MigrateTables,
	}
	return migrationWorker.Migration(ctx)
}

func (a *App) ExecuteDrop(ctx context.Context, dropQueries []string) error {
	defer a.MySQL.Close()
	client := a.MySQL.Client()
	for _, query := range dropQueries {
		if _, err := client.ExecContext(ctx, query); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

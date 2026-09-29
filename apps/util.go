package apps

import (
	db_client "github.com/RandySteven/go-cook/db"
	temporal_client "github.com/RandySteven/go-cook/temporal"
	"github.com/RandySteven/onboard-be/configs"
)

func prepareDBConfig(config *configs.Config) *db_client.DBConfig {
	dbConfig := config.Configs.Db
	return &db_client.DBConfig{
		Db:              dbConfig.Db,
		DbUser:          dbConfig.DbUser,
		DbPass:          dbConfig.DbPass,
		DbHost:          dbConfig.Host + ":" + dbConfig.Port,
		DbName:          dbConfig.DbName,
		SSLMode:         dbConfig.SSLMode,
		MaxIdleConns:    dbConfig.MaxIdleConns,
		MaxOpenConns:    dbConfig.MaxOpenConns,
		ConnMaxLifeTime: dbConfig.ConnMaxLifeTime,
		ConnMaxIdleTime: dbConfig.ConnMaxIdleTime,
	}
}

func prepareTemporalConfig(config *configs.Config) *temporal_client.TemporalConfig {
	temporalConfig := config.Configs.Temporal
	cfg := &temporal_client.TemporalConfig{
		Host:      temporalConfig.Host,
		Port:      temporalConfig.Port,
		Namespace: temporalConfig.Namespace,
		TaskQueue: temporalConfig.TaskQueue,
	}
	if temporalConfig.WorkerOptions != nil {
		cfg.WorkerOptions.MaxConcurrentActivityExecutionSize = temporalConfig.WorkerOptions.MaxConcurrentActivityExecutionSize
		cfg.WorkerOptions.MaxConcurrentLocalActivityExecutionSize = temporalConfig.WorkerOptions.MaxConcurrentLocalActivityExecutionSize
		cfg.WorkerOptions.WorkerActivitiesPerSecond = temporalConfig.WorkerOptions.WorkerActivitiesPerSecond
		cfg.WorkerOptions.WorkerLocalActivitiesPerSecond = temporalConfig.WorkerOptions.WorkerLocalActivitiesPerSecond
	}
	return cfg
}

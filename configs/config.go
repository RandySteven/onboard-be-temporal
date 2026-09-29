package configs

import "time"

type Config struct {
	Configs struct {
		Server struct {
			Host    string `yaml:"host"`
			Port    string `yaml:"port"`
			Timeout struct {
				Server time.Duration `yaml:"server"`
				Read   time.Duration `yaml:"read"`
				Write  time.Duration `yaml:"write"`
				Idle   time.Duration `yaml:"idle"`
			} `yaml:"timeout"`
		} `yaml:"server"`

		Db struct {
			Db              string `yaml:"db"`
			Host            string `yaml:"host"`
			Port            string `yaml:"port"`
			DbName          string `yaml:"dbname"`
			DbUser          string `yaml:"dbuser"`
			DbPass          string `yaml:"dbpass"`
			SSLMode         string `yaml:"ssl_mode"`
			MaxIdleConns    int    `yaml:"max_idle_conns"`
			MaxOpenConns    int    `yaml:"max_open_conns"`
			ConnMaxLifeTime int    `yaml:"conn_max_lifetime"`
			ConnMaxIdleTime int    `yaml:"conn_max_idletime"`
		} `yaml:"db"`

		Temporal struct {
			Host          string `yaml:"host"`
			Port          string `yaml:"port"`
			TaskQueue     string `yaml:"task_queue"`
			Namespace     string `yaml:"namespace"`
			WorkerOptions *struct {
				MaxConcurrentActivityExecutionSize      int     `yaml:"maxConcurrentActivityExecutionSize"`
				WorkerActivitiesPerSecond               float64 `yaml:"workerActivitiesPerSecond"`
				MaxConcurrentLocalActivityExecutionSize int     `yaml:"maxConcurrentLocalActivityExecutionSize"`
				WorkerLocalActivitiesPerSecond          float64 `yaml:"workerLocalActivitiesPerSecond"`
			} `yaml:"workerOptions"`
		} `yaml:"temporal"`
	} `yaml:"configs"`
}

func (c *Config) GetConfigs() any {
	return &c.Configs
}

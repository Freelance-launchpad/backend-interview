package jdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sqltrace "github.com/DataDog/dd-trace-go/contrib/database/sql/v2"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/lib/pq"
)

// Config represents a database configuration.
type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	Region   string `yaml:"region"`
	SSLMode  string `yaml:"ssl_mode"`

	MaxIdleConns    *int           `yaml:"max_idle_conns"`
	MaxOpenConns    *int           `yaml:"max_open_conns"`
	ConnMaxIdleTime *time.Duration `yaml:"conn_max_idle_time"`
	ConnMaxLifetime *time.Duration `yaml:"conn_max_lifetime"`
}

// OpenDatabase returns a database
func OpenDatabase(conf Config) (*sql.DB, error) {
	awsConf, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("config.LoadDefaultConfig error: %w", err)
	}

	if conf.Password == "" {
		conf.Password, err = auth.BuildAuthToken(
			context.Background(), fmt.Sprintf("%s:%d", conf.Host, conf.Port), conf.Region, conf.User, awsConf.Credentials)
		if err != nil {
			return nil, fmt.Errorf("auth.BuildAuthToken error: %w", err)
		}
	}

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Name)
	if conf.SSLMode != "" {
		psqlconn += fmt.Sprintf(" sslmode=%s", conf.SSLMode)
	}

	sqltrace.Register("postgres", &pq.Driver{})
	db, err := sqltrace.Open("postgres", psqlconn, sqltrace.WithDBMPropagation(tracer.DBMPropagationModeFull))
	if err != nil {
		return nil, fmt.Errorf("sqltrace.Open error: %w", err)
	}

	if conf.MaxIdleConns != nil {
		db.SetMaxIdleConns(*conf.MaxIdleConns)
	}
	if conf.ConnMaxIdleTime != nil {
		db.SetConnMaxIdleTime(*conf.ConnMaxIdleTime)
	}
	if conf.MaxOpenConns != nil {
		db.SetMaxOpenConns(*conf.MaxOpenConns)
	}
	if conf.ConnMaxLifetime != nil {
		db.SetConnMaxLifetime(*conf.ConnMaxLifetime)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

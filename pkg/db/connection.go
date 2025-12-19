package db

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	rdsutils "github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Anacardo89/lenic/config"
)

func Connect(cfg *config.Config, dsn, user string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}
	poolConfig.BeforeConnect = func(ctx context.Context, cc *pgx.ConnConfig) error {
		if cfg.AppEnv == "aws" {
			token, err := GetRDSToken(cfg, user)
			if err != nil {
				return fmt.Errorf("failed to get RDS token: %w", err)
			}
			cc.Password = token
		}
		return nil
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return pool, nil
}

func GetRDSToken(cfg *config.Config, user string) (string, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(cfg.AWS_Region))
	if err != nil {
		return "", err
	}
	return rdsutils.BuildAuthToken(context.TODO(), cfg.DB.Host, cfg.AWS_Region, user, awsCfg.Credentials)
}

func BuildDSN_URL(cfg *config.Config, user string) (string, error) {
	var (
		err           error
		host, portStr string
		port          uint16
	)
	if cfg.AppEnv == "aws" {
		host, portStr, err = net.SplitHostPort(cfg.DB.Host)
		if err != nil {
			return "", err
		}
		p, err := strconv.Atoi(portStr)
		if err != nil {
			return "", err
		}
		port = uint16(p)
	} else {
		host = cfg.DB.Host
		port = cfg.DB.Port
	}
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, cfg.DB.Pass),
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   cfg.DB.Name,
	}
	q := u.Query()
	q.Set("sslmode", cfg.DB.SSL)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

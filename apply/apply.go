package apply

import (
	"context"
	"fmt"

	"github.com/mattermost/morph"
	"github.com/mattermost/morph/drivers"
	"github.com/mattermost/morph/drivers/mysql"
	"github.com/mattermost/morph/drivers/postgres"
	"github.com/mattermost/morph/drivers/sqlite"
	"github.com/mattermost/morph/sources"
)

func Migrate(ctx context.Context, dsn, source, driverName, path string, options ...morph.EngineOption) error {
	engine, err := initializeEngine(ctx, dsn, source, driverName, path, options...)
	if err != nil {
		return err
	}
	defer engine.Close()

	return engine.ApplyAll()
}

func Up(ctx context.Context, limit int, dsn, source, driverName, path string, options ...morph.EngineOption) (int, error) {
	engine, err := initializeEngine(ctx, dsn, source, driverName, path, options...)
	if err != nil {
		return -1, err
	}
	defer engine.Close()

	return engine.Apply(limit)
}

func Down(ctx context.Context, limit int, dsn, source, driverName, path string, options ...morph.EngineOption) (int, error) {
	engine, err := initializeEngine(ctx, dsn, source, driverName, path, options...)
	if err != nil {
		return -1, err
	}
	defer engine.Close()

	return engine.ApplyDown(limit)
}

func initializeEngine(ctx context.Context, dsn, source, driverName, path string, options ...morph.EngineOption) (*morph.Morph, error) {
	src, err := sources.Open(source, path)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	var driver drivers.Driver
	switch driverName {
	case "mysql":
		driver, err = mysql.Open(dsn)
	case "postgresql", "postgres":
		driver, err = postgres.Open(dsn)
	case "sqlite":
		driver, err = sqlite.Open(dsn)
	default:
		err = fmt.Errorf("unsupported driver %s", driverName)
	}
	if err != nil {
		return nil, err
	}

	engine, err := morph.New(ctx, driver, src, options...)
	if err != nil {
		return nil, err
	}

	return engine, err
}

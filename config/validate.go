package config

import "github.com/pkg/errors"

func (c Config) Validate() error {
	switch {
	case c.Logger.Validate() != nil:
		return errors.New("Validate section 'logger' error: " + c.Logger.Validate().Error())
	case c.Database.Validate() != nil:
		return errors.New("Validate section 'database' error: " + c.Database.Validate().Error())
	default:
		return nil
	}
}

func (logger Logger) Validate() error {
	switch {
	case logger.LogLevel == "":
		return errors.New("Log level length may be greater than 0")

	default:
		return nil
	}
}

func (dbs DbSettingsStruct) Validate() error {
	switch {
	case dbs.Type != DbTypeMongo && dbs.Type != DbTypeSQLite3:
		return errors.New("Unknown db type: " + dbs.Type)
	default:
		return nil
	}
}

package config

import (
	"github.com/pkg/errors"
	"time"
)

func (c Config) Validate() error {
	switch {
	case c.Logger.Validate() != nil:
		return errors.New("Validate section 'logger' error: " + c.Logger.Validate().Error())
	case c.Notifier.Validate() != nil:
		return errors.New("Validate section 'notifier' error: " + c.Notifier.Validate().Error())
	case c.Motion.Validate() != nil:
		return errors.New("Validate section 'motion' error: " + c.Motion.Validate().Error())
	case c.Periods.Validate() != nil:
		return errors.New("Validate section 'periods' error: " + c.Periods.Validate().Error())
	case c.SentrySettings.Validate() != nil:
		return errors.New("Validate section 'sentry_settings' error: " + c.SentrySettings.Validate().Error())
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

func (n Notifier) Validate() error {
	switch {
	case !n.IsUsing:
		return nil
	case n.Type == "":
		return errors.New("Unknown notifier type")
	case n.Type != NotifierTypeTelegram:
		return errors.New("Unknown notifier type")
	case n.Options.ChatID == "":
		return errors.New("Empty chat ID")
	case n.Options.Token == "":
		return errors.New("Empty token")
	default:
		return nil
	}
}

func (ms MotionSettings) Validate() error {
	switch {
	case ms.FileExtension == "":
		return errors.New("Empty 'file_extension'")
	case ms.MoviesDirCam1 == "":
		return errors.New("Empty movies directory path")
	default:
		return nil
	}
}

func (ps Periods) Validate() error {
	switch {
	case ps.MainTickerTime == time.Duration(0):
		return errors.New("main ticker time is not define")
	default:
		return nil
	}
}

func (s SentrySettings) Validate() error {
	switch {
	case s.SentryUrl == "":
		return errors.New("Empty sentry api url")
	default:
		return nil
	}
}

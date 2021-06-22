package config

import (
	"errors"
	"fmt"
	"github.com/vasilpatelnya/rpi-home/tool/translate"
	"time"
)

func sectionError(section string, err error) error {
	msg := fmt.Sprintf("Validate section '%s' error: %s", section, err.Error())

	return errors.New(msg)
}

func (c Config) Validate() error {
	switch {
	case c.Logger.Validate() != nil:
		return sectionError("logger", c.Logger.Validate())
	case c.Notifier.Validate() != nil:
		return sectionError("notifier", c.Notifier.Validate())
	case c.Motion.Validate() != nil:
		return sectionError("motion", c.Motion.Validate())
	case c.Periods.Validate() != nil:
		return sectionError("periods", c.Periods.Validate())
	case c.SentrySettings.Validate() != nil:
		return sectionError("sentry_settings", c.SentrySettings.Validate())
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

func (as ApiSettings) Validate() error {
	switch {
	case as.Port == 0:
		return errors.New(translate.Txt(translate.ErrorValidatingAPIPort))
	case as.ApiKey == "":
		return errors.New(translate.Txt(translate.ErrorValidatingAPIKey))
	default:
		return nil
	}
}

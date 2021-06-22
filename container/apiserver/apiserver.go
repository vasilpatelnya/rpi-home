package apiserver

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/container/notification"
	"github.com/vasilpatelnya/rpi-home/dataservice"
	"github.com/vasilpatelnya/rpi-home/tool/errors"
	"github.com/vasilpatelnya/rpi-home/tool/translate"
	"net/http"
)

type ApiServer struct {
	Settings config.ApiSettings
	Repo     dataservice.EventData
	Logger   *logrus.Logger
	Notifier notification.Notifier
}

type ApiOpts struct {
	Settings config.ApiSettings
	Repo     dataservice.EventData
	Logger   *logrus.Logger
	Notifier notification.Notifier
}

func New(opts *ApiOpts) *ApiServer {
	return &ApiServer{
		Settings: opts.Settings,
		Repo:     opts.Repo,
		Logger:   opts.Logger,
		Notifier: opts.Notifier,
	}
}

func (s *ApiServer) Run() {
	http.HandleFunc("/api/v1/motioneye", MotionEyeHandler(s.Repo))

	addr := fmt.Sprintf(":%d", s.Settings.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		msg := errors.ErrorMsg(translate.ErrorApiServerCommon, err)
		s.Logger.Errorln(msg)
		sendErr := s.Notifier.SendText(msg)
		if sendErr != nil {
			s.Logger.Errorf(errors.ErrorMsg(translate.ErrorNotifierSendText, err))
		}
	}
}

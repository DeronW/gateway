package error_handle

import (
	"gateway/config"
	log "github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
	"sync"
)

var once sync.Once

func CheckError(err error) {
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Info()
	}
}

func FatalError(err error) {
	if err != nil {
		panic("ERROR: " + " " + err.Error())
	}
}

func ReportError(err error) {
	raven.CaptureError(err, nil, nil)
}

func init() {
	once.Do(func() {
		raven.SetDSN(config.GetSentryCfg().DSN)
		raven.CaptureMessage("Device Gateway server starting", nil)
	})
}

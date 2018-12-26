package password

import (
	"time"

	"github.com/labstack/echo"
	"github.com/johncoleman83/cerebrum/pkg/api/password"
	"github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// New creates new password logging service
func New(svc password.Service, logger cerebrum.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents password logging service
type LogService struct {
	password.Service
	logger cerebrum.Logger
}

const name = "password"

// Change logging
func (ls *LogService) Change(c echo.Context, id uint, oldPass, newPass string) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Change password request", err,
			map[string]interface{}{
				"req":  id,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Change(c, id, oldPass, newPass)
}

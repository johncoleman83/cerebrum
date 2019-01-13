package password

import (
	"time"

	"github.com/labstack/echo"

	"github.com/johncoleman83/cerebrum/pkg/api/password"
	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// packageName is the name of the package
const packageName = "password"

// LogService represents password logging service
type LogService struct {
	password.Service
	logger cerebrum.Logger
}

// New creates new password logging service
func New(svc password.Service, logger cerebrum.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// Change logging
func (ls *LogService) Change(c echo.Context, id uint, oldPass, newPass string) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			packageName, "Change password request", err,
			map[string]interface{}{
				"req":  id,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Change(c, id, oldPass, newPass)
}

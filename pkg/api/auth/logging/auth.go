package auth

import (
	"time"

	"github.com/labstack/echo"

	"github.com/johncoleman83/cerebrum/pkg/api/auth"
	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// packageName is the name of the package
const packageName = "auth"

// LogService represents auth logging service
type LogService struct {
	auth.Service
	logger cerebrum.Logger
}

// New creates new auth logging service
func New(svc auth.Service, logger cerebrum.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// Authenticate logging
func (ls *LogService) Authenticate(c echo.Context, user, password string) (resp *cerebrum.AuthToken, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			packageName, "Authenticate request", err,
			map[string]interface{}{
				"req":  user,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Authenticate(c, user, password)
}

// Refresh logging
func (ls *LogService) Refresh(c echo.Context, req string) (resp *cerebrum.RefreshToken, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			packageName, "Refresh request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Refresh(c, req)
}

// Me logging
func (ls *LogService) Me(c echo.Context) (resp *cerebrum.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			packageName, "Me request", err,
			map[string]interface{}{
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Me(c)
}

package user

import (
	"time"

	"github.com/labstack/echo"

	"github.com/johncoleman83/cerebrum/pkg/api/user"
	"github.com/johncoleman83/cerebrum/pkg/utl/models"
)

const packageName = "user"

// LogService represents user logging service
type LogService struct {
	user.Service
	logger models.Logger
}

// New creates new user logging service
func New(svc user.Service, logger models.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// Create logging
func (ls *LogService) Create(c echo.Context, req models.User) (resp *models.User, err error) {
	dupe := req
	dupe.Password = "xxx-redacted-xxx"
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			packageName, "Create user request", err,
			map[string]interface{}{
				"req":  dupe,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Create(c, req)
}

// List logging
func (ls *LogService) List(c echo.Context, req *models.Pagination) (resp []models.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			packageName, "List user request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.List(c, req)
}

// View logging
func (ls *LogService) View(c echo.Context, req uint) (resp *models.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			packageName, "View user request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.View(c, req)
}

// Delete logging
func (ls *LogService) Delete(c echo.Context, req uint) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			packageName, "Delete user request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Delete(c, req)
}

// Update logging
func (ls *LogService) Update(c echo.Context, req *user.Update) (resp *models.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			packageName, "Update user request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Update(c, req)
}

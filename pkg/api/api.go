package api

import (
	"crypto/sha1"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	// cerebrum/pkg/api
	"github.com/johncoleman83/cerebrum/pkg/api/auth"
	al "github.com/johncoleman83/cerebrum/pkg/api/auth/logging"
	at "github.com/johncoleman83/cerebrum/pkg/api/auth/transport"
	"github.com/johncoleman83/cerebrum/pkg/api/password"
	pl "github.com/johncoleman83/cerebrum/pkg/api/password/logging"
	pt "github.com/johncoleman83/cerebrum/pkg/api/password/transport"
	"github.com/johncoleman83/cerebrum/pkg/api/user"
	ul "github.com/johncoleman83/cerebrum/pkg/api/user/logging"
	ut "github.com/johncoleman83/cerebrum/pkg/api/user/transport"

	// cerebrum/pkg/utl
	"github.com/johncoleman83/cerebrum/pkg/utl/config"
	"github.com/johncoleman83/cerebrum/pkg/utl/datastore"
	jwtService "github.com/johncoleman83/cerebrum/pkg/utl/middleware/jsonwebtoken"
	rbacService "github.com/johncoleman83/cerebrum/pkg/utl/rbac"
	"github.com/johncoleman83/cerebrum/pkg/utl/secure"
	"github.com/johncoleman83/cerebrum/pkg/utl/server"
	"github.com/johncoleman83/cerebrum/pkg/utl/zlog"
)

// newServices initializes new services for API
func newServices(cfg *config.Configuration) (rbac *rbacService.Service, jwt *jwtService.Service, sec *secure.Service, log *zlog.Log, e *echo.Echo) {
	sec = secure.New(cfg.App.MinPasswordStr, sha1.New())
	rbac = rbacService.New()
	jwt = jwtService.New(cfg.JWT.Secret, cfg.JWT.SigningAlgorithm, cfg.JWT.Duration)
	log = zlog.New()

	e = server.New()
	e.Static("/swaggerui", cfg.App.SwaggerUIPath)
	return rbac, jwt, sec, log, e
}

// initializeControllers initializes new HTTP services for each controller
func initializeControllers(db *gorm.DB, rbac *rbacService.Service, jwt *jwtService.Service, sec *secure.Service, log *zlog.Log, e *echo.Echo) {
	at.NewHTTP(al.New(auth.Initialize(db, jwt, sec, rbac), log), e, jwt.MWFunc())

	v1 := e.Group("/v1")
	v1.Use(jwt.MWFunc())

	ut.NewHTTP(ul.New(user.Initialize(db, rbac, sec), log), v1)
	pt.NewHTTP(pl.New(password.Initialize(db, rbac, sec), log), v1)
}

// startServer starts HTTP server with correct config & initialized services
func startServer(e *echo.Echo, cfg *config.Configuration) {
	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})
}

// Start starts the API service
func Start(cfg *config.Configuration) error {
	db, err := datastore.NewMySQLGormDb(cfg.DB)
	if err != nil {
		return err
	}

	rbac, jwt, sec, log, e := newServices(cfg)

	initializeControllers(db, rbac, jwt, sec, log, e)

	startServer(e, cfg)

	return nil
}

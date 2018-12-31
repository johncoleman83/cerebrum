// CEREBRUM
//
// API Docs for CEREBRUM v1
//
// 	 Terms Of Service:  N/A
//     Schemes: http
//     Version: 1.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: David John Coleman II <me@davidjohncoleman.com> https://davidjohncoleman.com
//     Host: localhost:8080
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - Bearer: []
//
//     SecurityDefinitions:
//     Bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta

package api

import (
	"crypto/sha1"

	"github.com/johncoleman83/cerebrum/pkg/utl/zlog"

	"github.com/johncoleman83/cerebrum/pkg/api/auth"
	al "github.com/johncoleman83/cerebrum/pkg/api/auth/logging"
	at "github.com/johncoleman83/cerebrum/pkg/api/auth/transport"
	"github.com/johncoleman83/cerebrum/pkg/api/password"
	pl "github.com/johncoleman83/cerebrum/pkg/api/password/logging"
	pt "github.com/johncoleman83/cerebrum/pkg/api/password/transport"
	"github.com/johncoleman83/cerebrum/pkg/api/user"
	ul "github.com/johncoleman83/cerebrum/pkg/api/user/logging"
	ut "github.com/johncoleman83/cerebrum/pkg/api/user/transport"

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
	"github.com/johncoleman83/cerebrum/pkg/utl/datastore"
	"github.com/johncoleman83/cerebrum/pkg/utl/middleware/jwt"
	"github.com/johncoleman83/cerebrum/pkg/utl/rbac"
	"github.com/johncoleman83/cerebrum/pkg/utl/secure"
	"github.com/johncoleman83/cerebrum/pkg/utl/server"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {
	db, err := datastore.NewMySQLGormDb(cfg.DB)
	if err != nil {
		return err
	}

	sec := secure.New(cfg.App.MinPasswordStr, sha1.New())
	rbac := rbac.New()
	jwt := jwt.New(cfg.JWT.Secret, cfg.JWT.SigningAlgorithm, cfg.JWT.Duration)
	log := zlog.New()

	e := server.New()
	e.Static("/swaggerui", cfg.App.SwaggerUIPath)

	at.NewHTTP(al.New(auth.Initialize(db, jwt, sec, rbac), log), e, jwt.MWFunc())

	v1 := e.Group("/v1")
	v1.Use(jwt.MWFunc())

	ut.NewHTTP(ul.New(user.Initialize(db, rbac, sec), log), v1)
	pt.NewHTTP(pl.New(password.Initialize(db, rbac, sec), log), v1)

	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}

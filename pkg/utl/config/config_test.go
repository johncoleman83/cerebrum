package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
	"github.com/johncoleman83/cerebrum/pkg/utl/support"
)

func TestLoad(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		wantData *config.Configuration
		wantErr  bool
	}{
		{
			name:    "Fail on non-existing file",
			path:    "./path/does/not/Exists",
			wantErr: true,
		},
		{
			name:    "Fail on wrong file format",
			path:    "testdata/config.invalid.yaml",
			wantErr: true,
		},
		{
			name: "Success",
			path: support.TestingConfigPath(),
			wantData: &config.Configuration{
				DB: &config.Database{
					Dialect: 	"mysql",
					User: 		"mysql_test_user",
					Password: "mysql_test_password",
					Name: 		"cerebrum_mysql_test_db",
					Protocol: "tcp",
					Host: 		"localhost",
					Port: 		"3306",
					Settings: "tls=skip-verify&charset=utf8&parseTime=True&loc=Local&autocommit=true&timeout=20s",
				},
				Server: &config.Server{
					Port:         ":8080",
					Debug:        true,
					ReadTimeout:  15,
					WriteTimeout: 20,
				},
				JWT: &config.JWT{
					Secret:           "dsflaksdhflaksdhfalksjdhflasdfh",
					Duration:         10,
					RefreshDuration:  10,
					MaxRefresh:       144,
					SigningAlgorithm: "HS384",
				},
				App: &config.Application{
					MinPasswordStr: 3,
					SwaggerUIPath:  "assets/swagger",
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.LoadConfigFrom(tt.path)
			assert.Equal(t, tt.wantData, cfg)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

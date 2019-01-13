package mysqldb_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"github.com/johncoleman83/cerebrum/pkg/api/auth/platform/mysqldb"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock/mockdb"
	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
)

var (
	someTime   = time.Now().Round(time.Second)
	superAdmin = cerebrum.Role{
		ID:          cerebrum.AccessRole(100),
		AccessLevel: cerebrum.AccessRole(100),
		Name:        "SUPER_ADMIN",
	}
)

func TestView(t *testing.T) {
	cases := []struct {
		name         string
		expectedErr  bool
		id           uint
		expectedData *cerebrum.User
	}{
		{
			name:        "User should not not exist and not return error",
			expectedErr: true,
			id:          1000,
		},
		{
			name:        "Success",
			id:          2,
			expectedErr: false,
			expectedData: &cerebrum.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     cerebrum.AccessRole(100),
				CompanyID:  1,
				LocationID: 1,
				Password:   "newPass",
				Token:      "asdf",
				Base: cerebrum.Base{
					Model: gorm.Model{ID: 2},
				},
			},
		},
	}

	container := mockdb.NewMySQLDockerTestContainer(t)
	db, pool, resource := container.DB, container.Pool, container.Resource

	if err := mockdb.InsertMultiple(db, superAdmin, cases[1].expectedData); err != nil {
		t.Error(err)
	}

	udb := mysqldb.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := udb.View(db, tt.id)
			fmt.Println(err)
			fmt.Println(tt.expectedErr)
			assert.Equal(t, tt.expectedErr, err != nil)
			if tt.expectedData != nil {
				if user == nil {
					t.Errorf("response was nil due to: %v", err)
				} else {
					tt.expectedData.CreatedAt = user.CreatedAt
					tt.expectedData.UpdatedAt = user.UpdatedAt
					tt.expectedData.LastLogin = user.LastLogin
					tt.expectedData.LastPasswordChange = user.LastPasswordChange
					tt.expectedData.Role = superAdmin
					assert.Equal(t, tt.expectedData, user)
				}
			}
		})
	}
	db.Close()
	if err := pool.Purge(resource); err != nil {
		t.Fatal(fmt.Sprintf("Could not purge resource: %v", err))
	}
}

func TestFindByUsername(t *testing.T) {
	cases := []struct {
		name         string
		expectedErr  bool
		username     string
		expectedData *cerebrum.User
	}{
		{
			name:        "User does not exist",
			expectedErr: true,
			username:    "notExists",
		},
		{
			name:     "Success",
			username: "tomjones",
			expectedData: &cerebrum.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     cerebrum.AccessRole(100),
				CompanyID:  1,
				LocationID: 1,
				Password:   "newPass",
				Base: cerebrum.Base{
					Model: gorm.Model{ID: 2},
				},
			},
		},
	}

	container := mockdb.NewMySQLDockerTestContainer(t)
	db, pool, resource := container.DB, container.Pool, container.Resource

	if err := mockdb.InsertMultiple(db, superAdmin, cases[1].expectedData); err != nil {
		t.Error(err)
	}

	udb := mysqldb.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := udb.FindByUsername(db, tt.username)
			assert.Equal(t, tt.expectedErr, err != nil)

			if tt.expectedData != nil {
				tt.expectedData.CreatedAt = user.CreatedAt
				tt.expectedData.UpdatedAt = user.UpdatedAt
				tt.expectedData.LastLogin = user.LastLogin
				tt.expectedData.LastPasswordChange = user.LastPasswordChange
				tt.expectedData.Role = superAdmin
				assert.Equal(t, tt.expectedData, user)

			}
		})
	}
	db.Close()
	if err := pool.Purge(resource); err != nil {
		t.Fatal(fmt.Sprintf("Could not purge resource: %v", err))
	}
}

func TestFindByToken(t *testing.T) {
	cases := []struct {
		name         string
		expectedErr  bool
		token        string
		expectedData *cerebrum.User
	}{
		{
			name:        "User does not exist",
			expectedErr: true,
			token:       "notExists",
		},
		{
			name:  "Success",
			token: "loginrefresh",
			expectedData: &cerebrum.User{
				Email:      "johndoe@mail.com",
				FirstName:  "John",
				LastName:   "Doe",
				Username:   "johndoe",
				RoleID:     cerebrum.AccessRole(100),
				CompanyID:  1,
				LocationID: 1,
				Password:   "hunter2",
				Base: cerebrum.Base{
					Model: gorm.Model{ID: 1},
				},
				Token: "loginrefresh",
			},
		},
	}

	container := mockdb.NewMySQLDockerTestContainer(t)
	db, pool, resource := container.DB, container.Pool, container.Resource

	if err := mockdb.InsertMultiple(db, superAdmin, cases[1].expectedData); err != nil {
		t.Error(err)
	}

	udb := mysqldb.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := udb.FindByToken(db, tt.token)
			assert.Equal(t, tt.expectedErr, err != nil)

			if tt.expectedData != nil {
				tt.expectedData.CreatedAt = user.CreatedAt
				tt.expectedData.UpdatedAt = user.UpdatedAt
				tt.expectedData.LastLogin = user.LastLogin
				tt.expectedData.LastPasswordChange = user.LastPasswordChange
				tt.expectedData.Role = superAdmin
				assert.Equal(t, tt.expectedData, user)

			}
		})
	}
	db.Close()
	if err := pool.Purge(resource); err != nil {
		t.Fatal(fmt.Sprintf("Could not purge resource: %v", err))
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name         string
		expectedErr  bool
		usr          *cerebrum.User
		expectedData *cerebrum.User
	}{
		{
			name: "Success",
			usr: &cerebrum.User{
				Base: cerebrum.Base{
					Model: gorm.Model{
						ID: 2,
					},
				},
				FirstName: "Z",
				LastName:  "Freak",
				Address:   "Address",
				Phone:     "123456",
				Mobile:    "345678",
				Username:  "newUsername",
			},
			expectedData: &cerebrum.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Z",
				LastName:   "Freak",
				Username:   "tomjones",
				RoleID:     cerebrum.AccessRole(100),
				CompanyID:  1,
				LocationID: 1,
				Password:   "newPass",
				Address:    "Address",
				Phone:      "123456",
				Mobile:     "345678",
				Base: cerebrum.Base{
					Model: gorm.Model{ID: 2},
				},
			},
		},
	}

	container := mockdb.NewMySQLDockerTestContainer(t)
	db, pool, resource := container.DB, container.Pool, container.Resource

	if err := mockdb.InsertMultiple(db, superAdmin, cases[0].usr); err != nil {
		t.Error(err)
	}

	udb := mysqldb.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user := &cerebrum.User{}
			if err := db.First(user, tt.usr.ID).Error; err != nil {
				t.Error(err)
			}
			tt.expectedData.CreatedAt = user.CreatedAt
			tt.expectedData.LastLogin = user.LastLogin
			tt.expectedData.LastPasswordChange = user.LastPasswordChange
			err := udb.Update(db, tt.expectedData)
			assert.Equal(t, tt.expectedErr, err != nil)
			if tt.expectedData != nil {
				user = &cerebrum.User{}
				if err := db.First(user, tt.usr.ID).Error; err != nil {
					t.Error(err)
				}
				tt.expectedData.UpdatedAt = user.UpdatedAt
				assert.Equal(t, tt.expectedData, user)
			}
		})
	}
	db.Close()
	if err := pool.Purge(resource); err != nil {
		t.Fatal(fmt.Sprintf("Could not purge resource: %v", err))
	}
}

package mysqldb_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/johncoleman83/cerebrum/pkg/api/user/platform/mysqldb"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock/mockdb"
	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
	"github.com/stretchr/testify/assert"
)

var (
	someTime   = time.Now().Round(time.Second)
	superAdmin = cerebrum.Role{
		ID:          cerebrum.AccessRole(100),
		AccessLevel: cerebrum.AccessRole(100),
		Name:        "SUPER_ADMIN",
	}
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name         string
		expectedErr  bool
		req          cerebrum.User
		expectedData *cerebrum.User
	}{
		{
			name:        "Fail on insert duplicate ID",
			expectedErr: true,
			req: cerebrum.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     cerebrum.AccessRole(100),
				CompanyID:  1,
				LocationID: 1,
				Password:   "pass",
				Base: cerebrum.Base{
					Model: gorm.Model{
						ID: 1,
					},
				},
			},
		},
		{
			name:        "Fail on insert duplicate email but new id",
			expectedErr: true,
			req: cerebrum.User{
				Email:      "johndoe@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "asdf",
				RoleID:     cerebrum.AccessRole(100),
				CompanyID:  1,
				LocationID: 1,
				Password:   "pass",
				Base: cerebrum.Base{
					Model: gorm.Model{
						ID: 12,
					},
				},
			},
		},
		{
			name:        "Fail on insert duplicate username but new id",
			expectedErr: true,
			req: cerebrum.User{
				Email:      "asdf@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "johndoe",
				RoleID:     cerebrum.AccessRole(100),
				CompanyID:  1,
				LocationID: 1,
				Password:   "pass",
				Base: cerebrum.Base{
					Model: gorm.Model{
						ID: 13,
					},
				},
			},
		},
		{
			name:        "Success",
			expectedErr: false,
			req: cerebrum.User{
				Email:      "newtomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "newtomjones",
				RoleID:     cerebrum.AccessRole(100),
				CompanyID:  1,
				LocationID: 1,
				Password:   "pass",
				Base: cerebrum.Base{
					Model: gorm.Model{
						ID: 42,
					},
				},
			},
			expectedData: &cerebrum.User{
				Email:      "newtomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "newtomjones",
				RoleID:     cerebrum.AccessRole(100),
				CompanyID:  1,
				LocationID: 1,
				Password:   "pass",
				Base: cerebrum.Base{
					Model: gorm.Model{
						ID: 42,
					},
				},
			},
		},
	}

	container := mockdb.NewMySQLDockerTestContainer(t)
	db, pool, resource := container.DB, container.Pool, container.Resource

	duplicateUser := &cerebrum.User{
		Email:    "johndoe@mail.com",
		Username: "johndoe",
		Base: cerebrum.Base{
			Model: gorm.Model{
				ID: 1,
			},
		},
	}
	if err := mockdb.InsertMultiple(db, superAdmin, duplicateUser); err != nil {
		t.Error(err)
	}

	udb := mysqldb.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := udb.Create(db, tt.req)
			assert.Equal(t, tt.expectedErr, err != nil)
			if tt.expectedData != nil {
				if resp == nil {
					t.Error("Expected data, but received nil.")
					return
				}
				tt.expectedData.CreatedAt = resp.CreatedAt
				tt.expectedData.UpdatedAt = resp.UpdatedAt
				tt.expectedData.LastLogin = resp.LastLogin
				tt.expectedData.LastPasswordChange = resp.LastPasswordChange
				assert.Equal(t, tt.expectedData, resp)
			}
		})
	}
	db.Close()
	if err := pool.Purge(resource); err != nil {
		t.Fatal(fmt.Sprintf("Could not purge resource: %v", err))
	}
}

func TestView(t *testing.T) {
	cases := []struct {
		name         string
		expectedErr  bool
		id           uint
		expectedData *cerebrum.User
	}{
		{
			name:        "User does not exist",
			expectedErr: true,
			id:          1000,
		},
		{
			name: "Success",
			id:   3,
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
					Model: gorm.Model{
						ID: 3,
					},
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
						ID: 3,
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
					Model: gorm.Model{
						ID: 3,
					},
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
			tt.expectedData.DeletedAt = user.DeletedAt
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

// TODO fix these List tests
func TestList(t *testing.T) {
	cases := []struct {
		name         string
		expectedErr  bool
		qp           *cerebrum.ListQuery
		pg           *cerebrum.Pagination
		expectedData []cerebrum.User
	}{
		{
			name:        "Success, should return all 2 records",
			expectedErr: false,
			pg: &cerebrum.Pagination{
				Limit:  100,
				Offset: 0,
			},
			qp: &cerebrum.ListQuery{
				ID:    1,
				Query: "company_id = ?",
			},
			expectedData: []cerebrum.User{
				{
					Email:      "tomjones@mail.com",
					FirstName:  "Tom",
					LastName:   "Jones",
					Username:   "tomjones",
					RoleID:     cerebrum.AccessRole(100),
					CompanyID:  1,
					LocationID: 1,
					Password:   "newPass",
					Base: cerebrum.Base{
						Model: gorm.Model{
							ID: 1,
						},
					},
				},
				{
					Email:      "johnzone@mail.com",
					FirstName:  "John",
					LastName:   "Zone",
					Username:   "johnzone",
					RoleID:     cerebrum.AccessRole(100),
					CompanyID:  1,
					LocationID: 1,
					Password:   "hunter2",
					Token:      "loginrefresh",
					Base: cerebrum.Base{
						Model: gorm.Model{
							ID: 2,
						},
					},
				},
			},
		},
		{
			name:        "Success, should respect the limit and offset",
			expectedErr: false,
			pg: &cerebrum.Pagination{
				Limit:  1,
				Offset: 1,
			},
			qp: &cerebrum.ListQuery{
				ID:    1,
				Query: "company_id = ?",
			},
			expectedData: []cerebrum.User{
				{
					Email:      "johnzone@mail.com",
					FirstName:  "John",
					LastName:   "Zone",
					Username:   "johnzone",
					RoleID:     cerebrum.AccessRole(100),
					CompanyID:  1,
					LocationID: 1,
					Password:   "hunter2",
					Token:      "loginrefresh",
					Base: cerebrum.Base{
						Model: gorm.Model{
							ID: 2,
						},
					},
				},
			},
		},
	}

	container := mockdb.NewMySQLDockerTestContainer(t)
	db, pool, resource := container.DB, container.Pool, container.Resource

	if err := mockdb.InsertMultiple(db, superAdmin, &cases[0].expectedData[0], &cases[0].expectedData[1]); err != nil {
		t.Error(err)
	}

	udb := mysqldb.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			users, err := udb.List(db, tt.qp, tt.pg)
			assert.Equal(t, tt.expectedErr, err != nil)
			if tt.expectedData != nil {
				for i, v := range users {
					tt.expectedData[i].CreatedAt = v.CreatedAt
					tt.expectedData[i].UpdatedAt = v.UpdatedAt
					tt.expectedData[i].LastLogin = v.LastLogin
					tt.expectedData[i].LastPasswordChange = v.LastPasswordChange
					tt.expectedData[i].Role = superAdmin
				}
				assert.Equal(t, tt.expectedData, users)
			}
		})
	}
	db.Close()
	if err := pool.Purge(resource); err != nil {
		t.Fatal(fmt.Sprintf("Could not purge resource: %v", err))
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name         string
		id           uint
		expectedData *cerebrum.User
	}{
		{
			name: "Success",
			id:   2,
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
					Model: gorm.Model{
						ID: 2,
					},
				},
			},
		},
	}

	container := mockdb.NewMySQLDockerTestContainer(t)
	db, pool, resource := container.DB, container.Pool, container.Resource

	if err := mockdb.InsertMultiple(db, superAdmin, cases[0].expectedData); err != nil {
		t.Error(err)
	}

	udb := mysqldb.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			userBefore, err := udb.View(db, tt.id)
			assert.Equal(t, nil, err, "user should exist in db store")

			err = udb.Delete(db, userBefore)
			assert.Equal(t, nil, err, "should not error on delete")

			userAfter, err := udb.View(db, tt.id)
			emptyUser := new(cerebrum.User)
			assert.Equal(t, true, err != nil, "there should be an error when accessing deleted records")
			if err != nil {
				assert.Equal(t, "record not found", err.Error(), "error should be `record not found`")
			}
			assert.Equal(t, emptyUser, userAfter, "the response to find deleted user should be empty user")

			if err := db.Unscoped().Where("id = ?", tt.id).First(&emptyUser).Error; err != nil {
				t.Error(err)
			}
			actual := emptyUser.DeletedAt.UTC()
			actualFormatted := time.Date(actual.Year(), actual.Month(), actual.Day(), 0, 0, 0, 0, actual.Location())
			assert.Equal(t, 0, strings.Index(actualFormatted.String(), "2018-"), "the user should be have a time set for deleted_at")
		})
	}
	db.Close()
	if err := pool.Purge(resource); err != nil {
		t.Fatal(fmt.Sprintf("Could not purge resource: %v", err))
	}
}

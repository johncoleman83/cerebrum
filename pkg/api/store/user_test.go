package store_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"github.com/johncoleman83/cerebrum/pkg/api/store"
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

	udb := store.NewUser()

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
			name:        "User should not not exist and return a 404 not found error",
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

	udb := store.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := udb.View(db, tt.id)
			assert.Equal(t, tt.expectedErr, err != nil)
			if tt.expectedErr == true {
				assert.Equal(t, "code=404, message=user not found", err.Error(), "error should be `code=404, message=user not found`")
			}
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

	udb := store.NewUser()

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

	udb := store.NewUser()

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
		{
			name:        "Success, should return empty list if the query searches for non-existing id",
			expectedErr: false,
			pg: &cerebrum.Pagination{
				Limit:  100,
				Offset: 0,
			},
			qp: &cerebrum.ListQuery{
				ID:    99,
				Query: "company_id = ?",
			},
			expectedData: []cerebrum.User{},
		},
		{
			name:        "Success, should return all 3 records if no query is made",
			expectedErr: false,
			pg: &cerebrum.Pagination{
				Limit:  100,
				Offset: 0,
			},
			qp: nil,
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
				{
					Email:      "sarahsmith@mail.com",
					FirstName:  "Sarah",
					LastName:   "Smith",
					Username:   "sarahsmith",
					RoleID:     cerebrum.AccessRole(100),
					CompanyID:  3,
					LocationID: 3,
					Password:   "hunter2",
					Token:      "loginrefresh",
					Base: cerebrum.Base{
						Model: gorm.Model{
							ID: 3,
						},
					},
				},
			},
		},
	}

	container := mockdb.NewMySQLDockerTestContainer(t)
	db, pool, resource := container.DB, container.Pool, container.Resource

	if err := mockdb.InsertMultiple(db, superAdmin, &cases[3].expectedData[0], &cases[3].expectedData[1], &cases[3].expectedData[2]); err != nil {
		t.Error(err)
	}

	udb := store.NewUser()

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

	udb := store.NewUser()

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

	udb := store.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			userBefore := new(cerebrum.User)
			if err := db.Unscoped().Where("id = ?", tt.id).First(&userBefore).Error; err != nil {
				assert.Equal(t, nil, err, "user should exist in db store")
			}
			assert.Nil(t, userBefore.DeletedAt, "before user is deleted their deleted_at field should be set to NULL")

			err := udb.Delete(db, userBefore)
			assert.Nil(t, err, fmt.Sprintf("should not error on delete, error: %v", err))

			userAfter, err := udb.View(db, tt.id)
			emptyUser := new(cerebrum.User)
			assert.Equal(t, true, err != nil, "there should be an error when accessing deleted records")
			if err != nil {
				assert.Equal(t, "code=404, message=user not found", err.Error(), "error should be `code=404, message=user not found`")
			}
			assert.Equal(t, emptyUser, userAfter, "the response to find deleted user should be empty user")

			if err := db.Unscoped().Where("id = ?", tt.id).First(&emptyUser).Error; err != nil {
				assert.Nil(t, err, fmt.Sprintf("user should exist in db store and should be accessible with db.Unscopped(), error: %v", err))
			}
			assert.NotNil(t, emptyUser.DeletedAt, "the user should be have a time set for deleted_at at the same time as when the user was deleted")
		})
	}
	db.Close()
	if err := pool.Purge(resource); err != nil {
		t.Fatal(fmt.Sprintf("Could not purge resource: %v", err))
	}
}

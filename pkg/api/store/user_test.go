package store_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"github.com/johncoleman83/cerebrum/pkg/api/store"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock/mockstore"
	"github.com/johncoleman83/cerebrum/pkg/utl/models"
)

var (
	someTime   = time.Now().Round(time.Second)
	superAdmin = models.Role{
		ID:          1,
		AccessLevel: models.AccessRole(100),
		Name:        "SUPER_ADMIN",
	}
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name         string
		expectedErr  bool
		req          models.User
		expectedData *models.User
	}{
		{
			name:        "Fail on insert duplicate ID",
			expectedErr: true,
			req: models.User{
				Email:         "newname@mail.com",
				FirstName:     "New",
				LastName:      "Name",
				Username:      "newname",
				RoleID:        1,
				AccountID:     1,
				PrimaryTeamID: 1,
				Password:      "pass",
				Base: models.Base{
					Model: gorm.Model{
						ID: 1,
					},
				},
			},
		},
		{
			name:        "Fail on insert duplicate email but new id",
			expectedErr: true,
			req: models.User{
				Email:         "alreadyused@mail.com",
				FirstName:     "Never",
				LastName:      "Used",
				Username:      "neverused",
				RoleID:        1,
				AccountID:     1,
				PrimaryTeamID: 1,
				Password:      "pass",
				Base: models.Base{
					Model: gorm.Model{
						ID: 12,
					},
				},
			},
		},
		{
			name:        "Fail on insert duplicate username but new id",
			expectedErr: true,
			req: models.User{
				Email:         "brandnewmail@mail.com",
				FirstName:     "BrandNew",
				LastName:      "Name",
				Username:      "alreadyused",
				RoleID:        1,
				AccountID:     1,
				PrimaryTeamID: 1,
				Password:      "pass",
				Base: models.Base{
					Model: gorm.Model{
						ID: 13,
					},
				},
			},
		},
		{
			name:        "Success",
			expectedErr: false,
			req: models.User{
				Email:         "successfullyNew@mail.com",
				FirstName:     "Succeeding",
				LastName:      "Always",
				Username:      "successfullyNew",
				RoleID:        1,
				AccountID:     1,
				PrimaryTeamID: 1,
				Password:      "pass",
				Base: models.Base{
					Model: gorm.Model{
						ID: 42,
					},
				},
			},
			expectedData: &models.User{
				Email:         "successfullyNew@mail.com",
				FirstName:     "Succeeding",
				LastName:      "Always",
				Username:      "successfullyNew",
				RoleID:        1,
				AccountID:     1,
				PrimaryTeamID: 1,
				Password:      "pass",
				Base: models.Base{
					Model: gorm.Model{
						ID: 42,
					},
				},
			},
		},
	}

	db, err := mockstore.NewDataBaseConnection()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	duplicateUser := &models.User{
		Email:    "alreadyused@mail.com",
		Username: "alreadyused",
		Base: models.Base{
			Model: gorm.Model{
				ID: 1,
			},
		},
	}
	if err := mockstore.InsertRowsFor(db, superAdmin, duplicateUser); err != nil {
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
}

func TestView(t *testing.T) {
	cases := []struct {
		name         string
		expectedErr  bool
		id           uint
		expectedData *models.User
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
			expectedData: &models.User{
				Email:         "MrRogers@mail.com",
				FirstName:     "Mr",
				LastName:      "Rogers",
				Username:      "MrRogers",
				RoleID:        1,
				AccountID:     1,
				PrimaryTeamID: 1,
				Password:      "newPass",
				Token:         "asdf",
				Base: models.Base{
					Model: gorm.Model{ID: 2},
				},
			},
		},
	}

	db, err := mockstore.NewDataBaseConnection()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := mockstore.InsertRowsFor(db, superAdmin, cases[1].expectedData); err != nil {
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
}

func TestFindByUsername(t *testing.T) {
	cases := []struct {
		name         string
		expectedErr  bool
		username     string
		expectedData *models.User
	}{
		{
			name:        "User does not exist",
			expectedErr: true,
			username:    "notExists",
		},
		{
			name:     "Success",
			username: "indianajones",
			expectedData: &models.User{
				Email:         "indianajones@mail.com",
				FirstName:     "Indiana",
				LastName:      "Jones",
				Username:      "indianajones",
				RoleID:        1,
				AccountID:     1,
				PrimaryTeamID: 1,
				Password:      "newPass",
				Base: models.Base{
					Model: gorm.Model{ID: 2},
				},
			},
		},
	}

	db, err := mockstore.NewDataBaseConnection()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := mockstore.InsertRowsFor(db, superAdmin, cases[1].expectedData); err != nil {
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
}

func TestFindByToken(t *testing.T) {
	cases := []struct {
		name         string
		expectedErr  bool
		token        string
		expectedData *models.User
	}{
		{
			name:        "User does not exist",
			expectedErr: true,
			token:       "notExists",
		},
		{
			name:  "Success",
			token: "loginrefresh",
			expectedData: &models.User{
				Email:         "CharlieDarwin@mail.com",
				FirstName:     "Charles",
				LastName:      "Darwin",
				Username:      "CharlieDarwin",
				RoleID:        1,
				AccountID:     1,
				PrimaryTeamID: 1,
				Password:      "hunter2",
				Base: models.Base{
					Model: gorm.Model{ID: 1},
				},
				Token: "loginrefresh",
			},
		},
	}

	db, err := mockstore.NewDataBaseConnection()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := mockstore.InsertRowsFor(db, superAdmin, cases[1].expectedData); err != nil {
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
}

func TestList(t *testing.T) {
	cases := []struct {
		name         string
		expectedErr  bool
		qp           *models.ListQuery
		pg           *models.Pagination
		expectedData []models.User
	}{
		{
			name:        "Success, should return all 2 records",
			expectedErr: false,
			pg: &models.Pagination{
				Limit:  100,
				Offset: 0,
			},
			qp: &models.ListQuery{
				ID:    1,
				Query: "account_id = ?",
			},
			expectedData: []models.User{
				{
					Email:         "ElizabethSmart@mail.com",
					FirstName:     "Elizabeth",
					LastName:      "Smart",
					Username:      "ElizabethSmart",
					RoleID:        1,
					AccountID:     1,
					PrimaryTeamID: 1,
					Password:      "newPass",
					Base: models.Base{
						Model: gorm.Model{
							ID: 1,
						},
					},
				},
				{
					Email:         "amandacena@mail.com",
					FirstName:     "Amanda",
					LastName:      "Cena",
					Username:      "amandacena",
					RoleID:        1,
					AccountID:     1,
					PrimaryTeamID: 1,
					Password:      "hunter2",
					Token:         "loginrefresh",
					Base: models.Base{
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
			pg: &models.Pagination{
				Limit:  1,
				Offset: 1,
			},
			qp: &models.ListQuery{
				ID:    1,
				Query: "account_id = ?",
			},
			expectedData: []models.User{
				{
					Email:         "amandacena@mail.com",
					FirstName:     "Amanda",
					LastName:      "Cena",
					Username:      "amandacena",
					RoleID:        1,
					AccountID:     1,
					PrimaryTeamID: 1,
					Password:      "hunter2",
					Token:         "loginrefresh",
					Base: models.Base{
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
			pg: &models.Pagination{
				Limit:  100,
				Offset: 0,
			},
			qp: &models.ListQuery{
				ID:    99,
				Query: "account_id = ?",
			},
			expectedData: []models.User{},
		},
		{
			name:        "Success, should return all 3 records if no query is made",
			expectedErr: false,
			pg: &models.Pagination{
				Limit:  100,
				Offset: 0,
			},
			qp: nil,
			expectedData: []models.User{
				{
					Email:         "ElizabethSmart@mail.com",
					FirstName:     "Elizabeth",
					LastName:      "Smart",
					Username:      "ElizabethSmart",
					RoleID:        1,
					AccountID:     1,
					PrimaryTeamID: 1,
					Password:      "newPass",
					Base: models.Base{
						Model: gorm.Model{
							ID: 1,
						},
					},
				},
				{
					Email:         "amandacena@mail.com",
					FirstName:     "Amanda",
					LastName:      "Cena",
					Username:      "amandacena",
					RoleID:        1,
					AccountID:     1,
					PrimaryTeamID: 1,
					Password:      "hunter2",
					Token:         "loginrefresh",
					Base: models.Base{
						Model: gorm.Model{
							ID: 2,
						},
					},
				},
				{
					Email:         "sarahsmith@mail.com",
					FirstName:     "Sarah",
					LastName:      "Smith",
					Username:      "sarahsmith",
					RoleID:        1,
					AccountID:     3,
					PrimaryTeamID: 3,
					Password:      "hunter2",
					Token:         "loginrefresh",
					Base: models.Base{
						Model: gorm.Model{
							ID: 3,
						},
					},
				},
			},
		},
	}

	db, err := mockstore.NewDataBaseConnection()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := mockstore.InsertRowsFor(db, superAdmin, &cases[3].expectedData[0], &cases[3].expectedData[1], &cases[3].expectedData[2]); err != nil {
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
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name         string
		expectedErr  bool
		usr          *models.User
		expectedData *models.User
	}{
		{
			name: "Success",
			usr: &models.User{
				Base: models.Base{
					Model: gorm.Model{
						ID: 2,
					},
				},
				Email:     "iamold@village.com",
				FirstName: "OldName",
				LastName:  "Antiques",
				Address:   "1908 VintageHouse",
				Phone:     "123456",
				Mobile:    "345678",
				Username:  "OldSchool",
			},
			expectedData: &models.User{
				Email:         "refresh@mail.com",
				FirstName:     "Cool",
				LastName:      "Hip",
				Username:      "offthehook",
				RoleID:        1,
				AccountID:     1,
				PrimaryTeamID: 1,
				Password:      "newPass",
				Address:       "2020 forme",
				Phone:         "123456",
				Mobile:        "345678",
				Base: models.Base{
					Model: gorm.Model{ID: 2},
				},
			},
		},
	}

	db, err := mockstore.NewDataBaseConnection()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := mockstore.InsertRowsFor(db, superAdmin, cases[0].usr); err != nil {
		t.Error(err)
	}

	udb := store.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user := &models.User{}
			if err := db.First(user, tt.usr.ID).Error; err != nil {
				t.Error(err)
			}
			tt.expectedData.CreatedAt = user.CreatedAt
			tt.expectedData.LastLogin = user.LastLogin
			tt.expectedData.LastPasswordChange = user.LastPasswordChange
			err := udb.Update(db, tt.expectedData)
			assert.Equal(t, tt.expectedErr, err != nil)
			user = &models.User{}
			if err := db.First(user, tt.usr.ID).Error; err != nil {
				t.Error(err)
			}
			tt.expectedData.UpdatedAt = user.UpdatedAt
			assert.Equal(t, tt.expectedData, user)
		})
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name         string
		id           uint
		expectedData *models.User
	}{
		{
			name: "Success",
			id:   2,
			expectedData: &models.User{
				Email:         "tomjones@mail.com",
				FirstName:     "Tom",
				LastName:      "Jones",
				Username:      "tomjones",
				RoleID:        1,
				AccountID:     1,
				PrimaryTeamID: 1,
				Password:      "newPass",
				Base: models.Base{
					Model: gorm.Model{
						ID: 2,
					},
				},
			},
		},
	}

	db, err := mockstore.NewDataBaseConnection()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := mockstore.InsertRowsFor(db, superAdmin, cases[0].expectedData); err != nil {
		t.Error(err)
	}

	udb := store.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			userBefore := new(models.User)
			if err := db.Unscoped().Where("id = ?", tt.id).First(&userBefore).Error; err != nil {
				assert.Equal(t, nil, err, "user should exist in db store")
			}
			assert.Nil(t, userBefore.DeletedAt, "before user is deleted their deleted_at field should be set to NULL")

			err := udb.Delete(db, userBefore)
			assert.Nil(t, err, fmt.Sprintf("should not error on delete, error: %v", err))

			userAfter, err := udb.View(db, tt.id)
			emptyUser := new(models.User)
			assert.Equal(t, true, err != nil, "there should be an error when accessing deleted records")
			if err != nil {
				assert.Equal(t, "code=404, message=user not found", err.Error(), "error should be `code=404, message=user not found`")
			}
			assert.Equal(t, emptyUser, userAfter, "the response to find deleted user should be empty user")

			if err := db.Unscoped().Where("id = ?", tt.id).First(&emptyUser).Error; err != nil {
				assert.Nil(t, err, fmt.Sprintf("user should exist in db store and should be accessible with db.Unscopped(), error: %v", err))
			}
			assert.NotNil(t, emptyUser.DeletedAt, "the user should have a time set for deleted_at at the same time as when the user was deleted")
		})
	}
}

package mysqldb_test

import (
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/johncoleman83/cerebrum/pkg/utl/model"
	"github.com/johncoleman83/cerebrum/pkg/api/user/platform/mysqldb"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock/mockdb"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		req      cerebrum.User
		wantData *cerebrum.User
	}{
		{
			name:    "User already exists",
			wantErr: true,
			req: cerebrum.User{
				Email:    "johndoe@mail.com",
				Username: "johndoe",
			},
		},
		{
			name:    "Fail on insert duplicate ID",
			wantErr: true,
			req: cerebrum.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "pass",
				Base: cerebrum.Base{Model: gorm.Model{
					ID: 1,
				}},
			},
		},
		{
			name: "Success",
			req: cerebrum.User{
				Email:      "newtomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "newtomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "pass",
				Base: cerebrum.Base{Model: gorm.Model{
					ID: 2,
				}},
			},
			wantData: &cerebrum.User{
				Email:      "newtomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "newtomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "pass",
				Base: cerebrum.Base{Model: gorm.Model{
					ID: 2,
				}},
			},
		},
	}

	dbContainer, cfg := mockdb.MySqlTestContainerConfig(t)
	defer dbContainer.Shutdown()

	db := mockdb.NewDBConn(t, cfg, &cerebrum.Role{}, &cerebrum.User{})

	if err := mockdb.InsertMultiple(db, &cerebrum.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, &cases[1].req); err != nil {
		t.Error(err)
	}

	udb := mysqldb.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := udb.Create(db, tt.req)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				if resp == nil {
					t.Error("Expected data, but received nil.")
					return
				}
				tt.wantData.CreatedAt = resp.CreatedAt
				tt.wantData.UpdatedAt = resp.UpdatedAt
				assert.Equal(t, tt.wantData, resp)
			}
		})
	}
	db.Close()
}

func TestView(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		id       uint
		wantData *cerebrum.User
	}{
		{
			name:    "User does not exist",
			wantErr: true,
			id:      1000,
		},
		{
			name: "Success",
			id:   2,
			wantData: &cerebrum.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "newPass",
				Base: cerebrum.Base{Model: gorm.Model{
					ID: 2,
				}},
				Role: &cerebrum.Role{
					ID:          1,
					AccessLevel: 1,
					Name:        "SUPER_ADMIN",
				},
			},
		},
	}

	dbContainer, cfg := mockdb.MySqlTestContainerConfig(t)
	defer dbContainer.Shutdown()

	db := mockdb.NewDBConn(t, cfg, &cerebrum.Role{}, &cerebrum.User{})

	if err := mockdb.InsertMultiple(db, &cerebrum.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := mysqldb.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := udb.View(db, tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				if user == nil {
					t.Errorf("response was nil due to: %v", err)
				} else {
					tt.wantData.CreatedAt = user.CreatedAt
					tt.wantData.UpdatedAt = user.UpdatedAt
					assert.Equal(t, tt.wantData, user)
				}
			}
		})
	}
	db.Close()
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		usr      *cerebrum.User
		wantData *cerebrum.User
	}{
		{
			name: "Success",
			usr: &cerebrum.User{
				Base: cerebrum.Base{Model: gorm.Model{
					ID: 2,
				}},
				FirstName: "Z",
				LastName:  "Freak",
				Address:   "Address",
				Phone:     "123456",
				Mobile:    "345678",
				Username:  "newUsername",
			},
			wantData: &cerebrum.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Z",
				LastName:   "Freak",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "newPass",
				Address:    "Address",
				Phone:      "123456",
				Mobile:     "345678",
				Base: cerebrum.Base{Model: gorm.Model{
					ID: 2,
				}},
			},
		},
	}

	dbContainer, cfg := mockdb.MySqlTestContainerConfig(t)
	defer dbContainer.Shutdown()

	db := mockdb.NewDBConn(t, cfg, &cerebrum.Role{}, &cerebrum.User{})

	if err := mockdb.InsertMultiple(db, &cerebrum.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, cases[0].usr); err != nil {
		t.Error(err)
	}

	udb := mysqldb.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := udb.Update(db, tt.wantData)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				user := &cerebrum.User{}
				if err := db.First(user, tt.usr.ID).Error; err != nil {
					t.Error(err)
				}
				tt.wantData.UpdatedAt = user.UpdatedAt
				tt.wantData.CreatedAt = user.CreatedAt
				tt.wantData.LastLogin = user.LastLogin
				tt.wantData.DeletedAt = user.DeletedAt
				assert.Equal(t, tt.wantData, user)
			}
		})
	}
	db.Close()
}

func TestList(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		qp       *cerebrum.ListQuery
		pg       *cerebrum.Pagination
		wantData []cerebrum.User
	}{
		{
			name:    "Invalid pagination values",
			wantErr: true,
			pg: &cerebrum.Pagination{
				Limit: -100,
			},
		},
		{
			name: "Success",
			pg: &cerebrum.Pagination{
				Limit:  100,
				Offset: 0,
			},
			qp: &cerebrum.ListQuery{
				ID:    1,
				Query: "company_id = ?",
			},
			wantData: []cerebrum.User{
				{
					Email:      "tomjones@mail.com",
					FirstName:  "Tom",
					LastName:   "Jones",
					Username:   "tomjones",
					RoleID:     1,
					CompanyID:  1,
					LocationID: 1,
					Password:   "newPass",
					Base: cerebrum.Base{Model: gorm.Model{
						ID: 2,
					}},
					Role: &cerebrum.Role{
						ID:          1,
						AccessLevel: 1,
						Name:        "SUPER_ADMIN",
					},
				},
				{
					Email:      "johndoe@mail.com",
					FirstName:  "John",
					LastName:   "Doe",
					Username:   "johndoe",
					RoleID:     1,
					CompanyID:  1,
					LocationID: 1,
					Password:   "hunter2",
					Base: cerebrum.Base{Model: gorm.Model{
						ID: 1,
					}},
					Role: &cerebrum.Role{
						ID:          1,
						AccessLevel: 1,
						Name:        "SUPER_ADMIN",
					},
					Token: "loginrefresh",
				},
			},
		},
	}

	dbContainer, cfg := mockdb.MySqlTestContainerConfig(t)
	defer dbContainer.Shutdown()

	db := mockdb.NewDBConn(t, cfg, &cerebrum.Role{}, &cerebrum.User{})

	if err := mockdb.InsertMultiple(db, &cerebrum.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, &cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := mysqldb.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			users, err := udb.List(db, tt.qp, tt.pg)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				for i, v := range users {
					tt.wantData[i].CreatedAt = v.CreatedAt
					tt.wantData[i].UpdatedAt = v.UpdatedAt
				}
				assert.Equal(t, tt.wantData, users)
			}
		})
	}
	db.Close()
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		usr      *cerebrum.User
		wantData *cerebrum.User
	}{
		{
			name: "Success",
			usr: &cerebrum.User{
				Base: cerebrum.Base{Model: gorm.Model{
					ID:        2,
					DeletedAt: mock.TestTimePtr(2018),
				}},
			},
			wantData: &cerebrum.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "newPass",
				Base: cerebrum.Base{Model: gorm.Model{
					ID: 2,
				}},
			},
		},
	}

	dbContainer, cfg := mockdb.MySqlTestContainerConfig(t)
	defer dbContainer.Shutdown()

	db := mockdb.NewDBConn(t, cfg, &cerebrum.Role{}, &cerebrum.User{})

	if err := mockdb.InsertMultiple(db, &cerebrum.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, cases[0].wantData); err != nil {
		t.Error(err)
	}

	udb := mysqldb.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			err := udb.Delete(db, tt.usr)
			assert.Equal(t, tt.wantErr, err != nil)

			// TODO: See if below message means an updated needed
			// Check if the deleted_at was set
		})
	}
	db.Close()
}

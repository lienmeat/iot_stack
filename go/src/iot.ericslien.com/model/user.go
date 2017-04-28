package model

import (
	"fmt"
	"strconv"

	"iot.ericslien.com/database"
)

type User struct {
	ID       uint64   `json:"id" form:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Devices  []Device `gorm:"many2many:device_users;"`
}

func NewUser() CrudBase {
	return &User{}
}

func (u *User) GetId() uint64 {
	return u.ID
}

func (u *User) IdToStr(id uint64) string {
	return fmt.Sprintf("%d", id)
}

func (u *User) StrToId(id string) uint64 {
	n, _ := strconv.ParseUint(id, 10, 64)
	return n
}

func (u *User) Load(id uint64) error {
	return Load(u, id)
}

func (u *User) Search(params map[string]interface{}, limit int, offset int) ([]CrudBase, error) {
	us := []*User{}
	params = FilterSearchParams(u, params)
	err := Search(&us, params, limit, offset)
	ret := make([]CrudBase, len(us), len(us))
	for i := range us {
		ret[i] = us[i]
	}
	return ret, err
}

func (u *User) Create() error {
	return Create(u)
}

func (u *User) Update() error {
	return Update(u)
}

func (u *User) Delete() error {
	return Delete(u)
}

func (u *User) Migrate() error {
	db, _ := database.Open()
	return db.AutoMigrate(&User{}).Error
}

package model

import (
	"fmt"
	"strconv"
	"testing"

	"iot.ericslien.com/config"
	"iot.ericslien.com/database"
)

func init() {
	//you must customize these for your database to test.
	c := config.SQLDatabase{User: "root", Host: "localhost", Password: "123", Driver: "mysql", Port: "3306", Database: "test"}
	database.SetConfig(c)
	db, _ := database.Open()
	db.AutoMigrate(&Modeltests{})
}

type Modeltests struct {
	ID     uint64
	Name   string
	Number uint64
}

// func NewM() *Modeltests {
// 	return &Modeltests{}
// }

func (m *Modeltests) GetId() uint64 {
	return m.ID
}

func (m *Modeltests) IdToStr(id uint64) string {
	return fmt.Sprintf("%d", id)
}

func (m *Modeltests) StrToId(id string) uint64 {
	u, _ := strconv.ParseUint(id, 10, 64)
	return u
}

func (m *Modeltests) Load(id uint64) error {
	return Load(m, id)
}

func (m *Modeltests) Search(params map[string]interface{}, limit int, offset int) ([]CrudBase, error) {
	ms := []*Modeltests{}
	err := Search(&ms, params, limit, offset)
	ret := make([]CrudBase, len(ms), len(ms))
	for i := range ms {
		ret[i] = ms[i]
	}
	return ret, err
}

func (m *Modeltests) Create() error {
	return Create(m)
}

func (m *Modeltests) Update() error {
	return Update(m)
}

func (m *Modeltests) Delete() error {
	return Delete(m)
}

func TestCreate(t *testing.T) {
	m := Modeltests{Number: 1234}
	err := m.Create()
	if err != nil || m.GetId() == 0 {
		t.Error(fmt.Sprintf("Could not create: %s", err))
	}
	m.Delete()
}

func TestUpdate(t *testing.T) {
	m := Modeltests{Number: 1234}
	m.Create()
	m.Number = 123456
	err := m.Update()
	if err != nil {
		t.Error(fmt.Sprintf("Could not update: %s", err))
	}
	m.Delete()
}

func TestDelete(t *testing.T) {
	m := Modeltests{Number: 1234}
	m.Create()
	err := m.Delete()
	if err != nil {
		t.Error(fmt.Sprintf("Could not delete: %s", err))
	}
}

func TestSearch(t *testing.T) {
	m := Modeltests{Number: 12349}
	m.Create()
	m2 := Modeltests{Number: 12349, Name: "2"}
	m2.Create()
	m3 := Modeltests{Number: 12349}
	m3.Create()
	params := make(map[string]interface{})
	params["number"] = m.Number

	ms, err := m.Search(params, 1, 0)
	if err != nil || len(ms) != 1 {
		fmt.Printf("%s", ms)
		t.Error(fmt.Sprintf("Could not limit search: %s", err))
	}

	ms, err = m.Search(params, 1, 1)
	if err != nil || len(ms) != 1 {
		t.Error(fmt.Sprintf("Could not offset search: %s", err))
	} else {
		tmp := ms[0].(*Modeltests)
		if tmp.Name != "2" {
			t.Error(fmt.Sprintf("Could not offset search, got wrong object"))
		}
	}

	ms, err = m.Search(params, 0, 0)
	if err != nil || len(ms) < 3 {
		fmt.Printf("%d", len(ms))
		t.Error(fmt.Sprintf("Could not search: %s", err))
	} else {
		fmt.Printf("%d", len(ms))
	}
	m.Delete()
	m2.Delete()
	m3.Delete()
}

func TestLoad(t *testing.T) {
	m := Modeltests{Number: 1234}
	m.Create()
	mr := Modeltests{}
	err := mr.Load(m.GetId())
	if err != nil || mr.GetId() != m.GetId() {
		t.Error(fmt.Sprintf("Could not load: %s", err))
	}
	m.Delete()
}

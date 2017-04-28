package model

import (
	"fmt"
	"strconv"

	"iot.ericslien.com/database"
)

//DevicePassToken must be used by devices when registering
const DevicePassToken string = "VSpXJVt0NyJMYTZuVmEwLy1xR1tYKDc7SFo5YUslcFFaSGRCLjIpLEM="

type Device struct {
	ID         uint64 `json:"id" form:"id"`
	Type       string `json:"type" form:"type"`
	ExternalID string `json:"external_id" form:"external_id" gorm:"index"`
	Serial     string `json:"serial" form:"serial" gorm:"index"`
}

func NewDevice() CrudBase {
	return &Device{}
}

func (d *Device) GetId() uint64 {
	return d.ID
}

func (d *Device) IdToStr(id uint64) string {
	return fmt.Sprintf("%d", id)
}

func (d *Device) StrToId(id string) uint64 {
	u, _ := strconv.ParseUint(id, 10, 64)
	return u
}

func (d *Device) Load(id uint64) error {
	return Load(d, id)
}

func (d *Device) Search(params map[string]interface{}, limit int, offset int) ([]CrudBase, error) {
	ds := []*Device{}
	params = FilterSearchParams(d, params)
	err := Search(&ds, params, limit, offset)
	ret := make([]CrudBase, len(ds), len(ds))
	for i := range ds {
		ret[i] = ds[i]
	}
	return ret, err
}

func (d *Device) Create() error {
	return Create(d)
}

func (d *Device) Update() error {
	return Update(d)
}

func (d *Device) Delete() error {
	return Delete(d)
}

func (d *Device) Migrate() error {
	db, _ := database.Open()
	return db.AutoMigrate(&Device{}).Error
}

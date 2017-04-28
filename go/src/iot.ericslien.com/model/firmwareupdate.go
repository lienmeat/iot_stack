package model

import (
	"fmt"
	"strconv"

	"iot.ericslien.com/database"
)

type FirmwareUpdate struct {
	ID         uint64 `json:"id" form:"id"`
	DeviceType string `json:"device_type" form:"device_type" gorm:"index"`
	Version    string `json:"version" form:"version"`
	File       string `json:"file" form:"file"`
}

func NewFirmwareUpdate() CrudBase {
	return &FirmwareUpdate{}
}

func (o *FirmwareUpdate) GetId() uint64 {
	return o.ID
}

func (o *FirmwareUpdate) IdToStr(id uint64) string {
	return fmt.Sprintf("%d", id)
}

func (o *FirmwareUpdate) StrToId(id string) uint64 {
	u, _ := strconv.ParseUint(id, 10, 64)
	return u
}

func (o *FirmwareUpdate) Load(id uint64) error {
	return Load(o, id)
}

func (o *FirmwareUpdate) Search(params map[string]interface{}, limit int, offset int) ([]CrudBase, error) {
	os := []*FirmwareUpdate{}
	params = FilterSearchParams(o, params)
	err := Search(&os, params, limit, offset)
	ret := make([]CrudBase, len(os), len(os))
	for i := range os {
		ret[i] = os[i]
	}
	return ret, err
}

func (o *FirmwareUpdate) Create() error {
	return Create(o)
}

func (o *FirmwareUpdate) Update() error {
	return Update(o)
}

func (o *FirmwareUpdate) Delete() error {
	return Delete(o)
}

func (o *FirmwareUpdate) Migrate() error {
	db, _ := database.Open()
	return db.AutoMigrate(&FirmwareUpdate{}).Error
}

func (o *FirmwareUpdate) GetLatestForDevice(devicetype string, version string) error {
	db, _ := database.Open()
	return db.Where("device_type = ? AND version > ?", devicetype, version).Order("version desc").First(&o).Error
}

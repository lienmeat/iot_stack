package model

import (
	"fmt"
	"net/http"

	"iot.ericslien.com/config"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(e *gin.Engine) {
	//TODO make some form of auth tokens and middlewhere that uses them
	handlergroups := CrudHandlerGroups{}
	SetupCrudRest("id", "devices", NewDevice, e, &handlergroups)
	SetupCrudRest("id", "users", NewUser, e, &handlergroups)
	SetupCrudRest("id", "firmwareupdates", NewFirmwareUpdate, e, &handlergroups)

	auth := e.Group("/")

	auth.GET("/device/register/:token/:type/:external_id/:serial", registerDevice)
	auth.GET("/device/update/:type/:version", updateDevice)
	//yes normally this would be a post, but it's just not going to matter
	//and wget is easier to use from shell scripts this will be run with/against
	auth.GET("/firmwareupdate/new/:token/:type/:version/:filename", addFirmwareUpdate)
	// auth.GET("/vr/:id/voter/:vid", getByVid)
	// auth.GET("/vr/:id/voters/lookup", lookupVoters)
	// auth.POST("/vr/:id/voters/import", importVoters)
}

func registerDevice(c *gin.Context) {
	token := c.Param("token")
	t := c.Param("type")
	externalID := c.Param("external_id")
	serial := c.Param("serial")

	if token != DevicePassToken || t == "" || externalID == "" || serial == "" {
		c.JSON(http.StatusForbidden, "GET OUT OF HERE")
	}

	d := Device{}
	params := make(map[string]interface{}, 0)
	params["serial"] = serial
	params["external_id"] = externalID
	ds, err := d.Search(params, 1, 0)
	if len(ds) == 0 && err == nil {
		d.Serial = serial
		d.ExternalID = externalID
		d.Type = t
		err = d.Create()
		if err != nil {
			c.JSON(http.StatusInternalServerError, "")
			return
		}
		c.JSON(200, d)
	}
	c.JSON(http.StatusNotFound, "")
}

func addFirmwareUpdate(c *gin.Context) {
	token := c.Param("token")
	t := c.Param("type")
	filename := c.Param("filename")
	version := c.Param("version")

	if token != DevicePassToken || t == "" || filename == "" || version == "" {
		c.JSON(http.StatusForbidden, "GET OUT OF HERE")
		return
	}

	fw := FirmwareUpdate{DeviceType: t, File: filename, Version: version}
	if err := fw.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "OK")
}

func updateDevice(c *gin.Context) {
	cnf := config.Get()
	t := c.Param("type")
	ver := c.Param("version")

	if t == "" || ver == "" {
		c.JSON(http.StatusBadRequest, "MISSING URL PARAMS")
		return
	}

	fw := FirmwareUpdate{}
	if fw.GetLatestForDevice(t, ver) != nil {
		c.JSON(http.StatusNotFound, "No update available")
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+fw.File)
	c.Header("Content-Type", "application/octet-stream")
	file := fmt.Sprintf("%s/firmware/%s", cnf.FilesDir, fw.File)
	logrus.Debug(fmt.Sprintf("Sending firmware update file: %s", file))
	c.File(file)
}

package model

import (
	"errors"
	"fmt"

	"net/http"
	"reflect"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	// 	"strings"

	"strconv"
)

type CrudHandlerGroups struct {
	Search *gin.RouterGroup
	Get    *gin.RouterGroup
	Post   *gin.RouterGroup
	Update *gin.RouterGroup
	Delete *gin.RouterGroup
}

func SetupCrudRest(idParam string, route string, newModel NewModel, server *gin.Engine, handlerGroups *CrudHandlerGroups) {

	// Fill in the handlers with defaults if they aren't defined
	populateHandlerGroups(server, handlerGroups)

	handlerGroups.Search.GET(route, restSearch(newModel, server, handlerGroups))

	// Get request to load a single item
	handlerGroups.Get.GET(route+"/:"+idParam, restGet(idParam, newModel, server, handlerGroups))

	// Support post to create new items
	handlerGroups.Post.POST(route, restPost(newModel, server, handlerGroups))

	// Support updates
	handlerGroups.Update.PUT(route+"/:"+idParam, restPut(idParam, newModel, server, handlerGroups))

	// Support delete operations
	handlerGroups.Delete.DELETE(route+"/:"+idParam, restDelete(idParam, newModel, server, handlerGroups))
}

/**
* Returns closure to handle searches for an object
 */
func restSearch(newModel NewModel, server *gin.Engine, handlerGroups *CrudHandlerGroups) func(c *gin.Context) {
	// Support generic list requests to get all items
	return func(c *gin.Context) {
		logrus.Debug("SEARCH")
		// Get query params
		search, limit, offset := GetSearchFromGin(c)

		// Grab a new model to work with
		model := newModel()

		response, err := model.Search(search, limit, offset)
		if err != nil {
			c.Error(err)
			return
		}

		logrus.WithFields(logrus.Fields{
			"context": "Crud Access",
			"id":      response,
		})
		//need something similar once we get auth figured out!
		// @todo: #performance Calls to interface search have to copy the array over. Best way to go?
		// for i := 0; i < len(response); {
		// 	if ok := UserAccess(response[i].(AccessControlled), RequestPermView, c); !ok {
		// 		logrus.WithFields(logrus.Fields{
		// 			"context": "Crud Access",
		// 			"id":      response[i].GetId(),
		// 			"type":    reflect.TypeOf(response[i]),
		// 		}).Debug("Access denied. Removing from return list.")
		// 		response = append(response[:i], response[i+1:]...)
		// 	} else {
		// 		i++
		// 	}
		// }
		c.JSON(http.StatusOK, response)
	}
}

/**
* Returns closure to handle GET requests for an object
 */
func restGet(idParam string, newModel NewModel, server *gin.Engine, handlerGroups *CrudHandlerGroups) func(c *gin.Context) {
	return func(c *gin.Context) {
		logrus.Debug("GET")
		// Grab a new model to work with
		model := newModel()

		// Convert to a Model ID
		id := model.StrToId(c.Param(idParam))

		// Try to load the object
		err := model.Load(id)
		// if err != nil {
		// 	c.Error(err)
		// 	c.JSON(http.StatusNotFound, model)
		// 	return
		// }

		// if ok := UserAccess(model, RequestPermView, c); !ok {
		// 	GinExitWithAccessDenied(c)
		// 	return
		// }

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"context": "CRUD",
				// "userId": c.MustGet(common.RequestTokenClaimsKey).(*common.UserTokenClaims).Subject,
				"objectId":   model.GetId(),
				"objectType": reflect.TypeOf(model).String(),
			}).Error("Unable to get object by id for model.")
			c.Error(err)
			c.JSON(http.StatusNotFound, model)
			return
		}

		// Everything went well
		c.JSON(200, model)
	}
}

/**
* Returns closure to handle POST requests for an object
 */
func restPost(newModel NewModel, server *gin.Engine, handlerGroups *CrudHandlerGroups) func(c *gin.Context) {
	return func(c *gin.Context) {
		logrus.Debug("POST")
		// Grab a new model to work with
		model := newModel()

		// Map all incoming data to the object
		err := c.Bind(model)
		logrus.Debug(fmt.Sprintf("model after bind: %+v", model))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"context": "CRUD",
				// "userId": c.MustGet(common.RequestTokenClaimsKey).(*common.UserTokenClaims).Subject,
				"objectId":   model.GetId(),
				"objectType": reflect.TypeOf(model).String(),
			}).Error("Unable to bind post data to model.")
			c.Error(err)
			return
		}

		// if ok := UserAccess(model, RequestPermCreate, c); !ok {
		// 	logrus.WithFields(logrus.Fields{
		// 		"context": "CRUD",
		// 		// "userId": c.MustGet(common.RequestTokenClaimsKey).(*common.UserTokenClaims).Subject,
		// 		"objectId": model.GetId(),
		// 		"objectType": reflect.TypeOf(model).String(),
		// 	}).Warn("User attempted to CREATE object without permissions.")
		// 	GinExitWithAccessDenied(c)
		// 	return
		// }

		err = model.Create()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"context": "CRUD",
				// "userId": c.MustGet(common.RequestTokenClaimsKey).(*common.UserTokenClaims).Subject,
				"objectId":   model.GetId(),
				"objectType": reflect.TypeOf(model).String(),
			}).Error("Error creating object.")
			c.Error(err)
			return
		}

		// Everything went well, return the result
		logrus.WithFields(logrus.Fields{
			"context": "CRUD",
			// "userId": c.MustGet(common.RequestTokenClaimsKey).(*common.UserTokenClaims).Subject,
			"objectId":   model.GetId(),
			"objectType": reflect.TypeOf(model).String(),
		}).Info("User CREATED object.")
		c.JSON(http.StatusOK, model)
	}
}

/**
* Returns closure to handle PUT requests for an object
 */
func restPut(idParam string, newModel NewModel, server *gin.Engine, handlerGroups *CrudHandlerGroups) func(c *gin.Context) {
	return func(c *gin.Context) {
		logrus.Debug("PUT")
		// Grab a new model to work with
		model := newModel()

		id := model.StrToId(c.Param(idParam))

		// original := newModel()

		err := model.Load(id)
		if err != nil {
			c.Error(err)
			return
		}

		// Map all incoming data to the object
		err = c.Bind(model)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"context": "CRUD",
				// "userId": c.MustGet(common.RequestTokenClaimsKey).(*common.UserTokenClaims).Subject,
				"objectId":   model.GetId(),
				"objectType": reflect.TypeOf(model).String(),
			}).Error("unable to bind post data to model")
			c.Error(err)
			return
		}

		if id != model.GetId() {
			c.Error(errors.New("posted ID does not match parameter"))
			return
		}

		// if ok := UserAccess(model, RequestPermUpdate, c); !ok {
		// 	logrus.WithFields(logrus.Fields{
		// 		"context": "CRUD",
		// 		// "userId": c.MustGet(common.RequestTokenClaimsKey).(*common.UserTokenClaims).Subject,
		// 		"objectId": model.GetId(),
		// 		"objectType": reflect.TypeOf(model).String(),
		// 	}).Warn("User attempted to UPDATE object without permissions.")
		// 	GinExitWithAccessDenied(c)
		// 	return
		// }

		err = model.Update()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"context": "CRUD",
				// "userId": c.MustGet(common.RequestTokenClaimsKey).(*common.UserTokenClaims).Subject,
				"objectId":   model.GetId(),
				"objectType": reflect.TypeOf(model).String(),
			}).Error("Error updating object.")
			c.Error(err)
			return
		}

		// All is good
		logrus.WithFields(logrus.Fields{
			"context": "CRUD",
			// "userId": c.MustGet(common.RequestTokenClaimsKey).(*common.UserTokenClaims).Subject,
			"objectId":   model.GetId(),
			"objectType": reflect.TypeOf(model).String(),
		}).Info("User UPDATED object.")
		c.JSON(http.StatusOK, model)
	}
}

/**
* Returns closure to handle DELETE requests for an object
 */
func restDelete(idParam string, newModel NewModel, server *gin.Engine, handlerGroups *CrudHandlerGroups) func(c *gin.Context) {
	return func(c *gin.Context) {
		logrus.Debug("DELETE")
		// Grab a new model to work with
		model := newModel()

		id := model.StrToId(c.Param(idParam))

		err := model.Load(id)
		if err != nil {
			c.Error(err)
			return
		}

		// if ok := UserAccess(model, RequestPermDelete, c); !ok {
		// 	logrus.WithFields(logrus.Fields{
		// 		"context": "CRUD",
		// 		"userId": c.MustGet(common.RequestTokenClaimsKey).(*common.UserTokenClaims).Subject,
		// 		"objectId": model.GetId(),
		// 		"objectType": reflect.TypeOf(model).String(),
		// 	}).Warn("User attempted to DELETE object without permissions.")
		// 	GinExitWithAccessDenied(c)
		// 	return
		// }

		err = model.Delete()
		if err != nil {
			c.Error(err)
			return
		}

		// All is good
		logrus.WithFields(logrus.Fields{
			"context": "CRUD",
			// "userId": c.MustGet(common.RequestTokenClaimsKey).(*common.UserTokenClaims).Subject,
			"objectId":   model.GetId(),
			"objectType": reflect.TypeOf(model).String(),
		}).Warn("User DELETED object.")
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

func populateHandlerGroups(server *gin.Engine, handlerGroups *CrudHandlerGroups) {

	//@todo: authentication is a big todo!
	//add middlewhere when we get some that will work for this use case
	requireAuthenticated := server.Group("/")

	if handlerGroups == nil {
		handlerGroups = &CrudHandlerGroups{
			Search: requireAuthenticated,
			Get:    requireAuthenticated,
			Post:   requireAuthenticated,
			Update: requireAuthenticated,
			Delete: requireAuthenticated,
		}
	} else {
		if handlerGroups.Search == nil {
			handlerGroups.Search = requireAuthenticated
		}
		if handlerGroups.Get == nil {
			handlerGroups.Get = requireAuthenticated
		}
		if handlerGroups.Post == nil {
			handlerGroups.Post = requireAuthenticated
		}
		if handlerGroups.Update == nil {
			handlerGroups.Update = requireAuthenticated
		}
		if handlerGroups.Delete == nil {
			handlerGroups.Delete = requireAuthenticated
		}
	}
}

func GetSearchFromGin(c *gin.Context) (map[string]interface{}, int, int) {
	search := make(map[string]interface{}, 0)
	limit := 0
	offset := 0
	//map[string][]string of query args and vals
	query := c.Request.URL.Query()
	var val interface{}
	for k, v := range query {
		val = v[0]
		if k == "limit" {
			limit, _ = strconv.Atoi(val.(string))
			continue
		}
		if k == "offset" {
			offset, _ = strconv.Atoi(val.(string))
			continue
		}
		search[k] = val
	}
	return search, limit, offset
}

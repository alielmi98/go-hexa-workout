package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/alielmi98/go-hexa-workout/pkg/helper"
	"github.com/gin-gonic/gin"
)

// Create an entity
// Ti: Http request body
// To: Http response body (Usecase function output)
// usecaseCreate: usecase Create method
func Create[Ti any, To any](c *gin.Context,
	usecaseCreate func(ctx context.Context, req Ti) (To, error)) {

	// bind http request
	request := new(Ti)
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}

	// call use case method
	response, err := usecaseCreate(c, *request)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(response, true, 0))
}

// Update an entity
// Ti: Http request body
// To: Http response body (Usecase function output)
// usecaseUpdate: usecase Update method
func Update[Ti any, To any](c *gin.Context,
	usecaseUpdate func(ctx context.Context,
		id int, req Ti) (To, error)) {

	// bind http request
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	request := new(Ti)
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return

	}

	// call use case method
	response, err := usecaseUpdate(c, id, *request)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(response, true, 0))
}

func Delete(c *gin.Context, usecaseDelete func(ctx context.Context, id int) error) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound,
			helper.GenerateBaseResponse(nil, false, helper.ValidationError))
		return
	}

	err := usecaseDelete(c, id)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true, 0))
}

// Get an entity
// To: Http response body (Usecase function output)
// usecaseGet: usecase Get method
func GetById[To any](c *gin.Context,
	usecaseGet func(c context.Context, id int) (To, error)) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound,
			helper.GenerateBaseResponse(nil, false, helper.ValidationError))
		return
	}

	// call use case method
	response, err := usecaseGet(c, id)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(response, true, 0))
}

// Get entities by filter
// Ti: Http request body
// To: Http response body (Usecase function output)
// usecaseList: usecase GetByFilter method
func GetByFilter[Ti any, To any](c *gin.Context,
	usecaseList func(c context.Context, req Ti) (*To, error)) {

	req := new(Ti)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}

	res, err := usecaseList(c, *req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, 0))
}

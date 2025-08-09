package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/alielmi98/go-hexa-workout/internal/workout/port/filter"
	"github.com/alielmi98/go-hexa-workout/pkg/helper"
	"github.com/gin-gonic/gin"
)

// Create an entity
// TRequest: Http request body
// TUInput: Usecase method input that mapped from TRequest with TUInput := mapper(TRequest)
// TUOutput: Usecase function output
// TResponse: Http response body that mapped from TUOutput with TResponse := mapper(TUOutput)
// requestMapper: this function map endpoint input to usecase input
// responseMapper: this function map usecase output to endpoint output
// usecaseCreate: usecase Create method
func Create[TRequest any, TUInput any, TUOutput any, TResponse any](c *gin.Context,
	requestMapper func(req TRequest) (res TUInput),
	responseMapper func(req TUOutput) (res TResponse),
	usecaseCreate func(ctx context.Context, req TUInput) (TUOutput, error)) {

	// bind http request
	request := new(TRequest)
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}

	// map http request body to usecase input
	usecaseInput := requestMapper(*request)

	// call use case method
	usecaseResult, err := usecaseCreate(c, usecaseInput)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	// map usecase response to http response
	response := responseMapper(usecaseResult)

	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(response, true, 0))
}

// Update an entity
// TRequest: Http request body
// TUInput: Usecase method input that mapped from TRequest with TUInput := mapper(TRequest)
// TUOutput: Usecase function output
// TResponse: Http response body that mapped from TUOutput with TResponse := mapper(TUOutput)
// requestMapper: this function map endpoint input to usecase input
// responseMapper: this function map usecase output to endpoint output
// usecaseUpdate: usecase Update method
func Update[TRequest any, TUInput any, TUOutput any, TResponse any](c *gin.Context,
	requestMapper func(req TRequest) (res TUInput),
	responseMapper func(req TUOutput) (res TResponse),
	usecaseUpdate func(ctx context.Context,
		id int, req TUInput) (TUOutput, error)) {

	// bind http request
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, helper.ValidationError, err))
		return
	}
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, helper.ValidationError, errors.New("invalid id")))
		return
	}

	request := new(TRequest)
	err = c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return

	}
	// map http request body to usecase input
	usecaseInput := requestMapper(*request)

	// call use case method
	usecaseResult, err := usecaseUpdate(c, id, usecaseInput)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	// map usecase response to http response
	response := responseMapper(usecaseResult)

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(response, true, 0))
}

func Delete(c *gin.Context, usecaseDelete func(ctx context.Context, id int) error) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, helper.ValidationError, err))
		return
	}
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, helper.ValidationError, errors.New("invalid id")))
		return
	}

	err = usecaseDelete(c, id)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(nil, true, 0))
}

// Get an entity
// TUOutput: Usecase function output
// TResponse: Http response body that mapped from TUOutput with TResponse := mapper(TUOutput)
// responseMapper: this function map usecase output to endpoint output
// usecaseGet: usecase Get method
func GetById[TUOutput any, TResponse any](c *gin.Context,
	responseMapper func(req TUOutput) (res TResponse),
	usecaseGet func(c context.Context, id int) (TUOutput, error)) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, helper.ValidationError, err))
		return
	}
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, helper.ValidationError, errors.New("invalid id")))
		return
	}

	// call use case method
	usecaseResult, err := usecaseGet(c, id)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	// map usecase response to http response
	response := responseMapper(usecaseResult)

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(response, true, 0))
}

// Get entities by filter
// TUOutput: Usecase function output
// TResponse: Http response body that mapped from TUOutput with TResponse := mapper(TUOutput)
// responseMapper: this function map usecase output to endpoint output
// usecaseList: usecase GetByFilter method
func GetByFilter[TUOutput any, TResponse any](c *gin.Context,
	responseMapper func(req TUOutput) (res TResponse),
	usecaseList func(c context.Context, req filter.PaginationInputWithFilter) (*filter.PagedList[TUOutput], error)) {

	req := new(filter.PaginationInputWithFilter)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}

	// call use case method
	usecaseResult, err := usecaseList(c, *req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	response := filter.PagedList[TResponse]{
		PageNumber:      usecaseResult.PageNumber,
		PageSize:        usecaseResult.PageSize,
		TotalRows:       usecaseResult.TotalRows,
		TotalPages:      usecaseResult.TotalPages,
		HasPreviousPage: usecaseResult.HasPreviousPage,
		HasNextPage:     usecaseResult.HasNextPage,
	}

	// map usecase response to http response
	items := []TResponse{}
	for _, item := range *usecaseResult.Items {

		items = append(items, responseMapper(item))
	}
	response.Items = &items

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(response, true, 0))
}

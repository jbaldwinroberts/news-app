// Package server implements the API endpoints
package server

import (
	"net/http"

	"github.com/josephroberts/esqimo/store"
	Swagger "github.com/josephroberts/esqimo/swagger"
	"github.com/labstack/echo/v4"
)

// Server contains a Store interface for reading data.
type Server struct {
	Store store.Store
}

// GetItems returns a JSON response containing a list of items.
func (s *Server) GetItems(ctx echo.Context, params Swagger.GetItemsParams) error {
	titles := checkParamSlice(params.Titles)
	categories := checkParamSlice(params.Categories)
	limit := checkParamInt(params.Limit)

	return ctx.JSON(http.StatusOK, s.Store.GetItems(titles, categories, limit))
}

// GetTitles returns a JSON response containing a list of titles.
func (s *Server) GetTitles(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, s.Store.GetTitles())
}

// GetCategories returns a JSON reponse containing a list of categories.
func (s *Server) GetCategories(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, s.Store.GetCategories())
}

// Handle null pointer.
func checkParamSlice(param *[]string) []string {
	if param == nil {
		return make([]string, 0)
	}
	return *param
}

// Handle null pointer.
func checkParamInt(param *int32) int32 {
	if param == nil {
		return 0
	}
	return *param
}

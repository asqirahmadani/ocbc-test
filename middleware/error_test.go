package middleware

import (
	"boilerplate/dto/response"
	"boilerplate/utils"
	"context"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type MockField struct {
	Name string `validate:"required"`
}

func TestErrorMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Should handle NotFoundError", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)

		r.Use(ErrorMiddleware())
		r.GET("/test", func(c *gin.Context) {
			c.Error(utils.ErrNotFound("Resource not found"))
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(recorder, req)

		var res response.ErrorDTO
		json.Unmarshal(recorder.Body.Bytes(), &res)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.Equal(t, "Resource not found", res.Message)
	})

	t.Run("Should handle Context DeadlineExceeded", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)

		r.Use(ErrorMiddleware())
		r.GET("/timeout", func(c *gin.Context) {
			c.Error(context.DeadlineExceeded)
		})

		req, _ := http.NewRequest("GET", "/timeout", nil)
		r.ServeHTTP(recorder, req)

		var res response.ErrorDTO
		json.Unmarshal(recorder.Body.Bytes(), &res)

		assert.Equal(t, http.StatusRequestTimeout, recorder.Code)
		assert.Equal(t, "Request has timeout", res.Message)
	})

	t.Run("Should handle generic Internal Server Error", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)

		r.Use(ErrorMiddleware())
		r.GET("/error", func(c *gin.Context) {
			c.Error(assert.AnError) // generic error
		})

		req, _ := http.NewRequest("GET", "/error", nil)
		r.ServeHTTP(recorder, req)

		var res response.ErrorDTO
		json.Unmarshal(recorder.Body.Bytes(), &res)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, "Uh oh, something went wrong!", res.Message)
	})

	t.Run("Should pass through if no errors occur", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)

		r.Use(ErrorMiddleware())
		r.GET("/success", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req, _ := http.NewRequest("GET", "/success", nil)
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("Coverage: context.Canceled", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)
		r.Use(ErrorMiddleware())
		r.GET("/test", func(c *gin.Context) {
			c.Error(context.Canceled)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(recorder, req)

		var res response.ErrorDTO
		json.Unmarshal(recorder.Body.Bytes(), &res)
		assert.Equal(t, http.StatusRequestTimeout, recorder.Code)
		assert.Equal(t, "Request has been cancelled", res.Message)
	})

	t.Run("Coverage: validator.ValidationErrors", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)
		r.Use(ErrorMiddleware())
		r.GET("/test", func(c *gin.Context) {
			validate := validator.New()
			err := validate.Struct(MockField{}) // Trigger missing 'Name'
			c.Error(err)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(recorder, req)

		var res response.ErrorDTO
		json.Unmarshal(recorder.Body.Bytes(), &res)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, res.Message, "Missing or incorrect type for Name")
	})

	t.Run("Coverage: json.UnmarshalTypeError", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)
		r.Use(ErrorMiddleware())
		r.GET("/test", func(c *gin.Context) {
			err := &json.UnmarshalTypeError{Field: "age", Type: nil}
			c.Error(err)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(recorder, req)

		var res response.ErrorDTO
		json.Unmarshal(recorder.Body.Bytes(), &res)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, "Incorrect type for age", res.Message)
	})

	t.Run("Coverage: BadRequestError", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)
		r.Use(ErrorMiddleware())
		r.GET("/test", func(c *gin.Context) {
			c.Error(utils.ErrBadRequest("Bad input"))
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Coverage: InternalServerError", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)
		r.Use(ErrorMiddleware())
		r.GET("/test", func(c *gin.Context) {
			c.Error(utils.ErrInternalServer("DB down"))
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("Coverage: UnprocessableEntityError", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)
		r.Use(ErrorMiddleware())
		r.GET("/test", func(c *gin.Context) {
			c.Error(utils.ErrUnprocessableEntity("Invalid logic"))
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	})
}
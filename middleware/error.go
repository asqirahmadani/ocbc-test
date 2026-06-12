package middleware

import (
	"boilerplate/dto/response"
	"boilerplate/utils"
	"context"
	"encoding/json"
	"errors"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgconn"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		errResponse := response.JSONResponse{
			Message: "",
			Data:    nil,
		}

		if errors.Is(err, context.DeadlineExceeded) {
			errResponse.Message = "waktu request telah habis (timeout)"
			c.JSON(http.StatusRequestTimeout, errResponse)

			return
		}

		if errors.Is(err, context.Canceled) {
			errResponse.Message = "request telah dibatalkan"
			c.JSON(http.StatusRequestTimeout, errResponse)

			return
		}

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			var errMsgs []string
			for _, fe := range ve {
				errMsgs = append(errMsgs, "kolom '"+strings.ToLower(fe.Field())+"' tidak valid atau kosong")
			}
			errResponse.Message = strings.Join(errMsgs, ", ")

			c.JSON(http.StatusBadRequest, errResponse)
			return
		}

		var ute *json.UnmarshalTypeError
		if errors.As(err, &ute) {
			errResponse.Message = "format data untuk kolom '" + strings.ToLower(ute.Field) + "' tidak valid."
			c.JSON(http.StatusBadRequest, errResponse)
			return
		}

		var nf utils.NotFoundError
		if errors.As(err, &nf) {
			errResponse.Message = err.Error()
			c.JSON(http.StatusNotFound, errResponse)

			return
		}

		var br utils.BadRequestError
		if errors.As(err, &br) {
			errResponse.Message = err.Error()
			c.JSON(http.StatusBadRequest, errResponse)

			return
		}

		var ae utils.AuthorizationError
		if errors.As(err, &ae) {
			errResponse.Message = err.Error()
			c.JSON(http.StatusUnauthorized, errResponse)

			return
		}

		var is utils.InternalServerError
		if errors.As(err, &is) {
			errResponse.Message = err.Error()
			c.JSON(http.StatusInternalServerError, errResponse)

			return
		}

		var ue utils.UnprocessableEntityError
		if errors.As(err, &ue) {
			errResponse.Message = err.Error()
			c.JSON(http.StatusUnprocessableEntity, errResponse)

			return
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				errResponse.Message = "data yang dimasukkan sudah terdaftar di sistem. silakan gunakan data lain."
				c.JSON(http.StatusConflict, errResponse)
				return

			case "23503": // foreign_key_violation
				errResponse.Message = "permintaan ditolak karena data ini masih terkait dengan informasi lain di dalam sistem."
				c.JSON(http.StatusBadRequest, errResponse)
				return

			case "23502": // not_null_violation
				errResponse.Message = "terdapat data wajib yang belum terisi dengan benar."
				c.JSON(http.StatusBadRequest, errResponse)
				return

			case "22P02": // invalid_text_representation
				errResponse.Message = "format data yang dikirimkan tidak valid atau tidak dapat diproses oleh sistem."
				c.JSON(http.StatusBadRequest, errResponse)
				return

			default:
				// error database lainnya
				errResponse.Message = "terjadi kendala saat memproses data. silakan coba beberapa saat lagi."
				c.AbortWithStatusJSON(http.StatusInternalServerError, errResponse)
				return
			}
		}

		errResponse.Message = "mohon maaf, terjadi gangguan pada server kami. silakan coba beberapa saat lagi."
		c.AbortWithStatusJSON(http.StatusInternalServerError, errResponse)
	}
}

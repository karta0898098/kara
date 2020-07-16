package view

import (
	"github.com/gin-gonic/gin"
	appError "github.com/karta0898098/kara/errors"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

func Error(c *gin.Context, err error) {
	if err != nil {
		errorOfApp, ok := errors.Cause(err).(*appError.AppError)
		if ok {
			if errorOfApp.Status >= http.StatusInternalServerError {
				log.Error().Msgf("%+v", err)
			} else if errorOfApp.Status >= http.StatusBadRequest {
				log.Info().Msgf("%+v", err)
			} else {
				log.Info().Msgf("%+v", err)
			}

			c.JSON(errorOfApp.Status, gin.H{
				"code":    errorOfApp.Code,
				"message": errorOfApp.Message,
				"data":    nil,
			})
		} else {
			log.Error().Msgf("%+v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500000,
				"message": http.StatusText(http.StatusInternalServerError),
				"data":    nil,
			})
		}
	}
}

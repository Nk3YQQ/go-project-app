package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func ParseDateToTime(date string, parsedDate *time.Time) error {
	layout := "02.01.2006"

	birthDate, err := time.Parse(layout, date)

	if err != nil {
		return err
	}

	*parsedDate = birthDate

	return nil
}

func RaiseBadRequestError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, err.Error())
	c.Abort()
}

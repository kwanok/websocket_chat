package utils

import (
	"github.com/gin-gonic/gin"
	"log"
)

type FatalError struct {
	Error error
}

type PanicError struct {
	Error error
}

type HttpError struct {
	Error   error
	Context *gin.Context
	Status  int
}

func (fe FatalError) Handle() {
	if fe.Error != nil {
		log.Fatal(fe.Error)
	}
}

func (pe PanicError) Handle() {
	if pe.Error != nil {
		log.Panic(pe.Error)
	}
}

func (he HttpError) Handle(messages ...string) {
	if he.Error != nil {
		if len(messages) == 0 {
			he.Context.JSON(he.Status, he.Error.Error())
		} else {
			he.Context.JSON(he.Status, messages[0])
		}

		return
	}
}

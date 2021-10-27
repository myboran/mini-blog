package app

import (
	"gin-blog/pkg/logging"
	"github.com/astaxie/beego/validation"
)

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors{
		println(err.Key, err.Message)
		logging.Info(err.Key, err.Message)
	}
	return
}
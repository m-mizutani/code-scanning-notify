package interfaces

import "github.com/m-mizutani/code-scanning-notify/pkg/domain/model"

type Usecase interface {
	Notify(model.InputNotify) error

	SetConfig(*model.Config)
}

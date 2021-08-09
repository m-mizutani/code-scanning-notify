package interfaces

import "github.com/m-mizutani/cs-alert-notify/pkg/domain/model"

type Usecase interface {
	Notify(model.InputNotify) error

	SetConfig(*model.Config)
}

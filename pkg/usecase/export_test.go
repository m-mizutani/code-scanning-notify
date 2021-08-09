package usecase

import (
	"github.com/m-mizutani/cs-alert-notify/pkg/domain/interfaces"
)

func (x *Usecase) InjectFactories(fac *interfaces.Factories) {
	x.factories = fac
}

package usecase

import (
	"github.com/m-mizutani/code-scanning-notify/pkg/domain/interfaces"
)

func (x *Usecase) InjectFactories(fac *interfaces.Factories) {
	x.factories = fac
}

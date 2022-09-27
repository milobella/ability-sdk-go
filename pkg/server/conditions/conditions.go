package conditions

import "github.com/milobella/ability-sdk-go/pkg/model"

func IfIntents(intents ...string) func(request *model.Request) (result bool) {
	return func(request *model.Request) (result bool) {
		for _, intent := range intents {
			if request.Nlu.BestIntent == intent {
				return true
			}
		}
		return false
	}
}

func IfInSlotFilling(id string) func(request *model.Request) (result bool) {
	return func(request *model.Request) (result bool) {
		return request.IsInSlotFillingAction(id)
	}
}

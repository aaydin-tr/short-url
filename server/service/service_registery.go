package service

import (
	"reflect"

	"go.uber.org/zap"
)

type Services struct {
	ShortURLService *URLService
}

func RegisterServices(services ...interface{}) *Services {

	newServices := &Services{}

	for _, service := range services {
		serviceType := getServiceType(service)
		switch serviceType {
		case "URLService":
			newServices.ShortURLService = service.(*URLService)
			zap.S().Info("URLService registered successfully")
		default:
			zap.S().Error("Unknown service: " + serviceType)
		}
	}
	return newServices
}

func getServiceType(service interface{}) string {
	valueOf := reflect.ValueOf(service)
	var name string
	if valueOf.Type().Kind() == reflect.Ptr {
		name = reflect.Indirect(valueOf).Type().Name()
	} else {
		name = valueOf.Type().Name()
	}
	return name
}

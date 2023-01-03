package formats

import (
	"fmt"
	"github.com/cocov-ci/coverage-reporter/meta"
	"go.uber.org/zap"
	"reflect"
)

type BaseFormat interface {
	Wants(diffs map[string]string) *string
	Name() string
	Parse(path string) (map[string]string, error)
	SetMeta(meta *meta.Metadata)
}

var covers = []reflect.Type{
	reflect.TypeOf(&GoCov{}),
	reflect.TypeOf(&SimpleCov{}),
	reflect.TypeOf(&Lcov{}),
}

func makeCover(kind reflect.Type) BaseFormat {
	return reflect.New(kind.Elem()).Interface().(BaseFormat)
}

func AutoFind(diffs map[string]string, meta *meta.Metadata) (map[string]string, error) {
	log := zap.L()
	log.Debug("AutoFind called")
	for _, kind := range covers {
		handler := makeCover(kind)
		log := log.With(zap.String("handler", handler.Name()))
		log.Debug("Trying handler")
		handler.SetMeta(meta)
		which := handler.Wants(diffs)
		if which != nil {
			log.Debug("Handler indicates it may be able to handle report", zap.Stringp("report_path", which))
			if data, err := handler.Parse(*which); err == nil {
				log.Debug("Successful parse from handler. Returning it.")
				return data, nil
			} else {
				log.Debug("Handler failed to process report. Trying next.", zap.Stringp("report_path", which), zap.Error(err))
				continue
			}
		}
	}

	return nil, fmt.Errorf("could not auto-detect coverage data")
}

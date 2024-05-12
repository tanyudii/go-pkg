package go_validator

import (
	"context"
	"errors"
	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	goerr "pkg.tanyudii.me/go-pkg/go-err"
	gotex "pkg.tanyudii.me/go-pkg/go-tex"
	"reflect"
	"regexp"
	"strings"
)

type Service interface {
	Struct(val interface{}) error
	StructWithCtx(ctx context.Context, val interface{}) error
	StructWithLang(lang string, val interface{}) error
}

type service struct {
	cfg      *Config
	validate *validator.Validate
	uni      *ut.UniversalTranslator
}

func NewValidator(args ...ConfigFunc) Service {
	v := &service{cfg: generateConfig(args...)}
	v.init()
	return v
}

func (s *service) init() {
	s.initUni()

	s.validate = validator.New()
	s.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("label")
	})

	for name, fn := range s.cfg.mapRegisterTranslation {
		translator, _ := s.uni.GetTranslator(name)
		_ = fn(s.validate, translator)
	}
}

func (s *service) initUni() {
	if s.uni != nil {
		return
	}
	var translations []locales.Translator
	for _, trans := range s.cfg.mapLocalesTranslator {
		translations = append(translations, trans)
	}
	if len(translations) > 1 {
		s.uni = ut.New(translations[0], translations[1:]...)
	} else {
		s.uni = ut.New(translations[0])
	}
}

func (s *service) Struct(val interface{}) error {
	return s.StructWithLang(DefaultLocaleName, val)
}

func (s *service) StructWithCtx(ctx context.Context, val interface{}) error {
	return s.StructWithLang(gotex.GetAcceptLanguage(ctx, DefaultLocaleName), val)
}

func (s *service) StructWithLang(lang string, val interface{}) error {
	err := s.validate.Struct(val)
	if err == nil {
		return nil
	}
	fields := make(goerr.ErrorField)
	translatorFn := s.getTranslator(lang)
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		for _, e := range errs {
			realNs := strings.Split(e.StructNamespace(), ".")[1:]
			fields[s.getFieldName(val, realNs)] = e.Translate(translatorFn)
		}
	}
	return goerr.NewBadRequestErrorUsingFieldsOrNil(fields)
}

func (s *service) getTranslator(localeName string) ut.Translator {
	transFn, ok := s.uni.GetTranslator(localeName)
	if !ok {
		transFn, _ = s.uni.GetTranslator(DefaultLocaleName)
	}
	return transFn
}

func (s *service) getFieldName(val any, fields []string) string {
	valClone := val
	result := s.constructFieldName(valClone, fields)
	return result
}

func (s *service) constructFieldName(val any, fields []string) string {
	var fieldName, result string
	nFields := len(fields)
	for i := 0; i < nFields; i++ {
		field := fields[i]
		if isSlice, parent, pos := s.isSliceField(field); isSlice {
			fieldName, val = s.extractFieldName(val, parent)
			fieldName += "." + pos
		} else {
			fieldName, val = s.extractFieldName(val, field)
		}
		result += fieldName
		if nFields > i+1 {
			result += "."
		}
	}
	return result
}

func (s *service) isSliceField(input string) (bool, string, string) {
	rgx := regexp.MustCompile(`([^\[]+)\[(\d+)\]`)
	match := rgx.FindStringSubmatch(input)
	if len(match) == 3 {
		return true, match[1], match[2]
	}
	return false, input, ""
}

func (s *service) extractFieldName(i interface{}, original string) (string, any) {
	reflected := reflect.ValueOf(i)
	switch reflected.Kind() {
	case reflect.Ptr:
		reflected = reflected.Elem()
	case reflect.Struct:
		break
	case reflect.Slice:
		return s.extractFieldName(reflected.Index(0).Interface(), original)
	default:
		return strcase.ToLowerCamel(original), i
	}
	if field, ok := reflected.Type().FieldByName(original); ok {
		i = reflected.FieldByName(original).Interface()
		if tag := field.Tag.Get("field"); tag != "" {
			return tag, i
		}
	}
	return strcase.ToLowerCamel(original), i
}

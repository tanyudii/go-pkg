package go_validator

import (
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
)

var (
	DefaultLocaleName  = "en_US"
	DefaultLocaleTrans = en_US.New()
)

type Config struct {
	mapLocalesTranslator   map[string]locales.Translator
	mapRegisterTranslation map[string]RegisterTranslationFunc
}

type ConfigFunc func(c *Config)

type RegisterTranslationFunc func(v *validator.Validate, trans ut.Translator) (err error)

func generateConfig(args ...ConfigFunc) *Config {
	c := &Config{
		mapLocalesTranslator:   make(map[string]locales.Translator),
		mapRegisterTranslation: make(map[string]RegisterTranslationFunc),
	}
	RegisterLocaleTranslator(DefaultLocaleName, DefaultLocaleTrans, en.RegisterDefaultTranslations)(c)
	for i := range args {
		args[i](c)
	}
	return c
}

func RegisterLocaleTranslator(name string, trans locales.Translator, fn RegisterTranslationFunc) ConfigFunc {
	return func(c *Config) {
		c.mapLocalesTranslator[name] = trans
		c.mapRegisterTranslation[name] = fn
	}
}

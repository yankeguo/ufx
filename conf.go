package ufx

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
)

// Conf is the configuration type
type Conf map[string]any

// Bind binds the configuration to the given data structure, supports json tags
func (c Conf) Bind(data any, keys ...string) (err error) {
	m := map[string]any(c)

	for _, key := range keys {
		if v := m[key]; v == nil {
			m = map[string]any{}
			break
		} else {
			switch v := v.(type) {
			case map[string]any:
				m = v
			case Conf:
				m = v
			default:
				err = errors.New("ufx.Conf#Bind: invalid key: " + strings.Join(keys, "."))
				return
			}
		}
	}

	var buf []byte
	if buf, err = json.Marshal(m); err != nil {
		return
	}
	if err = json.Unmarshal(buf, data); err != nil {
		return
	}
	if err = defaults.Set(data); err != nil {
		return
	}
	if err = validator.New().Struct(data); err != nil {
		return err
	}
	return err
}

// ProvideEmptyConf provides an empty Conf
func ProvideEmptyConf() fx.Option {
	return fx.Provide(func() Conf { return Conf{} })
}

// ProvideConfFromYAMLFile provides a Conf from a YAML file
func ProvideConfFromYAMLFile(name string) fx.Option {
	return fx.Provide(func() (conf Conf, err error) {
		var buf []byte
		if buf, err = os.ReadFile(name); err != nil {
			return
		}
		if err = yaml.Unmarshal(buf, &conf); err != nil {
			return
		}
		return
	})
}

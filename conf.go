package ufx

import (
	"encoding/json"
	"errors"
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

// Conf is the configuration type
type Conf map[string]any

// Bind binds the configuration to the given data structure, supports json tags
func (c Conf) Bind(data interface{}, keys ...string) (err error) {
	m := map[string]any(c)

	for _, key := range keys {
		var ok bool
		if m, ok = m[key].(map[string]any); !ok {
			if m == nil {
				m = map[string]any{}
				break
			} else {
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

func LoadConf() (c Conf, err error) {
	var buf []byte
	if buf, err = os.ReadFile("config.yaml"); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	if err = yaml.Unmarshal(buf, &c); err != nil {
		return
	}
	return
}

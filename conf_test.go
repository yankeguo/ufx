package ufx

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

type TestParams struct {
	World string `json:"world"`
}

func TestConf_Bind(t *testing.T) {
	c := Conf(
		map[string]any{
			"hello": map[string]any{
				"world": "world",
			},
		},
	)

	var p TestParams
	require.NoError(t, c.Bind(&p, "hello"))

	require.Equal(t, "world", p.World)
}

func TestProvideConfFromYAMLFile(t *testing.T) {
	buf, err := yaml.Marshal(map[string]any{
		"hello": map[string]any{
			"world": "world",
		},
	})
	require.NoError(t, err)
	err = os.WriteFile("test.yaml", buf, 0644)
	require.NoError(t, err)
	defer os.RemoveAll("test.yaml")

	a := fx.New(
		ProvideConfFromYAMLFile("test.yaml"),
		fx.Invoke(
			func(conf Conf) {
				var p TestParams
				require.NoError(t, conf.Bind(&p, "hello"))
				require.Equal(t, "world", p.World)
			},
		),
	)
	require.NoError(t, a.Err())
}

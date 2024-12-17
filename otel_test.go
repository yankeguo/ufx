package ufx

import (
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"testing"
)

func TestSetupOTEL(t *testing.T) {
	require.NoError(t, SetupOTEL())
	_ = otel.GetTracerProvider()
}

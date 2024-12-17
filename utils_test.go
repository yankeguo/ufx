package ufx

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type M = map[string]any

func TestExtractStringSlice(t *testing.T) {
	m := M{}
	extractStringSlice(m, "b_", "hello", []string{"a"})
	require.Equal(t, M{
		"hello":         "a",
		"hello_array":   []string{"a"},
		"b_hello":       "a",
		"b_hello_array": []string{"a"},
	}, m)
}

func TestExtractRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "https://example.com/get?aaa=bbb", nil)

	fm := map[string][]*multipart.FileHeader{}

	m := map[string]any{}
	err := extractRequest(m, fm, req)
	require.NoError(t, err)
	require.Equal(
		t,
		map[string]interface{}{
			"aaa":               "bbb",
			"aaa_array":         []string{"bbb"},
			"header_host":       "example.com",
			"header_host_array": []string{"example.com"},
			"host":              "example.com",
			"host_array":        []string{"example.com"},
			"query_aaa":         "bbb",
			"query_aaa_array":   []string{"bbb"}},
		m,
	)

	req = httptest.NewRequest("POST", "https://example.com/post?aaa=bbb", bytes.NewReader([]byte(`{"hello":"world"}`)))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	m = map[string]any{}
	err = extractRequest(m, fm, req)
	require.NoError(t, err)
	require.Equal(
		t,
		map[string]interface{}{
			"aaa":                       "bbb",
			"aaa_array":                 []string{"bbb"},
			"content_type":              "application/json;charset=utf-8",
			"content_type_array":        []string{"application/json;charset=utf-8"},
			"header_content_type":       "application/json;charset=utf-8",
			"header_content_type_array": []string{"application/json;charset=utf-8"},
			"header_host":               "example.com",
			"header_host_array":         []string{"example.com"},
			"hello":                     "world",
			"host":                      "example.com",
			"host_array":                []string{"example.com"},
			"query_aaa":                 "bbb",
			"query_aaa_array":           []string{"bbb"}},
		m,
	)

	req = httptest.NewRequest("POST", "https://example.com/post?aaa=bbb", bytes.NewReader([]byte(`hello=world`)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	m = map[string]any{}
	err = extractRequest(m, fm, req)
	require.NoError(t, err)
	require.Equal(
		t,
		map[string]interface{}{
			"aaa":                       "bbb",
			"aaa_array":                 []string{"bbb"},
			"content_type":              "application/x-www-form-urlencoded;charset=utf-8",
			"content_type_array":        []string{"application/x-www-form-urlencoded;charset=utf-8"},
			"form_hello":                "world",
			"form_hello_array":          []string{"world"},
			"header_content_type":       "application/x-www-form-urlencoded;charset=utf-8",
			"header_content_type_array": []string{"application/x-www-form-urlencoded;charset=utf-8"},
			"header_host":               "example.com",
			"header_host_array":         []string{"example.com"},
			"hello":                     "world",
			"hello_array":               []string{"world"},
			"host":                      "example.com",
			"host_array":                []string{"example.com"},
			"query_aaa":                 "bbb",
			"query_aaa_array":           []string{"bbb"}},
		m,
	)

	req = httptest.NewRequest("POST", "https://example.com/post?aaa=bbb", bytes.NewReader([]byte(`hello=world`)))
	req.Header.Set("Content-Type", "text/plain;charset=utf-8")

	m = map[string]any{}
	err = extractRequest(m, fm, req)
	require.NoError(t, err)
	require.Equal(
		t,
		map[string]interface{}{
			"aaa":                       "bbb",
			"aaa_array":                 []string{"bbb"},
			"body":                      "hello=world",
			"content_type":              "text/plain;charset=utf-8",
			"content_type_array":        []string{"text/plain;charset=utf-8"},
			"header_content_type":       "text/plain;charset=utf-8",
			"header_content_type_array": []string{"text/plain;charset=utf-8"},
			"header_host":               "example.com",
			"header_host_array":         []string{"example.com"},
			"host":                      "example.com",
			"host_array":                []string{"example.com"},
			"query_aaa":                 "bbb",
			"query_aaa_array":           []string{"bbb"}},
		m,
	)

	req = httptest.NewRequest("POST", "https://example.com/post?aaa=bbb", bytes.NewReader([]byte(`hello=world`)))
	req.Header.Set("Content-Type", "application/x-custom")

	m = map[string]any{}
	err = extractRequest(m, fm, req)
	require.NoError(t, err)
	require.Equal(
		t,
		map[string]interface{}{
			"aaa":                       "bbb",
			"aaa_array":                 []string{"bbb"},
			"body":                      []uint8{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x3d, 0x77, 0x6f, 0x72, 0x6c, 0x64},
			"content_type":              "application/x-custom",
			"content_type_array":        []string{"application/x-custom"},
			"header_content_type":       "application/x-custom",
			"header_content_type_array": []string{"application/x-custom"},
			"header_host":               "example.com",
			"header_host_array":         []string{"example.com"},
			"host":                      "example.com",
			"host_array":                []string{"example.com"},
			"query_aaa":                 "bbb",
			"query_aaa_array":           []string{"bbb"}},
		m,
	)
}

package ufx

import (
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

func extractStringSlice(m map[string]any, pfx string, key string, vs []string) {
	var v any
	if len(vs) == 1 {
		v = vs[0]
	} else {
		v = vs
	}
	m[key] = v
	m[key+"_array"] = vs
	if pfx != "" {
		m[pfx+key] = v
		m[pfx+key+"_array"] = vs
	}
}

func extractRequest(m map[string]any, f map[string][]*multipart.FileHeader, req *http.Request) (err error) {
	// header
	for k, vs := range req.Header {
		k = strings.ToLower(strings.ReplaceAll(k, "-", "_"))
		extractStringSlice(m, "header_", k, vs)
	}

	// query
	for k, vs := range req.URL.Query() {
		extractStringSlice(m, "query_", k, vs)
	}

	// body
	var buf []byte

	contentType, _, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))

	if contentType != ContentTypeMultipart {
		if buf, err = io.ReadAll(req.Body); err != nil {
			return
		}

		if len(buf) == 0 {
			return
		}
	}

	switch contentType {
	case ContentTypeTextPlain:
		m["body"] = string(buf)
	case ContentTypeApplicationJSON:
		var j map[string]any
		if err = json.Unmarshal(buf, &j); err != nil {
			return
		}
		for k, v := range j {
			m[k] = v
		}
	case ContentTypeFormURLEncoded:
		var q url.Values
		if q, err = url.ParseQuery(string(buf)); err != nil {
			return
		}
		for k, vs := range q {
			extractStringSlice(m, "form_", k, vs)
		}
	case ContentTypeMultipart:
		if err = req.ParseMultipartForm(1024 * 1024 * 10); err != nil {
			return
		}
		for k, vs := range req.MultipartForm.Value {
			extractStringSlice(m, "form_", k, vs)
		}
		for k, v := range req.MultipartForm.File {
			f[k] = v
		}
	default:
		m["body"] = buf
		return
	}

	return
}

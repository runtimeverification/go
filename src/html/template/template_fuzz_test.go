// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package template_test

import (
	"bytes"
	"testing"

	"html/template"
)

func FuzzHTMLTemplateParse(f *testing.F) {
	f.Add([]byte("{{.}}"))
	f.Add([]byte("{{if .}}x{{end}}"))
	f.Add([]byte("{{range .}}{{.}}{{end}}"))
	f.Add([]byte("{{define \"t\"}}x{{end}}{{template \"t\" .}}"))
	f.Add([]byte("{{/* comment */}}"))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 1<<20 {
			t.Skip()
		}
		s := string(data)
		tmpl, err := template.New("fuzz").Parse(s)
		if err != nil {
			return
		}
		var buf bytes.Buffer
		_ = tmpl.Execute(&buf, map[string]any{})
	})
}

func TestFuzzHTMLTemplateParseAnchor(t *testing.T) {
	const anchor = "{{.}}"
	tmpl, err := template.New("anchor").Parse(anchor)
	if err != nil {
		t.Fatalf("parse anchor: %v", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, "ok"); err != nil {
		t.Fatalf("execute anchor: %v", err)
	}
	if got, want := buf.String(), "ok"; got != want {
		t.Fatalf("anchor output: got %q want %q", got, want)
	}
}

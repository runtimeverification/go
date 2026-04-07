// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package template_test

import (
	"html/template"
	"testing"
)

func TestFuzzHTMLTemplateParseAnchor(t *testing.T) {
	_, err := template.New("anchor").Parse("{{.}}")
	if err != nil {
		t.Fatal(err)
	}
}

func FuzzHTMLTemplateParse(f *testing.F) {
	f.Add([]byte("{{.}}"))
	f.Add([]byte(""))
	f.Add([]byte("{{if .X}}{{.Y}}{{end}}"))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 256*1024 {
			return
		}
		_, err := template.New("fuzz").Parse(string(data))
		if err != nil {
			return
		}
	})
}

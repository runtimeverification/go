// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cookiejar

import (
	"bytes"
	"net/http"
	"testing"
)

func TestFuzzCookieJarSetCookiesAnchor(t *testing.T) {
	u := mustParseURL("https://example.org/path")
	j := newTestJar()
	c, err := http.ParseSetCookie("a=b; Path=/")
	if err != nil {
		t.Fatal(err)
	}
	j.SetCookies(u, []*http.Cookie{c})
}

func FuzzCookieJarSetCookies(f *testing.F) {
	u := mustParseURL("https://example.org/path")
	f.Add([]byte("a=b; Path=/"))
	f.Add([]byte("session=xyz; Path=/; HttpOnly\nlang=en; Path=/"))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 8*1024 {
			return
		}
		jar := newTestJar()
		for _, line := range bytes.Split(data, []byte("\n")) {
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			c, err := http.ParseSetCookie(string(line))
			if err != nil {
				continue
			}
			jar.SetCookies(u, []*http.Cookie{c})
		}
	})
}

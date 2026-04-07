// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cookiejar

import (
	"bytes"
	"cmp"
	"net/http"
	"net/url"
	"slices"
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

// fuzzNormalizeCookies returns a sorted copy for (Name, Value, Quoted) comparison.
// jar.Cookies only populates those fields (see jar.go).
func fuzzNormalizeCookies(cs []*http.Cookie) []http.Cookie {
	out := make([]http.Cookie, len(cs))
	for i, c := range cs {
		if c != nil {
			out[i] = *c
		}
	}
	slices.SortFunc(out, func(a, b http.Cookie) int {
		if r := cmp.Compare(a.Name, b.Name); r != 0 {
			return r
		}
		if r := cmp.Compare(a.Value, b.Value); r != 0 {
			return r
		}
		if a.Quoted == b.Quoted {
			return 0
		}
		if !a.Quoted {
			return -1
		}
		return 1
	})
	return out
}

// fuzzAssertJarCookiesSemanticallyConsistent checks that Cookies is idempotent
// and returns only well-formed names — stronger than “no panic” alone.
func fuzzAssertJarCookiesSemanticallyConsistent(t *testing.T, jar *Jar, u *url.URL) {
	t.Helper()
	got1 := jar.Cookies(u)
	got2 := jar.Cookies(u)
	n1 := fuzzNormalizeCookies(got1)
	n2 := fuzzNormalizeCookies(got2)
	if len(n1) != len(n2) {
		t.Fatalf("Cookies length unstable between reads: %d vs %d", len(n1), len(n2))
	}
	for i := range n1 {
		a, b := n1[i], n2[i]
		if a.Name != b.Name || a.Value != b.Value || a.Quoted != b.Quoted {
			t.Fatalf("Cookies mismatch on re-read at %d: %+v vs %+v", i, a, b)
		}
	}
	for _, c := range got1 {
		if c.Name == "" {
			t.Fatalf("jar returned cookie with empty Name")
		}
	}
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
		applied := 0
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
			applied++
		}
		fuzzAssertJarCookiesSemanticallyConsistent(t, jar, u)
		// Trivial bound: jar cannot return more cookies than Set-Cookie lines we applied.
		if applied > 0 {
			if n := len(jar.Cookies(u)); n > applied {
				t.Fatalf("jar returned %d cookies but only %d Set-Cookie lines applied", n, applied)
			}
		}
	})
}

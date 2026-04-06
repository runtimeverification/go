// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestFuzzReadRequestAnchor(t *testing.T) {
	const raw = "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n"
	req, err := ReadRequest(bufio.NewReader(strings.NewReader(raw)))
	if err != nil {
		t.Fatal(err)
	}
	if req.Method != "GET" {
		t.Fatalf("Method = %q", req.Method)
	}
	if req.Body != nil {
		io.Copy(io.Discard, req.Body)
		req.Body.Close()
	}
}

func FuzzReadRequest(f *testing.F) {
	f.Add([]byte("GET / HTTP/1.1\r\nHost: x\r\n\r\n"))
	f.Add([]byte("GET http://x/ HTTP/1.1\r\nHost: x\r\n\r\n"))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 64*1024 {
			return
		}
		req, err := ReadRequest(bufio.NewReader(bytes.NewReader(data)))
		if err != nil {
			return
		}
		if req.Body != nil {
			_, _ = io.Copy(io.Discard, req.Body)
			req.Body.Close()
		}
	})
}

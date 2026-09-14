// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

//go:build linux && !android

package tailssh

import (
	"reflect"
	"testing"

	"github.com/godbus/dbus/v5"
)

func TestParseCreateSessionResponse(t *testing.T) {
	validBody := []any{
		"session-1",
		dbus.ObjectPath("/org/freedesktop/login1/session/_1"),
		"/run/user/1000",
		dbus.UnixFD(3),
		uint32(1000),
		"seat0",
		uint32(1),
		false,
	}
	want := createSessionResp{
		sessionID:   "session-1",
		objectPath:  "/org/freedesktop/login1/session/_1",
		runtimePath: "/run/user/1000",
		fifoFD:      3,
		uid:         1000,
		seatID:      "seat0",
		vtnr:        1,
	}

	got, err := parseCreateSessionResponse(validBody)
	if err != nil {
		t.Fatalf("parseCreateSessionResponse(valid) error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseCreateSessionResponse(valid) = %#v, want %#v", got, want)
	}

	tests := []struct {
		name string
		body []any
	}{
		{name: "empty"},
		{name: "short", body: validBody[:7]},
		{name: "long", body: append(append([]any(nil), validBody...), "extra")},
		{name: "wrong_type", body: append(append([]any(nil), validBody[:7]...), "false")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := parseCreateSessionResponse(tt.body); err == nil {
				t.Fatal("parseCreateSessionResponse() error = nil, want malformed response error")
			}
		})
	}
}

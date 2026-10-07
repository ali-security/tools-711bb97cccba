// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"
)

const implicitAllInterfaces = "implicitly binds all network interfaces"

// runServe runs the serve command with the given -listen address and returns
// its error. It fails the test if the command does not return promptly, which
// means it is serving.
func runServe(t *testing.T, address string) error {
	t.Helper()
	app := New("", nil)
	app.Serve.app = app
	app.Serve.Address = address
	errc := make(chan error, 1)
	go func() { errc <- app.Serve.Run(context.Background()) }()
	select {
	case err := <-errc:
		return err
	case <-time.After(30 * time.Second):
		t.Fatalf("-listen=%s did not return: gopls is serving on it", address)
		return nil
	}
}

// TestServeListenImplicitAllInterfaces checks that a -listen address with no
// host, which would implicitly bind all network interfaces, is rejected before
// any listener is created.
func TestServeListenImplicitAllInterfaces(t *testing.T) {
	for _, address := range []string{":0", ":37374"} {
		err := runServe(t, address)
		if err == nil || !strings.Contains(err.Error(), implicitAllInterfaces) {
			t.Errorf("-listen=%s: got error %v, want %q", address, err, implicitAllInterfaces)
		}
	}
}

// TestServeListenExplicitHost checks that a -listen address with an explicit
// host is not rejected. An invalid port makes the listener fail without
// serving.
func TestServeListenExplicitHost(t *testing.T) {
	for _, address := range []string{"localhost:-1", "0.0.0.0:-1"} {
		err := runServe(t, address)
		if err == nil {
			t.Errorf("-listen=%s: got no error, want invalid port", address)
		} else if strings.Contains(err.Error(), implicitAllInterfaces) {
			t.Errorf("-listen=%s: explicit host rejected: %v", address, err)
		}
	}
}

// TestServeNoPortFlag checks that gopls no longer has the -port flag, which
// listened on all network interfaces.
func TestServeNoPortFlag(t *testing.T) {
	var check func(typ reflect.Type, path string)
	check = func(typ reflect.Type, path string) {
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			if f.Tag.Get("flag") == "port" {
				t.Errorf("%s.%s defines the -port flag", path, f.Name)
			}
			if f.PkgPath == "" && f.Type.Kind() == reflect.Struct {
				check(f.Type, path+"."+f.Name)
			}
		}
	}
	check(reflect.TypeOf(Application{}), "Application")
}

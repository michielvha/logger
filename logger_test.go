package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

func records(t *testing.T, out string) []map[string]any {
	t.Helper()
	var recs []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line == "" {
			continue
		}
		var r map[string]any
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			t.Fatalf("not a JSON record: %q: %v", line, err)
		}
		recs = append(recs, r)
	}
	return recs
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()
	fn()
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestInitWithOptionsWritesToOutput(t *testing.T) {
	var buf bytes.Buffer
	stdout := captureStdout(t, func() {
		InitWithOptions(Options{Level: "info", Output: &buf})
		Info("started", "key", "value")
	})
	if stdout != "" {
		t.Fatalf("stdout = %q, want nothing", stdout)
	}
	recs := records(t, buf.String())
	if len(recs) != 1 || recs[0]["msg"] != "started" || recs[0]["key"] != "value" || recs[0]["level"] != "INFO" {
		t.Fatalf("records = %v", recs)
	}
}

func TestInitWithOptionsLevel(t *testing.T) {
	tests := []struct {
		name, level, env string
		want             []string
	}{
		{"explicit warn", "warn", "", []string{"w", "e"}},
		{"warning alias", "WARNING", "", []string{"w", "e"}},
		{"env debug", "", "debug", []string{"d", "i", "w", "e"}},
		{"explicit wins over env", "error", "debug", []string{"e"}},
		{"unknown is info", "loud", "", []string{"i", "w", "e"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("LOG_LEVEL", tt.env)
			var buf bytes.Buffer
			InitWithOptions(Options{Level: tt.level, Output: &buf})
			Debug("d")
			Info("i")
			Warn("w")
			Error("e")
			var got []string
			for _, r := range records(t, buf.String()) {
				got = append(got, r["msg"].(string))
			}
			if strings.Join(got, ",") != strings.Join(tt.want, ",") {
				t.Fatalf("messages = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitDefaultsToStdout(t *testing.T) {
	t.Setenv("LOG_LEVEL", "")
	out := captureStdout(t, func() {
		Init("")
		Info("hello")
		InitWithOptions(Options{})
		Info("again")
	})
	recs := records(t, out)
	if len(recs) != 2 || recs[0]["msg"] != "hello" || recs[1]["msg"] != "again" {
		t.Fatalf("records = %v", recs)
	}
}

package main

import (
	"bytes"
	"testing"
)

func TestRun(t *testing.T) {
	for _, tt := range []struct {
		name   string
		args   []string
		status int
		stdout string
		stderr string
	}{
		{
			name:   "historical code",
			args:   []string{"110103"},
			stdout: `{"province":"北京市","city":"","district":"崇文区"}` + "\n",
		},
		{
			name:   "unknown code",
			args:   []string{"999999"},
			stdout: "未知\n",
		},
		{
			name:   "invalid prefix",
			args:   []string{"12345x"},
			status: 2,
			stderr: "usage: cn-idcard-area <6 ASCII digits>\n",
		},
		{
			name:   "missing argument",
			status: 2,
			stderr: "usage: cn-idcard-area <6 ASCII digits>\n",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			status := run(tt.args, &stdout, &stderr)
			if status != tt.status || stdout.String() != tt.stdout || stderr.String() != tt.stderr {
				t.Fatalf("run(%q) = (%d, %q, %q), want (%d, %q, %q)",
					tt.args, status, stdout.String(), stderr.String(), tt.status, tt.stdout, tt.stderr)
			}
		})
	}
}

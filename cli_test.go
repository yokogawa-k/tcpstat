package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_global(t *testing.T) {
	tests := []struct {
		desc           string
		arg            string
		expectedStatus int
		expectedSubOut string
		expectedSubErr string
	}{
		{
			desc:           "undefined flag",
			arg:            "tcpstat --undefined",
			expectedStatus: exitCodeErr,
			expectedSubErr: "flag provided but not defined",
		},
		{
			desc:           "help",
			arg:            "tcpstat --help",
			expectedStatus: exitCodeErr,
			expectedSubErr: "Usage: tcpstat",
		},
		{
			desc:           "version",
			arg:            "tcpstat --version",
			expectedStatus: exitCodeOK,
			expectedSubErr: "tcpstat version",
		},
		{
			desc:           "all",
			arg:            "tcpstat --all",
			expectedStatus: exitCodeOK,
			expectedSubOut: "Local Address",
		},
		{
			desc:           "normal",
			arg:            "tcpstat",
			expectedStatus: exitCodeOK,
			expectedSubOut: "Local Address",
		},
		{
			desc:           "--numeric",
			arg:            "tcpstat -n",
			expectedStatus: exitCodeOK,
			expectedSubOut: "Local Address",
		},
		{
			desc:           "--json",
			arg:            "tcpstat --json",
			expectedStatus: exitCodeOK,
			expectedSubOut: "{\"proto\":",
		},
	}
	for _, tc := range tests {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		cli := &CLI{outStream: outStream, errStream: errStream}
		args := strings.Split(tc.arg, " ")

		status := cli.Run(args)
		if status != tc.expectedStatus {
			t.Errorf("desc: %q, status should be %v, not %v", tc.desc, tc.expectedStatus, status)
		}

		if !strings.Contains(outStream.String(), tc.expectedSubOut) {
			t.Errorf("desc: %q, subout should contain %q, got %q", tc.desc, tc.expectedSubOut, outStream.String())
		}
		if !strings.Contains(errStream.String(), tc.expectedSubErr) {
			t.Errorf("desc: %q, subout should contain %q, got %q", tc.desc, tc.expectedSubErr, errStream.String())
		}
	}
}

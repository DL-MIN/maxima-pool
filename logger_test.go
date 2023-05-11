/*******************************************************************************
 * Test: Logger
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-05-11
 ******************************************************************************/

package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"os/exec"
	"testing"
)

func TestLogger_Print(t *testing.T) {
	type fields struct {
		level int
	}
	type args struct {
		format string
		v      []any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"debug one string", fields{Debug}, args{"%s", []any{"TEST"}}, "[\x1b[0;37mDEBUG\x1b[0m] TEST\n"},
		{"debug two strings", fields{Debug}, args{"%s%s", []any{"TE", "ST"}}, "[\x1b[0;37mDEBUG\x1b[0m] TEST\n"},
		{"debug one string and integer", fields{Debug}, args{"%s%d", []any{"TEST", 123}}, "[\x1b[0;37mDEBUG\x1b[0m] TEST123\n"},
		{"info one string", fields{Info}, args{"%s", []any{"TEST"}}, "[\x1b[0;32mINFO\x1b[0m] TEST\n"},
		{"info two strings", fields{Info}, args{"%s%s", []any{"TE", "ST"}}, "[\x1b[0;32mINFO\x1b[0m] TEST\n"},
		{"info one string and integer", fields{Info}, args{"%s%d", []any{"TEST", 123}}, "[\x1b[0;32mINFO\x1b[0m] TEST123\n"},
		{"warn one string", fields{Warn}, args{"%s", []any{"TEST"}}, "[\x1b[0;33mWARNING\x1b[0m] TEST\n"},
		{"warn two strings", fields{Warn}, args{"%s%s", []any{"TE", "ST"}}, "[\x1b[0;33mWARNING\x1b[0m] TEST\n"},
		{"warn one string and integer", fields{Warn}, args{"%s%d", []any{"TEST", 123}}, "[\x1b[0;33mWARNING\x1b[0m] TEST123\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			l := NewLogger()
			l.loggerList[tt.fields.level].SetOutput(&buf)
			l.loggerList[tt.fields.level].SetFlags(0)

			switch tt.fields.level {
			case Debug:
				l.Debug(tt.args.v...)
			case Info:
				l.Info(tt.args.v...)
			case Warn:
				l.Warn(tt.args.v...)
			}

			assert.Equal(t, tt.want, buf.String())
			buf.Reset()

			switch tt.fields.level {
			case Debug:
				l.Debugf(tt.args.format, tt.args.v...)
			case Info:
				l.Infof(tt.args.format, tt.args.v...)
			case Warn:
				l.Warnf(tt.args.format, tt.args.v...)
			}

			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestLogger_Fatal(t *testing.T) {
	if os.Getenv("CRASHTEST") == "1" {
		l := NewLogger()
		l.loggerList[Fatal].SetFlags(0)
		l.Fatalf("%s%d", "TEST", 123)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestLogger_Fatal")
	cmd.Env = append(os.Environ(), "CRASHTEST=1")
	bufReader, _ := cmd.StderrPipe()
	err := cmd.Start()
	bufOut, _ := io.ReadAll(bufReader)
	err = cmd.Wait()

	e, ok := err.(*exec.ExitError)
	assert.Equal(t, true, ok)
	assert.Equal(t, false, e.Success())

	want := "[\x1b[0;31mFATAL\x1b[0m] TEST123\n"
	assert.EqualValues(t, want, string(bufOut))
}

func TestLogger_Level(t *testing.T) {
	type fields struct {
		level int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"debug level", fields{Debug}, Debug},
		{"info level", fields{Info}, Info},
		{"warn level", fields{Warn}, Warn},
		{"fatal level", fields{Fatal}, Fatal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				level: tt.fields.level,
			}
			if got := l.Level(); got != tt.want {
				t.Errorf("Level() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogger_SetLevel(t *testing.T) {
	type fields struct {
		level int
	}
	type args struct {
		level int
	}
	tests := []struct {
		name  string
		args  args
		wants fields
	}{
		{"debug level", args{Debug}, fields{Debug}},
		{"info level", args{Info}, fields{Info}},
		{"warn level", args{Warn}, fields{Warn}},
		{"fatal level", args{Fatal}, fields{Fatal}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{}
			l.SetLevel(tt.args.level)
		})
		assert.Equal(t, tt.wants.level, tt.args.level)
	}
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name string
		want *Logger
	}{
		{name: "valid logger", want: &Logger{level: 0, loggerList: [4]*log.Logger{log.Default(), log.Default(), log.Default()}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLogger()
			assert.Equal(t, tt.want.level, got.level)
			for i := range got.loggerList {
				assert.IsType(t, tt.want.loggerList[i], got.loggerList[i])
			}
		})
	}
}

func TestLogger_logWithLevel(t *testing.T) {
	type args struct {
		level  int
		format *string
		v      []any
	}
	tests := []struct {
		name     string
		args     args
		setLevel int
		want     string
	}{
		{"debug level", args{Debug, nil, []any{"TEST"}}, Debug, "[\x1b[0;37mDEBUG\x1b[0m] TEST\n"},
		{"info level", args{Info, nil, []any{"TEST"}}, Info, "[\x1b[0;32mINFO\x1b[0m] TEST\n"},
		{"warn level", args{Warn, nil, []any{"TEST"}}, Warn, "[\x1b[0;33mWARNING\x1b[0m] TEST\n"},
		{"debug at warn level", args{Debug, nil, []any{"TEST"}}, Warn, ""},
		{"info at warn level", args{Info, nil, []any{"TEST"}}, Warn, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			l := NewLogger()
			l.SetLevel(tt.setLevel)
			l.loggerList[tt.setLevel].SetOutput(&buf)
			l.loggerList[tt.setLevel].SetFlags(0)

			l.logWithLevel(tt.args.level, tt.args.format, tt.args.v...)

			assert.Equal(t, tt.want, buf.String())
		})
	}
}

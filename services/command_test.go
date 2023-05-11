/*******************************************************************************
 * Test: Service: command
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-05-11
 ******************************************************************************/

package services

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"syscall"
	"testing"
	"time"
)

func Test_commandCreateWorkspace(t *testing.T) {
	err := os.MkdirAll("/tmp/Test_commandCreateWorkspace", 0755)
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll("/tmp/Test_commandCreateWorkspace")
	}()

	type args struct {
		uid int64
		gid int64
	}
	tests := []struct {
		name         string
		args         args
		setWorkspace string
		wantUid      int
		wantGid      int
		wantErr      error
	}{
		{"valid permissions", args{
			uid: -1,
			gid: -1,
		}, "/tmp/Test_commandCreateWorkspace", os.Getuid(), os.Getgid(), nil},
		{"invalid permissions", args{
			uid: 1001,
			gid: 1001,
		}, "/tmp/Test_commandCreateWorkspace", 1001, 1001, &fs.PathError{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("storage.workspace", tt.setWorkspace)
			gotWorkspace, gotErr, gotClean := commandCreateWorkspace(tt.args.uid, tt.args.gid)

			if tt.wantErr != nil {
				assert.Error(t, gotErr)
				assert.IsType(t, tt.wantErr, gotErr)
			} else {
				assert.NoError(t, gotErr)
				assert.DirExists(t, gotWorkspace)

				info, err := os.Stat(gotWorkspace)
				require.NoError(t, err)
				stat := info.Sys().(*syscall.Stat_t)

				assert.EqualValues(t, tt.wantUid, stat.Uid)
				assert.EqualValues(t, tt.wantGid, stat.Gid)
			}

			gotClean()
			assert.NoDirExists(t, gotWorkspace)
		})
	}
}

func Test_commandGetUser(t *testing.T) {
	tests := []struct {
		name    string
		user    string
		wantUid int64
		wantGid int64
		wantErr bool
	}{
		{"valid user", "root", 0, 0, false},
		{"invalid user", "1234", 0, 0, true},
		{"invalid empty user", "", 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("job.user", tt.user)

			gotUid, gotGid, gotErr := commandGetUser()
			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, tt.wantUid, gotUid)
				assert.Equal(t, tt.wantGid, gotGid)
			}
		})
	}
}

func TestCommandCreate(t *testing.T) {
	err := os.MkdirAll("/tmp/TestCommandCreate", 0755)
	require.NoError(t, err)
	viper.Set("storage.workspace", "/tmp/TestCommandCreate")
	defer func() {
		_ = os.RemoveAll("/tmp/TestCommandCreate")
	}()

	type args struct {
		user    string
		timeout time.Duration
		stdIn   string
		command string
		args    []string
	}
	tests := []struct {
		name       string
		args       args
		wantStdOut []byte
		wantStdErr []byte
		wantErr    bool
	}{
		{
			name:       "valid command",
			args:       args{"", time.Second, "", "echo", []string{"-n", "TEST"}},
			wantStdOut: []byte("TEST"),
			wantStdErr: []byte{},
		},
		{
			name:    "invalid command",
			args:    args{"", time.Second, "", "INVALID", nil},
			wantErr: true,
		},
		{
			name:       "deadline exceeded",
			args:       args{"", 3 * time.Second, "", "sleep", []string{"10"}},
			wantStdOut: []byte{},
			wantStdErr: []byte{},
			wantErr:    true,
		},
		{
			name:    "invalid permission",
			args:    args{"root", time.Second, "", "pwd", nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if len(tt.args.user) > 0 {
			viper.Set("job.user", tt.args.user)
		} else {
			viper.Set("job.user", nil)
		}

		t.Run(tt.name, func(t *testing.T) {
			gotStdOut, gotStdErr, _, _, gotErr := CommandCreate(tt.args.timeout, tt.args.stdIn, tt.args.command, tt.args.args...)
			assert.Equal(t, tt.wantErr, gotErr != nil)
			assert.Equal(t, tt.wantStdOut, gotStdOut)
			assert.Equal(t, tt.wantStdErr, gotStdErr)
		})
	}
}

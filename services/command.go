/*******************************************************************************
 * Service: command
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-05-11
 ******************************************************************************/

package services

import (
	"context"
	"github.com/spf13/viper"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

func CommandCreate(timeout time.Duration, stdIn string, command string, args ...string) (stdOut []byte, stdErr []byte, workspace string, clean func(), err error) {
	uid, gid, err := commandGetUser()
	if err != nil {
		return
	}

	workspace, err, clean = commandCreateWorkspace(uid, gid)
	if err != nil {
		return
	}

	stdOut, stdErr, err = commandRun(timeout, uid, gid, workspace, stdIn, command, args...)
	return
}

func commandRun(timeout time.Duration, uid int64, gid int64, workspace string, stdIn string, command string, args ...string) (stdOut []byte, stdErr []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create command context
	cmdCtx := exec.CommandContext(ctx, command, args...)
	cmdCtx.Dir = workspace

	// User credentials
	if uid >= 0 {
		cmdCtx.SysProcAttr = &syscall.SysProcAttr{}
		cmdCtx.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	}

	// I/O setup
	stdInWriter, err := cmdCtx.StdinPipe()
	if err != nil {
		return
	}

	stdOutReader, err := cmdCtx.StdoutPipe()
	if err != nil {
		return
	}

	stdErrReader, err := cmdCtx.StderrPipe()
	if err != nil {
		return
	}

	// Start command
	if err = cmdCtx.Start(); err != nil {
		return
	}

	// I/O interaction
	if _, err = stdInWriter.Write([]byte(stdIn)); err != nil {
		return
	}
	if err = stdInWriter.Close(); err != nil {
		return
	}

	stdOut, err = io.ReadAll(stdOutReader)
	if err != nil {
		return
	}

	stdErr, err = io.ReadAll(stdErrReader)
	if err != nil {
		return
	}

	// Process handling
	errCmd := cmdCtx.Wait()
	if err = ctx.Err(); err != nil {
		return
	}
	err = errCmd

	return
}

func commandGetUser() (uid int64, gid int64, err error) {
	if !viper.IsSet("job.user") {
		return -1, -1, nil
	}

	localUser, err := user.Lookup(viper.GetString("job.user"))
	if err != nil {
		return
	}

	uid, err = strconv.ParseInt(localUser.Uid, 10, 32)
	gid, err = strconv.ParseInt(localUser.Gid, 10, 32)

	return
}

func commandCreateWorkspace(uid int64, gid int64) (workspace string, err error, clean func()) {
	clean = func() {}
	if workspace, err = os.MkdirTemp(viper.GetString("storage.workspace"), "maxima-"); err != nil {
		return
	}
	clean = func() {
		_ = os.RemoveAll(workspace)
	}

	if uid >= 0 {
		err = os.Chown(workspace, int(uid), int(gid))
	}

	return
}

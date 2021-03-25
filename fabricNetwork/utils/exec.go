package utils

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
)

// 暂时用不上
//ExecuteCommand example ls -al ==> ls al
func ExecuteCommand(name string, opts ...string) (stdoutBuf, stderrBuf bytes.Buffer, err error) {
	cmd := exec.Command(name, opts...)
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err = cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}
	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()
	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()
	err = cmd.Wait()
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
	return stdoutBuf, stderrBuf, err
}

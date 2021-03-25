package utils

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"

	scp "github.com/bramvdbogaerde/go-scp"

	"golang.org/x/crypto/ssh"
)

var nilBytes = *bytes.NewBufferString("")
var shellCmdReplace = "/bin/bash %s"

//Dial connect remote machine
func Dial(user, password, ipPort string) (*ssh.Client, error) {
	cfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	return ssh.Dial("tcp", ipPort, cfg)
}

//RunCommand exec remote cmd and shell
func RunCommand(client *ssh.Client, cmd string, isShell bool) (bytes.Buffer, error) {
	session, err := client.NewSession()
	defer session.Close()
	if err != nil {
		return nilBytes, err
	}
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf
	// the shell of the remote machine
	if isShell {
		cmd = fmt.Sprintf(shellCmdReplace, cmd)
	}
	if err := session.Run(cmd); err != nil {
		return nilBytes, err
	}
	fmt.Println(stderrBuf.String())
	// if err := session.Wait(); err != nil {
	// 	return nilBytes, err
	// }
	return stdoutBuf, nil
}

//Scp local file to remote machine ;dPath 目标 需要处理子文件夹
func Scp(client *ssh.Client, file io.Reader, size int64, mode os.FileMode, dPath, dName string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	var stdoutBuf, errStderr bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &errStderr
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		fmt.Fprintln(w, "C"+fmt.Sprintf("%04o", mode), size, dName)
		io.Copy(w, file)
		fmt.Fprint(w, "\x00")
		// content := "123456789\n"
		// fmt.Fprintln(w, "D0755", 0, "testdir/test") // mkdir
		// fmt.Fprintln(w, "C0644", len(content), "testfile1")
		// fmt.Fprint(w, content)
		// fmt.Fprint(w, "\x00") // transfer end with \x00
		// fmt.Fprintln(w, "C0644", len(content), "testfile2")
		// fmt.Fprint(w, content)
		// fmt.Fprint(w, "\x00")
	}()

	if err := session.Run(fmt.Sprintf("/usr/bin/scp -qt %s", dPath)); err != nil {
		fmt.Println("执行scp命令失败:", err)
		fmt.Println(string(stdoutBuf.Bytes()), string(errStderr.Bytes()))
		return err
	}
	// if err := session.Wait(); err != nil {
	// 	fmt.Println(err)
	// }
	return nil
}

//ScpUseThirdLib
func ScpUseThirdLib(client *ssh.Client, from io.Reader, dist string, mode os.FileMode) {
	newClient, err := scp.NewClientBySSH(client)
	defer newClient.Close()
	if err != nil {
		fmt.Println("Error creating new SSH session from existing connection", err)
	}
	err = newClient.CopyFile(from, dist, fmt.Sprintf("%04o", mode))
	if err != nil {
		fmt.Println(err)
	}
}

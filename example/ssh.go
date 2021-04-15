package main

import (
	"Data_Bank/fabric-deploy-tools/fabricNetwork/utils"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {
	client, err := dial("root", "LEIlei1985!@#", "8.131.95.158:22")
	defer client.Close()
	if err != nil {
		panic(err)
	}

	file, err := os.Open("/Users/leixw/Documents/go-work/src/Data_Bank/fabric-deploy-tools/README.md")
	if err != nil {
		panic(err)
	}
	info, _ := file.Stat()
	// os.FileMode
	fmt.Println(info.Name())
	fmt.Println(info.Mode())
	f2, err := os.Open("/Users/leixw/Documents/go-work/src/Data_Bank/fabric-deploy-tools/.")
	f2Info, err := f2.Stat()
	if err != nil {
		fmt.Println(err)
	}
	// cmd.ScpUseThirdLib(client, file, "./test.md", info.Mode())
	err = utils.Scp(client, file, info.Size(), info.Mode(), "/root", info.Name())
	if err != nil {
		fmt.Println(err)
	}

	err = utils.Scp(client, f2, f2Info.Size(), f2Info.Mode(), "/root", f2Info.Name())
	if err != nil {
		fmt.Println(err)
	}
	// defer client.Close()
	// session, err := client.NewSession()
	// defer session.Close()
	// if err != nil {
	// 	panic(err)
	// }
	// println("ok")
	// var stdoutBuf bytes.Buffer
	// session.Stdout = &stdoutBuf

	// if err := session.Run("/bin/bash ./test.sh"); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(stdoutBuf.Bytes()))
	// stdout, stderr, err := cmd.ExecuteCommand("test.sh")
	// fmt.Println(string(stdout.Bytes()), string(stderr.Bytes()), err)
}

//dial
func dial(user, password, ipPort string) (*ssh.Client, error) {
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

//faultTolerant 兼容
func faultTolerant(session *ssh.Session) (err error) {
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm", 25, 80, modes)
	return err
}

//scp 定制 基于run
func scp(client *ssh.Client, file io.Reader, size int64, path string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	// session.StdinPipe()
	return nil
}

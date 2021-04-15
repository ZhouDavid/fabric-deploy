package main

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {
	client, err := dial("root", "LEIlei1985!@#", "8.131.95.158:22")
	if err != nil {
		panic(err)
	}
	sftpClient, err := sftp.NewClient(client)
	defer sftpClient.Close()
	println("...")
	srcFile, _ := os.Open("/Users/leixw/Documents/go-work/src/Data_Bank/fabric-deploy-tools/example/test.yaml")
	dstFile, _ := sftpClient.Create("/opt/test")
	defer func() {
		_ = dstFile.Close()
		srcFile.Close()
	}()
	buf := make([]byte, 1024)
	for {
		n, err := srcFile.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatalln("error occurred:", err)
			} else {
				break
			}
		}
		_, _ = dstFile.Write(buf[:n])
	}
}
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

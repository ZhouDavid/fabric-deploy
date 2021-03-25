package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/pkg/sftp"
)

//Connect sftp
func Connect(user, password, ipPort string) (*sftp.Client, error) {
	var sftpClient *sftp.Client
	sshClient, err := Dial(user, password, ipPort)

	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

//UploadFile sftp
func UploadFile(sftpClient *sftp.Client, localFilePath string, remotePath string) {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("os.Open error : ", localFilePath)
		log.Fatal(err)
		return
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(localFilePath)

	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		fmt.Println("sftpClient.Create error : ", path.Join(remotePath, remoteFileName))
		return
	}
	defer dstFile.Close()

	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("ReadAll error : ", localFilePath)
		return
	}
	if _, err := dstFile.Write(ff); err != nil {
		fmt.Println("sftp copy error", err)
	}
	stat, _ := srcFile.Stat()
	sftpClient.Chmod(path.Join(remotePath, remoteFileName), stat.Mode())
	// fmt.Println(localFilePath + "  copy file to remote server finished!")
}

//UploadDirectory sftp
func UploadDirectory(sftpClient *sftp.Client, localPath string, remotePath string) {
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		log.Fatal("read dir list fail ", err)
	}

	sftpClient.MkdirAll(remotePath)

	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		if backupDir.IsDir() {
			sftpClient.Mkdir(remoteFilePath)
			UploadDirectory(sftpClient, localFilePath, remoteFilePath)
		} else {
			UploadFile(sftpClient, path.Join(localPath, backupDir.Name()), remotePath)
		}
	}
	// fmt.Println(localPath + "  copy directory to remote server finished!")
}

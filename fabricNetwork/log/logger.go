package log

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
)

//InitLogger logger
func InitLogger(filename, filePath, level string) {
	fmt.Println(filename, filePath, level)
	logFilePath := filePath
	logFileName := filename

	// 日志文件
	fileName := path.Join(logFilePath, logFileName)
	// 写入文件
	src, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		fmt.Println("err", err)
	}

	// 实例化
	// logger := logrus.New()

	// 设置输出
	logrus.SetOutput(src)

	l, _ := logrus.ParseLevel(level)
	// 设置日志级别
	logrus.SetLevel(l)

	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d.log",

		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),

		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	fmt.Println(err)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// 新增 Hook
	logrus.AddHook(lfHook)
	// return logger
}

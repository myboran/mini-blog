package logging

import (
	"fmt"
	"gin-blog/pkg/setting"
	"time"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

func getLogFileFullPath() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
		)
}

//func openLogFile(filePath string) *os.File {
//	_, err := os.Stat(filePath)
//	switch {
//		case os.IsNotExist(err):
//			mkDir(getLogFilePath())
//		case os.IsPermission(err):
//			log.Fatalf("Permission :%v", err)
//	}
//	handle, err := os.OpenFile(filePath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
//	if err != nil {
//		log.Fatalf("Fail to OpenFile :%v", err)
//	}
//	return handle
//}
//
//func mkDir(string2 string) {
//	dir, _ := os.Getwd()
//	err := os.MkdirAll(dir + "/" + string2, os.ModePerm)
//	if err != nil {
//		panic(err)
//	}
//}
package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Log 全局logger
var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.Out = os.Stdout
	Log.Formatter = &logrus.JSONFormatter{}
	Log.AddHook(&DefaultFiledsHook{})
}

// DefaultFiledsHook 将应用默认信息添加为Fields的Hook
type DefaultFiledsHook struct {
}

// Levels 哪些 level 触发 Fire
func (hook *DefaultFiledsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire 触发时执行的操作
func (hook *DefaultFiledsHook) Fire(entry *logrus.Entry) error {
	entry.Data["serviceName"] = "archivist"
	entry.Data["serviceInstance"] = ""
	return nil
}

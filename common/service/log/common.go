package serviceLog

import (
	"github.com/sirupsen/logrus"
	// "gorm.io/gorm"
	// repo "go-web-app/common/repository"
)

// Constant values
const ()

var (
	DebugMode = false
)

// LevelThreshold - Returns every logging level above and including the given parameter.
func LevelThreshold(l logrus.Level) []logrus.Level {
	for i := range logrus.AllLevels {
		if logrus.AllLevels[i] == l {
			return logrus.AllLevels[:i+1]
		}
	}
	return []logrus.Level{}
}

package config

import (
	"os"
	"os/user"
	"path"

	utilslog "github.com/chengyumeng/khadijah/pkg/utils/log"
)

var Version string = ""

var GlobalOption *Option = new(Option)

var BaseURL string = "http://127.0.0.1:4200"

var UserConfigDir string = ""
var ConfigFile string = "config"

func init() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	UserConfigDir = path.Join(user.HomeDir, ".khadijah")
	if _, err := os.Stat(UserConfigDir); os.IsNotExist(err) {
		err = os.Mkdir(UserConfigDir, 0744)
		if err != nil {
			utilslog.AppLogger.Error(err)
		}
	}
}
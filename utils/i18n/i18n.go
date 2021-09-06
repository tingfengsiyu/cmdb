package i18n

import (
	"cmdb/utils"
	"path"

	"github.com/leonelquinteros/gotext"
)

func Initial() {
	localePath := path.Join(utils.RootPath, "locale")
	//gotext.Configure(localePath, "en_US", "koko")
	gotext.Configure(localePath, "zh_CN", "koko")

}

func T(s string) string {
	return gotext.Get(s)
}

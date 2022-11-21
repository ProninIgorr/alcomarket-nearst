package formatter

import (
	"regexp"
	"sync"
)

var formatterInstance *Formatter
var formatterOnce sync.Once

type Formatter struct {
	PhoneRegexp *regexp.Regexp
}

func GetFormatter() *Formatter {
	formatterOnce.Do(func() {
		formatterInstance = &Formatter{}
		phoneRegexp, _ := regexp.Compile(`\D`)
		formatterInstance.PhoneRegexp = phoneRegexp

	})

	return formatterInstance
}

func (f *Formatter) FormatPhone(phone string) string {
	processedString := f.PhoneRegexp.ReplaceAllString(phone, "")
	return processedString
}

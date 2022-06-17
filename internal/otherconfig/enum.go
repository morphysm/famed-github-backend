package otherconfig

import (
	"reflect"

	"github.com/phuslu/log"
)

type Enum int

const (
	APIServerListenAddress Enum = iota
	LogLevel
)

func EnumString() []string {
	return []string{
		"apiserver.listenaddress", // APIServerListenAddress
		"log.level",               // LogLevel
	}
}

func EnumsDefault() []any {
	return []any{
		":1337",       // APIServerListenAddress
		log.InfoLevel, // LogLevel
	}
}

func (e Enum) String() string {
	return EnumString()[e]
}

func (e Enum) Default() any {
	return EnumsDefault()[e]
}

func (e Enum) Type() reflect.Kind {
	return reflect.TypeOf(EnumsDefault()[e]).Kind()
}

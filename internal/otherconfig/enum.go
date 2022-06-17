package otherconfig

import (
	"reflect"

	"github.com/phuslu/log"
)

type Enum int

const (
	ServerListenAddress Enum = iota
	LogLevel
)

func EnumString() []string {
	return []string{
		"server.listenaddress", // ServerListenAddress
		"log.level",            // LogLevel
	}
}

func EnumsDefault() []any {
	return []any{
		":1337",       // ServerListenAddress
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

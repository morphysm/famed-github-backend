package otherconfig

import (
	"reflect"

	"github.com/knadh/koanf/parsers/json"

	"github.com/knadh/koanf/providers/rawbytes"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/rotisserie/eris"
)

func (c *Configuration) Set(path Enum, value interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	wantedKind := path.Type()
	valueKind := reflect.TypeOf(value).Kind()

	if wantedKind != valueKind {
		return eris.New(path.String() + " should be of type `" + wantedKind.String() + "` and not of type `" + valueKind.String())
	}

	err := c.koanf.Load(confmap.Provider(map[string]interface{}{
		path.String(): value,
	}, c.koanf.Delim()), nil)
	if err != nil {
		return eris.Wrap(err, "failed to set new value to "+path.String())
	}

	return nil
}

func (c *Configuration) SetRaw(raw []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.koanf.Load(rawbytes.Provider(raw), json.Parser())
	if err != nil {
		return eris.Wrap(err, "failed to load raw bytes")
	}

	return nil
}

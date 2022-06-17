package otherconfig

import "github.com/knadh/koanf/parsers/json"

func (c *Configuration) GetStringMap(path Enum) map[string]string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.koanf.StringMap(path.String())
}

func (c *Configuration) GetString(path Enum) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.koanf.String(path.String())
}

func (c *Configuration) GetBool(path Enum) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.koanf.Bool(path.String())
}

func (c *Configuration) GetInt(path Enum) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.koanf.Int(path.String())
}

func (c *Configuration) JSONMarshal() ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	b, _ := c.koanf.Marshal(json.Parser())

	return b, nil
}

package protoconf_go

import (
	"fmt"
)

//KVReader kv reader for configuration
type KVReader interface {
	GetValues(appName string, keys []string) map[string]*string
	WatchApp(appName string, callback func(k string, v string))
}

//ConfigurationReader configuration reader
type ConfigurationReader struct {
	reader KVReader
}

//NewConfigurationReader create new configuration reader
func NewConfigurationReader(r KVReader) *ConfigurationReader {
	if r == nil {
		panic(fmt.Errorf("kv reader can not be nil"))
	}
	result := &ConfigurationReader{reader: r}
	return result
}

//Config read value needed by the configuration object
func (p *ConfigurationReader) Config(data Configuration) bool {
	appName := data.GetApplicationName()
	keys := data.GetValidKeys()
	kv := p.reader.GetValues(appName, keys)
	for k, v := range kv {
		if v == nil {
			defValue := data.GetDefaultValue(k)
			v = &defValue
		}
		if v == nil {
			panic(fmt.Errorf("No value for %s is found", k))
		}
		if !data.SetValue(k, *v) {
			panic(fmt.Errorf("Invalid value %s for %s is found", *v, k))
		}
	}
	return true
}

//WatchKeys watch specified keys
func (p *ConfigurationReader) WatchKeys(data Configuration) {
	p.reader.WatchApp(data.GetApplicationName(), data.NotifyValueChange)
}

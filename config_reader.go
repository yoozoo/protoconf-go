package protoconf

import (
	"fmt"
	"strings"
)

//NotifyInterface notification interface
type NotifyInterface interface {
	AddKey(k string, v string)
	UpdateKey(k string, v string)
	DeleteKey(k string)
}

//KVReader kv reader for configuration
type KVReader interface {
	GetValues(appName string) map[string]string
	WatchApp(appName string, callback NotifyInterface)
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

type configStruct struct {
	keys map[string]bool
	maps map[string]*configStruct
}

func getMapKey(k string, cfg *configStruct) (string, map[string]bool) {
	for prefix, sub := range cfg.maps {
		if strings.HasPrefix(k, prefix) {
			k = strings.TrimPrefix(k, prefix)
			s := strings.SplitN(k, "/", 2)

			// must be map of basic type
			if len(s) == 1 {
				if sub == nil && len(k) > 0 {
					return prefix + k, nil
				}
				return "", nil
			} else if sub != nil {
				// a item of the map exists
				if _, ok := sub.keys["/"+s[1]]; ok {
					return prefix + s[0], sub.keys
				}
				p, keylist := getMapKey("/"+s[1], sub)
				if len(p) > 0 {
					return prefix + s[0] + p, keylist
				}
			}
			return "", nil
		}
	}
	return "", nil
}

func createConfigStruct(keys []string) *configStruct {

	cfgs := &configStruct{keys: make(map[string]bool), maps: make(map[string]*configStruct)}

	for _, k := range keys {
		if strings.Contains(k, mapKeyPlaceHolder) {
			subs := strings.Split(k, mapKeyPlaceHolder)
			isBasicType := strings.HasSuffix(k, mapKeyPlaceHolder)
			if isBasicType {
				subs = subs[:len(subs)-1]
			}
			root := cfgs
			for i, s := range subs {
				if i == (len(subs) - 1) {
					if isBasicType {
						root.maps[s] = nil
					} else {
						root.keys[s] = true
					}
				} else {
					if _, ok := root.maps[s]; !ok {
						root.maps[s] = &configStruct{keys: make(map[string]bool), maps: make(map[string]*configStruct)}
					}
					root = root.maps[s]
					if root == nil {
						fmt.Println("Invalid config key")
						break
					}
				}
			}
		} else {
			cfgs.keys[k] = true
		}
	}
	return cfgs
}

//Config read value needed by the configuration object
func (p *ConfigurationReader) Config(data Configuration) error {
	appName := data.ApplicationName()
	keys := data.ValidKeys()
	kv := p.reader.GetValues(appName)

	cfgs := createConfigStruct(keys)
	// normal values
	for k := range cfgs.keys {
		value, ok := kv[k]
		if !ok {
			v := data.DefaultValue(k)
			if v == nil {
				return fmt.Errorf("No value for %s is found", k)
			}
			value = *v
		} else {
			delete(kv, k)
		}
		err := data.SetValue(k, value)
		if err != nil {
			return fmt.Errorf("Invalid value %s for %s is found : %s", value, k, err)
		}
	}
	// map values
	for k, v := range kv {
		prefix, keys := getMapKey(k, cfgs)
		if len(prefix) > 0 {
			if keys == nil {
				err := data.SetValue(k, v)
				if err != nil {
					return fmt.Errorf("Invalid value %s for %s is found : %s", v, k, err)
				}
			} else {
				// make sure the map object exists
				err := data.SetValue(k, v)
				if err != nil {
					return fmt.Errorf("Invalid value %s for %s is found : %s", v, k, err)
				}
				for subkey := range keys {
					key := prefix + subkey
					value, ok := kv[key]
					if !ok {
						v := data.DefaultValue(key)
						if v == nil {
							return fmt.Errorf("No value for %s is found", k)
						}
						value = *v
					} else {
						delete(kv, key)
					}
					err := data.SetValue(key, value)
					if err != nil {
						return fmt.Errorf("Invalid value %s for %s is found : %s", value, k, err)
					}
				}
			}
		}
	}
	return nil
}

type notifycationObject struct {
	data Configuration
	cfgs *configStruct
}

func (p *notifycationObject) AddKey(k string, v string) {
	p.UpdateKey(k, v)
}
func (p *notifycationObject) UpdateKey(k string, v string) {
	if _, ok := p.cfgs.keys[k]; !ok {
		prefix, _ := getMapKey(k, p.cfgs)
		if len(prefix) == 0 {
			return
		}
	}
	p.data.NotifyValueChange(k, v)
}
func (p *notifycationObject) DeleteKey(k string) {
	if _, ok := p.cfgs.keys[k]; !ok {
		prefix, _ := getMapKey(k, p.cfgs)
		if len(prefix) == 0 {
			return
		}
	}
	DeleteKey(p.data, k)
}

//WatchKeys watch specified keys
func (p *ConfigurationReader) WatchKeys(data Configuration) {
	p.reader.WatchApp(data.ApplicationName(), &notifycationObject{data: data, cfgs: createConfigStruct(data.ValidKeys())})
}

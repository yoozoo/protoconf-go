package main

import (
	"fmt"
	"time"

	protoconf "github.com/yoozoo/protoconf_go"
)

func main() {
	etcd := protoconf.NewEtcdReader("default", "root", "root", []string{"192.168.115.57:2379"})
	reader := protoconf.NewConfigurationReader(etcd)
	var t ConfigurationTest
	reader.Config(&t)
	reader.WatchKeys(&t)

	time.Sleep(1000 * time.Second)
}

// ConfigurationTest protoconf configuration object interface
type ConfigurationTest struct {
}

//GetApplicationName retrieve application name
func (p *ConfigurationTest) GetApplicationName() string {
	return "服务1"
}

//GetValidKeys retrieve all keys
func (p *ConfigurationTest) GetValidKeys() []string {
	return []string{"msg/name", "msg/def", "msg/id", "name"}
}

//SetValue set values inside the java config class
func (p *ConfigurationTest) SetValue(key string, value string) error {
	fmt.Printf("setting value %s for key %s\n", value, key)
	return nil
}

//GetDefaultValue get default values from the java config class
func (p *ConfigurationTest) GetDefaultValue(key string) *string {
	fmt.Println("Getting default value of key ", key)
	v := ""
	return &v
}

//NotifyValueChange add key change to the change list
func (p *ConfigurationTest) NotifyValueChange(key string, newValue string) {
	fmt.Printf("key %s value changed to %s\n", key, newValue)
}

package main

import (
	"fmt"
	"time"

	"github.com/yoozoo/protoconf_go"
)

func main() {
	etcd := protoconf.NewEtcdReader("default")
	etcd.SetUser("root", "root")
	etcd.SetEndpoints([]string{"192.168.115.57:2379"})

	reader := protoconf.NewConfigurationReader(etcd)
	var t ConfigurationTest
	reader.Config(&t)
	reader.WatchKeys(&t)

	time.Sleep(1000 * time.Second)
}

// ConfigurationTest protoconf configuration object interface
type ConfigurationTest struct {
}

//ApplicationName retrieve application name
func (p *ConfigurationTest) ApplicationName() string {
	return "服务1"
}

//ValidKeys retrieve all keys
func (p *ConfigurationTest) ValidKeys() []string {
	return []string{"msg/name", "msg/def", "msg/id", "name"}
}

//SetValue set values inside the java config class
func (p *ConfigurationTest) SetValue(key string, value string) error {
	fmt.Printf("setting value %s for key %s\n", value, key)
	return nil
}

//DefaultValue get default values from the java config class
func (p *ConfigurationTest) DefaultValue(key string) *string {
	fmt.Println("Getting default value of key ", key)
	v := ""
	return &v
}

//NotifyValueChange add key change to the change list
func (p *ConfigurationTest) NotifyValueChange(key string, newValue string) {
	fmt.Printf("key %s value changed to %s\n", key, newValue)
}

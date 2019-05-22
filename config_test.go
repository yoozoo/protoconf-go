package protoconf

import (
	"fmt"
	"reflect"
	"testing"
)

type testCfg struct {
	TestKey    []string
	TestValue  map[string]string
	Result     map[string]string
	Expected   map[string]string
	DefaultKey map[string]string
}

func (m *testCfg) GetValues(appName string) map[string]string {
	return m.TestValue
}
func (m *testCfg) WatchApp(appName string, callback NotifyInterface) {

}

func TestReadMap1(t *testing.T) {
	cfg := testCfg{}
	cfg.TestKey = []string{"test_int/MAP_ENTRY", "test_string/MAP_ENTRY", "test_obj/MAP_ENTRY/id", "test_obj/MAP_ENTRY/id2", "key1", "key2/abc"}
	cfg.TestValue = make(map[string]string)
	cfg.Result = make(map[string]string)
	cfg.Expected = make(map[string]string)
	cfg.DefaultKey = make(map[string]string)

	cfg.TestValue["test_int/abc"] = "123"
	cfg.TestValue["test_int/def"] = "123"
	cfg.TestValue["test_string/abc"] = "str"
	cfg.TestValue["key1"] = "key1"
	cfg.TestValue["key2/abc"] = "key2abc"
	cfg.TestValue["test_obj/abc/id"] = "id1"
	cfg.TestValue["test_obj/abc/id2"] = "id2"

	cfg.Expected["test_int/abc"] = "123"
	cfg.Expected["test_int/def"] = "123"
	cfg.Expected["test_string/abc"] = "str"
	cfg.Expected["key1"] = "key1"
	cfg.Expected["key2/abc"] = "key2abc"
	cfg.Expected["test_obj/abc/id"] = "id1"
	cfg.Expected["test_obj/abc/id2"] = "id2"

	reader := NewConfigurationReader(&cfg)
	reader.Config(&cfg)

	if !reflect.DeepEqual(cfg.Expected, cfg.Result) {
		t.Fail()
	}

}

func TestReadMultiLevelMap(t *testing.T) {
	cfg := testCfg{}
	cfg.TestKey = []string{"test/MAP_ENTRY/proj/MAP_ENTRY/setting", "test/MAP_ENTRY/proj/MAP_ENTRY/setting2"}
	cfg.TestValue = make(map[string]string)
	cfg.Result = make(map[string]string)
	cfg.Expected = make(map[string]string)
	cfg.DefaultKey = make(map[string]string)

	cfg.TestValue["test/MAP_ENTRY/proj/MAP_ENTRY/setting"] = "1"
	cfg.TestValue["test/MAP_ENTRY/proj/MAP_ENTRY/setting2"] = "2"
	cfg.TestValue["test_string/abc"] = "str"
	cfg.TestValue["key1"] = "key1"
	cfg.TestValue["key2/abc"] = "key2abc"
	cfg.TestValue["test_obj/abc/id"] = "id1"
	cfg.TestValue["test_obj/abc/id2"] = "id2"

	cfg.Expected["test/MAP_ENTRY/proj/MAP_ENTRY/setting"] = "1"
	cfg.Expected["test/MAP_ENTRY/proj/MAP_ENTRY/setting2"] = "2"

	reader := NewConfigurationReader(&cfg)
	reader.Config(&cfg)

	if !reflect.DeepEqual(cfg.Expected, cfg.Result) {
		t.Fail()
	}
}

func TestReadMultiLevelMap2(t *testing.T) {
	cfg := testCfg{}
	cfg.TestKey = []string{"test/MAP_ENTRY/proj/MAP_ENTRY/setting", "test/MAP_ENTRY/proj/MAP_ENTRY/setting2"}
	cfg.TestValue = make(map[string]string)
	cfg.Result = make(map[string]string)
	cfg.Expected = make(map[string]string)
	cfg.DefaultKey = make(map[string]string)

	cfg.TestValue["test/MAP_ENTRY/proj/MAP_ENTRY/setting"] = "1"
	cfg.TestValue["test/MAP_ENTRY/proj/MAP_ENTRY/setting2"] = "2"
	cfg.TestValue["test/MAP_ENTRY/proj/MAP_ENTRY2/setting"] = "1"
	cfg.TestValue["test/MAP_ENTRY/proj/MAP_ENTRY2/setting2"] = "2"

	cfg.Expected["test/MAP_ENTRY/proj/MAP_ENTRY/setting"] = "1"
	cfg.Expected["test/MAP_ENTRY/proj/MAP_ENTRY/setting2"] = "2"
	cfg.Expected["test/MAP_ENTRY/proj/MAP_ENTRY2/setting"] = "1"
	cfg.Expected["test/MAP_ENTRY/proj/MAP_ENTRY2/setting2"] = "2"

	reader := NewConfigurationReader(&cfg)
	reader.Config(&cfg)

	if !reflect.DeepEqual(cfg.Expected, cfg.Result) {
		t.Fail()
	}
}

func TestReadDefaultValue(t *testing.T) {
	cfg := testCfg{}
	cfg.TestKey = []string{"test/MAP_ENTRY/proj/MAP_ENTRY/setting", "test/MAP_ENTRY/proj/MAP_ENTRY/setting2"}
	cfg.TestValue = make(map[string]string)
	cfg.Result = make(map[string]string)
	cfg.Expected = make(map[string]string)
	cfg.DefaultKey = make(map[string]string)

	cfg.TestValue["test/MAP_ENTRY/proj/MAP_ENTRY/setting"] = "1"
	cfg.TestValue["test_string/abc"] = "str"
	cfg.TestValue["key1"] = "key1"
	cfg.TestValue["key2/abc"] = "key2abc"
	cfg.TestValue["test_obj/abc/id"] = "id1"
	cfg.TestValue["test_obj/abc/id2"] = "id2"

	cfg.Expected["test/MAP_ENTRY/proj/MAP_ENTRY/setting"] = "1"
	cfg.Expected["test/MAP_ENTRY/proj/MAP_ENTRY/setting2"] = "2"

	cfg.DefaultKey["test/MAP_ENTRY/proj/MAP_ENTRY/setting2"] = "2"
	reader := NewConfigurationReader(&cfg)
	reader.Config(&cfg)

	if !reflect.DeepEqual(cfg.Expected, cfg.Result) {
		t.Fail()
	}
}

func TestReadDefaultValue2(t *testing.T) {
	cfg := testCfg{}
	cfg.TestKey = []string{"test/MAP_ENTRY/proj/MAP_ENTRY/setting", "test/MAP_ENTRY/proj/MAP_ENTRY/setting2", "key1", "key2/abc"}
	cfg.TestValue = make(map[string]string)
	cfg.Result = make(map[string]string)
	cfg.Expected = make(map[string]string)
	cfg.DefaultKey = make(map[string]string)

	cfg.TestValue["test/MAP_ENTRY/proj/MAP_ENTRY/setting"] = "1"
	cfg.TestValue["test_string/abc"] = "str"

	cfg.TestValue["key2/abc"] = "key2abc"
	cfg.TestValue["test_obj/abc/id"] = "id1"
	cfg.TestValue["test_obj/abc/id2"] = "id2"

	cfg.Expected["test/MAP_ENTRY/proj/MAP_ENTRY/setting"] = "1"
	cfg.Expected["test/MAP_ENTRY/proj/MAP_ENTRY/setting2"] = "2"
	cfg.Expected["key1"] = "key1"
	cfg.Expected["key2/abc"] = "key2abc"

	cfg.DefaultKey["test/MAP_ENTRY/proj/MAP_ENTRY/setting2"] = "2"
	cfg.DefaultKey["key1"] = "key1"
	reader := NewConfigurationReader(&cfg)
	reader.Config(&cfg)

	if !reflect.DeepEqual(cfg.Expected, cfg.Result) {
		t.Fail()
	}
}

//ApplicationName retrieve application name
func (m *testCfg) ApplicationName() string {
	return "服务1"
}

//ValidKeys retrieve all keys
func (m *testCfg) ValidKeys() []string {
	return m.TestKey
}

//SetValue set values inside the java config class
func (m *testCfg) SetValue(key string, value string) error {
	fmt.Printf("setting value %s for key %s\n", value, key)
	m.Result[key] = value
	return nil
}

//DefaultValue get default values from the java config class
func (m *testCfg) DefaultValue(key string) *string {
	if v, ok := m.DefaultKey[key]; ok {
		return &v
	}
	return nil
}

//NotifyValueChange add key change to the change list
func (m *testCfg) NotifyValueChange(key string, newValue string) {
	fmt.Printf("key %s value changed to %s\n", key, newValue)
}

package protoconf

const (
	mapKeyPlaceHolder = "MAP_ENTRY"
)

// Configuration protoconf configuration object interface
type Configuration interface {
	//retrieve application name
	ApplicationName() string
	//retrieve all keys
	ValidKeys() []string
	//set values inside the java config class
	SetValue(key string, value string) error
	//get default values from the java config class
	DefaultValue(key string) *string
	// add key change to the change list
	NotifyValueChange(key string, newValue string)
}

// MapConfiguration map related func
type MapConfiguration interface {
	//delete the related object
	DeleteKey(key string)
}

// DeleteKey delete related map object
func DeleteKey(cfg Configuration, key string) {
	//check if has delete key interface
	if withOptions, ok := cfg.(MapConfiguration); ok {
		withOptions.DeleteKey(key)
	}
}

package protoconf_go

import (
	"fmt"
	"os"
	"strings"

	"github.com/coreos/etcd/clientv3"
)

const (
	defaultEndpoints = "192.168.115.57:2379"
	defaultUsername  = "root"
	defaultPassword  = "root"
	defaultEnv       = "default"

	endpointsEnvKey = "etcd_endpoints"
	userpassEnvKey  = "etcd_user"
	envEnvKey       = "etcd_envkey"
)

// EtcdReader etcd reader
type EtcdReader struct {
	client    *clientv3.Client
	endpoints []string
	user      string
	pass      string
	env       string
}

//NewEtcdReader create etcd reader
func NewEtcdReader(env string, user string, pass string, endpoints []string) *EtcdReader {
	if len(endpoints) == 0 {
		s := os.Getenv(endpointsEnvKey)
		if len(s) > 0 {
			endpoints = strings.Split(s, ",")
		} else {
			fmt.Println("invalid env value for endpoint:", os.Getenv(s))
		}
	}
	if len(user) == 0 || len(pass) == 0 {
		s := strings.Split(os.Getenv(userpassEnvKey), ":")
		if len(s) == 2 && len(s[0]) > 0 && len(s[1]) > 0 {
			user = s[0]
			pass = s[1]
		} else {
			fmt.Println("invalid env value for user/pass:", os.Getenv(userpassEnvKey))
		}
	}
	if len(env) == 0 {
		s := os.Getenv(etcd_envkey)
		if len(s) > 0 {
			env = s
		}
	}
	// todo add protoagent support

	if len(endpoints) == 0 {
		panic(fmt.Errorf("empty etcd endpoints"))
	}
	if len(user) == 0 {
		panic(fmt.Errorf("empty username"))
	}
	if len(pass) == 0 {
		panic(fmt.Errorf("empty password"))
	}
	if len(env) == 0 {
		panic(fmt.Errorf("empty envkey"))
	}

	result := &EtcdReader{
		endpoints: endpoints,
		user:      user,
		pass:      pass,
		env:       env,
	}

	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   endpoints,
		Username:    username,
		Password:    password,
	})
	if err != nil {
		panic(err)
	}

	result.client = cli
	return cli
}

//GetValues get values of the key list
func (p *EtcdReader) GetValues(appName string, keys []string) map[string]*string {

}

//WatchApp watch keys of the app
func (p *EtcdReader) WatchApp(appName string, callback func(k string, v string)) {

}

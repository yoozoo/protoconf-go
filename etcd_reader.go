package protoconf_go

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

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

	dialTimeout = 2 * time.Second
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
		s := os.Getenv(envEnvKey)
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
		Username:    user,
		Password:    pass,
	})

	if err != nil {
		panic(err)
	}
	result.client = cli
	return result
}

//GetValues get values of the key list
func (p *EtcdReader) GetValues(appName string, keys []string) map[string]*string {

	prefix := "/" + p.env + "/" + appName + "/"
	result := make(map[string]*string)

	txn := p.client.Txn(context.TODO())
	for _, k := range keys {
		txn = txn.Then(clientv3.OpGet(prefix + k))
		result[k] = nil
	}

	txnResp, err := txn.Commit()

	if err != nil {
		fmt.Println("error to retrieve config values: ", err)
		return result
	}

	if !txnResp.Succeeded {
		fmt.Println("Failed to retrieve config values")

	}
	resp := txnResp.OpResponse().Get()
	for _, kv := range resp.Kvs {
		v := string(kv.Value)
		result[strings.TrimPrefix(prefix, string(kv.Key))] = &v
	}

	return result

}

//WatchApp watch keys of the app
func (p *EtcdReader) WatchApp(appName string, callback func(k string, v string)) {
	go p.watchingApp(appName, callback)
}
func (p *EtcdReader) watchingApp(appName string, callback func(k string, v string)) {
	prefix := "/" + p.env + "/" + appName + "/"
RESTART_POINT:
	for {
		ch := p.client.Watch(context.Background(), prefix, clientv3.WithPrefix())

		for {
			select {
			case s, ok := <-ch:
				if ok {
					if len(s.Events) > 0 {
						cache := make(map[string]string)
						for _, e := range s.Events {
							cache[strings.TrimPrefix(string(e.Kv.Key), prefix)] = string(e.Kv.Value)
						}
						for k, v := range cache {
							callback(k, v)
						}
					}
					if s.Canceled {
						err := s.Err()
						if err != nil {
							fmt.Println("error :", err)
						}
						break RESTART_POINT
					}
				} else {
					fmt.Println("channel closed.")
					break RESTART_POINT
				}
			}
		}

	}
}

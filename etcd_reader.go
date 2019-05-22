package protoconf

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/yoozoo/protoconf_go/agentApplicationService"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	defaultEnv = "default"

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
	appToken  string
}

//SetUser set etcd user/pass
func (p *EtcdReader) SetUser(user string, pass string) {
	p.user = user
	p.pass = pass
}

//SetEndpoints set etcd endpionts
func (p *EtcdReader) SetEndpoints(endpoints []string) {
	p.endpoints = endpoints
}

//SetToken set proto agent token
func (p *EtcdReader) SetToken(token string) {
	p.appToken = token
}

//NewEtcdReader create etcd reader
func NewEtcdReader(env string) *EtcdReader {
	s := os.Getenv(endpointsEnvKey)
	var result EtcdReader
	if len(s) > 0 {
		result.endpoints = strings.Split(s, ",")
	}
	userpass := strings.Split(os.Getenv(userpassEnvKey), ":")
	if len(userpass) == 2 && len(userpass[0]) > 0 && len(userpass[1]) > 0 {
		result.user = userpass[0]
		result.pass = userpass[1]
	}
	result.env = env
	if len(env) == 0 {
		s = os.Getenv(envEnvKey)
		if len(s) > 0 {
			result.env = s
		} else {
			result.env = defaultEnv
		}
	}
	return &result
}
func (p *EtcdReader) getClient(appName string) *clientv3.Client {
	// todo add protoagent support
	if p.client != nil {
		return p.client
	}

	if len(p.endpoints) == 0 || len(p.user) == 0 || len(p.pass) == 0 {
		p.getSettingFromAgent(appName)
	}

	if len(p.endpoints) == 0 {
		panic(fmt.Errorf("empty etcd endpoints"))
	}
	if len(p.user) == 0 {
		panic(fmt.Errorf("empty username"))
	}
	if len(p.pass) == 0 {
		panic(fmt.Errorf("empty password"))
	}

	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   p.endpoints,
		Username:    p.user,
		Password:    p.pass,
	})

	if err != nil {
		panic(err)
	}
	p.client = cli
	return cli
}

//GetValues get values of the key list
func (p *EtcdReader) GetValues(appName string) map[string]string {

	prefix := "/" + p.env + "/" + appName + "/"
	result := make(map[string]string)

	cli := p.getClient(appName)
	if cli == nil {
		fmt.Println("Failed to get etcd client connection")
		return result
	}

	resp, err := cli.Get(context.TODO(), prefix, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("error to retrieve config values: ", err)
		return result
	}

	for _, kv := range resp.Kvs {
		v := string(kv.Value)
		key := strings.TrimPrefix(string(kv.Key), prefix)
		result[key] = v
	}

	return result

}

//WatchApp watch keys of the app
func (p *EtcdReader) WatchApp(appName string, callback NotifyInterface) {
	go p.watchingApp(appName, callback)
}
func (p *EtcdReader) watchingApp(appName string, callback NotifyInterface) {
	prefix := "/" + p.env + "/" + appName + "/"
RESTART_POINT:
	for {
		ch := p.client.Watch(context.Background(), prefix, clientv3.WithPrefix())

		for {
			select {
			case s, ok := <-ch:
				if ok {
					if len(s.Events) > 0 {
						cache := make(map[string]*clientv3.Event)
						for _, e := range s.Events {
							cache[strings.TrimPrefix(string(e.Kv.Key), prefix)] = e
						}
						for k, e := range cache {
							if e.Type == mvccpb.DELETE {
								callback.DeleteKey(k)
							} else if e.IsCreate() {
								callback.AddKey(k, string(e.Kv.Value))
							} else {
								callback.UpdateKey(k, string(e.Kv.Value))
							}
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

func (p *EtcdReader) getSettingFromAgent(appName string) {

	if len(p.appToken) == 0 {
		return
	}
	agentaddr := "127.0.0.1:57581"
	agentConn, err := grpc.Dial(agentaddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer agentConn.Close()

	c := agentApplicationService.NewAgentApplicationServiceClient(agentConn)
	req := &agentApplicationService.LogonInfoRequest{
		AppToken: p.appToken,
		Env:      p.env,
	}
	resp, err := c.GetLogonInfo(context.TODO(), req)
	if err != nil {
		st := status.Convert(err)
		for _, t := range st.Details() {
			switch t := t.(type) {
			case *agentApplicationService.LogonError:
				fmt.Println(t.Detail)
			}
		}
		return
	}

	if resp.AppName != appName {
		panic(fmt.Errorf("protoagent response app name is %s which is different from our app %s", resp.AppName, appName))
	}
	p.endpoints = strings.Split(resp.Endpoints, ",")
	p.user = resp.User
	p.pass = resp.Password
}

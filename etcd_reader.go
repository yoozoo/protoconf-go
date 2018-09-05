package protoconf_go

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
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
			panic(fmt.Errorf("env value is empty"))
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
func (p *EtcdReader) GetValues(appName string, keys []string) map[string]*string {

	prefix := "/" + p.env + "/" + appName + "/"
	result := make(map[string]*string)

	var ops []clientv3.Op
	for _, k := range keys {
		ops = append(ops, clientv3.OpGet(prefix+k))

		result[k] = nil
	}
	cli := p.getClient(appName)
	if cli == nil {
		fmt.Println("Failed to get etcd client connection")
		return result
	}

	txnResp, err := cli.Txn(context.TODO()).Then(ops...).Commit()

	if err != nil {
		fmt.Println("error to retrieve config values: ", err)
		return result
	}

	if !txnResp.Succeeded {
		fmt.Println("Failed to retrieve config values")

	}
	for _, resp := range txnResp.Responses {
		r := resp.GetResponseRange()
		if r != nil {
			for _, kv := range r.Kvs {
				v := string(kv.Value)
				result[strings.TrimPrefix(string(kv.Key), prefix)] = &v
			}
		}
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
	}

	if resp.AppName != appName {
		panic(fmt.Errorf("protoagent response app name is %s which is different from our app %s", resp.AppName, appName))
	}
	p.endpoints = strings.Split(resp.Endpoints, ",")
	p.user = resp.User
	p.pass = resp.Password
}

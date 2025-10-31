package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.etcd.io/etcd/client/v3"
)

// EtcdSource is a configuration source for etcd.
type EtcdSource struct {
	client   *clientv3.Client
	prefix   string
	user     string
	password string
}

// NewEtcdSource creates a new EtcdSource.
func NewEtcdSource(endpoints []string, prefix string, user string, password string) (*EtcdSource, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
		Username:    user,
		Password:    password,
	})
	if err != nil {
		return nil, fmt.Errorf("creating etcd client: %w", err)
	}

	return &EtcdSource{
		client:   cli,
		prefix:   prefix,
		user:     user,
		password: password,
	}, nil
}

// Load loads the configuration from etcd.
func (s *EtcdSource) Load() (map[string]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := s.client.Get(ctx, s.prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("getting config from etcd: %w", err)
	}

	data := make(map[string]any)
	for _, kv := range resp.Kvs {
		key := strings.TrimPrefix(string(kv.Key), s.prefix)
		key = strings.TrimPrefix(key, "/")
		keys := strings.Split(key, "/")
		m := data
		for i, k := range keys {
			if i == len(keys)-1 {
				var val any
				if err := json.Unmarshal(kv.Value, &val); err != nil {
					m[k] = string(kv.Value)
				} else {
					m[k] = val
				}
			} else {
				if _, ok := m[k]; !ok {
					m[k] = make(map[string]any)
				}
				m = m[k].(map[string]any)
			}
		}
	}

	return data, nil
}

// Watch watches for changes in the configuration in etcd.
func (s *EtcdSource) Watch(onChange func(map[string]any)) error {
	rch := s.client.Watch(context.Background(), s.prefix, clientv3.WithPrefix())
	go func() {
		for range rch {
			data, err := s.Load()
			if err != nil {
				// TODO: log error
				continue
			}
			onChange(data)
		}
	}()
	return nil
}

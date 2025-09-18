package sources

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/client/v3"
)

func TestEtcdSource(t *testing.T) {
	endpoints := []string{"localhost:2379"}
	prefix := "/config"

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	require.NoError(t, err)
	defer cli.Close()

	t.Cleanup(func() {
		_, err := cli.Delete(context.Background(), prefix, clientv3.WithPrefix())
		require.NoError(t, err)
	})

	_, err = cli.Put(context.Background(), "/config/foo", "bar")
	require.NoError(t, err)
	_, err = cli.Put(context.Background(), "/config/baz/qux", "123")
	require.NoError(t, err)

	source, err := NewEtcdSource(endpoints, prefix, "", "")
	require.NoError(t, err)

	data, err := source.Load()
	require.NoError(t, err)

	expected := map[string]any{
		"foo": "bar",
		"baz": map[string]any{
			"qux": float64(123),
		},
	}
	assert.Equal(t, expected, data)
}

func TestEtcdSource_Watch(t *testing.T) {
	endpoints := []string{"localhost:2379"}
	prefix := "/config"

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	require.NoError(t, err)
	defer cli.Close()

	t.Cleanup(func() {
		_, err := cli.Delete(context.Background(), prefix, clientv3.WithPrefix())
		require.NoError(t, err)
	})

	_, err = cli.Put(context.Background(), "/config/foo", "bar")
	require.NoError(t, err)

	source, err := NewEtcdSource(endpoints, prefix, "", "")
	require.NoError(t, err)

	ch := make(chan map[string]any)
	err = source.Watch(func(data map[string]any) {
		ch <- data
	})
	require.NoError(t, err)

	_, err = cli.Put(context.Background(), "/config/foo", "new-bar")
	require.NoError(t, err)

	select {
	case data := <-ch:
		expected := map[string]any{
			"foo": "new-bar",
		}
		assert.Equal(t, expected, data)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for config change")
	}
}

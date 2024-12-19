package cache

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/AndreiGoStorm/previewer/internal/config"
	"github.com/AndreiGoStorm/previewer/internal/logger"
	"github.com/AndreiGoStorm/previewer/internal/service"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	conf := initConfig()
	logg := logger.New(conf.Log.Level)
	storage := service.NewStorage(logg)

	t.Run("simple", func(t *testing.T) {
		conf.Capacity = 5
		c := New(conf.Cache, storage)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("dequeue the oldest element", func(t *testing.T) {
		conf.Capacity = 2
		c := New(conf.Cache, storage)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		val, ok = c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 300, val)
	})

	t.Run("dequeue the oldest element after get", func(t *testing.T) {
		conf.Capacity = 3
		c := New(conf.Cache, storage)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		c.Set("ddd", 400)

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("dequeue the oldest element for more quantity elements", func(t *testing.T) {
		conf.Capacity = 4
		c := New(conf.Cache, storage)

		for i := 0; i <= 1000; i++ {
			c.Set(strconv.Itoa(i), i)
		}

		val, ok := c.Get("997")
		require.True(t, ok)
		require.Equal(t, 997, val)

		val, ok = c.Get("998")
		require.True(t, ok)
		require.Equal(t, 998, val)

		val, ok = c.Get("999")
		require.True(t, ok)
		require.Equal(t, 999, val)

		val, ok = c.Get("1000")
		require.True(t, ok)
		require.Equal(t, 1000, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		conf.Capacity = 3
		c := New(conf.Cache, storage)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		c.Clear()

		_, ok = c.Get("aaa")
		require.False(t, ok)
		_, ok = c.Get("bbb")
		require.False(t, ok)
		_, ok = c.Get("ccc")
		require.False(t, ok)

		c.Set("ddd", 1000)

		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 1000, val)
	})

	err := os.RemoveAll(storage.Dir)
	require.NoError(t, err)
}

func initConfig() *config.Config {
	file, err := filepath.Abs("main.go")
	if err != nil {
		panic(err)
	}
	configFile := "/../../configs/config.yml"
	conf := config.New(path.Dir(file) + configFile)

	return conf
}

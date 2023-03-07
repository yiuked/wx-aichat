package lib

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// 运行测试用例
	exitCode := m.Run()

	// 跳过项目中的 init 函数
	os.Exit(exitCode)
}

func TestCache(t *testing.T) {
	cache := NewCache(time.Second)
	cache.Set("aaa", "1111", 5*time.Second)
	get, err := cache.Get("aaa")
	assert.NoError(t, err, "err should nil ")
	assert.Equal(t, get, "1111", "they should equal")
	time.Sleep(6 * time.Second)
	assert.NoError(t, err, "err should not nil ")
	assert.Equal(t, get, "1111", "they should empty")
}

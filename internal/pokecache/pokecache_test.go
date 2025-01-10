package pokecache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewCache(5 * time.Second)

	key := "test-key"
	val := []byte("test-value")

	cache.Add(key, val)

	got, exists := cache.Get(key)
	if !exists {
		t.Errorf("Expected to find key")
		return
	}

	if string(got) != string(val) {
		t.Errorf("Expected to get %s, got %s", string(val), string(got))
	}
}

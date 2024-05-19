package cachery

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"math"
	"testing"
	"time"
)

type SuiteLRU struct {
	suite.Suite

	cache Cache

	key string
	val any
}

func TestLRU(t *testing.T) {
	c := new(SuiteLRU)
	c.cache = NewCache(LRU)

	c.key = "key"
	c.val = []byte("values")

	suite.Run(t, c)
}

func (s *SuiteLRU) TestSetAndGet() {
	s.cache.Set(context.Background(), s.key, s.val)

	ok, value := s.cache.Get(context.Background(), s.key)

	require.Equal(s.T(), s.val, value)
	require.True(s.T(), ok)
}

func (s *SuiteLRU) TestMultiSetAndGet() {
	test := []struct {
		key string
		val any
	}{
		{
			key: "1",
			val: nil,
		},
		{
			key: "2",
			val: "oh my god",
		},
		{
			key: "3",
			val: []byte("really?"),
		},
		{
			key: "4",
			val: struct{}{},
		},
		{
			key: "5",
			val: math.MaxFloat64,
		},
	}
	for _, v := range test {
		s.T().Run(fmt.Sprintf("Set: %s", v.key), func(t *testing.T) {
			s.cache.Set(context.Background(), v.key, v.val)
		})
	}
	for _, v := range test {
		s.T().Run(fmt.Sprintf("Get: %s", v.key), func(t *testing.T) {
			ok, _ := s.cache.Get(context.Background(), v.key)
			require.True(s.T(), ok)
		})
	}
}

func (s *SuiteLRU) TestExtrusion() {
	c := NewCache(LRU)

	tests := []struct {
		key string
	}{
		{
			key: "1",
		},
		{
			key: "2",
		},
		{
			key: "3",
		},
		{
			key: "4",
		},
		{
			key: "5",
		},
		{
			key: "6",
		},
		{
			key: "7",
		},
		{
			key: "8",
		},
		{
			key: "9",
		},
		{
			key: "10",
		},
		{
			key: "11",
		},
	}

	for _, tc := range tests {
		c.Set(context.Background(), tc.key, nil)
	}

	// giving time to delete least recently used value
	time.Sleep(time.Millisecond * 10)

	ok, val := c.Get(context.Background(), tests[0].key)
	require.False(s.T(), ok)
	require.Nil(s.T(), val)
}

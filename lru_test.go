package cachery

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"math"
	"testing"
)

func TestLRU(t *testing.T) {
	c := new(SuiteLRU)
	c.cache = NewCache()

	c.key = "key"
	c.val = []byte("values")

	suite.Run(t, c)
}

type SuiteLRU struct {
	suite.Suite

	cache Cache

	key string
	val any
}

func (s *SuiteLRU) TestSetAndGet() {
	s.cache.Set(context.Background(), s.key, s.val)

	ok, val := s.cache.Get(context.Background(), s.key)

	require.True(s.T(), ok)
	require.Equal(s.T(), s.val, val)
}

func (s *SuiteLRU) TestMultiSetGet() {
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

// TODO: Extrusion test

package Cachery

import (
	"context"
	"github.com/stretchr/testify/suite"
	"math"
	"testing"
)

func TestSuiteLFU(t *testing.T) {
	suite.Run(t, new(SuiteLRU))
}

type SuiteLRU struct {
	suite.Suite

	cache Cache

	key string
	val any
}

func (s *SuiteLRU) SetupTest() {
	s.cache = NewLRU()

	s.key = "key"
	s.val = []byte("values")
}

func (s *SuiteLRU) TestSetGet() {
	s.cache.Set(context.Background(), s.key, s.val)
	val, ok := s.cache.Get(context.Background(), s.key)
	s.Require().True(ok)
	s.Require().Equal(s.val, val)

	s.cache.Set(context.Background(), s.key, "newVal")
	s.val = "newVal"
	val, ok = s.cache.Get(context.Background(), s.key)
	s.Require().True(ok)
	s.Require().Equal(s.val, val)
}

func (s *SuiteLRU) TestExtrusion() {
	tests := []struct {
		key string
		val any
	}{
		{
			key: "1",
			val: "ok",
		},
		{
			key: "2",
			val: math.Inf(1),
		},
		{
			key: "3",
			val: math.Float32bits(32.3),
		},
		{
			key: "4",
			val: nil,
		},
		{
			key: "5",
			val: struct {
			}{},
		},
		{
			key: "6",
			val: "omaha",
		},
		{
			key: "7",
			val: 54,
		},
		{
			key: "8",
			val: 300,
		},
		{
			key: "9",
			val: []int{1, 2, 3, 4},
		},
		{
			key: "10",
			val: []string{"h", "e", "l", "l", "o"},
		},
	}

	for _, v := range tests {
		s.cache.Set(context.Background(), v.key, v.val)
	}
	for _, v := range tests {
		val, ok := s.cache.Get(context.Background(), v.key)
		s.Require().True(ok)
		s.Require().Equal(v.val, val)
	}

	s.cache.Set(context.Background(), "11", "extra")
	v, ok := s.cache.Get(context.Background(), "1")
	s.Require().False(ok)
	s.Require().Nil(v)
}

func (s *SuiteLRU) TestReplace() {
	c := NewCache(LRU, WithCapacity(3))

	c.Set(context.Background(), "1", 23)
	c.Set(context.Background(), "2", 43)

	c.Get(context.Background(), "1")
	c.Set(context.Background(), "3", 90)
	c.Set(context.Background(), "4", 5)

	val, ok := c.Get(context.Background(), "1")
	s.Require().True(ok)
	s.Require().NotNil(val)

	val, ok = c.Get(context.Background(), "2")
	s.Require().False(ok)
	s.Require().Nil(val)

	val, ok = c.Get(context.Background(), "3")
	s.Require().True(ok)
	s.Require().NotNil(val)
}

package random

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/suite"
)

func (s *randomSuite) TestLenght() {
	given := int(8)
	got := NewRandomString(given)
	s.Len(got, 8, got)
}

func (s *randomSuite) TestRightSymbols() {
	got := NewRandomString(8)
	r := regexp.MustCompile("^[A-Za-z0-9]{8}$")
	s.True(r.Match([]byte(got)), got)
}

func (s *randomSuite) TestDifferents() {
	got := NewRandomString(8)
	got2 := NewRandomString(8)
	s.NotEqual(got, got2)
}

type randomSuite struct {
	suite.Suite
}

func TestRandomSuite(t *testing.T) {
	suite.Run(t, new(randomSuite))
}

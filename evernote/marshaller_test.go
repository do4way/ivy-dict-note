package evernote

import (
	"testing"

	. "gopkg.in/check.v1"
)

type SinkerTestSuite struct {
}

var (
	_ = Suite(&SinkerTestSuite{})
)

func TestEvernoteSinkerSuite(t *testing.T) {
	TestingT(t)
}

func (s *SinkerTestSuite) SetUpSuite(c *C) {

}

func (s *SinkerTestSuite) TestWriteTo(c *C) {
}

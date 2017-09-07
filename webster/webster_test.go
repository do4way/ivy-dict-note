package webster

import . "gopkg.in/check.v1"

type WebsterTestSuite struct {
}

var _ = Suite(&WebsterTestSuite{})

func (w *WebsterTestSuite) TestParseWebster(c *C) {
	words, err := ParseWebster("word")
	c.Check(err, Equals, nil)
	c.Check(len(words), Equals, 2)
}

package paramconvert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Test struct {
	Start  string
	Expect string
}

func TestBraceToColon(t *testing.T) {
	for _, v := range []Test{
		{"/{user}", "/:user"},
		{"/{user}/", "/:user/"},
		{"/{user}/another", "/:user/another"},
		{"/test/{user}", "/test/:user"},
		{"/test/{user}/asdf", "/test/:user/asdf"},
		{"/{user1}/{user2}", "/:user1/:user2"},
		{"/{user1}/something/{user2}", "/:user1/something/:user2"},
	} {
		assert.Equal(t, v.Expect, BraceToColon(v.Start))
	}
}

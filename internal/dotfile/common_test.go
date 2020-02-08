package dotfile

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type (
	RequireSuite struct {
		suite.Suite
		*require.Assertions
	}
)

func (rs *RequireSuite) SetT(t *testing.T) {
	rs.Assertions = rs.Suite.Require()
	rs.Suite.SetT(t)
}


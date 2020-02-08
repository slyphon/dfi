package testhelper

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type (
	RequireSuite struct {
		suite.Suite
		*require.Assertions
		beforeTests []HookFn
		afterTests []HookFn
	}

	HookFn func (suiteName, testName string)
)

func (rs *RequireSuite) AddBeforeHook(f HookFn) {
	rs.beforeTests = append(rs.beforeTests, f)
}

func (rs *RequireSuite) AddAfterHook(f HookFn) {
	rs.afterTests = append(rs.afterTests, f)
}

func (rs *RequireSuite) BeforeTest(suiteName, testName string) {
	for _, hook := range rs.beforeTests {
		hook(suiteName, testName)
	}
}

func (rs *RequireSuite) AfterTest(suiteName, testName string) {
	for _, hook := range rs.afterTests {
		hook(suiteName, testName)
	}
}

func (rs *RequireSuite) SetT(t *testing.T) {
	rs.Assertions = rs.Suite.Require()
	rs.Suite.SetT(t)
}

var _ suite.BeforeTest = &RequireSuite{}
var _ suite.AfterTest = &RequireSuite{}

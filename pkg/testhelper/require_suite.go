package testhelper

import (
	"runtime/debug"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type (
	DFISuite struct {
		suite.Suite
		beforeTests []HookFn
		afterTests []HookFn
	}

	HookFn func (suiteName, testName string)
)


func (rs *DFISuite) AddBeforeHook(f HookFn) {
	rs.beforeTests = append(rs.beforeTests, f)
}

func (rs *DFISuite) AddAfterHook(f HookFn) {
	rs.afterTests = append(rs.afterTests, f)
}

func (rs *DFISuite) BeforeTest(suiteName, testName string) {
	for _, hook := range rs.beforeTests {
		func() {
			defer func() {
				if p := recover(); p != nil {
					log.Errorf("stopped panic %s", string(debug.Stack()))
				}
			}()
			hook(suiteName, testName)
		}()
	}
}

func (rs *DFISuite) AfterTest(suiteName, testName string) {
	for _, hook := range rs.afterTests {
		func() {
			defer func() {
				if p := recover(); p != nil {
					log.Errorf("stopped panic %s", string(debug.Stack()))
				}
			}()
			hook(suiteName, testName)
		}()
	}
}

func (rs *DFISuite) SetT(t *testing.T) {
	rs.Suite.SetT(t)
}

var _ suite.BeforeTest = &DFISuite{}
var _ suite.AfterTest = &DFISuite{}

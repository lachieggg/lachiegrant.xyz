package temporalio

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestUnitTestSuite
func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
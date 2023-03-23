package temporalio

import (
	"context"
	"errors"

	// "fmt"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
)

const testsuccess = "test_success"
const testfailure = "test_failure"

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	Env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.Env = s.NewTestWorkflowEnvironment()
	s.Env.RegisterWorkflow(SimpleWorkflow)
	s.Env.RegisterActivity(SimpleActivity)
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.Env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) Test_SimpleWorkflow_Success() {
	s.Env.ExecuteWorkflow(SimpleWorkflow, WorkflowArgs{Status: testsuccess})
	s.True(s.Env.IsWorkflowCompleted())
	s.NoError(s.Env.GetWorkflowError())
}

func (s *UnitTestSuite) Test_SimpleWorkflow_ActivityParamCorrect() {
	s.Env.OnActivity(SimpleActivity, mock.Anything, mock.AnythingOfType("temporalio.WorkflowArgs")).Return(
		func(ctx context.Context, args WorkflowArgs) error {
			// TODO: fixme
			return nil
		})
	args := WorkflowArgs{Status: testsuccess}
	s.Env.ExecuteWorkflow(SimpleWorkflow, args)

	s.True(s.Env.IsWorkflowCompleted())
	s.NoError(s.Env.GetWorkflowError())
}

func (s *UnitTestSuite) Test_SimpleWorkflow_ActivityFails() {
	s.Env.OnActivity(SimpleActivity, mock.Anything, mock.Anything).Return(
		errors.New("SimpleActivityFailure"))
	s.Env.ExecuteWorkflow(SimpleWorkflow, WorkflowArgs{Status: testfailure})
	s.True(s.Env.IsWorkflowCompleted())

	err := s.Env.GetWorkflowError()
	s.Error(err)
	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("SimpleActivityFailure", applicationErr.Error())
}

// WorkflowArgs
type WorkflowArgs struct {
	Status string
}

// InitialiseWorkflow
func InitialiseWorkflow(ctx workflow.Context) workflow.Context {
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
		StartToCloseTimeout:    time.Second * 30,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumAttempts:    5,
		},
	}

	return workflow.WithActivityOptions(ctx, ao)
}

// SimpleWorkflow
func SimpleWorkflow(ctx workflow.Context, args WorkflowArgs) error {
	// Call ExecuteActivity to run SimpleActivity
	ctx = InitialiseWorkflow(ctx)
	return workflow.ExecuteActivity(ctx, SimpleActivity, nil, args).Get(ctx, nil)
}

// SimpleActivity
func SimpleActivity(ctx context.Context, args WorkflowArgs) error {
	return nil
}

package loanprocess

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/worker"
)

func TestReplayWorkflowHistoryFromFile(t *testing.T) {
	replayer := worker.NewWorkflowReplayer()

	replayer.RegisterWorkflow(LoanProcessingWorkflow)

	err := replayer.ReplayWorkflowHistoryFromJSONFile(nil, "history_for_original_execution.json")

	assert.NoError(t, err)
}

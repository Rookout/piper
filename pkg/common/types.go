package common

import (
	"github.com/rookout/piper/pkg/git"
)

type WorkflowsBatch struct {
	OnStart    []*git.CommitFile
	OnExit     []*git.CommitFile
	Templates  []*git.CommitFile
	Parameters *git.CommitFile
}

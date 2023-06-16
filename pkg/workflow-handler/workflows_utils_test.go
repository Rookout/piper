package workflow_handler

import (
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rookout/piper/pkg/git"
	assertion "github.com/stretchr/testify/assert"
	"testing"
)

func TestAddFilesToTemplates(t *testing.T) {
	assert := assertion.New(t)

	template := make([]v1alpha1.Template, 0)
	files := make([]*git.CommitFile, 0)

	content := `
- name: local-step
  inputs:
    parameters:
      - name: message
  script:
    image: alpine
    command: [ sh ]
    source: |
      echo "wellcome to {{ workflow.parameters.global }}
      echo "{{ inputs.parameters.message }}"
`
	path := "path"
	files = append(files, &git.CommitFile{
		Content: &content,
		Path:    &path,
	})

	template, err := AddFilesToTemplates(template, files)

	assert.Nil(err)
	assert.Equal("alpine", template[0].Script.Container.Image)
	assert.Equal([]string{"sh"}, template[0].Script.Command)
	expectedScript := "echo \"wellcome to {{ workflow.parameters.global }}\necho \"{{ inputs.parameters.message }}\"\n"
	assert.Equal(expectedScript, template[0].Script.Source)
}

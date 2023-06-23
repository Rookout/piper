package workflow_handler

import (
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rookout/piper/pkg/git_provider"
	assertion "github.com/stretchr/testify/assert"
	"testing"
)

func TestAddFilesToTemplates(t *testing.T) {
	assert := assertion.New(t)

	template := make([]v1alpha1.Template, 0)
	files := make([]*git_provider.CommitFile, 0)

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
	files = append(files, &git_provider.CommitFile{
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

func TestCreateDAGTemplate(t *testing.T) {
	// Test case 1: Empty file list
	assert := assertion.New(t)

	fileList := []*git_provider.CommitFile{}
	name := "template1"
	template, err := CreateDAGTemplate(fileList, name)
	assert.Nil(template)
	assert.Nil(err)

	// Test case 2: Missing content or path
	file := &git_provider.CommitFile{
		Content: nil,
		Path:    nil,
	}
	fileList = []*git_provider.CommitFile{file}
	name = "template2"
	template, err = CreateDAGTemplate(fileList, name)
	assert.Nil(template)
	assert.NotNil(err)

	// Test case 3: Valid file list
	path := "some-path"
	content := `- name: local-step1
  template: local-step
  arguments:
    parameters:
      - name: message
        value: step-1
- name: local-step2
  templateRef:
    name: common-toolkit
    template: versioning
    clusterScope: true
  arguments:
    parameters:
      - name: message
        value: step-2
  dependencies:
    - local-step1`
	file = &git_provider.CommitFile{
		Content: &content,
		Path:    &path,
	}
	fileList = []*git_provider.CommitFile{file}
	name = "template3"
	template, err = CreateDAGTemplate(fileList, name)

	assert.Nil(err)
	assert.NotNil(template)

	assert.Equal(name, template.Name)
	assert.NotNil(template.DAG)
	assert.Equal(2, len(template.DAG.Tasks))

	assert.NotNil(template.DAG.Tasks[0])
	assert.Equal("local-step1", template.DAG.Tasks[0].Name)
	assert.Equal("local-step", template.DAG.Tasks[0].Template)
	assert.Equal(1, len(template.DAG.Tasks[0].Arguments.Parameters))
	assert.Equal("message", template.DAG.Tasks[0].Arguments.Parameters[0].Name)
	assert.Equal("step-1", template.DAG.Tasks[0].Arguments.Parameters[0].Value.String())

	assert.NotNil(template.DAG.Tasks[1])
	assert.Equal("local-step2", template.DAG.Tasks[1].Name)
	assert.Equal(1, len(template.DAG.Tasks[1].Dependencies))
	assert.Equal("local-step1", template.DAG.Tasks[1].Dependencies[0])
	assert.NotNil(template.DAG.Tasks[1].TemplateRef)
	assert.Equal("common-toolkit", template.DAG.Tasks[1].TemplateRef.Name)
	assert.Equal("versioning", template.DAG.Tasks[1].TemplateRef.Template)
	assert.True(template.DAG.Tasks[1].TemplateRef.ClusterScope)
}

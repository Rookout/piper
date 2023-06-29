package workflow_handler

import (
	"fmt"
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
func TestValidateDAGTasks(t *testing.T) {
	assert := assertion.New(t)
	// Define test cases
	tests := []struct {
		name  string
		tasks []v1alpha1.DAGTask
		want  error
	}{
		{
			name: "Valid tasks",
			tasks: []v1alpha1.DAGTask{
				{Name: "Task1", Template: "Template1"},
				{Name: "Task2", TemplateRef: &v1alpha1.TemplateRef{Name: "Template2"}},
			},
			want: nil,
		},
		{
			name: "Empty task name",
			tasks: []v1alpha1.DAGTask{
				{Name: "", Template: "Template1"},
			},
			want: fmt.Errorf("task name cannot be empty"),
		},
		{
			name: "Empty template and templateRef",
			tasks: []v1alpha1.DAGTask{
				{Name: "Task1"},
			},
			want: fmt.Errorf("task template or templateRef cannot be empty"),
		},
	}

	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Call the function being tested
			got := ValidateDAGTasks(test.tasks)

			// Use assert to check the equality of the error
			if test.want != nil {
				assert.Error(got)
				assert.NotNil(got)
			} else {
				assert.NoError(got)
				assert.Nil(got)
			}
		})
	}
}

func TestCreateDAGTemplate(t *testing.T) {
	assert := assertion.New(t)
	t.Run("Empty file list", func(t *testing.T) {
		fileList := []*git_provider.CommitFile{}
		name := "template1"
		template, err := CreateDAGTemplate(fileList, name)
		assert.Nil(template)
		assert.Nil(err)
	})

	t.Run("Missing content or path", func(t *testing.T) {
		file := &git_provider.CommitFile{
			Content: nil,
			Path:    nil,
		}
		fileList := []*git_provider.CommitFile{file}
		name := "template2"
		template, err := CreateDAGTemplate(fileList, name)
		assert.Nil(template)
		assert.NotNil(err)
	})

	t.Run("Valid file list", func(t *testing.T) {
		path3 := "some-path"
		content3 := `- name: local-step1
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
		file := &git_provider.CommitFile{
			Content: &content3,
			Path:    &path3,
		}
		fileList := []*git_provider.CommitFile{file}
		name := "template3"
		template, err := CreateDAGTemplate(fileList, name)

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
	})

	t.Run("Invalid configuration", func(t *testing.T) {
		path4 := "some-path"
		content4 := `- noName: local-step1
  wrongkey2: local-step
- noName: local-step2
  wrongkey: something
  dependencies:
    - local-step1`
		file := &git_provider.CommitFile{
			Content: &content4,
			Path:    &path4,
		}
		fileList := []*git_provider.CommitFile{file}
		name := "template4"
		template, err := CreateDAGTemplate(fileList, name)

		assert.Nil(template)
		assert.NotNil(err)
	})

	t.Run("YAML syntax error", func(t *testing.T) {
		path5 := "some-path"
		content5 := `- noName: local-step1
  wrongkey2: local-step
error: should be list`
		file := &git_provider.CommitFile{
			Content: &content5,
			Path:    &path5,
		}
		fileList := []*git_provider.CommitFile{file}
		name := "template5"
		template, err := CreateDAGTemplate(fileList, name)

		assert.Nil(template)
		assert.NotNil(err)
	})
}

func TestConvertToValidString(t *testing.T) {
	assert := assertion.New(t)

	tests := []struct {
		input    string
		expected string
	}{
		{"A@bC!-123.def", "abc-123.def"},
		{"Hello World!", "helloworld"},
		{"123$%^", "123"},
		{"abc_123.xyz", "abc-123.xyz"}, // Underscore (_) should be converted to hyphen (-)
		{"..--..", "..--.."},           // Only dots (.) and hyphens (-) should remain
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			converted := ConvertToValidString(test.input)
			assert.Equal(converted, test.expected)
		})
	}
}

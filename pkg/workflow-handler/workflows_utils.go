package workflow_handler

import (
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rookout/piper/pkg/git"
	"gopkg.in/yaml.v3"
)

func CreateDAGTemplate(fileList []*git.CommitFile, name string) (*v1alpha1.Template, error) {
	DAGs := make([]v1alpha1.DAGTask, 0)
	for _, file := range fileList {
		DAGTask := &v1alpha1.DAGTask{}
		err := yaml.Unmarshal([]byte(*file.Content), DAGTask)
		if err != nil {
			return nil, err
		}
		DAGs = append(DAGs, *DAGTask)
	}

	template := &v1alpha1.Template{
		Name: name,
		DAG: &v1alpha1.DAGTemplate{
			Tasks: DAGs,
		},
	}

	return template, nil
}

func AddFilesToTemplate(templates []v1alpha1.Template, files []*git.CommitFile) ([]v1alpha1.Template, error) {
	for _, f := range files {
		t := &v1alpha1.Template{}
		err := yaml.Unmarshal([]byte(*f.Content), t)
		if err != nil {
			return nil, err
		}
		templates = append(templates, *t)
	}
	return templates, nil
}

func GetParameters(paramsFile *git.CommitFile) ([]v1alpha1.Parameter, error) {
	var params []v1alpha1.Parameter
	err := yaml.Unmarshal([]byte(*paramsFile.Content), &params)
	if err != nil {
		return nil, err
	}
	return params, nil
}

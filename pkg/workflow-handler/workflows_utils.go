package workflow_handler

import (
	"encoding/json"
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rookout/piper/pkg/git"
	"github.com/rookout/piper/pkg/utils"
	"gopkg.in/yaml.v3"
)

func CreateDAGTemplate(fileList []*git.CommitFile, name string) (*v1alpha1.Template, error) {
	if len(fileList) == 0 {
		return nil, fmt.Errorf("empty file list for %s", name)
	}
	DAGs := make([]v1alpha1.DAGTask, 0)
	for _, file := range fileList {
		if file.Content == nil || file.Path == nil {
			return nil, fmt.Errorf("missing content or path for %s", name)
		}
		DAGTask := make([]v1alpha1.DAGTask, 0)
		err := yaml.Unmarshal([]byte(*file.Content), &DAGTask)
		if err != nil {
			return nil, err
		}
		DAGs = append(DAGs, DAGTask...)
	}

	template := &v1alpha1.Template{
		Name: name,
		DAG: &v1alpha1.DAGTemplate{
			Tasks: DAGs,
		},
	}

	return template, nil
}

func AddFilesToTemplates(templates []v1alpha1.Template, files []*git.CommitFile) ([]v1alpha1.Template, error) {
	for _, f := range files {
		t := make([]v1alpha1.Template, 0)
		jsonBytes, err := utils.ConvertYAMLListToJSONList(*f.Content)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(jsonBytes, &t)
		templates = append(templates, t...)
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

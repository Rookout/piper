package workflow_handler

import (
	"encoding/json"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rookout/piper/pkg/git"
	"gopkg.in/yaml.v3"
	"log"
)

func CreateDAGTemplate(fileList []*git.CommitFile, name string) (*v1alpha1.Template, error) {
	DAGs := make([]v1alpha1.DAGTask, 0)
	for _, file := range fileList {
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

func AddFilesToTemplate(templates []v1alpha1.Template, files []*git.CommitFile) ([]v1alpha1.Template, error) {
	for _, f := range files {
		t := make([]v1alpha1.Template, 0)
		yamlData := make([]map[string]interface{}, 0)
		err := yaml.Unmarshal([]byte(*f.Content), &yamlData)
		if err != nil {
			return nil, err
		}

		jsonBytes, err := json.Marshal(&yamlData)
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %v", err)
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

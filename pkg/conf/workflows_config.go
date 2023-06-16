package conf

import (
	"encoding/json"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rookout/piper/pkg/utils"
)

type WorkflowConfig struct {
	Configs map[string]*ConfigInstance
}

type ConfigInstance struct {
	Spec   v1alpha1.WorkflowSpec `yaml:"spec"`
	OnExit []v1alpha1.DAGTask    `yaml:"onExit"`
}

func (wfc *WorkflowConfig) WorkflowsSpecLoad(configPath string) error {
	var jsonBytes []byte
	wfc.Configs = make(map[string]*ConfigInstance)

	configs, err := utils.GetFilesData(configPath)
	if err != nil {
		return err
	}

	for key, config := range configs {
		tmp := &ConfigInstance{}
		jsonBytes, err = utils.ConvertYAMToJSON(config)
		if err != nil {
			return err
		}
		err = json.Unmarshal(jsonBytes, tmp)
		if err != nil {
			return err
		}
		wfc.Configs[key] = tmp
	}

	return nil
}

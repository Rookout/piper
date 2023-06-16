package conf

import (
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rookout/piper/pkg/utils"
	"gopkg.in/yaml.v3"
)

type WorkflowConfig struct {
	Configs map[string]*ConfigInstance
}

type ConfigInstance struct {
	Spec   v1alpha1.WorkflowSpec `yaml:"spec"`
	OnExit []v1alpha1.DAGTask    `yaml:"onExit"`
}

func (wfc *WorkflowConfig) WorkflowsSpecLoad() error {
	wfc.Configs = make(map[string]*ConfigInstance)

	configs, err := utils.GetFilesData("/piper-config/..data")
	if err != nil {
		return err
	}

	for key, config := range configs {
		tmp := &ConfigInstance{}
		err = yaml.Unmarshal(config, tmp)
		if err != nil {
			return err
		}
		wfc.Configs[key] = tmp
	}

	return nil
}

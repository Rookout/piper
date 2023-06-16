package conf

import (
	"github.com/rookout/piper/pkg/utils"
	"gopkg.in/yaml.v3"
)

type WorkflowConfig struct {
	Configs map[string]*ConfigInstance
}

type ConfigInstance struct {
	Spec   []byte `yaml:"spec"`
	OnExit []byte `yaml:"onExit"`
}

func (wfc *WorkflowConfig) WorkflowsSpecLoad() error {
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

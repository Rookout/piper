package conf

import "github.com/rookout/piper/pkg/utils"

type WorkflowConfig struct {
	Configs map[string][]byte
}

func (cfg *WorkflowConfig) WorkflowsSpecLoad() error {
	configs, err := utils.GetFilesData("/piper-config/..data")
	if err != nil {
		return err
	}
	cfg.Configs = configs
	return nil
}

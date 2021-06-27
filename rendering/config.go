package rendering

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

func (s *Service) loadVariables(syncFile string) {
	// Load the config from config_defaults.yml
	err := s.k.Merge(s.cfg.ConfigDefaults)
	s.p.CheckIfError(err)
	s.p.DebugF("Loading sync config %s", syncFile)
	// Load the config from .sync.yml
	err = s.k.Load(file.Provider(syncFile), yaml.Parser())
	if err != nil {
		s.p.WarnF("file %s not loaded: %s", syncFile, err)
	}
}

func (s *Service) loadDataForFile(fileName string) Data {
	// Load the global variables into exposed values
	data := make(Data)
	err := s.k.Unmarshal(":globals", &data)
	s.p.CheckIfError(err)
	// Load the file-specific variables into values
	err = s.k.Unmarshal(fileName, &data)
	s.p.CheckIfError(err)
	return data
}

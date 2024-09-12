package module

import (
	"plugin"

	"github.com/trinitytechnology/ebrick/config"
	"github.com/trinitytechnology/ebrick/utils"
	"go.uber.org/zap"
)

type ModuleManager struct {
	options *Options
	modules map[string]Module
}

func NewModuleManager(options ...Option) *ModuleManager {

	return &ModuleManager{
		modules: make(map[string]Module),
		options: newOptions(options...),
	}
}

func (mm *ModuleManager) RegisterModule(m Module) error {
	log := mm.options.Logger
	log.Info("Registering module", zap.String("id", m.Id()))
	if err := m.Initialize(mm.options); err != nil {
		log.Error("Initialize module error", zap.Error(err))
		return err
	}
	mm.modules[m.Id()] = m
	log.Info("Module registered", zap.String("id", m.Id()), zap.String("name", m.Name()), zap.String("version", m.Version()))
	return nil
}

func (mm *ModuleManager) LoadDynamicModules() {
	log := mm.options.Logger
	log.Info("Loading dynamic modules")
	modules := config.GetConfig().Modules
	for _, module := range modules {
		if module.Enable {
			if module.Id == "" {
				log.Error("Module id is required", zap.String("name", module.Name))
				continue
			}
			err := mm.RegisterModuleById(module.Id)
			if err != nil {
				log.Error("Failed to load module", zap.String("id", module.Id), zap.String("name", module.Name), zap.Error(err))
			}
		}

	}
}

func (mm *ModuleManager) RegisterModuleById(moduleId string) error {

	// if module already installed
	path := MODULES_DIR + "/" + moduleId + ".so"

	if !utils.FileExists(path) {
		return ErrModuleNotFound
	}

	module, err := mm.LoadModule(path)
	if err != nil {
		return err
	}

	return mm.RegisterModule(module)
}

func (mm *ModuleManager) LoadModule(modulePath string) (Module, error) {
	log := mm.options.Logger
	// Open the module file.

	plug, err := plugin.Open(modulePath)
	if err != nil {
		log.Error("Failed to open module", zap.String("path", modulePath), zap.Error(err))
		return nil, err
	}

	// Look up the 'Module' symbol.
	symModule, err := plug.Lookup("Module")
	if err != nil {
		log.Error("Failed to find 'Module' symbol in module", zap.String("path", modulePath), zap.Error(err))
		return nil, err
	}

	// Assert the type to module.Module.
	loadedModule, ok := symModule.(Module)
	if !ok {
		err := ErrInvalidModuleType
		log.Error("Module type assertion failed", zap.String("path", modulePath), zap.String("expectedType", "module.Module"), zap.Error(err))
		return nil, err
	}
	return loadedModule, nil
}

func (mm *ModuleManager) GetModule(moduleId string) Module {
	return mm.modules[moduleId]
}

func (mm *ModuleManager) GetModules() map[string]Module {
	return mm.modules
}

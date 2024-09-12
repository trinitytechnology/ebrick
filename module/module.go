package module

import (
	"errors"
	"plugin"

	"github.com/trinitytechnology/ebrick/logger"
	"github.com/trinitytechnology/ebrick/utils"
	"go.uber.org/zap"
)

// Module
const (
	MODULES_DIR = "modules"
)

var (
	ErrInvalidModuleType  = errors.New("invalid plugin type")
	ErrModulePathNotFound = errors.New("Module path not found")
	ErrModuleNotFound     = errors.New("Module not found")
)

type Initializer interface {
	Initialize(options *Options) error
}

type MetaDataProvider interface {
	Id() string
	Name() string
	Version() string
	Description() string
}

type Module interface {
	Initializer
	MetaDataProvider
}

type ModuleManager struct {
	options *Options
	modules map[string]Module
	logger  *zap.Logger
}

func NewModuleManager(options ...Option) *ModuleManager {

	return &ModuleManager{
		modules: make(map[string]Module),
		options: newOptions(options...),
		logger:  logger.DefaultLogger,
	}
}

func (mm *ModuleManager) RegisterModule(m Module) error {
	mm.logger.Info("Registering module", zap.String("id", m.Id()))
	if err := m.Initialize(mm.options); err != nil {
		mm.logger.Error("Initialize module error", zap.Error(err))
		return err
	}
	mm.modules[m.Id()] = m
	mm.logger.Info("Module registered", zap.String("id", m.Id()))
	return nil
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
	// Open the module file.

	plug, err := plugin.Open(modulePath)
	if err != nil {
		mm.logger.Error("Failed to open module", zap.String("path", modulePath), zap.Error(err))
		return nil, err
	}

	// Look up the 'Module' symbol.
	symModule, err := plug.Lookup("Module")
	if err != nil {
		mm.logger.Error("Failed to find 'Module' symbol in module", zap.String("path", modulePath), zap.Error(err))
		return nil, err
	}

	// Assert the type to module.Module.
	loadedModule, ok := symModule.(Module)
	if !ok {
		err := ErrInvalidModuleType
		mm.logger.Error("Module type assertion failed", zap.String("path", modulePath), zap.String("expectedType", "module.Module"), zap.Error(err))
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

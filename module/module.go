package module

// Module
const (
	MODULES_DIR = "modules"
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

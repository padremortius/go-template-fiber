package storage

type (
	StorageCfg struct {
		Path string `yaml:"path" json:"path" validate:"required"`
	}
)

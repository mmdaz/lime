package models

import "github.com/mmdaz/lime/config"

type Module struct {
	Name string `gorm:"primary_key" json:"name"`
}

// SaveModule is a ...
func (m *Module) SaveModule() (*Module, error) {
	err := config.DB.Create(&m).Error
	if err != nil {
		return &Module{}, err
	}
	return m, nil
}

// FindModule is a ...
func (m *Module) FindModule(name string) (*Module, error) {
	err := config.DB.Model(Module{}).Where("name = ?", name).Take(&m).Error
	if err != nil {
		return &Module{}, err
	}
	return m, err
}

// DeleteModule is a ...
func (m *Module) DeleteModule(name string) error {
	err := config.DB.Model(Module{}).Where("name = ?", name).Delete(&m).Error
	if err != nil {
		return err
	}
	return nil
}

// FindAllModules is a ...
func FindAllModules() (*[]Module, error) {
	var modules []Module
	err := config.DB.Model(&Module{}).Find(&modules).Error
	if err != nil {
		return &[]Module{}, err
	}
	return &modules, err
}

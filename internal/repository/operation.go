package repository

import (
	"github.com/joao-fontenele/s-rank-fiis/internal/model"
	"gorm.io/gorm"
)

const operationTableName = "operations"

type Operation struct {
	Conn *gorm.DB
}

func (o Operation) Save(operation *model.Operation) error {
	return o.Conn.Table(operationTableName).Create(operation).Error
}

func (o Operation) FindByTicker(code string) ([]model.Operation, error) {
	var ops []model.Operation
	err := o.Conn.Table(operationTableName).Where("code=?", code).Find(&ops).Error
	return ops, err
}

func (o Operation) FindAll() ([]model.Operation, error) {
	var ops []model.Operation
	err := o.Conn.Table(operationTableName).Find(&ops).Error
	return ops, err
}

func (o Operation) DeleteByTicker(code string) error {
	return o.Conn.Table(operationTableName).Where("code=?", code).Delete(&model.Operation{}).Error
}

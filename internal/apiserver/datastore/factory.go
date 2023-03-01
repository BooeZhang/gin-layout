package datastore

import (
	"github.com/BooeZhang/gin-layout/internal/apiserver/datastore/datainterface"
	"gorm.io/gorm"
)

var _datastore Factory

// Factory 数据集工场
type Factory interface {
	SysUser() datainterface.SysUserData
	Close() error
	GetDB() *gorm.DB
}

// Client return the store client instance.
func Client() Factory {
	return _datastore
}

// SetClient set the iam store client.
func SetClient(factory Factory) {
	_datastore = factory
}

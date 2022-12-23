package datastore

import (
	"github.com/BooeZhang/gin-layout/internal/apiserver/datastore/datainterface"
)

var _datastore Factory

// Factory 数据集工场
type Factory interface {
	SysUser() datainterface.ISysUser
	Close() error
}

// Client return the store client instance.
func Client() Factory {
	return _datastore
}

// SetClient set the iam store client.
func SetClient(factory Factory) {
	_datastore = factory
}

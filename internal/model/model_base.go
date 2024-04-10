package model

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"

	"gorm.io/plugin/soft_delete"
)

// Model base model
type Model struct {
	ID        uint                  `json:"id" gorm:"primaryKey"`
	CreatedAt MysqlTimestamp        `json:"created_at" gorm:"column:created_at;index;default=0;autoCreateTime:milli;comment:创建时间"`
	UpdatedAt MysqlTimestamp        `json:"updated_at" gorm:"column:updated_at;default=0;autoUpdateTime:milli;comment:更新时间"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index;default=0;comment:删除时间"`
}

func (m *Model) SetModifyTime(t time.Time) {
	m.UpdatedAt = MysqlTimestamp(t.Unix())
}

func (m *Model) GetModifyTime() time.Time {
	length := len(fmt.Sprintf("%d", m.UpdatedAt))
	if length > 10 {
		return time.UnixMilli(int64(m.UpdatedAt))
	}

	return time.Unix(int64(m.UpdatedAt), 0)
}

type MysqlTimestamp int64

// UnmarshalJSON implements json unmarshal interface.
func (t *MysqlTimestamp) UnmarshalJSON(data []byte) (err error) {
	v, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = MysqlTimestamp(v)

	return
}

// MarshalJSON implements json marshal interface.
func (t MysqlTimestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(t), 10)), nil
}

// Value ...
func (t MysqlTimestamp) Value() (driver.Value, error) {

	return int64(t), nil
}

// Scan valueof time.Time 注意是指针类型 method
func (t *MysqlTimestamp) Scan(v interface{}) error {
	value, ok := v.(int64)
	if ok {
		s := fmt.Sprintf("%d", t)
		if len(s) > 10 {
			value = value / 1000
		}

		*t = MysqlTimestamp(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t *MysqlTimestamp) Timstamp() int64 {
	return int64(*t)
}

func (t *MysqlTimestamp) Time() time.Time {
	return time.Unix(int64(*t), 0).In(time.Local)
}

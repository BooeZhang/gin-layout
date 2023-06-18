package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Model base model
type Model struct {
	ID        uint32         `json:"id" gorm:"primary_key"`
	CreatedAt JSONTime       `json:"created_at" gorm:"column:created_at;index;type:datetime;comment:创建时间"`
	UpdatedAt JSONTime       `json:"updated_at" gorm:"column:updated_at;type:datetime;comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index;type:datetime;comment:删除时间"`
}

// JSON 自定义数据库json数据类型
type JSON []byte

func (j *JSON) Scan(value interface{}) error {
	_bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %s", value)
	}

	result := json.RawMessage{}
	err := json.Unmarshal(_bytes, &result)
	*j = JSON(result)
	return err
}

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

func (j *JSON) MarshalJSON() ([]byte, error) {
	if *j == nil {
		return []byte("null"), nil
	}
	return *j, nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("null point exception")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

func (j JSON) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

func (j JSON) Equals(j1 JSON) bool {
	return bytes.Equal([]byte(j), []byte(j1))
}

// JSONTime 自定义数据库时间数据类型
type JSONTime struct {
	time.Time
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

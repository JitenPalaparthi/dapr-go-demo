package models

import (
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
)

type Tip struct {
	ID           uint              `json:"id" gorm:"primarykey"`
	Type         string            `json:"type"`
	SubType      string            `json:"sub_type" gorm:"column:sub_type"`
	Data         datatypes.JSONMap `json:"data" gorm:"column:data"`
	Tags         string            `json:"tags"`
	Status       string            `json:"status" gorm:"default:active"`
	LastModified int64             `json:"last_modified" gorm:"column:last_modified"`
}

func (t *Tip) ToByte() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Tip) ToJSONString() (string, error) {
	buf, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (t *Tip) ToString() string {
	return fmt.Sprintln(*t)
}

func (t *Tip) Validate() error {
	//TODO
	return nil
}

func (t *Tip) ToTags() string {
	if t != nil {
		return fmt.Sprintf("%s,%s,%s,%s", t.Type, t.SubType, t.Data, t.Status)
	}
	return ""
}

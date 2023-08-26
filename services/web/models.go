package web

import (
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	DocumentContent []byte `gorm:"type: LONGBLOB; not null" json:"document_content"`
	FileName        string `gorm:"size:256; not null" json:"file_name"`
	Extension       string `gorm:"size:6; not null" json:"extension"`
	UserId          uint   `gorm:"not null" json:"user_id"`
	ProcessedDoc    ProcessedDocument
}

type ProcessedDocument struct {
	gorm.Model
	DocumentContent []byte `gorm:"type: LONGBLOB; not null" json:"document_content"`
	DocumentId      uint   `gorm:"not null"`
}

type DocumentRes struct {
	FileName       string `json:"file_name"`
	Extension      string `json:"extension"`
	Encoding       string `json:"encoding"`
	EncodedContent string `json:"encoded_content"`
}

type DocumentsRes struct {
	Documents []DocumentRes `json:"documents"`
}

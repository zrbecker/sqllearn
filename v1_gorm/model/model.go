package model

import (
	"gorm.io/gorm"
)

type Chain struct {
	gorm.Model

	Name    string
	ChainID string `gorm:"unique;"`
}

type Token struct {
	gorm.Model

	ChainID  uint  `gorm:"uniqueIndex:ux_tokens_denom_chain_id;"`
	Chain    Chain `gorm:"constraint:OnDelete:CASCADE;"`
	Name     string
	Denom    string `gorm:"uniqueIndex:ux_tokens_denom_chain_id;"`
	Decimals int
}

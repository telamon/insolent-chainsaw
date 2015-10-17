package wharfmaster

import (
	. "github.com/telamon/wharfmaster/models"
)

type WharfMaster struct {
	Services []Service
}

func New() *WharfMaster {

	return &WharfMaster{}
}

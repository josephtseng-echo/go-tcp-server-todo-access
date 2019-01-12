package service

import (
	"sync"
)

type ManagerStruct struct {
	Tasks map[uint64]*Task
	M *sync.RWMutex
}

var Manager = &ManagerStruct {
	Tasks : make(map[uint64]*Task),
	M : new(sync.RWMutex),
}
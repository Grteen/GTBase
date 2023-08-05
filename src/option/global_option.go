package option

import "sync"

type GlobalOpt struct {
	// figuare out how to store the data
	// 1 is cache storage
	// 0 is persistent storage
	StoreWay int
	// figure out whether server need to redo log when server run
	// 1 means need
	// 0 means not need
	NeedRedo int
}

var instance *GlobalOpt
var once sync.Once

func GetGlobalOpt() *GlobalOpt {
	once.Do(func() {
		if instance == nil {
			instance = &GlobalOpt{
				StoreWay: 0,
				NeedRedo: 1,
			}
		}
	})
	return instance
}

func SetGlobalOpt(opt *GlobalOpt) {
	instance = opt
}

func IsCache() bool {
	return GetGlobalOpt().StoreWay == 1
}

func NeedRedo() bool {
	return GetGlobalOpt().NeedRedo == 1
}

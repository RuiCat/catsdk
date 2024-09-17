package logs

import (
	"device/boltdb"
	"device/boltdb/bbolt"
	"runtime"
	"time"
)

// ConfingDB 全局数据库
var ConfingDB = func() *boltdb.DB {
	confing, err := boltdb.Open("./confing.db", boltdb.BoltOptions(0600, &bbolt.Options{Timeout: 1 * time.Second}))
	if err != nil {
		panic(err)
	}
	runtime.SetFinalizer(confing, func(db *boltdb.DB) {
		db.Close()
	})
	return confing
}()

// Confing 配置
var Confing = ConfingDB.From("Confing")

type confing struct {
	ID    int    `storm:"id,increment"`
	Key   string `storm:"unique"`
	Value any
}

// GetConfing 得到配置
func GetConfing[T any](key string, value T) T {
	var to confing
	if Confing.One("Key", key, &to) == nil {
		return (to.Value).(T)
	}
	Confing.Save(&confing{Key: key, Value: value})
	return value
}

// SetConfing 设置配置
func SetConfing[T any](key string, value T) error {
	var to confing
	err := Confing.One("Key", key, &to)
	if err != nil {
		return err
	}
	to.Value = value
	return Confing.Update(&to)
}

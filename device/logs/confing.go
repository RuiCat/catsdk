package logs

import (
	"device/boltdb"
	"device/boltdb/bbolt"
	"encoding/json"
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
	Value []byte
}

// GetConfing 得到配置
func GetConfing(key string, value any) error {
	to := new(confing)
	err := Confing.One("Key", key, to)
	if err != nil {
		to.Key = key
		to.Value, err = json.Marshal(value)
		if err != nil {
			return err
		}
		err = Confing.Save(to)
		if err != nil {
			return err
		}
	} else {
		err = json.Unmarshal(to.Value, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetConfing 设置配置
func SetConfing(key string, value any) error {
	var to confing
	err := Confing.One("Key", key, &to)
	if err != nil {
		return err
	}
	to.Value, err = json.Marshal(value)
	if err != nil {
		return err
	}
	return Confing.Update(&to)
}

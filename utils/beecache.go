package utils

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"time"

	"github.com/beego/beego/v2/core/logs"

	"github.com/beego/beego/v2/client/cache"
)

var cc cache.Cache

func InitCache() {
	//cacheConfig, _ := config.String("cache")
	cc = nil
	// if "redis" == cacheConfig {
	// 	initRedis()
	// } else if "memory" == cacheConfig {
	initMemory()
	// } else {
	// 	initFile()
	// }
	if cc != nil {
		logs.Info("Init cache success!!")
	}
}

// func initFile() {
// 	var err error
// 	cc, err = cache.NewCache("file", `{"CachePath":"./tmp/cache","FileSuffix":".cache","DirectoryLevel":"2","EmbedExpiry":"180"}`)
// 	if err != nil {
// 		logs.Error("New file cache error", err)
// 	}
// }

func initMemory() {
	//var err error
	mc := cache.NewMemoryCache()
	// mc, err := cache.NewCache("memory", `{"interval":"180"}`)
	// use the default strategy which will generate random time offset (range: [3s,8s)) expired
	cc = cache.NewRandomExpireCache(mc)
	// so the expiration will be [1m3s, 1m8s)
	// if err != nil {
	// 	logs.Error("New memory cache error", err)
	// }
}

// func initRedis() {
// 	var err error
// 	defer func() {
// 		if r := recover(); r != nil {
// 			cc = nil
// 		}
// 	}()
// 	key, _ := config.String("cacherediskey")
// 	conn, _ := config.String("cacheredishost")
// 	password, _ := config.String("cacheredispass")
// 	cc, err = cache.NewCache("redis", `{"key":"`+key+`","conn":"`+conn+`","password":"`+password+`"}`)

// 	if err != nil {
// 		logs.Error("New redis cache error", err)
// 	}
// }

func SetCache(key string, value interface{}, timeoutSecond int) error {
	data, err := Encode(value)
	if err != nil {
		logs.Error("Set cache error:", err)
		return err
	}
	if cc == nil {
		logs.Error("Set cache error cache is nil")
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			logs.Error("recover cache error:", r)
			//cc = nil
		}
	}()
	timeouts := time.Duration(timeoutSecond) * time.Second
	err = cc.Put(context.Background(), key, data, timeouts)
	if err != nil {
		logs.Error("Set cache error:", err)
		return err
	} else {
		return nil
	}
}

func GetCache(key string, to interface{}) error {
	if cc == nil {
		logs.Error("Get cache error cache is nil")
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			logs.Error("recover cache error:", r)
			//cc = nil
		}
	}()

	data, _ := cc.Get(context.Background(), key)
	if data == nil {
		logs.Warn("Get cache warn Cache不存在")
		return nil
	}
	err := Decode(data.([]byte), to)
	if err != nil {
		logs.Error(err)
	}

	return err
}

func DelCache(key string) error {
	if cc == nil {
		logs.Error("Delete cache error cache is nil")
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			logs.Error("recover cache error:", r)
			//cc = nil
		}
	}()

	err := cc.Delete(context.Background(), key)
	if err != nil {
		logs.Error("Delete cache error", err)
		return err
	} else {
		return nil
	}
}

// --------------------
// Encode
// 用gob进行数据编码
func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// -------------------
// Decode
// 用gob进行数据解码
func Decode(data []byte, to interface{}) error {
	if len(data) == 0 {
		logs.Info("Decode data is empty")
		return nil
	}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

package aliyun

import (
	"fmt"
	"reflect"
	"syscall"
	"testing"
)

func TestSyncLocalToOSS(t *testing.T) {
	metaKey := "Content-Md5sum"
	config := &OSSConfig{
		//Endpoint:        "oss-cn-hangzhou.aliyuncs.com",
		//BucketName:      "dev-blog-xiexianbin-cn",
		Endpoint:        "oss-cn-shanghai.aliyuncs.com",
		BucketName:      "blog-xiexianbin-cn",
		AccessKeyID:     "",
		AccessKeySecret: "",
	}
	sourceDir := "/Users/xiexianbin/work/code/github/xiexianbin/xiexianbin.github.io/public"
	//sourceDir := "/Users/xiexianbin/work/code/github/xiexianbin/docs.xiexianbin.cn/public"

	err := SyncLocalToOSS(config, sourceDir, metaKey, "", []string{})
	if err != nil {
		fmt.Println("err", err)
	}
}

func TestCacheWrite(t *testing.T) {
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	filePath := "./dev-blog-xiexianbin-cn.js"
	m := make(map[string]interface{})
	m["a"] = "a"
	m["b"] = "b"
	m["c"] = 3

	err := CacheWrite(m, filePath)
	if err != nil {
		fmt.Println("write", m, "to filepath", filePath, "err", err)
	}
}

func TestCacheRead(t *testing.T) {
	filePath := "./persistence_test.json"
	m, err := CacheRead(filePath)

	if err != nil {
		fmt.Println("read from", filePath, "err", err)
	}

	fmt.Println(m)
	fmt.Println(reflect.TypeOf(m))
}

package aliyun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/xiexianbin/gsync/utils"
)

func ReadDir(filesMap map[string]interface{}, sourceDir, subDir string) {
	if strings.HasSuffix(sourceDir, "/") == false {
		sourceDir += "/"
	}

	currentDir := sourceDir
	if subDir != "" {
		if strings.HasSuffix(subDir, "/") == false {
			subDir += "/"
		}
		currentDir += subDir
	}

	dirInfos, err := ioutil.ReadDir(currentDir)
	if err != nil {
		utils.Println("Read Source Dir error:", err)
	}

	for _, dirInfo := range dirInfos {
		if dirInfo.IsDir() {
			newSubDir := subDir + dirInfo.Name()
			ReadDir(filesMap, sourceDir, newSubDir)
		} else {
			filePath := currentDir + dirInfo.Name()
			fileByte, err := ioutil.ReadFile(filePath)
			if err != nil {
				utils.Println("read file err:", err)
			}

			fileContent := string(fileByte)
			re, err := regexp.Compile("<li>Build <small>&copy; .*</small></li>")
			if err != nil {
				utils.Println("init regexp err:", err)
			}

			fileContent = re.ReplaceAllString(fileContent, "")
			md5sum := utils.Md5sum(fileContent)

			filesMap[strings.Replace(filePath, sourceDir, "", 1)] = md5sum
		}
	}
}

func CacheWrite(m map[string]interface{}, cacheFile string) error {
	_, err := os.Stat(cacheFile)
	if err != nil {
		_, _ = os.Create(cacheFile)
	}

	_j, err := json.Marshal(m)
	if err != nil {
		utils.Println("json.Marshal failed:", err)
		return err
	}

	var j bytes.Buffer
	err = json.Indent(&j, _j, "", "  ")

	err = ioutil.WriteFile(cacheFile, []byte(j.String()), 0644)
	if err != nil {
		panic(err)
		return err
	}

	return nil
}

func CacheRead(filename string) (map[string]interface{}, error) {
	cacheBytes, err := ioutil.ReadFile(filename)
	var m map[string]interface{}
	err = json.Unmarshal(cacheBytes, &m)
	if err != nil {
		utils.Println("Unmarshal failed, ", err)
		return nil, err
	}
	return m, nil
}

func isStartWitch(str string, excludeList []string) bool {
	for _, e := range excludeList {
		if e != "" && strings.HasPrefix(str, e) {
			return true
		}
	}
	return false
}

// syncFiles Do upload new files to aliyun oss
func syncFiles(m map[string]interface{}, metaKey, sourceDir, action string, config *OSSConfig, ch chan bool) {
	if action == "" {
		action = "new"
	}
	// new ConcurrentMap
	cMap := utils.NewConcurrentMap()

	// set key and value to cMap
	for k, v := range m {
		cMap.Set(k, v)
	}

	// new goroutine
	wg := sync.WaitGroup{}
	wg.Add(utils.ShareCount)

	utils.Println("Do upload", action, "files")
	for i := 0; i < utils.ShareCount; i++ {
		// pre shared map, new goroutine to statics
		go func(ms *utils.ConcurrentMapShared, index int) {
			count := 1
			sum := len(ms.Items)
			ms.Mu.RLock() // read locak
			for k := range ms.Items {
				process := fmt.Sprintf("(map shard id %d, [%d/%d])", index, count, sum)
				metasMap := make(map[string]interface{})
				v, _ := cMap.Get(k)
				metasMap[metaKey] = v
				err := PutObjectFromFile(config, k, sourceDir+"/"+k, metasMap)
				if err != nil {
					HandleError(err)
					utils.Println("Upload", action, "OSS Object", process, k, "Error:", err)
				} else {
					utils.Println("Upload", action, "OSS Object", process, k, "Done.")
				}
				count++
			}
			ms.Mu.RUnlock() // unlock
			wg.Done()
		}((*cMap)[i], i)
	}

	// wait all goroutine stop
	wg.Wait()
	ch <- true
	utils.Println("upload", action, "files Done")
}

// syncDelFiles Do delete files from aliyun oss
func syncDelFiles(m map[string]interface{}, config *OSSConfig, ch chan bool) {
	// new ConcurrentMap
	cMap := utils.NewConcurrentMap()

	// set key and value to cMap
	for k, v := range m {
		cMap.Set(k, v)
	}

	// new goroutine
	wg := sync.WaitGroup{}
	wg.Add(utils.ShareCount)

	utils.Println("Do delete files:")
	for i := 0; i < utils.ShareCount; i++ {
		// pre shared map, new goroutine to statics
		go func(ms *utils.ConcurrentMapShared, index int) {
			count := 1
			sum := len(ms.Items)
			ms.Mu.RLock() // read locak
			for k := range ms.Items {
				process := fmt.Sprintf("(map shard id %d, [%d/%d])", index, count, sum)
				err := DeleteObject(config, k)
				if err != nil {
					HandleError(err)
					utils.Println("Delete OSS Object", process, k, "Error:", err)
				} else {
					utils.Println("Delete OSS Object", process, "]", k, "Done.")
				}
				count++
			}
			ms.Mu.RUnlock() // unlock
			wg.Done()
		}((*cMap)[i], i)
	}

	// wait goroutine done
	wg.Wait()
	ch <- true
	utils.Println("delete files Done")
}

func SyncLocalToOSS(config *OSSConfig, sourceDir, metaKey, cacheFile string, excludeList []string) error {
	if metaKey == "" {
		metaKey = "Content-Md5sum"
	}
	if cacheFile == "" {
		cacheFile = "/tmp/" + config.BucketName + ".json"
	}
	utils.Println("Begin to sync", sourceDir, "files, metaKey is", metaKey, ", cacheFile is", cacheFile, ", exclude file or direct is", excludeList)

	// read local files
	_filesMap := make(map[string]interface{})
	ReadDir(_filesMap, sourceDir, "")

	filesMap := make(map[string]interface{})
	for k := range _filesMap {
		if isStartWitch(k, excludeList) {
			utils.Println("Skip", k, "by exclude rule.")
			continue
		}
		filesMap[k] = _filesMap[k]
	}

	// list oss object metadata
	objectsMap := make(map[string]interface{})
	_, err := os.Stat(cacheFile)
	if err != nil {
		objectsMap, err = ListObjects(config, metaKey)
		if err != nil {
			HandleError(err)
		}
	} else {
		objectsMap, err = CacheRead(cacheFile)
		if err != nil {
			HandleError(err)
		}
	}

	// get diff map
	justM1, justM2, diffM1AndM2, err := utils.DiffMap(filesMap, objectsMap)
	if err != nil {
		HandleError(err)
	}

	// signal channel
	newFileChan, diffFileChan, deleteFileChan := make(chan bool), make(chan bool), make(chan bool)

	// do upload new file
	go syncFiles(justM1, metaKey, sourceDir, "new", config, newFileChan)

	// do update diff files
	go syncFiles(diffM1AndM2, metaKey, sourceDir, "update", config, diffFileChan)

	// do delete files
	go syncDelFiles(justM2, config, deleteFileChan)

	<-newFileChan
	<-diffFileChan
	<-deleteFileChan

	// cache new map to file
	_, err = os.Stat(cacheFile)
	if err == nil {
		_ = os.Truncate(cacheFile, 0)
	}

	err = CacheWrite(filesMap, cacheFile)
	if err != nil {
		utils.Println("cache file map to file fail.")
	} else {
		utils.Println("write cache success! path:", cacheFile)
	}

	utils.Println("Sync done! files is", filesMap)

	return nil
}

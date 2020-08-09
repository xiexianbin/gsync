package aliyun

import (
	"github.com/xiexianbin/gsync/utils"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSSConfig struct {
	Endpoint        string
	BucketName      string
	AccessKeyID     string
	AccessKeySecret string
}

func HandleError(err error) {
	utils.Println("Error:", err)
}

func ListObjects(config *OSSConfig, metaKey string) (map[string]interface{}, error) {
	objectsMap := make(map[string]interface{})
	client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		HandleError(err)
		return nil, err
	}

	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		HandleError(err)
		return nil, err
	}

	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			HandleError(err)
		}
		for _, object := range lsRes.Objects {
			headers, err := bucket.GetObjectDetailedMeta(object.Key)
			if err != nil {
				HandleError(err)
			}

			objectsMap[object.Key] = headers.Get("X-Oss-Meta-" + metaKey)
		}

		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}

	return objectsMap, nil
}

func PutObjectFromFile(config *OSSConfig, objectKey, filePath string, metasMap map[string]interface{}) error {
	client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		HandleError(err)
		return err
	}

	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		HandleError(err)
		return err
	}

	utils.Println("Begin to put objectKey:", objectKey, "filePath:", filePath, "metasMap:", metasMap)
	err = bucket.PutObjectFromFile(objectKey, filePath)
	if err != nil {
		HandleError(err)
		return err
	}

	for k, v := range metasMap {
		switch v.(type) {
		case string:
			err = bucket.SetObjectMeta(objectKey, oss.Meta(k, v.(string)))
			if err != nil {
				HandleError(err)
				return err
			}
		default:
			break
		}
	}
	utils.Println("--> put object", objectKey, "done.")

	return nil
}

func DeleteObject(config *OSSConfig, objectKey string) error {
	client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		HandleError(err)
		return err
	}

	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		HandleError(err)
		return err
	}

	err = bucket.DeleteObject(objectKey)
	if err != nil {
		HandleError(err)
		return err
	}

	return nil
}

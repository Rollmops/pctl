package app

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/davidscholberg/go-durationfmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func WaitUntilTrue(testFunction func() (bool, error), timeout time.Duration, interval time.Duration) (bool, error) {
	if interval > timeout {
		return false, fmt.Errorf("interval %s is greater than timeout %s", interval, timeout)
	}
	start := time.Now()
	for {
		result, err := testFunction()
		if err != nil {
			return false, err
		}
		if result == true {
			return true, nil
		}
		if time.Since(start) > timeout {
			return false, nil
		}
		time.Sleep(interval)
	}
}

func CompareStringSlices(slice1 []string, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for index := range slice1 {
		if slice1[index] != slice2[index] {
			return false
		}
	}
	return true
}

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", os.Getenv("HOME"), 1)
	}
	path = os.ExpandEnv(path)
	return path
}

func ByteCountIEC(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func DurationToString(d time.Duration) string {
	var format string
	if d > time.Hour*24 {
		format = "%dd %hh"
	} else if d > time.Hour {
		format = "%hh %mm"
	} else if d > time.Minute {
		format = "%mm %ss"
	} else {
		format = "%ss"
	}
	durationString, _ := durationfmt.Format(d, format)
	return durationString
}

func HashMd5File(filePath string) (string, error) {
	var returnMD5String string

	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

func CreateFileHashesFromCommand(command []string) (*map[string]string, error) {
	md5hashes := make(map[string]string)
	for _, arg := range command {
		fullPath, err := GetFullPathFromFile(arg)
		if err == nil {
			hash, err := HashMd5File(fullPath)
			if err != nil {
				return nil, err
			}
			md5hashes[arg] = hash
		}
	}
	return &md5hashes, nil
}

func GetFullPathFromFile(path string) (string, error) {
	logrus.Tracef("Getting full Path of %s", path)
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	logrus.Tracef("Trying absolute Path %s", fullPath)
	if FileExists(fullPath) {
		return fullPath, nil
	}
	fullPath, err = exec.LookPath(path)
	logrus.Tracef("Trying lookup Path %s", fullPath)
	if err != nil {
		return "", err
	}
	if FileExists(fullPath) {
		return fullPath, nil
	}
	return "", fmt.Errorf("unable to find Path %s", path)
}

func FileExists(path string) bool {
	logrus.Tracef("Checking Path %s exists", path)
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}

func IsInStringList(stringList []string, string string) bool {
	for _, s := range stringList {
		if s == string {
			return true
		}
	}
	return false
}

func MergeStringMap(map1 map[string]string, map2 map[string]string) map[string]string {
	returnMap := make(map[string]string)
	for k, v := range map1 {
		returnMap[k] = v
	}
	for k, v := range map2 {
		returnMap[k] = v
	}
	return returnMap
}

package filez

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// PathExists is used to determine whether the path folder exists
// True if it exists, false otherwise
// 判断文件是否存在，如果是链接文件或文件夹，检查链接的目标是否存在
func PathExists(path string) bool {
	// 获取文件信息
	fileInfo, err := os.Lstat(path)
	if err != nil {
		// 如果文件不存在
		if os.IsNotExist(err) {
			return false
		}
		// 如果出现其他错误
		return false
	}

	// 如果是符号链接文件
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		// 获取链接的目标
		target, err := os.Readlink(path)
		if err != nil {
			// 读取链接目标出错
			return false
		}
		if !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(path), target)
			//target2 := filepath.Join(path, target)
			//if e == nil {
			//target = target2
			//}
		}
		// 检查链接的目标是否存在
		if _, err := os.Stat(target); err != nil {
			// 如果目标不存在
			if os.IsNotExist(err) {
				return false
			}
			// 其他错误
			return false
		}
		return true // 链接的目标存在
	}

	// 如果是文件夹
	if fileInfo.IsDir() {
		return true
	}

	// 其他情况都认为文件存在
	return true
}
func SiblingPath(path string, sibFile string) string {
	parent := filepath.Dir(filepath.Clean(path))
	return filepath.Join(parent, filepath.Clean(sibFile))
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return true
	}
	return false
}
func FindFiles(folderPath string, ext []string) []string {
	var filesList []string
	filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && slices.Contains(ext, filepath.Ext(path)) {
			filesList = append(filesList, path)
		}
		return nil
	})
	return filesList
}

// a/b/c/d.txt   return d
func FileName(filePath string) string {
	return strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
}

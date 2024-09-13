package filez

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func FirstExistPath(list []string) string {
	for _, v := range list {
		if PathExists(v) {
			return v
		}
	}
	return ""
}

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

// 查找文件和子目录下的文件，不会返回文件夹路径
func FindFiles(folderPath string, ext []string) []string {
	var filesList []string
	filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if len(ext) == 0 {
				filesList = append(filesList, path)
			} else if slices.Contains(ext, filepath.Ext(path)) {
				filesList = append(filesList, path)
			}
		}
		return nil
	})
	return filesList
}

// a/b/c/d.txt   return d
func FileName(filePath string) string {
	return strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
}

// ReadDir lists all the file or dir names in the directory.not recursively
func ListDir(dir string) ([]string, error) {
	infos, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(infos))
	for i, info := range infos {
		names[i] = info.Name()
	}
	return names, nil
}

// 匹配模式*.txt这种unix风格的文件名匹配调用filepath.Match
func MatchesPattern(pattern, name string) bool {
	matched, err := filepath.Match(pattern, name)
	if err != nil {
		fmt.Println("Error matching pattern:", err)
		return false
	}
	return matched
}

// 删除符合规则的文件或文件夹
func DeleteFilesOrDirsByPattern(rootDir, pattern string, isFile bool) error {
	return filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 根据 isFile 判断是否处理文件
		if isFile && !info.IsDir() {
			if MatchesPattern(pattern, info.Name()) {
				fmt.Println("Deleting file:", path)
				if err := os.Remove(path); err != nil {
					return err
				}
			}
		} else if !isFile && info.IsDir() {
			if MatchesPattern(pattern, info.Name()) {
				fmt.Println("Deleting directory:", path)
				// 删除目录前必须先删除目录中的所有文件和子目录
				if err := os.RemoveAll(path); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

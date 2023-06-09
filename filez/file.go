package filez

import (
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"io"
	"os"
	"path/filepath"
)

// PathExists is used to determine whether the path folder exists
// True if it exists, false otherwise
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
func CopyFile(src string, dst string) (err error) {
	// 打开源文件
	inFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer inFile.Close()

	// 创建并打开目标文件
	outFile, err := os.Create(dst)
	if err != nil {
		return
	}
	defer outFile.Close()

	// 将源文件复制到目标文件
	_, err = io.Copy(outFile, inFile)
	return
}

// CopyDir is used to copy the folder and its contents recursively
func CopyDir(src string, dst string) (err error) {
	// 判断源文件夹是否存在
	if !PathExists(src) {
		return errors.New("source folder does not exist")
	}

	// 判断目标文件夹是否存在，不存在则创建
	if !PathExists(dst) {
		err = os.MkdirAll(dst, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// 打开源文件夹
	dir, err := os.Open(src)
	if err != nil {
		return err
	}
	defer dir.Close()

	// 读取源文件夹下的所有文件和文件夹
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	// 遍历所有文件和文件夹
	for _, fileInfo := range fileInfos {
		srcPath := filepath.Join(src, fileInfo.Name())
		dstPath := filepath.Join(dst, fileInfo.Name())

		// 如果是文件夹，则递归复制
		if fileInfo.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// 如果是文件，则直接复制
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
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

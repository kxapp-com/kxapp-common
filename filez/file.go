package filez

import (
	"archive/zip"
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"io"
	"os"
	"path/filepath"
	"strings"
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
func SiblingPath(path string, sibFile string) string {
	parent := filepath.Dir(filepath.Clean(path))
	return filepath.Join(parent, filepath.Clean(sibFile))
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

func ZipDir(source string, zipFilePath string) (string, error) {
	// Create a new zip archive
	//	zipFilePath := source + ".zip"
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through the directory and add all files to the zip archive
	err = filepath.Walk(source, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if fileInfo.IsDir() {
			return nil
		}

		// Open the file
		rpath, _ := filepath.Rel(source, filePath)
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
		rpath = filepath.ToSlash(rpath)
		//fmt.Printf("zip file %s to %s\n",filePath,rpath)

		//rpath = strings.ReplaceAll(rpath, "\\", "/")
		// Create a new file in the zip archive
		zipFileT, err := zipWriter.Create(rpath)
		if err != nil {
			return err
		}

		// Copy the file to the zip archive
		_, err = io.Copy(zipFileT, file)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	// Return the path to the zip file
	return zipFilePath, nil
}
func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			continue
		}

		fpath := filepath.Dir(path)
		if _, err := os.Stat(fpath); err != nil {
			if err := os.MkdirAll(fpath, 0755); err != nil {
				return err
			}
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

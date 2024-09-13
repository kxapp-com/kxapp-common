package filez

import (
	"encoding/json"
	"fmt"
	"github.com/kxapp-com/kxapp-common/shellz"
	"runtime"

	//"golang.org/x/net/html/charset"
	"io"
	"os"
	"path/filepath"
)

const CopyOptionSkipSymlink = 0
const CopyOptionCopySymlink = 1
const CopyOptionFollowLink = 2

// CopyFile 复制单个文件
func CopyFile(src, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}

// CopySymlink 复制符号链接本身
func CopySymlink(src, dst string) error {
	linkTarget, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(linkTarget, dst)
}

// CopyDir 递归复制目录及其内容
func CopyDir(src, dst string, symlinkHandling int) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", src)
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(srcPath, dstPath, symlinkHandling); err != nil {
				return err
			}
		} else {
			srcInfo, err := os.Lstat(srcPath)
			if err != nil {
				return err
			}

			if srcInfo.Mode()&os.ModeSymlink != 0 {
				// 处理符号链接
				if symlinkHandling == CopyOptionSkipSymlink {
					// 跳过符号链接
					continue
				} else if symlinkHandling == CopyOptionCopySymlink {
					// 复制符号链接本身
					if err := CopySymlink(srcPath, dstPath); err != nil {
						return err
					}
				} else if symlinkHandling == CopyOptionFollowLink {
					// 复制符号链接的目标
					linkTarget, err := os.Readlink(srcPath)
					if err != nil {
						return err
					}
					if err := CopyPath(linkTarget, dstPath, symlinkHandling); err != nil {
						return err
					}
				}
			} else {
				// 复制普通文件
				if err := CopyFile(srcPath, dstPath); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// CopyPath 复制文件或目录
func CopyPath(src, dst string, symlinkHandling int) error {
	srcInfo, err := os.Lstat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return CopyDir(src, dst, symlinkHandling)
	} else {
		if srcInfo.Mode()&os.ModeSymlink != 0 && symlinkHandling == 2 {
			linkTarget, err := os.Readlink(src)
			if err != nil {
				return err
			}
			src = linkTarget
		}
		return CopyFile(src, dst)
	}
}

/*
	func CopyFile(src string, dst string) (err error) {
		// 打开源文件
		inFile, err := os.Open(src)
		if err != nil {
			return
		}
		defer inFile.Close()

		if !PathExists(filepath.Dir(dst)) {
			os.MkdirAll(filepath.Dir(dst), os.ModePerm)
		}
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
*/
const LINK_PATCH_FILE_NAME = ".LinkPatch.json"

// recreateAll 如果是true，则会重建所有链接，如果是false，则遇到一个链接是有效存在就停止后面的链接创建
func RestoreSymlinks(linksFile string, destDir string, recreateAll bool) error {
	// 读取链接文件信息
	linksJSON, err := os.ReadFile(linksFile)
	if err != nil {
		return err
	}

	// 解析链接文件信息
	var linkInfos map[string]string
	if err := json.Unmarshal(linksJSON, &linkInfos); err != nil {
		return err
	}

	// 恢复链接文件
	for linkFile, linkFileTarget := range linkInfos {
		newName := filepath.Join(destDir, linkFile)
		if !recreateAll && PathExists(newName) {
			break
		} else {
			os.RemoveAll(newName)
			if runtime.GOOS == "windows" {
				shellz.CreateLinkWindows(linkFileTarget, newName)
			} else {
				// 创建链接文件
				if err := os.Symlink(linkFileTarget, newName); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func WriteFileAppend(name string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}

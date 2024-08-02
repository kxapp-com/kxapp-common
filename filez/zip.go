package filez

import (
	"archive/zip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// IsZipFile 检查文件是否是zip文件
func IsZipFile(filename string) bool {
	if !PathExists(filename) {
		return false
	}
	if IsDir(filename) {
		return false
	}
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	// 读取前四个字节
	magic := make([]byte, 4)
	_, err = file.Read(magic)
	if err != nil {
		return false
	}

	// 判断是否是zip文件的魔数
	return magic[0] == 0x50 && magic[1] == 0x4B && magic[2] == 0x03 && magic[3] == 0x04
}

func ZipDir(source string, zipFilePath string) (string, error) {
	return CompressFolder(source, zipFilePath, false)
}
func CompressFolder(srcDir, destZip string, patchLink bool) (string, error) {
	zipFile, err := os.Create(destZip)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	linkInfos := make(map[string]string)

	err = filepath.Walk(srcDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileInfo, err := os.Lstat(filePath)
		if err != nil {
			return err
		}

		if fileInfo.Name() == ".DS_Store" {
			return nil
		}

		zipPath, err := filepath.Rel(srcDir, filePath)
		if err != nil {
			return err
		}

		zipPath = filepath.ToSlash(zipPath) // 确保路径使用斜杠'/'

		if fileInfo.IsDir() {
			if fileInfo.Mode()&os.ModeSymlink != 0 {
				if linkTarget, err := os.Readlink(filePath); err == nil {
					linkInfos[zipPath] = linkTarget
				} else {
					return err
				}
				return nil
			}

			_, err = zipWriter.CreateHeader(&zip.FileHeader{
				Name:     zipPath + "/",
				Method:   zip.Store,
				Modified: fileInfo.ModTime(),
				Flags:    0x800, // 设置EFS，标记文件名和评论采用UTF-8编码
			})
		} else {
			if fileInfo.Mode()&os.ModeSymlink != 0 {
				if linkTarget, err := os.Readlink(filePath); err == nil {
					linkInfos[zipPath] = linkTarget
				} else {
					return err
				}
				return nil
			}

			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			zipFileHeader := &zip.FileHeader{
				Name:     zipPath,
				Method:   zip.Deflate,
				Modified: fileInfo.ModTime(),
				Flags:    0x800, // 设置EFS
			}

			zipFile, err := zipWriter.CreateHeader(zipFileHeader)
			if err != nil {
				return err
			}

			_, err = io.Copy(zipFile, file)
		}

		return err
	})

	if err != nil {
		return "", err
	}

	if patchLink {
		linksJSON, err := json.MarshalIndent(linkInfos, "", "    ")
		if err != nil {
			return "", err
		}
		linkFile, err := zipWriter.CreateHeader(&zip.FileHeader{
			Name:  LINK_PATCH_FILE_NAME,
			Flags: 0x800, // 同理，使用UTF-8编码
		})
		if err != nil {
			return "", err
		}
		_, err = linkFile.Write(linksJSON)
		if err != nil {
			return "", err
		}
	}

	return destZip, nil
}

func Unzip(zipFile, destDir string) error {
	// 打开 zip 文件进行读取
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	// 遍历 zip 文件中的每个文件
	for _, file := range r.File {
		// 解压缩文件的路径
		filePath := filepath.Join(destDir, file.Name)

		// 检查是否为目录
		if file.FileInfo().IsDir() {
			// 如果是目录，创建对应的目录
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		// 如果不是目录，创建对应的文件
		if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		// 创建文件
		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		// 打开 zip 文件中的文件
		zippedFile, err := file.Open()
		if err != nil {
			return err
		}
		defer zippedFile.Close()

		// 将 zip 文件中的文件内容拷贝到目标文件中
		_, err = io.Copy(outFile, zippedFile)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
	func ZipDir(source string, zipFilePath string) (string, error) {
		// Create a new zip archive
		//	zipFilePath := source + ".zip"
		zipFile, err1 := os.Create(zipFilePath)
		if err1 != nil {
			return "", err1
		}
		defer zipFile.Close()

		// Create a new zip writer
		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()

		if IsDir(source) {
			// Walk through the directory and add all files to the zip archive
			err := filepath.Walk(source, func(filePath string, fileInfo os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// Skip directories
				if fileInfo.IsDir() {
					return nil
				}

				if filepath.Base(filePath) == ".DS_Store" {
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
		} else {
			// Open the file
			file, err := os.Open(source)
			if err != nil {
				return "", err
			}
			defer file.Close()
			// Create a new file in the zip archive
			zipFileT, err := zipWriter.Create(filepath.Base(source))
			if err != nil {
				return "", err
			}
			// Copy the file to the zip archive
			_, err = io.Copy(zipFileT, file)
			if err != nil {
				return "", err
			}
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
*/

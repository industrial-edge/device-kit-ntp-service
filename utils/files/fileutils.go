/*
 * Copyright © Siemens 2023 - 2025. ALL RIGHTS RESERVED.
 * Licensed under the MIT license
 * See LICENSE file in the top-level directory
 */

package files

import (
	"bufio"
	"errors"
	"io"
	"io/fs"
	"log"
	"path/filepath"
)

type FileUtil interface {
	Copy(source string, target string) error
	CreateOrUpdateFile(path string, content string) error
	IsFileExist(path string) (bool, error)
}

type OsFileUtils struct {
	FileSystemOperations
}

func (fileUtils *OsFileUtils) Copy(source string, target string) error {
	src, err := fileUtils.Open(source)
	if err != nil {
		log.Printf("Copying failed cannot access source file %s: %s\n", source, err.Error())
		return err
	}
	defer src.Close()

	if err := fileUtils.MkdirAll(filepath.Dir(target), 0666); err != nil {
		log.Printf("Copying file failed while creating parent directory for %s: %s\n", target, err.Error())
		return err
	}

	dst, err := fileUtils.Create(target)
	if err != nil {
		log.Printf("Copying file failed while create: %s\n", err.Error())
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		log.Printf("Copying file failed while copy: %s\n", err.Error())
		return err
	}

	log.Printf("File %s copied to file %s successfully.", source, target)
	return nil
}

func (fileUtils *OsFileUtils) CreateOrUpdateFile(path string, content string) error {
	if targetFile, err := fileUtils.Create(path); err == nil {
		writer := bufio.NewWriter(targetFile)
		if _, err := writer.WriteString(content); err != nil {
			return err
		}

		if err = writer.Flush(); err != nil {
			return err
		}
		defer targetFile.Close()
	}
	return nil
}

func (fileUtils *OsFileUtils) IsFileExist(path string) (bool, error) {
	stat, err := fileUtils.Stat(path)

	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return false, err
		}

		return false, nil
	}

	if stat != nil && stat.IsDir() {
		//There might be a folder name same as with the target file
		log.Printf("Expected file, found directory: %s ", path)
		return false, nil
	}

	return true, nil
}

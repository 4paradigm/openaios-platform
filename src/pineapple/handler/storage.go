/*
 * Copyright © 2021 peizhaoyou <peizhaoyou@4paradigm.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handler

import (
	"code.cloudfoundry.org/bytefmt"
	"flag"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var root = flag.String("storage-root", os.Getenv("PINEAPPLE_STORAGE_ROOT"),
	"storage root")

func (handler *Handler) GetDirectory(ctx echo.Context, params apigen.GetDirectoryParams) error {
	userID := ctx.Get("userID").(string)

	var path string
	if params.Path == nil {
		path = ""
	} else {
		path = string(*params.Path)
	}
	cPath, err := abPath(userID, path)
	if err != nil {
		return response.BadRequestWithMessage(ctx, "Invalid path")
	}

	if fileInfo, err := os.Stat(cPath); os.IsNotExist(err) {
		return response.BadRequestWithMessage(ctx, "No such Directory")
	} else if !fileInfo.IsDir() {
		return response.BadRequestWithMessage(ctx, "Not a directory")
	}

	files, err := ioutil.ReadDir(cPath)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "read directory failed.").SetInternal(
			errors.Wrap(err, "Cannot read directory: "+cPath+utils.GetRuntimeLocation()))
	}

	fileInfoList := []apigen.FileInfo{}
	for _, v := range files {
		isDir := v.IsDir()
		modTime := v.ModTime()
		name := v.Name()
		size := bytefmt.ByteSize(uint64(v.Size()))
		fileInfoList = append(fileInfoList,
			apigen.FileInfo{IsDir: &isDir,
				ModificationTime: &modTime,
				Name:             &name,
				Size:             &size})
	}
	return ctx.JSON(http.StatusOK, fileInfoList)
}

func (handler *Handler) CreateDirectory(ctx echo.Context, params apigen.CreateDirectoryParams) error {
	userID := ctx.Get("userID").(string)

	path := string(params.Path)
	cPath, err := abPath(userID, path)
	if err != nil {
		return response.BadRequestWithMessage(ctx, "Invalid path")
	}

	err = os.Mkdir(cPath, 0777)
	if err != nil {
		if strings.Contains(err.Error(), "file exists") {
			return response.BadRequestWithMessage(ctx, "File exists")
		} else {
			return echo.NewHTTPError(
				http.StatusInternalServerError, "creation failed.").SetInternal(
				errors.Wrap(err, "Cannot create directory: "+cPath+". "+utils.GetRuntimeLocation()))
		}
	}
	return response.StatusOKNoContent(ctx)
}

func (handler *Handler) DeleteDirectoryOrFile(ctx echo.Context, params apigen.DeleteDirectoryOrFileParams) error {
	userID := ctx.Get("userID").(string)

	path := string(params.Path)
	cPath, err := abPath(userID, path)
	if err != nil {
		return response.BadRequestWithMessage(ctx, "invalid path")
	}

	err = os.RemoveAll(cPath)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "delete failed.").SetInternal(
			errors.Wrap(err, "Cannot delete directory or file: "+cPath+utils.GetRuntimeLocation()))
	}
	return response.StatusOKNoContent(ctx)
}

func (handler *Handler) PostStorageUpload(ctx echo.Context, params apigen.PostStorageUploadParams) error {
	userID := ctx.Get("userID").(string)

	// Source
	file, err := ctx.FormFile("file")
	if err != nil {
		return response.BadRequestWithMessage(ctx, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "upload failed.").SetInternal(
			errors.Wrap(err, "open source file error."+utils.GetRuntimeLocation()))
	}
	defer func() {
		if err = src.Close(); err != nil {
			ctx.Logger().Error(errors.Wrap(err, "close source file failed. "))
		}
	}()

	// Destination
	var path string
	if params.Path == nil {
		path = file.Filename
	} else {
		path = filepath.Join(string(*params.Path), file.Filename)
	}
	cPath, err := abPath(userID, path)
	if err != nil {
		return response.BadRequestWithMessage(ctx, "invalid path")
	}

	dst, err := os.Create(cPath)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "upload failed.").SetInternal(
			errors.Wrap(err, "cannot create file "+cPath+utils.GetRuntimeLocation()))
	}
	defer func() {
		if err = dst.Close(); err != nil {
			ctx.Logger().Error(errors.Wrap(err, "cannot close file: "+cPath))
		}
	}()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "upload failed.").SetInternal(
			errors.Wrap(err, "cannot copy file "+file.Filename+utils.GetRuntimeLocation()))
	}
	return response.StatusOKNoContent(ctx)
}

func (handler *Handler) GetStorageDownload(ctx echo.Context, params apigen.GetStorageDownloadParams) error {
	userID := ctx.Get("userID").(string)

	path := string(params.Path)
	cPath, err := abPath(userID, path)
	if err != nil {
		return response.BadRequestWithMessage(ctx, "invalid path")
	}

	if fileInfo, err := os.Lstat(cPath); os.IsNotExist(err) {
		return response.BadRequestWithMessage(ctx, "no such file")
	} else if fileInfo.IsDir() {
		return response.BadRequestWithMessage(ctx, "is a directory")
	}
	_, fileName := filepath.Split(cPath)
	return ctx.Attachment(cPath, fileName)
}

// MkUserDir create user directory when init user
func MkUserDir(userID string) error {
	if root == nil {
		return errors.New("flag storage root is not set " + utils.GetRuntimeLocation())
	}
	quotaBytes := conf.GetUserStorageQuotaBytes()
	userDir := filepath.Join(*root, userID)

	// Create new user directory
	e := os.MkdirAll(userDir, 0777)
	if e != nil {
		return errors.New("cannot create user directory " + utils.GetRuntimeLocation())
	}
	// set quota
	if quotaBytes != "" {
		cmd := exec.Command("setfattr", "-n", "ceph.quota.max_bytes", "-v", quotaBytes, userDir)
		if err := cmd.Run(); err != nil {
			return errors.New("cannot set quota for user " + utils.GetRuntimeLocation())
		}
	}
	return nil
}

//func GetUserRoot(userID string) (string, error) {
//	if root == nil {
//		return "", errors.New("flag storage root is not set")
//	}
//	userDir := filepath.Join(*root, userID)
//	if _, err := os.Stat(userDir); os.IsNotExist(err) {
//		return "", errors.New("user directory not exists")
//	}
//	return userDir, nil
//}

// join userid with path and validation
func abPath(userID string, path string) (string, error) {
	if root == nil {
		return "", errors.New("flag storage root is not set " + utils.GetRuntimeLocation())
	}
	rootPath := filepath.Clean(filepath.Join(*root, userID, "./"+path))
	if !strings.HasPrefix(rootPath, *root) {
		return "", errors.New("invalid file or directory path " + utils.GetRuntimeLocation())
	}
	return rootPath, nil
}

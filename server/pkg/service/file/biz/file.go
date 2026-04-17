package biz

import (
	"errors"
	"gitee.com/liujit/shop/server/api/file"
	"gitee.com/liujit/shop/server/lib/utils/str"
	"go.newcapec.cn/nctcommon/nmslib"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/oss"
	"go.newcapec.cn/ncttools/nmskit/log"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"slices"
)

type FileCase struct {
	oss.OSS
}

// NewFileCase new a File use case.
func NewFileCase(
	oss oss.OSS,
) *FileCase {
	nmslib.Runtime.SetOSS(oss)
	return &FileCase{
		OSS: oss,
	}
}

func (c *FileCase) MultiUploadFile(req *file.MultiUploadFileRequest) (*file.MultiUploadFileResponse, error) {
	files := make([]*file.FileInfo, 0)
	uploadFiles := req.GetFiles()
	if len(uploadFiles) == 0 {
		return nil, errors.New("no upload file")
	}
	for _, item := range uploadFiles {
		url, err := c.UploadByByte(item.GetName(), item.GetPath(), item.GetContent())
		if err != nil {
			log.Error("MultiUploadFile err:", err.Error())
			return nil, errors.New("文件上传失败")
		}
		files = append(files, &file.FileInfo{
			Url:     url,
			Name:    item.GetName(),
			Extname: item.GetExtname(),
		})
	}
	return &file.MultiUploadFileResponse{Files: files}, nil
}

func (c *FileCase) UploadFile(req *file.UploadFileInfo) (*file.FileInfo, error) {
	url, err := c.UploadByByte(req.GetName(), req.GetPath(), req.GetContent())
	if err != nil {
		log.Error("UploadFile err:", err.Error())
		return nil, errors.New("文件上传失败")
	}
	return &file.FileInfo{
		Url:     url,
		Name:    req.GetName(),
		Extname: req.GetExtname(),
	}, nil
}

func (c *FileCase) DownloadFile(req *file.DownloadFileRequest) (*wrapperspb.BytesValue, error) {
	fileByte, err := c.GetFileByte(req.GetPath())
	if err != nil {
		log.Error("DownloadFile err:", err.Error())
		return nil, errors.New("文件下载失败")
	}
	return &wrapperspb.BytesValue{Value: fileByte}, nil
}

func (c *FileCase) MultiDeleteFileByString(oldFile string, newFile []string) {
	c.MultiDeleteFile(str.ConvertJsonStringToStringArray(oldFile), newFile)
}

func (c *FileCase) MultiDeleteFile(oldFile, newFile []string) {
	for _, item := range oldFile {
		if len(newFile) == 0 || !slices.Contains(newFile, item) {
			err := c.OSS.DeleteFile(item)
			if err != nil {
				log.Error("MultiDeleteFile err:", err.Error())
			}
		}
	}
}

func (c *FileCase) DeleteFile(oldFile string, newFile string) {
	if newFile == "" || oldFile != newFile {
		// 删除旧文件
		err := c.OSS.DeleteFile(oldFile)
		if err != nil {
			log.Error("DeleteFile err:", err.Error())
		}
	}
}

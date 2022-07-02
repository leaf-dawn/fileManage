package manage

import (
	"FileManagement/constant"
	"FileManagement/obejct"
	"strings"
)

/**
 * 目录管理
 */

type DirectoryManage struct {
}

func NewDirectoryManage() *DirectoryManage {
	return &DirectoryManage{}
}

/**
 * 创建目录
 */
func (this *DirectoryManage) MakeDirectory(directoryName string) obejct.Result {
	if directoryName == "" {
		return obejct.OperateFailWithMessage("文件名不能为空")
	}
	directoryName = strings.TrimSpace(directoryName)
	//检验重名
	for _, directory := range obejct.Memory.CurrentDirectory.ChildDirectory {
		if directory.Fcb.FileName == directoryName {
			return obejct.OperateFailWithMessage("文件名重复")
		}
	}
	//创建fcb
	fcb := &obejct.Fcb{
		IsDirectory: true,
		FileName:    directoryName,
		StartBlock:  0,
		BlockNum:    0,
	}

	// 将文件控制块存储到磁盘中
	obejct.Cache.Disk.FcbList = append(obejct.Cache.Disk.FcbList, fcb)

	// 新建一个目录项
	directory := &obejct.Directory{
		Fcb:            fcb,
		Index:          0,
		ChildDirectory: []*obejct.Directory{},
		ParentIndex:    this.indexOf(obejct.Cache.Disk.DirectoryStruct, obejct.Memory.CurrentDirectory),
	}
	// 将目录项保存到树形结构目录中
	obejct.Cache.Disk.DirectoryStruct = append(obejct.Cache.Disk.DirectoryStruct, directory)
	directory.Index = len(obejct.Cache.Disk.DirectoryStruct) - 1
	obejct.Memory.CurrentDirectory.ChildDirectory = append(obejct.Memory.CurrentDirectory.ChildDirectory, directory)
	return obejct.OperateSuccess()
}

/**
 * 改变目录
 */
func (this *DirectoryManage) ChangeDirectory(path string) obejct.Result {
	if path == "" {
		// 如果路径为空，直接返回当前目录
		return obejct.OperateSuccess()
	}
	result := this.PathResolve(path)
	if result.IsSuccess {
		if result.Data.(*obejct.Directory).Fcb.IsDirectory {
			// 切换到目标目录
			obejct.Memory.CurrentDirectory = result.Data.(*obejct.Directory)
			return obejct.OperateSuccess()
		} else {
			// 如果最终找到的目录项不是目录文件类型，则报错
			return obejct.OperateFailWithMessage("找不到对应目录")
		}
	} else {
		return obejct.OperateFailWithMessage("找不到对应目录")
	}
}

/**
 *  列出文件目录
 */
func (this *DirectoryManage) ShowDirectory(directory *obejct.Directory) obejct.Result {
	if directory == nil {
		return obejct.OperateFailWithMessage("系统出错!!")
	}

	if !directory.Fcb.IsDirectory {
		return obejct.OperateFailWithMessage(directory.Fcb.FileName + "不是目录")
	}

	fcbs := []*obejct.Fcb{}
	for _, d := range directory.ChildDirectory {
		fcbs = append(fcbs, d.Fcb)
	}
	return obejct.OperateSuccess2(fcbs)
}

/**
* 解析路径并得到目录项
 */
func (this *DirectoryManage) PathResolve(path string) obejct.Result {
	if path == "" {
		return obejct.OperateFailWithMessage("路径不能为空")
	}

	currentDirectory := obejct.Memory.CurrentDirectory
	pathArray := strings.Split(strings.TrimSpace(path), constant.PATH_SEPARATOR)
	for i := 0; i < len(pathArray); i++ {
		if this.IsBackToPrevious(pathArray[i]) {
			if currentDirectory.ParentIndex != 0 {
				// 回退到上一级目录，知道当前目录为用户的根目录
				currentDirectory = obejct.Cache.Disk.DirectoryStruct[currentDirectory.ParentIndex]
			}
			continue
		}

		if pathArray[i] == "" {
			if i != 0 && i != len(pathArray)-1 {
				// 如果路径中间有空格则报错找不到目录（因为文件名不能为空）
				return obejct.OperateFailWithMessage("找不到对应目录")
			}
			continue
		}

		// 寻找子目录项
		childDirectory := this.SearchChildDirectory(currentDirectory, pathArray[i])
		if childDirectory == nil {
			return obejct.OperateFailWithMessage("找不到对应文件")
		} else {
			if (i+1) < len(pathArray) && !childDirectory.Fcb.IsDirectory {
				// 说明是路径的中间部分，路径的中间部分文件名应该都是对应文件夹类型才是正确的，否则报错
				return obejct.OperateFailWithMessage("找不到对应文件")
			}
			currentDirectory = childDirectory
		}
	}

	return obejct.OperateSuccess2(currentDirectory)
}

/**
 * 判断是否返回上一级目录
 */
func (this *DirectoryManage) IsBackToPrevious(path string) bool {
	if len(path) >= 2 {
		// 规定: 形如".."、"../"、".. "、"  .."都是表示上一级目录的路径
		return (len(path) == 2 && constant.BACK_TO_PREVIOUS_ONE == path) ||
			(len(path) > 2 && constant.BACK_TO_PREVIOUS_ONE == strings.TrimSpace(path))
	}

	return false
}

/**
 * 在一个目录下面寻找子目录项
 */
func (this *DirectoryManage) SearchChildDirectory(directory *obejct.Directory, directoryName string) *obejct.Directory {
	if directory == nil {
		return nil
	}
	for _, childDirectory := range directory.ChildDirectory {
		if childDirectory.Fcb.FileName == directoryName {
			return childDirectory
		}
	}
	return nil
}

//todo
func (this *DirectoryManage) indexOf(l []*obejct.Directory, d *obejct.Directory) int {
	for i, v := range l {
		if v == d {
			return i
		}
	}
	return -1
}

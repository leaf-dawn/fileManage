package manage

import (
	"FileManagement/obejct"
	"strings"
)

type FileManage struct {
	diskManage      *DiskManage
	directoryManage *DirectoryManage
}

func NewFileManage() *FileManage {
	return &FileManage{
		diskManage:      NewDiskManage(),
		directoryManage: NewDirectoryManage(),
	}
}

/**
 * 创建文件
 */
func (this *FileManage) CreateFile(fileName string) obejct.Result {
	if fileName == "" {
		return obejct.OperateFailWithMessage("文件名不能为空")
	}

	// 去除前后空格
	fileName = strings.TrimSpace(fileName)
	for _, directory := range obejct.Memory.CurrentDirectory.ChildDirectory {
		if directory.Fcb.FileName == fileName {
			return obejct.OperateFailWithMessage("文件名不可重复")
		}
	}

	// 新建数据文件的文件控制块。注意：创建新文件的时候，文件没有记录，不为它分配盘块
	fcb := &obejct.Fcb{
		IsDirectory: false,
		FileName:    fileName,
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
 * 打开文件
 */
func (this *FileManage) OpenFile(path string) obejct.Result {
	if path == "" {
		return obejct.OperateFailWithMessage("文件名不能为空")
	}

	if obejct.Memory.ActiveFile != nil {
		return obejct.OperateFailWithMessage("当前存在未关闭的文件")
	}

	resolveResult := this.directoryManage.PathResolve(path)
	if resolveResult.IsSuccess && !resolveResult.Data.(*obejct.Directory).Fcb.IsDirectory {
		// 将文件的记录从磁盘中读取出来
		result := this.diskManage.ReadRecord(resolveResult.Data.(*obejct.Directory).Fcb.StartBlock, resolveResult.Data.(*obejct.Directory).Fcb.BlockNum)
		if result.IsSuccess {
			//添加一个activeFile
			activeFile := &obejct.ActiveFile{
				Fcb:        resolveResult.Data.(*obejct.Directory).Fcb,
				FileRecord: result.Data.([]byte),
				ReadPtr:    0,
				WritePtr:   len(result.Data.([]byte)),
			}

			// 加载进内存中
			obejct.Memory.ActiveFile = activeFile
			return obejct.OperateSuccess2(activeFile)
		} else {
			return obejct.OperateFailWithMessage("打开文件失败: " + result.Message)
		}
	}

	return obejct.OperateFailWithMessage("文件不存在")
}

/**
 * 读取文件
 */
func (this *FileManage) ReadFile(recordNum int) obejct.Result {
	if recordNum == 0 {
		return obejct.OperateSuccess2("")
	}

	activeFile := obejct.Memory.ActiveFile
	if activeFile == nil {
		return obejct.OperateFailWithMessage("请先打开文件再进行读取")
	}

	readResult := []byte{}
	if recordNum > 0 {
		// 如果超过时改为最后一个
		if activeFile.ReadPtr+recordNum > len(activeFile.FileRecord) {
			recordNum = len(activeFile.FileRecord) - activeFile.ReadPtr
		}

		// 读取记录
		for i := 0; i < recordNum; i++ {
			// 读取一个记录
			readResult = append(readResult, activeFile.FileRecord[activeFile.ReadPtr])
			// 读指针向前移动
			activeFile.ReadPtr++
		}
	} else {
		// 向后读取记录时
		if activeFile.ReadPtr+recordNum < 0 {
			return obejct.OperateFailWithMessage("读取文件失败")
		}

		// 读取记录
		for i := 0; i > recordNum; i-- {
			// 读取一个记录
			readResult = append([]byte{activeFile.FileRecord[activeFile.ReadPtr]}, readResult...)
			// 读指针向后移动
			activeFile.ReadPtr--
		}
	}

	return obejct.OperateSuccess2(string(readResult))
}

/**
 * 向文件中填写记录
 */
func (this *FileManage) WriteToFile(record string) obejct.Result {
	if record == "" {
		return obejct.OperateSuccess()
	}

	activeFile := obejct.Memory.ActiveFile
	if activeFile == nil {
		return obejct.OperateFailWithMessage("请先打开文件再进行写入")
	}

	// 添加记录
	originalRecord := activeFile.FileRecord
	for i := 0; i < len(record); i++ {
		originalRecord = append(originalRecord, record[i])
	}
	// 存储记录
	resultTwo := this.diskManage.StoreRecord(activeFile.Fcb, originalRecord)
	if resultTwo.IsSuccess {
		// 设置写指针
		obejct.Memory.ActiveFile.WritePtr = obejct.Memory.ActiveFile.WritePtr + len(record)
		// 设置内存中的文件记录
		obejct.Memory.ActiveFile.FileRecord = originalRecord
		obejct.Memory.ActiveFile.Fcb = resultTwo.Data.(*obejct.Fcb)
		return obejct.OperateSuccess()
	} else {
		return obejct.OperateFailWithMessage("写入文件失败: " + resultTwo.Message)
	}
}

/**
 * 关闭文件
 */
func (this *FileManage) CloseFile() obejct.Result {
	if obejct.Memory.ActiveFile == nil {
		return obejct.OperateFailWithMessage("当前没有文件被打开")
	}
	obejct.Memory.ActiveFile = nil
	return obejct.OperateSuccess()
}

func (this *FileManage) DeleteFile(path string) obejct.Result {
	if path == "" {
		return obejct.OperateFailWithMessage("文件名不能为空")
	}
	//解析路径
	resolveResult := this.directoryManage.PathResolve(path)
	if resolveResult.IsSuccess && !resolveResult.Data.(*obejct.Directory).Fcb.IsDirectory {
		if obejct.Memory.ActiveFile != nil && resolveResult.Data.(*obejct.Directory).Fcb == obejct.Memory.ActiveFile.Fcb {
			return obejct.OperateFailWithMessage("[删除文件失败]: 该文件已被打开，请先关闭")
		}

		// 如果找到目录项且是数据文件类型
		parentDirectory := obejct.Cache.Disk.DirectoryStruct[resolveResult.Data.(*obejct.Directory).ParentIndex]
		index := this.indexOf(parentDirectory.ChildDirectory, resolveResult.Data.(*obejct.Directory))
		parentDirectory.ChildDirectory = append(parentDirectory.ChildDirectory[0:index], parentDirectory.ChildDirectory[index+1:]...)
		obejct.Cache.Disk.DirectoryStruct[resolveResult.Data.(*obejct.Directory).Index] = nil
		// 释放空间
		this.diskManage.FreeSpace(resolveResult.Data.(*obejct.Directory).Fcb.StartBlock, resolveResult.Data.(*obejct.Directory).Fcb.BlockNum)
		return obejct.OperateSuccess()
	} else {
		return obejct.OperateFailWithMessage("找不到对应文件")
	}
}

/**
 * 重命名文件
 */
func (this *FileManage) RenameFile(path string, newName string) obejct.Result {
	if path == "" || newName == "" {
		return obejct.OperateFailWithMessage("文件名不能为空")
	}

	resolveResult := this.directoryManage.PathResolve(path)
	if resolveResult.IsSuccess {
		parentDirectory := obejct.Cache.Disk.DirectoryStruct[resolveResult.Data.(*obejct.Directory).ParentIndex]
		//检验重名
		for _, childDirectory := range parentDirectory.ChildDirectory {
			if childDirectory.Fcb.FileName == newName {
				return obejct.OperateFailWithMessage("文件名重复")
			}
		}
		fcb := resolveResult.Data.(*obejct.Directory).Fcb
		fcb.FileName = newName
		return obejct.OperateSuccess()
	}
	return obejct.OperateFailWithMessage("文件不存在")
}

func (this *FileManage) indexOf(l []*obejct.Directory, d *obejct.Directory) int {
	for i, v := range l {
		if v == d {
			return i
		}
	}
	return -1
}

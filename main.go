package main

import (
	"FileManagement/constant"
	"FileManagement/manage"
	"FileManagement/obejct"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var userManage = &manage.UserManage{}
var directoryManage = &manage.DirectoryManage{}
var fileManage = manage.NewFileManage()

func main() {

	initDisk()

	for {
		showUserAndDirectory()
		fmt.Println()
		command := input()
		switch command[0] {
		// 注册用户
		case "register":
			if len(command) < 3 {
				fmt.Println("输入有误")
				break
			}
			registerResult := userManage.Register(command[1], command[2])
			fmt.Println(registerResult.Message)
		// 用户登录
		case "login":
			if len(command) < 3 {
				fmt.Println("输入有误")
				break
			}
			loginResult := userManage.Login(command[1], command[2])
			if !loginResult.IsSuccess {
				fmt.Println(loginResult.Message)
			}
		// 用户注销
		case "logout":

			logoutResult := userManage.Logout()
			if !logoutResult.IsSuccess {
				fmt.Println(logoutResult.Message)
			}
		// 创建目录
		case "mkdir":
			if obejct.Memory.CurrentUser == nil {
				fmt.Println("请先登录")
				break
			}
			if len(command) < 2 {
				fmt.Println("输入有误")
				break
			}
			mkdirResult := directoryManage.MakeDirectory(command[1])
			if !mkdirResult.IsSuccess {
				fmt.Println(mkdirResult.Message)
			}
		// 切换目录
		case "cd":
			if obejct.Memory.CurrentUser == nil {
				fmt.Println("请先登录")
				break
			}
			if len(command) < 2 {
				fmt.Println("输入有误")
				break
			}
			cdResult := directoryManage.ChangeDirectory(command[1])
			if !cdResult.IsSuccess {
				fmt.Println(cdResult.Message)
			}
		// 查看目录
		case "dir":
			if obejct.Memory.CurrentUser == nil {
				fmt.Println("请先登录")
				break
			}
			dirResult := directoryManage.ShowDirectory(obejct.Memory.CurrentDirectory)
			if dirResult.IsSuccess {
				printFileList(dirResult.Data.([]*obejct.Fcb))
			} else {
				fmt.Println(dirResult.Message)
			}
		// 创建文件
		case "create":
			if obejct.Memory.CurrentUser == nil {
				fmt.Println("请先登录")
				break
			}
			if len(command) < 2 {
				fmt.Println("输入有误")
				break
			}
			createResult := fileManage.CreateFile(command[1])
			if !createResult.IsSuccess {
				fmt.Println(createResult.Message)
			}
		// 打开文件
		case "open":
			if obejct.Memory.CurrentUser == nil {
				fmt.Println("请先登录")
			}
			if len(command) < 2 {
				fmt.Println("输入有误")
				break
			}
			openResult := fileManage.OpenFile(command[1])
			if openResult.IsSuccess {
				printOpenedFile(openResult.Data.(*obejct.ActiveFile))
			} else {
				fmt.Println(openResult.Message)
			}
		// 读取文件
		case "read":
			if obejct.Memory.CurrentUser == nil {
				fmt.Println("请先登录")
				break
			}
			if len(command) < 2 {
				fmt.Println("输入有误")
				break
			}
			in, e := strconv.Atoi(command[1])
			if e != nil {
				fmt.Println("输入有误")
				break
			}
			readResult := fileManage.ReadFile(in)
			if readResult.IsSuccess {
				fmt.Println(readResult.Data)
			} else {
				fmt.Println(readResult.Message)
			}
		// 写入文件
		case "write":
			if obejct.Memory.CurrentUser == nil {
				fmt.Println("[写入文件失败]: 请先登录")
				break
			}

			// 循环读取用户的输入，直到输入的内容以"$"结尾
			record := input2()

			writeResult := fileManage.WriteToFile(record)
			if writeResult.IsSuccess {
				// 显示修改后的文件记录
				printOpenedFile(obejct.Memory.ActiveFile)
			} else {
				fmt.Println(writeResult.Message)
			}
		// 删除文件
		case "delete":
			if obejct.Memory.CurrentUser == nil {
				fmt.Println("请先登录")
			}
			if len(command) < 2 {
				fmt.Println("输入有误")
				break
			}
			deleteResult := fileManage.DeleteFile(command[1])
			if !deleteResult.IsSuccess {
				fmt.Println(deleteResult.Message)
			}
		// 关闭文件
		case "close":
			if obejct.Memory.CurrentUser == nil {
				fmt.Println("请先登录")
				break
			}
			closeResult := fileManage.CloseFile()
			if !closeResult.IsSuccess {
				fmt.Println(closeResult.Message)
			}
		// 重命名文件
		case "rename":
			if obejct.Memory.CurrentUser == nil {
				fmt.Println("请先登录")
				break
			}
			if len(command) < 3 {
				fmt.Println("输入有误")
			}
			renameResult := fileManage.RenameFile(command[1], command[2])
			if !renameResult.IsSuccess {
				fmt.Println(renameResult.Message)
			}
			break
		// 退出系统
		case "exit":
			return
		}
	}
}

/**
 * 获取输入
 */
func input() []string {
	var msg string
	reader := bufio.NewReader(os.Stdin) // 标准输入输出
	msg, _ = reader.ReadString('\n')    // 回车结束
	msg = strings.TrimSpace(msg)        // 去除最后一个空格
	spaceRe, _ := regexp.Compile(`\s+`)
	answer := spaceRe.Split(msg, -1)
	return answer
}

/**
 * 获取输入
 */
func input2() string {
	var msg string
	reader := bufio.NewReader(os.Stdin) // 标准输入输出
	msg, _ = reader.ReadString('$')     // 钱字符结束
	return msg[0 : len(msg)-1]
}

/**
 * 显示当前用户和当前目录
 */
func showUserAndDirectory() {
	username := ""
	if obejct.Memory.CurrentUser != nil {
		username = obejct.Memory.CurrentUser.Username
	}
	path := ""
	if obejct.Memory.CurrentDirectory.ParentIndex == -1 {
		path = "/"
	} else {
		path = obejct.Memory.CurrentDirectory.Fcb.FileName
	}
	answer := "\n[" + username + "]" + path + ":"
	fmt.Print(answer)
}

/**
 * 初始化一个新的磁盘
 */
func initDisk() {

	newDisk := &obejct.Disk{}
	{
		// 初始化盘块
		disk := make([][]byte, 1024)
		newDisk.D = disk

		// 初始化位示图
		bitmap := make([][]int, constant.BITMAP_ROW_LENGTH)
		for i := 0; i < len(bitmap); i++ {
			bitmap[i] = make([]int, constant.BITMAP_LINE_LENGTH)
		}
		for i := 0; i < constant.BITMAP_ROW_LENGTH; i++ {
			for j := 0; j < constant.BITMAP_LINE_LENGTH; j++ {
				bitmap[i][j] = constant.BITMAP_FREE
			}
		}
		newDisk.BitMap = bitmap

		for i := 0; i < constant.BLOCK_NUM; i++ {
			block := make([]byte, 1024)
			disk = append(disk, block)
		}
		//todo
		newDisk.D = disk

		// 初始化系统用户数据盘块区
		for i := constant.USER_START_BLOCK; i < constant.USER_START_BLOCK+constant.USER_BLOCK_NUM; i++ {
			for j := 0; j < constant.BLOCK_SIZE; j++ {
				newDisk.D[i] = append(newDisk.D[i], 'U') //todo
				bitmap[0][i] = constant.BITMAP_BUSY
			}
		}

		// 初始化文件控制块数据盘块区
		for i := constant.FCB_START_BLOCK; i < constant.FCB_START_BLOCK+constant.FCB_BLOCK_NUM; i++ {
			for j := 0; j < constant.BLOCK_SIZE; j++ {
				newDisk.D[i] = append(newDisk.D[i], 'F') //todo
				bitmap[0][i] = constant.BITMAP_BUSY
			}
		}

		// 初始化树形目录结构数据盘块区
		for i := constant.DIR_START_BLOCK; i < constant.DIR_START_BLOCK+constant.DIR_BLOCK_NUM; i++ {
			for j := 0; j < constant.BLOCK_SIZE; j++ {
				newDisk.D[i] = append(newDisk.D[i], 'D')
				bitmap[0][i] = constant.BITMAP_BUSY
			}
		}

		// 初始化位示图数据盘块区
		for i := constant.BITMAP_START_BLOCK; i < constant.BITMAP_START_BLOCK+constant.BITMAP_BLOCK_NUM; i++ {
			for j := 0; j < constant.BLOCK_SIZE; j++ {
				newDisk.D[i] = append(newDisk.D[i], 'U')
				bitmap[0][i] = constant.BITMAP_BUSY
			}
		}
	}

	{
		// 初始化系统用户集
		userMap := map[string]*obejct.User{}
		// 新增一个管理员用户
		userMap["admin"] = &obejct.User{
			Username: "admin",
			Password: "123456",
		}
		newDisk.UserMap = userMap
	}

	{
		// 初始化文件控制块集
		fileControlBlockList := []*obejct.Fcb{}
		newDisk.FcbList = fileControlBlockList
		// 初始化根目录文件控制块
		root := &obejct.Fcb{
			IsDirectory: true,
			FileName:    "root",
			StartBlock:  0,
			BlockNum:    0,
		}
		newDisk.FcbList = append(newDisk.FcbList, root)

		// 初始化Administrator目录文件控制块
		administrator := &obejct.Fcb{
			IsDirectory: true,
			FileName:    "admin",
			StartBlock:  0,
			BlockNum:    0,
		}
		newDisk.FcbList = append(newDisk.FcbList, administrator)

		newDisk.DirectoryStruct = []*obejct.Directory{}
		// 初始化根目录文件项
		rootDirectory := &obejct.Directory{
			Fcb:            root,
			Index:          0,
			ChildDirectory: []*obejct.Directory{},
			ParentIndex:    -1,
		}
		newDisk.DirectoryStruct = append(newDisk.DirectoryStruct, rootDirectory)
		rootDirectory.Index = 0

		// 初始化Administrator目录文件项
		administratorDirectory := &obejct.Directory{
			Fcb:            administrator,
			Index:          0,
			ChildDirectory: []*obejct.Directory{},
			ParentIndex:    0,
		}
		newDisk.DirectoryStruct = append(newDisk.DirectoryStruct, administratorDirectory)
		administratorDirectory.Index = len(newDisk.DirectoryStruct) - 1
		newDisk.DirectoryStruct[0].ChildDirectory = append(newDisk.DirectoryStruct[0].ChildDirectory, administratorDirectory)
	}

	obejct.Cache.Disk = newDisk
	obejct.Memory.CurrentDirectory = obejct.Cache.Disk.DirectoryStruct[0]
	obejct.Memory.BitMap = obejct.Cache.Disk.BitMap
}

func printFileList(fileControlBlockList []*obejct.Fcb) obejct.Result {
	if fileControlBlockList == nil || len(fileControlBlockList) == 0 {
		return obejct.OperateSuccess()
	}

	fmt.Println()
	fmt.Printf("%-20s", "文件名")
	fmt.Printf("%-20s", "文件长度")
	fmt.Println()

	for _, fcb := range fileControlBlockList {
		fmt.Printf("%-20s", fcb.FileName)
		if fcb.IsDirectory {
			fmt.Printf("%-20s", "文件夹")
		} else {
			fmt.Printf("%-20s", strconv.Itoa(fcb.BlockNum)+" KB")
		}
		fmt.Println()
	}
	return obejct.OperateSuccess()
}

func printOpenedFile(activeFile *obejct.ActiveFile) obejct.Result {
	if activeFile == nil {
		return obejct.OperateFailWithMessage("文件为空")
	}
	fmt.Println("[" + activeFile.Fcb.FileName + "]: ")
	record := []byte{}
	for _, ch := range activeFile.FileRecord {
		record = append(record, ch)
	}
	fmt.Println(string(record))
	fmt.Println("[end]")

	return obejct.OperateSuccess()
}

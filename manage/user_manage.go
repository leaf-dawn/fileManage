package manage

import (
	"FileManagement/obejct"
)

type UserManage struct {
}

func NewUserManage() *UserManage {
	return &UserManage{}
}

/**
 * 登录
 */
func (*UserManage) Register(username string, password string) obejct.Result {
	//检验是否当前有用户
	if obejct.Memory.CurrentUser != nil {
		return obejct.OperateFailWithMessage("注册失败，请先退出当前用户")
	}
	if username == "" || password == "" {
		return obejct.OperateFailWithMessage("注册失败，用户名和密码不能为空")
	}
	//添加用户
	if _, ok := obejct.Cache.Disk.UserMap[username]; ok {
		return obejct.OperateFailWithMessage("注册失败,该用户已存在")

	}
	// 新建一个系统用户并保存
	obejct.Cache.Disk.UserMap[username] = &obejct.User{Username: username, Password: password}
	// 新建一个用户文件夹
	fcb := &obejct.Fcb{
		IsDirectory: true,
		FileName:    username,
		StartBlock:  0,
		BlockNum:    0,
	}
	obejct.Cache.Disk.FcbList = append(obejct.Cache.Disk.FcbList, fcb)

	// 新建一个用户目录项
	directory := &obejct.Directory{
		Fcb:            fcb,
		ChildDirectory: []*obejct.Directory{},
		ParentIndex:    0,
	}

	// 保存
	obejct.Cache.Disk.DirectoryStruct = append(obejct.Cache.Disk.DirectoryStruct, directory)
	directory.Index = len(obejct.Cache.Disk.DirectoryStruct) - 1
	// 添加到根目录下
	obejct.Cache.Disk.DirectoryStruct[0].ChildDirectory = append(obejct.Cache.Disk.DirectoryStruct[0].ChildDirectory, directory)

	return obejct.OperateSuccessWithMessage("[注册用户成功]")

}

func (*UserManage) Login(username string, password string) obejct.Result {
	if obejct.Memory.CurrentUser != nil {
		return obejct.OperateFailWithMessage("登录失败，请先退出当前用户")
	}
	user, ok := obejct.Cache.Disk.UserMap[username]
	if ok {
		if user.Password != password {
			return obejct.OperateFailWithMessage("登录失败: 密码错误")
		}
	} else {
		return obejct.OperateFailWithMessage("登录失败: 该用户不存在")

	}

	for _, directory := range obejct.Cache.Disk.DirectoryStruct[0].ChildDirectory {
		if directory.Fcb.IsDirectory && directory.Fcb.FileName == username {
			obejct.Memory.CurrentUser = user
			obejct.Memory.CurrentDirectory = directory
			// 切换到用户目录
			return obejct.OperateSuccess()
		}
	}
	return obejct.OperateFailWithMessage("系统错误")
}

func (*UserManage) Logout() obejct.Result {
	if obejct.Memory.CurrentUser == nil {
		return obejct.OperateFailWithMessage("当前没有登录用户")
	}
	obejct.Memory.CurrentUser = nil
	// 切换到根目录
	obejct.Memory.CurrentDirectory = obejct.Cache.Disk.DirectoryStruct[0]
	return obejct.OperateSuccess()
}

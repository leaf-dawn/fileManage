package obejct

/**
 * 磁盘
 */

type Disk struct {
	/** 磁盘空间 */
	D [][]byte
	/** 用户 */
	UserMap map[string]*User
	/** 文件控制块 */
	FcbList []*Fcb
	/** 目录 */
	DirectoryStruct []*Directory
	/** 位视图 */
	BitMap [][]int
}

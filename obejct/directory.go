package obejct

/**
 * 目录结构
 */
type Directory struct {

	/**文件控制块*/
	Fcb *Fcb
	/**目录的位置*/
	Index int
	/** 子目录*/
	ChildDirectory []*Directory
	/**父目录项的位置*/
	ParentIndex int
}

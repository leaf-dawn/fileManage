package obejct

/**
 * 文件控制块
 */

type Fcb struct {
	/** 是否使目录 */
	IsDirectory bool
	/** 文件名称 */
	FileName string
	/** 开始盘块 */
	StartBlock int
	/** 占用盘块数 */
	BlockNum int
}

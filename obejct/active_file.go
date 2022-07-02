package obejct

/**
 * 打开的文件
 */

type ActiveFile struct {
	Fcb        *Fcb
	FileRecord []byte
	ReadPtr    int
	WritePtr   int
}

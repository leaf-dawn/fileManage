package obejct

/**
 * 内存
 */

var Memory *memory = &memory{}

type memory struct {
	CurrentUser      *User
	CurrentDirectory *Directory
	BitMap           [][]int
	ActiveFile       *ActiveFile
}

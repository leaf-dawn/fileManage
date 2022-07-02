package manage

import (
	"FileManagement/constant"
	"FileManagement/obejct"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type DiskManage struct {
}

func NewDiskManage() *DiskManage {
	return &DiskManage{}
}

func (this *DiskManage) StoreRecord(fcb *obejct.Fcb, record []byte) obejct.Result {
	if record == nil || len(record) == 0 {
		return obejct.OperateSuccess()
	}

	// 修改位示图，为重新分配盘块做准备
	this.changeBitmapStatus(fcb.StartBlock, fcb.BlockNum, true)
	// 计算存储该长度记录所需要的盘块数量
	requiredNum := this.ceilDivide(len(record), constant.BLOCK_SIZE)

	count := 0
	startBlockId := constant.RECORD_START_BLOCK

	for i := 0; i < constant.BITMAP_ROW_LENGTH; i++ {
		for j := 0; j < constant.BITMAP_LINE_LENGTH; j++ {
			// 如果该盘块空闲
			if constant.BITMAP_FREE == obejct.Memory.BitMap[i][j] {
				if count == 0 {
					// 记下起始盘块号: 当前行 * 总列数 + 当前列
					startBlockId = i*constant.BITMAP_LINE_LENGTH + j
				}
				count++
				if count == requiredNum {
					// 如果有足够的连续盘区供存储，则进行存储，并改变位示图的相应状态
					this.StoreToDisk(startBlockId, record)
					//todo
					this.changeBitmapStatus(startBlockId, requiredNum, false)
					fcb.StartBlock = startBlockId
					fcb.BlockNum = requiredNum
					return obejct.OperateSuccess2(fcb)
				}
			} else {
				// 因为是连续分配，如果该盘块不空闲则要重新计数
				count = 0
			}
		}
	}

	return obejct.OperateFailWithMessage("分配盘块失败 磁盘空间不足")
}

func (this *DiskManage) FreeSpace(startBlockId int, blockNum int) obejct.Result {
	this.changeBitmapStatus(startBlockId, blockNum, true)
	return obejct.OperateSuccess()
}

/**
 *获取指定盘块上的文件记录
 */
func (this *DiskManage) ReadRecord(startBlockId int, blockNum int) obejct.Result {
	record := []byte{}

	if startBlockId == 0 {
		return obejct.OperateSuccess2(record)
	}

	if startBlockId < constant.RECORD_START_BLOCK || startBlockId >= constant.BLOCK_NUM ||
		blockNum <= 0 || blockNum > constant.BLOCK_NUM {
		return obejct.OperateFailWithMessage("系统错误")
	}

	disk := obejct.Cache.Disk
	// 盘块指针，用于从盘块中读取记录，加载进内存
	blockPtr := 0
	i := 0
	for i < blockNum {
		if i == blockNum-1 {
			if blockPtr >= len(disk.D[startBlockId+i]) || disk.D[startBlockId+i][blockPtr] == 0 {
				// 说明所有记录都读取完成了
				break
			}
		}
		record = append(record, disk.D[startBlockId+i][blockPtr])
		if blockPtr >= constant.BLOCK_SIZE {
			blockPtr = 0
			i++
		} else {
			blockPtr++
		}
	}

	return obejct.OperateSuccess2(record)
}

/**
 * 保存虚拟磁盘数据到本地文件中
 */
func (this *DiskManage) SaveDisk(savePath string) obejct.Result {
	file, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	data, _ := json.Marshal(obejct.Cache.Disk)
	_, e := write.Write(data)
	//Flush将缓存的文件真正写入到文件中
	e = write.Flush()
	if e != nil {
		return obejct.OperateFailWithMessage("保存磁盘失败")
	} else {
		return obejct.OperateSuccessWithMessage("保存磁盘成功")

	}
}

/**
 * 将记录存储到磁盘中
 */
func (this *DiskManage) StoreToDisk(startBlockId int, record []byte) {
	index := 0
	blockId := startBlockId
	disk := obejct.Cache.Disk

	// 可以直接覆盖掉原来的磁盘中的记录
	for _, ch := range record {
		if index >= constant.BLOCK_SIZE {
			// 如果一个盘块的空间被用完了，则使用下一个盘块来进行存储
			blockId++
			index = 0
		}
		//todo
		disk.D[blockId] = append(append(append([]byte{}, disk.D[blockId][0:index]...), ch), disk.D[blockId][index:]...)
		index++
	}

	for index < constant.BLOCK_SIZE && len(disk.D[blockId]) > index {
		// 擦除最后一个盘块中没有用到的空间
		disk.D[blockId][index] = 0
		index++
	}
}

/**
 * 更改相应的位示图的状态
 * 仅用于连续分配
 */
func (this *DiskManage) changeBitmapStatus(startBlockId int, blockNum int, changeToFree bool) {
	if startBlockId < constant.RECORD_START_BLOCK || startBlockId >= constant.BLOCK_NUM || blockNum <= 0 {
		return
	}

	// 解析该盘块在位示图中的第几行
	row := startBlockId / constant.BITMAP_LINE_LENGTH
	// 解析该盘块在位示图中的第几列
	line := startBlockId % constant.BITMAP_LINE_LENGTH
	for i := 0; i < blockNum; i++ {
		if changeToFree {
			obejct.Memory.BitMap[row][line] = constant.BITMAP_FREE
		} else {
			obejct.Memory.BitMap[row][line] = constant.BITMAP_BUSY
		}
		if line >= constant.BITMAP_LINE_LENGTH-1 {
			line = 0
			row++
		} else {
			line++
		}
	}
}

func (this *DiskManage) ceilDivide(dividend int, divisor int) int {
	if dividend%divisor == 0 {
		// 整除
		return dividend / divisor
	} else {
		// 不整除，向上取整
		return (dividend + divisor) / divisor
	}
}

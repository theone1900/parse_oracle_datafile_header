//要修改 block[4:8] 位置的内容并将修改后的数据刷新到物理文件上，
//你需要使用 os 包的 OpenFile 函数打开文件，
//并使用 WriteAt 方法将修改后的数据写入指定的偏移位置。以下是一个示例：
// oracle dbname 修改 ；改成SHITAN
package main

import (
	"flag"
	//"encoding/binary"
	"fmt"
	"os"
)

func main() {
// 假设有一个名为 filePath 的文件路径
//filePath := "C:\\Users\\lixora\\SYSTEM.dbf"
	// 定义一个命令行标志
	var filePath string
	flag.StringVar(&filePath, "dbfile", "", "指定数据库文件路径")

	// 解析命令行参数
	flag.Parse()

	// 检查是否提供了数据库文件路径
	if filePath == "" {
		fmt.Println("[Info]: 请提供数据库文件路径--Pls provide path of Oracle datafile ")
		os.Exit(1)
	}


	// Open the Oracle data file
	// local dev test demo
	//file, err := os.Open("C:\\Users\\ZMI\\Desktop\\asm-diskb\\system01_recovered.dbf")


// 打开文件，使用 os.O_RDWR 标志表示读写模式
file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
if err != nil {
fmt.Println("Error opening file:", err)
return
}

// 在 main 函数结束时关闭文件
defer file.Close()

// 从文件中读取原始数据
blockSize := 8
originalData := make([]byte, blockSize)
fmt.Printf("originalData  : %#08x\n",originalData)

// dbname 偏移位置block0+offset
offset1 := int64(8192+32)

_, err = file.ReadAt(originalData, offset1) // 从偏移位置4开始读取8个字节
if err != nil {
fmt.Println("Error reading original data:", err)
return
}

// 修改 block[0:7] 位置的内容
//newData := uint64(0x48454C4F57494E00)


// 字节反转  0x48454C4F57494E00  --> 0x004e49574f4c4548
//binary.LittleEndian.PutUint64(originalData[0:], newData)

//fmt.Printf("originalData2 : %#08x\n:",originalData)

// 一个包含中文字符的字符串
dbname := "lixora"
//fmt.Printf("dbname :%s\n",dbname)

// 将字符串转换为字节切片(隐藏)
//byteSlice := []byte(dbname)
//fmt.Printf("dbname_byteSlice : %#08x\n:",byteSlice)


// 创建8字节的切片并初始化为零值
data := make([]byte, 8)

// 将用户输入的字符串复制到切片中
copy(data, dbname)


// 将修改后的数据写入文件的相应位置
var nn int
nn, err = file.WriteAt(data, offset1) // 从偏移位置4开始写入8个字节
if err != nil {
fmt.Println("Error writing modified data:", err)
return
}
fmt.Println(nn)

// 刷新文件以确保数据被写入磁盘
err = file.Sync()
if err != nil {
fmt.Println("Error syncing file:", err)
return
}

fmt.Println("[info] : Data modification and file update successful. v2023-11-22-By huanglinjie 17767151782")
}
//在这个示例中，文件被以读写模式打开，并使用 ReadAt 方法读取原始数据，
//然后修改 block[4:8] 位置的内容，最后使用 WriteAt 方法将修改后的数据写入文件的相应位置。
//最后，通过调用 Sync 方法，确保修改被刷新到磁盘。请根据实际需要调整文件路径、偏移位置、原始数据和修改数据。

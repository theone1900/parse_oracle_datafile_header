package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// KCBlockStruct represents the structure of the Oracle block
/*
KCVFH means kernel cache recovery file header.
The length of datafile header for different version is different.

9i  kcvfh
10g kcvfh   676 bytes
11g kcvfh   860 bytes
12c kcvfh   1112 bytes

BBED> p kcvfh
struct kcvfh, 676 bytes @0
struct kcvfhbfh, 20 bytes @0
ub1 type_kcbh @0 0x0b -----数据的块类型 11可以看出是否是数据文件头
ub1 frmt_kcbh @1 0xa2 -----数据块的格式。1=oracle7 ,2=oracle8+
ub1 spare1_kcbh @2 0x00
ub1 spare2_kcbh @3 0x00
ub4 rdba_kcbh @4 0x00400001
ub4 bas_kcbh @8 0x00000000 ---SCN BASE
ub2 wrp_kcbh @12 0x0000 ---SCN WRAP
ub1 seq_kcbh @14 0x01 ---SCN***
ub1 flg_kcbh @15 0x04 (KCBHFCKV) ----块属性
ub2 chkval_kcbh @16 0x5064 ---检验值
ub2 spare3_kcbh @18 0x0000
struct kcvfhhdr, 76 bytes @20 ---此结构存储这个数据文件的属性
ub4 kccfhswv @20 0x00000000
ub4 kccfhcvn @24 0x0a200500 ---文件创建的版本号
ub4 kccfhdbi @28 0x783cfa8c ---数据库的DBID
text kccfhdbn[0] @32 Q ---所属实例的名字
text kccfhdbn[1] @33 X
text kccfhdbn[2] @34 P
text kccfhdbn[3] @35 T
text kccfhdbn[4] @36 F
text kccfhdbn[5] @37 H
text kccfhdbn[6] @38 0
text kccfhdbn[7] @39 1
ub4 kccfhcsq @40 0x00003db5 ---控制序列，控制文件事务会增加此值
ub4 kccfhfsz @44 0x0000f000 ---文件当前所包含数据块的个数
s_blkz kccfhbsz @48 0x00 ---文件存放的块大小，关闭数据库有值
ub2 kccfhfno @52 0x0001 ---文件号
ub2 kccfhtyp @54 0x0003 ---文件类型，03代表数据文件，06表示undo文件
ub4 kccfhacid @56 0x00000000 ---活动ID
ub4 kccfhcks @60 0x00000000 ---创建检查点的SCN
text kccfhtag[0] @64
text kccfhtag[1] @65
text kccfhtag[2] @66
text kccfhtag[3] @67
text kccfhtag[4] @68
text kccfhtag[5] @69
text kccfhtag[6] @70
text kccfhtag[7] @71
text kccfhtag[8] @72
text kccfhtag[9] @73
text kccfhtag[10] @74
text kccfhtag[11] @75
text kccfhtag[12] @76
text kccfhtag[13] @77
text kccfhtag[14] @78
text kccfhtag[15] @79
text kccfhtag[16] @80
text kccfhtag[17] @81
text kccfhtag[18] @82
text kccfhtag[19] @83
text kccfhtag[20] @84
text kccfhtag[21] @85
text kccfhtag[22] @86
text kccfhtag[23] @87
text kccfhtag[24] @88
text kccfhtag[25] @89
text kccfhtag[26] @90
text kccfhtag[27] @91
text kccfhtag[28] @92
text kccfhtag[29] @93
text kccfhtag[30] @94
text kccfhtag[31] @95
ub4 kcvfhrdb @96 0x00400179 ---ROOT DBA
struct kcvfhcrs, 8 bytes @100 ---文件创建的SCN
ub4 kscnbas @100 0x00000007 ---SCN BASE
ub2 kscnwrp @104 0x0000 ---SCN WRAP
ub4 kcvfhcrt @108 0x2ab9923a ---文件创建的时间戳
ub4 kcvfhrlc @112 0x30f3d1cf ---resetlogs的次数
struct kcvfhrls, 8 bytes @116 ---resetlogs的SCN
ub4 kscnbas @116 0x0005eca9 ---SCN BASE
ub2 kscnwrp @120 0x0000 ---SCN WRAP
ub4 kcvfhbti @124 0x00000000
struct kcvfhbsc, 8 bytes @128 ---备份的SCN
ub4 kscnbas @128 0x00000000 ---SCN BASE
ub2 kscnwrp @132 0x0000 ---SCN WRAP
ub2 kcvfhbth @136 0x0000
ub2 kcvfhsta @138 0x2004 (KCVFHOFZ) ---数据文件状态：04为正常，00为关闭，01为begin backup
struct kcvfhckp, 36 bytes @484 ---检查点checkpoint
struct kcvcpscn, 8 bytes @484 ---数据文件改变的检查点SCN
ub4 kscnbas @484 0x01a947ff --SCN BASE
ub2 kscnwrp @488 0x0000 --SCN WRAP
ub4 kcvcptim @492 0x338a07e7 --最后改变的时间
ub2 kcvcpthr @496 0x0001 --resetlogs的线程号
union u, 12 bytes @500
struct kcvcprba, 12 bytes @500
ub4 kcrbaseq @500 0x000005a0 --***
ub4 kcrbabno @504 0x00000002 --块号
ub2 kcrbabof @508 0x0010 --偏移量offset
ub1 kcvcpetb[0] @512 0x02 --最大线程数
ub1 kcvcpetb[1] @513 0x00
ub1 kcvcpetb[2] @514 0x00
ub1 kcvcpetb[3] @515 0x00
ub1 kcvcpetb[4] @516 0x00
ub1 kcvcpetb[5] @517 0x00
ub1 kcvcpetb[6] @518 0x00
ub1 kcvcpetb[7] @519 0x00
ub4 kcvfhcpc @140 0x00000619 --数据文件发生checkpoint的次数
ub4 kcvfhrts @144 0x3348a98a --resetlogs的次数
ub4 kcvfhccc @148 0x00000618 --控制文件记录的检查点次数
struct kcvfhbcp, 36 bytes @152
struct kcvcpscn, 8 bytes @152
ub4 kscnbas @152 0x00000000
ub2 kscnwrp @156 0x0000
ub4 kcvcptim @160 0x00000000
ub2 kcvcpthr @164 0x0000
union u, 12 bytes @168
struct kcvcprba, 12 bytes @168
ub4 kcrbaseq @168 0x00000000
ub4 kcrbabno @172 0x00000000
ub2 kcrbabof @176 0x0000
ub1 kcvcpetb[0] @180 0x00
ub1 kcvcpetb[1] @181 0x00
ub1 kcvcpetb[2] @182 0x00
ub1 kcvcpetb[3] @183 0x00
ub1 kcvcpetb[4] @184 0x00
ub1 kcvcpetb[5] @185 0x00
ub1 kcvcpetb[6] @186 0x00
ub1 kcvcpetb[7] @187 0x00
ub4 kcvfhbhz @312 0x00000000
struct kcvfhxcd, 16 bytes @316
ub4 space_kcvmxcd[0] @316 0x00000000
ub4 space_kcvmxcd[1] @320 0x00000000
ub4 space_kcvmxcd[2] @324 0x00000000
ub4 space_kcvmxcd[3] @328 0x00000000
word kcvfhtsn @332 0 --表空间号
ub2 kcvfhtln @336 0x0006
text kcvfhtnm[0] @338 S --表空间的名字，最长为30字符
text kcvfhtnm[1] @339 Y
text kcvfhtnm[2] @340 S
text kcvfhtnm[3] @341 T
text kcvfhtnm[4] @342 E
text kcvfhtnm[5] @343 M
text kcvfhtnm[6] @344
text kcvfhtnm[7] @345
text kcvfhtnm[8] @346
text kcvfhtnm[9] @347
text kcvfhtnm[10] @348
text kcvfhtnm[11] @349
text kcvfhtnm[12] @350
text kcvfhtnm[13] @351
text kcvfhtnm[14] @352
text kcvfhtnm[15] @353
text kcvfhtnm[16] @354
text kcvfhtnm[17] @355
text kcvfhtnm[18] @356
text kcvfhtnm[19] @357
text kcvfhtnm[20] @358
text kcvfhtnm[21] @359
text kcvfhtnm[22] @360
text kcvfhtnm[23] @361
text kcvfhtnm[24] @362
text kcvfhtnm[25] @363
text kcvfhtnm[26] @364
text kcvfhtnm[27] @365
text kcvfhtnm[28] @366
text kcvfhtnm[29] @367
ub4 kcvfhrfn @368 0x00000001 --相对文件号
struct kcvfhrfs, 8 bytes @372 --文件SCN
ub4 kscnbas @372 0x00000000 --SCN BASE
ub2 kscnwrp @376 0x0000 --SCN WRAP
ub4 kcvfhrft @380 0x00000000
struct kcvfhafs, 8 bytes @384 --绝对文件号
ub4 kscnbas @384 0x00000000 --SCN BASE
ub2 kscnwrp @388 0x0000 --SCN WRAP
ub4 kcvfhbbc @392 0x00000000
ub4 kcvfhncb @396 0x00000000
ub4 kcvfhmcb @400 0x00000000
ub4 kcvfhlcb @404 0x00000000
ub4 kcvfhbcs @408 0x00000000
ub2 kcvfhofb @412 0x000a
ub2 kcvfhnfb @414 0x000a
ub4 kcvfhprc @416 0x2ab99238 --上个resetlogs的次数
struct kcvfhprs, 8 bytes @420 --上个resetlogs的SCN
ub4 kscnbas @420 0x00000001
ub2 kscnwrp @424 0x0000
struct kcvfhprfs, 8 bytes @428
ub4 kscnbas @428 0x00000000
ub2 kscnwrp @432 0x0000
ub4 kcvfhtrt @444 0x00000000"
*/

type KCBlockStruct struct {
	TypeKCBH    byte
	FrmtKCBH    byte
	RDBAKCBH    uint32
	ChkvalKCBH  uint16
	KCCFHDBI    uint32
	KCCFHDBNX   []byte
	KCCFHCSQ    uint32
	KCCFHFSZ    uint32
	KCCFHFNO    uint16
	KCVFHRFN    uint32
	KCCFHTYP    uint16
	KCVFHRDB    uint32
	KSCNBAS     uint32
	KSCNWRP     uint16
	KCVFHCRT    uint32
	KCVFHRLC    uint32
	KCVFHRLS    struct {
		KSCNBAS uint32
		KSCNWRP uint16
	}
	KCVFHBSBSC  struct {
		KSCNBAS uint32
		KSCNWRP uint16
	}
	KCVFHSTA    uint16
	KCVFHCPC    uint32
	KCVFHCCC    uint32
	KCVFHTSN    uint32
	KCVFHTLN    uint16
	KCVFHTNM    []byte
	KCVFHPRC    uint32
	KCVFHPRS    struct {
		KSCNBAS uint32
		KSCNWRP uint16
	}
	KCVCPSCN     struct {
		KSCNBAS uint32
		KSCNWRP uint16
	}
	KCVCPTime   uint32
	KCVCPThr    uint16
	KCVCPRA     struct {
		KCRBASEQ uint32
		KCRBABNO uint32
		KCRBABOF uint32
	}
}

func main() {
	fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	fmt.Println("##powered by ：黄林杰_Huanglinjie\n##version : 2023-v11\n##联系方式：17767151782\n##blog: https://blog.csdn.net/lixora/\n##info: Oracle 11g datafile  header parse")
	fmt.Println("##demo : parseOracleKcvfh.exe -dbfile c:\\lixora.dbf")
	fmt.Println("##demo : parseOracleKcvfh.exe -dbfile C:\\system01.dbf -modify -bas 123456 -wrp 000\n")

	// 定义命令行标志
	var dbFilePath string
	var modifySCN bool
	var newKSCNBAS uint64
	var newKSCNWRP uint64
	
	flag.StringVar(&dbFilePath, "dbfile", "", "指定数据库文件路径")
	flag.BoolVar(&modifySCN, "modify", false, "是否修改SCN值")
	flag.Uint64Var(&newKSCNBAS, "bas", 0, "新的KSCNBAS值")
	flag.Uint64Var(&newKSCNWRP, "wrp", 0, "新的KSCNWRP值")

	// 解析命令行参数
	flag.Parse()

	// 检查是否提供了数据库文件路径
	if dbFilePath == "" {
		fmt.Println("[Info]: 请提供数据库文件路径--Pls provide path of Oracle datafile ")
		os.Exit(1)
	}

	// Open the Oracle data file
	file, err := os.OpenFile(dbFilePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the first block (assuming 8192 bytes)
	blockSize := 8192
	block := make([]byte, blockSize)
	blockNumber := 1

	// 计算块的偏移量
	blockOffset := blockSize * blockNumber

	// 移动文件指针到块的起始位置
	_, err = file.Seek(int64(blockOffset), io.SeekStart)
	if err != nil {
		fmt.Println("Error seeking to block:", err)
		return
	}

	_, err = file.Read(block)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// 如果需要修改SCN值
	if modifySCN {
		if newKSCNBAS > 0xFFFFFFFF || newKSCNWRP > 0xFFFF {
			fmt.Println("Error: KSCNBAS must be <= 0xFFFFFFFF and KSCNWRP must be <= 0xFFFF")
			return
		}
		
		// 修改SCN值
		binary.LittleEndian.PutUint32(block[484:488], uint32(newKSCNBAS))
		binary.LittleEndian.PutUint16(block[488:490], uint16(newKSCNWRP))
		
		// 写回文件
		_, err = file.Seek(int64(blockOffset), io.SeekStart)
		if err != nil {
			fmt.Println("Error seeking to block for write:", err)
			return
		}
		
		_, err = file.Write(block)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		
		fmt.Printf("Successfully modified SCN values: KSCNBAS=%d, KSCNWRP=%d\n", newKSCNBAS, newKSCNWRP)
	}

	// Parse the block using the defined structure
	kcBlock := parseKCBlock(block)

	// Print the extracted information
	fmt.Println("========================================> BLOCK SUMMARY <========================================")
	fmt.Printf("TypeKCBH: %#02x\n", kcBlock.TypeKCBH)
	fmt.Printf("FrmtKCBH: %#x\n", kcBlock.FrmtKCBH)
	fmt.Printf("RDBAKCBH: %#x\n", kcBlock.RDBAKCBH)
	fmt.Printf("KCCFHDBI: %d\n", kcBlock.KCCFHDBI)
	fmt.Printf("KCCFHDBNX: %s\n", strings.ReplaceAll(string(kcBlock.KCCFHDBNX), "\x00", ""))
	fmt.Printf("KCVCPSCN_KSCNBAS: %#08x,ckp scn-base:%d\n", kcBlock.KSCNBAS, kcBlock.KSCNBAS)
	fmt.Printf("KCVCPSCN_KSCNWRP: %#04x,ckp scn-wrap:%d\n", kcBlock.KSCNWRP,kcBlock.KSCNWRP)
	fmt.Printf("KCVFHSTA: %#x\n", kcBlock.KCVFHSTA)
	fmt.Printf("KCVFHTNM: %s\n", strings.ReplaceAll(string(kcBlock.KCVFHTNM), "\x00", ""))
	fmt.Printf("KCVFHTSN: %#d\n", kcBlock.KCVFHTSN)
	fmt.Printf("KCVFHCRT: %#x,CREATION TIME:%s\n", kcBlock.KCVFHCRT, time.Unix(int64(kcBlock.KCVFHCRT), 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("KCVFHCCC: %#08x,controlfile-chkpoint-count:%d\n", kcBlock.KCVFHCCC,kcBlock.KCVFHCCC)
	fmt.Printf("KCVFHCPC: %#08x,chkpoint-count:%d\n", kcBlock.KCVFHCPC,kcBlock.KCVFHCPC)
	fmt.Printf("KCVCPTime: %#x,chkpoint-TIME:%s\n", kcBlock.KCVCPTime, time.Unix(int64(kcBlock.KCVCPTime), 0).Format("2006-01-02 15:04:05"))

	hexNumber := fmt.Sprintf("%X", kcBlock.KCCFHFSZ)
	decimalNumber, err := strconv.ParseInt(hexNumber, 16, 64)
	fmt.Printf("KCCFHFSZ: %#x, FILE SIZE(bytes):%d\n", kcBlock.KCCFHFSZ, decimalNumber*int64(blockSize)+1*int64(blockSize))
	fmt.Printf("KCCFHFNO: %d\n", kcBlock.KCCFHFNO)
	fmt.Printf("KCVFHRFN: %d\n", kcBlock.KCVFHRFN)
}

func parseKCBlock(block []byte) KCBlockStruct {
	kcBlock := KCBlockStruct{
		TypeKCBH:   block[0],
		FrmtKCBH:   block[1],
		RDBAKCBH:   binary.LittleEndian.Uint32(block[4:8]),
		KCCFHDBI:   binary.LittleEndian.Uint32(block[28:32]),
		KCCFHDBNX:  block[32:40],
		KSCNBAS:    binary.LittleEndian.Uint32(block[484:488]),
		KSCNWRP:    binary.LittleEndian.Uint16(block[488:490]),
		KCVFHCRT:   binary.LittleEndian.Uint32(block[108:112]),
		KCVFHSTA:   binary.LittleEndian.Uint16(block[138:140]),
		KCVFHTNM:   block[338:368],
		KCVFHCCC:   binary.LittleEndian.Uint32(block[148:152]),
		KCVFHCPC:   binary.LittleEndian.Uint32(block[140:144]),
		KCCFHFSZ:   binary.LittleEndian.Uint32(block[44:48]),
		KCCFHFNO:   binary.LittleEndian.Uint16(block[52:54]),
		KCVFHRFN:   binary.LittleEndian.Uint32(block[368:372]),
		KCVFHTSN:   binary.LittleEndian.Uint32(block[332:336]),
		KCVCPTime:  binary.LittleEndian.Uint32(block[492:496]),
	}
	return kcBlock
}

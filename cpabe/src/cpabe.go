package src

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fentec-project/gofe/abe"
)

var inst = abe.NewFAME()

func Cpabe_setup() (*abe.FAMEPubKey, *abe.FAMESecKey) {
	// inst := abe.NewFAME()
	pubKey, secKey, err := abe.NewFAME().GenerateMasterKeys()
	if err != nil {
		panic(err)
	}
	return pubKey, secKey
}

func Gen_secKey(att []string, secKey *abe.FAMESecKey) *abe.FAMEAttribKeys {
	// inst := abe.NewFAME()
	// att := strings.Fields(arr)
	keys, err := inst.GenerateAttribKeys(att, secKey)
	// fmt.Println(keys)
	if err != nil {
		panic(err)
	}
	return keys
}

func Cpabe_Enc(filename string, policy string, pubKey *abe.FAMEPubKey) *abe.FAMECipher {
	content, err := os.ReadFile(filename)
	msp, err := abe.BooleanToMSP(string(policy), false)
	// fmt.Println(msp, '\n')
	if err != nil {
		panic(err)
	}
	cipher, err := inst.Encrypt(string(content), msp, pubKey)
	// fmt.Println(cipher)
	if err != nil {
		panic(err)
	}
	var str string
	str = cipher.CtPrime.String()

	// 将加密后数据写入cpabe文件

	file_cpabe_Path := filename + ".cpabe"
	file_cpabe, err := os.OpenFile(file_cpabe_Path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	file_cpabe.Write([]byte(str))
	file_cpabe.Close()

	// 删除初始文件
	del_origin := os.Remove(filename)
	if del_origin != nil {
		fmt.Println(del_origin)
	}
	return cipher
}

func Cpabe_Dec(file_cpabe_Path string, cipher *abe.FAMECipher, keys *abe.FAMEAttribKeys, pubKey *abe.FAMEPubKey) {
	msgCheck, err := inst.Decrypt(cipher, keys, pubKey)
	if err != nil {
		// panic(err)
		fmt.Println("无法解密")
	} else {
		fmt.Println("解密成功")
		//fmt.Println(msgCheck)
		// 删除cpabe文件

		del_cpabe := os.Remove(file_cpabe_Path)
		if del_cpabe != nil {
			fmt.Println(del_cpabe)
		}
		place := strings.Index(file_cpabe_Path, ".cpabe")
		file_output_path := file_cpabe_Path[:place]
		// 写入数据
		file_output, err := os.OpenFile(file_output_path, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("文件打开失败", err)
		}
		//及时关闭file句柄
		defer file_output.Close()
		//写入文件时，使用带缓存的 *Writer
		write1 := bufio.NewWriter(file_output)
		write1.Write([]byte(msgCheck))
		//Flush将缓存的文件真正写入到文件中
		write1.Flush()
	}

}

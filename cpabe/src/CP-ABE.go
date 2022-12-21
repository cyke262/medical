package src

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fentec-project/gofe/abe"
)

func main1() {

	// 首先输入想要加密的文件
	var file_origin string
	fmt.Println("输入您想加密的文件名称：")
	fmt.Scanln(&file_origin)
	content, err := os.ReadFile(file_origin)
	// fmt.Println(string(content))

	inst := abe.NewFAME()

	// 生成系统主密钥和公钥
	pubKey, secKey, err := inst.GenerateMasterKeys()
	// fmt.Println(pubKey, '\n')
	// fmt.Println(secKey, '\n')

	if err != nil {
		panic(err)
	}

	// 构造策略信息
	// "((0 AND 1) OR (2 AND 3)) AND 5",

	// 输入访问策略
	fmt.Println("请输入访问策略")
	var a string
	policy := bufio.NewScanner(os.Stdin)
	for policy.Scan() {
		a = policy.Text()
		//fmt.Println(a, '\n')
		break
	}
	// fmt.Scan(&policy)

	// fmt.Println(a, '\n')
	// msp, err := abe.BooleanToMSP("((医生1 AND 病种2 AND 科室2 AND 机构1) OR (研究员2 AND 病种1 AND 科室2 AND 机构3))", false)
	//((张三 AND 精神病 AND 神经科 AND 北大六院) OR (李四 AND 心脏病 AND 外科 AND 北大研究所))
	msp, err := abe.BooleanToMSP(string(a), false)
	// fmt.Println(msp, '\n')
	if err != nil {
		panic(err)
	}

	// 生成密文数据
	cipher, err := inst.Encrypt(string(content), msp, pubKey)
	// fmt.Println(cipher)
	if err != nil {
		panic(err)
	}
	// fmt.Println(cipher)
	var str string
	str = cipher.CtPrime.String()

	// 将加密后数据写入cpabe文件

	file_cpabe_Path := file_origin + ".cpabe"
	file_cpabe, err := os.OpenFile(file_cpabe_Path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	file_cpabe.Write([]byte(str))
	file_cpabe.Close()

	// 删除初始文件
	del_origin := os.Remove(file_origin)
	if del_origin != nil {
		fmt.Println(del_origin)
	}

	// 解密时构造 属性
	fmt.Println("请输入解密时的属性：")
	var s string
	sca := bufio.NewScanner(os.Stdin)
	for sca.Scan() {
		s = sca.Text()
		break
	}
	arr := strings.Fields(s)

	//gamma := []string{"医生1",  "病种2",  "科室2",  "机构1"}
	//fmt.Printf(gamma)

	keys, err := inst.GenerateAttribKeys(arr, secKey)
	// fmt.Println(keys)
	if err != nil {
		panic(err)
	}

	//解密
	msgCheck, err := inst.Decrypt(cipher, keys, pubKey)
	if err != nil {
		// panic(err)
		fmt.Println("无法解密")
	} else {
		fmt.Println("解密成功")
		//fmt.Println(msgCheck)
		// 删除cpabe文件
		defer file_cpabe.Close()
		del_cpabe := os.Remove(file_cpabe_Path)
		if del_cpabe != nil {
			fmt.Println(del_cpabe)
		}
		// 写入数据
		file_output_path := file_origin
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

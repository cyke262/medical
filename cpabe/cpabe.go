package main

import (
	"bufio"
	"fmt"
	"goproject/awesomeProject1/src"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// 读取文件中的Pubkey或Seckey
func ReadBybuffio(path string) string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("打开文件失败,err:%v\n", err)
		return "1"
	}
	defer file.Close() //关闭文件,为了避免文件泄露和忘记写关闭文件

	//使用buffio读取文件内容
	reader := bufio.NewReader(file) //创建新的读的对象
	var line string
	for {
		line, _ := reader.ReadString('\n')
		//fmt.Println(line)
		return line
	}
	return line
}

// 将Pubkey和Seckey写入文件中保存
func Write_PubSec() {
	pubkey, seckey := src.Cpabe_setup()
	pubkey2, seckey2 := src.Cpabe_setup()
	pub, _ := src.ToJson_Pub(pubkey)
	pub2, _ := src.ToJson_Pub(pubkey2)
	sec, _ := src.ToJson_Sec(seckey)
	sec2, _ := src.ToJson_Sec(seckey2)
	// pubkey1
	file, err := os.OpenFile("enc/pubkey_1.txt", os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()
	file.WriteString(pub)
	if err != nil {
		fmt.Println(err)
	}
	// pubkey2
	file2, err := os.OpenFile("enc/pubkey_2.txt", os.O_WRONLY|os.O_CREATE, 0644)
	defer file2.Close()
	file2.WriteString(pub2)
	if err != nil {
		fmt.Println(err)
	}
	// seckey1
	file3, err := os.OpenFile("enc/seckey_1.txt", os.O_WRONLY|os.O_CREATE, 0644)
	defer file3.Close()
	file3.WriteString(sec)
	if err != nil {
		fmt.Println(err)
	}
	// seckey2
	file4, err := os.OpenFile("enc/seckey_2.txt", os.O_WRONLY|os.O_CREATE, 0644)
	defer file4.Close()
	file4.WriteString(sec2)
	if err != nil {
		fmt.Println(err)
	}
}

// 通过编号获取对应编号组的秘钥
func Get_key(num string) (string, string) {
	//fmt.Println(num)
	return ReadBybuffio(fmt.Sprintf("enc/pubkey_%s.txt", num)), ReadBybuffio(fmt.Sprintf("enc/seckey_%s.txt", num))
}

// 生成1到2的随机数
func Gen_rand() string {
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(2) + 1
	key_num := strconv.Itoa(random)
	return key_num
}

// 将用户和该用户的秘钥编码存到文件中保存
// Todo: 存节点
func Write_user(user string) {
	file, err := os.OpenFile("enc/user_list.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644) // 追加模式写入文件
	defer file.Close()
	// 生成随机数
	key_num := Gen_rand()
	info := user + "," + key_num // 用户和编码中间逗号隔开
	fmt.Println(info)
	fmt.Fprintln(file, info)
	//file.WriteString(info)
	if err != nil {
		fmt.Println(err)
	}
}

// 读取用户和用户的秘钥编码文件，并匹配用户是否存在，若用户存在返回对应编码，若不存在则返回0
func Read_user(user string) string {
	//打开文件
	file, err := os.Open("enc/user_list.txt")
	if err != nil {
		fmt.Printf("打开文件失败,err:%v\n", err)
		return "0"
	}
	defer file.Close() //关闭文件,为了避免文件泄露和忘记写关闭文件

	//使用buffio读取文件内容
	reader := bufio.NewReader(file) //创建新的读的对象
	for {
		line, err := reader.ReadString('\n')
		//fmt.Println(line)
		place := strings.Index(line, ",")
		user_name := line[:place]
		num := line[place+1:]
		if user_name == user {
			fmt.Println("用户:", user_name, "的秘钥编码为 ", num)
			num = Trim(num)
			//fmt.Println(len(num))
			return num
		}
		if err == io.EOF { //如果读到末尾就会进入
			return "0"
		}
	}
}

func main() {
	// 加密的目标文件
	filename := "11.txt"
	// 第一次使用此函数，之后注释掉
	// Write_PubSec()

	// 用户名
	user := "张三"

	// 将用户信息保存
	// Write_user(user, key_num)
	key_num := Read_user(user)
	if key_num == "0" {
		fmt.Println("该用户不存在")
	}
	//fmt.Println("----------------------------------------------------------------")

	pub, sec := Get_key(key_num)
	pubkey, _ := src.ParseJson_Pub(pub)
	seckey, _ := src.ParseJson_Sec(sec)
	//fmt.Println(*&pubkey)
	//fmt.Println(*&seckey)

	//fmt.Println(pubkey, seckey, pubkey2, seckey2)
	gamma := []string{"医生1", "病种2", "科室2", "机构1"}

	keys := src.Gen_secKey(gamma, seckey)
	policy := "((医生1 AND 病种2) OR (研究员2 AND 病种1))"

	cipher := src.Cpabe_Enc(filename, policy, pubkey)
	file_cpabe_Path := filename + ".cpabe"
	time.Sleep(8 * time.Second)
	src.Cpabe_Dec(file_cpabe_Path, cipher, keys, pubkey)
	fmt.Println("成功")
}

// 删除字符串中的换行和空格
func Trim(src string) (dist string) {
	if len(src) == 0 {
		return
	}
	r, distR := []rune(src), []rune{}
	for i := 0; i < len(r); i++ {
		if r[i] == 10 || r[i] == 32 {
			continue
		}
		distR = append(distR, r[i])
	}
	dist = string(distR)
	return
}

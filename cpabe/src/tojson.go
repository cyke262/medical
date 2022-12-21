package src

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/fentec-project/gofe/abe"
)

var parseJsonError = errors.New("json parse error")
var toJsonError = errors.New("to json error")

func ToJson_Pub(c *abe.FAMEPubKey) (string, error) {
	//fmt.Printf("原始结构体: %v\n", c)
	if jsonStr, err := json.Marshal(c); err != nil {
		fmt.Println("Error =", err)
		return "", parseJsonError
	} else {
		return string(jsonStr), nil
	}
}

func ParseJson_Pub(a string) (*abe.FAMEPubKey, error) {
	//fmt.Printf("原始字符串: %s\n", a)
	var c *abe.FAMEPubKey
	if err := json.Unmarshal([]byte(a), &c); err != nil {
		fmt.Println("Error =", err)
		return c, parseJsonError
	}
	//fmt.Println(c)
	return c, nil
}

func ToJson_Sec(c *abe.FAMESecKey) (string, error) {
	//fmt.Printf("原始结构体: %v\n", c)
	if jsonStr, err := json.Marshal(c); err != nil {
		fmt.Println("Error =", err)
		return "", parseJsonError
	} else {
		return string(jsonStr), nil
	}
}

func ParseJson_Sec(a string) (*abe.FAMESecKey, error) {
	//fmt.Printf("原始字符串: %s\n", a)
	var c *abe.FAMESecKey
	if err := json.Unmarshal([]byte(a), &c); err != nil {
		fmt.Println("Error =", err)
		return c, parseJsonError
	}
	//fmt.Println(c)
	return c, nil
}

func main3() {
	fmt.Println("hello world")
	file, err := os.Create("11.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//根据file创建内存文件句柄
	write := bufio.NewWriter(file)
	defer write.Flush()
	pubKey, secKey, err := abe.NewFAME().GenerateMasterKeys()
	pub, err := ToJson_Pub(pubKey)
	fmt.Println(secKey)
	fmt.Println("----------------------------------------------------------------------------------")
	fmt.Println(pubKey.PartG2)
	fmt.Println("----------------------------------------------------------------------------------")
	pub1, err := ParseJson_Pub(pub)
	fmt.Println(pub1.PartG2)
	if pub1.PartG2 == pubKey.PartG2 {
		fmt.Println("success")
	} else {
		fmt.Println("Error")
	}

	//将内容写到内存
	fmt.Fprintln(write, pubKey)
	//fmt.Fprintln(write, secKey)

}

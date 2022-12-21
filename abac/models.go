package abac

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type Date2DB struct {
	Groups       string
	SubjectMark  string
	Diseases     string
	Researcher   string
	Organization string
}

// 主体
type Sub struct {
	UID        string `json:"uid"`        // 作为某个人的外键，进而可以指到某个具体的人，和那个样本标识符是一个样的 具有唯一性
	Department string `json:"department"` //TODO： 目前可设置 为空，后期可能会删除。
	Role       string `json:"role"`       // role of the subject, 可取值为："Adminstor"，“u1”，“u2”，“u3”
	Group      string `json:"group"`      // 疾病种类，也就算研究方向,eg：精神病 or 心脏病；
	Org        string `json:"org"`        //隶属机构, e.g.: 北大六院
}

// sub转码为json
func (p *Sub) ToBytes() []byte {
	b, err := json.Marshal(*p)
	if err != nil {
		fmt.Println("Obj转码json字符串错误: ", err.Error())
		return nil
	}
	return b
}

// 客体
type Obj struct {
	Caseumber string `json:"casenumber"` // ID of the object, e.g. E01,这个应该是系统生成，而不是用户赋值
	Owner     string `json:"owner"`      // Owner of the object, e.g. the EHR E01 belongs to the patient P01.
	Time      string `json:"time"`       // 时间
}

// Obj转码为json
func (p *Obj) ToBytes() []byte {
	b, err := json.Marshal(*p)
	if err != nil {
		fmt.Println("Obj转码json字符串错误: ", err.Error())
		return nil
	}
	return b
}

// []byte => obj
func NewResource(b []byte) (Obj, error) {
	r := Obj{}
	err := json.Unmarshal(b, &r)
	return r, err
}

// 操作
//  TODO: 目前设计为单操作，所以没有设计结构体

// 环境
type Env struct {
	AllowOrg    string `json:"allowOrg"`
	CreatedTime string `json:"createdTime"`
	EndTime     string `json:"endTime"` // 代表有效期
}

// 访问策略
type Policy struct {
	Obj      string `json:"obj"`   //仅用写id
	Owner    string `json:"owner"` //
	Env      Env
	SubRules []string `json:"subRules"`
}

// 请求
type ABACRequest struct {
	Sub     string `json:"sub"` //仅用写id
	Obj     string `json:"obj"` //仅用写id
	CurTime int64  `json:"curTime"`
	Op      string `json:"op"` //
}

// Policy 的ID生成
func (p *Policy) GetID() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(p.Obj)))
}

// Policy转码为json
func (p *Policy) ToBytes() []byte {
	b, err := json.Marshal(*p)
	if err != nil {
		fmt.Println("Policy转码json字符串错误: ", err.Error())
		return nil
	}
	return b
}

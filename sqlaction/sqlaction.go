package sqlaction

// 参考https://blog.csdn.net/naiwenw/article/details/79281220
// https://blog.csdn.net/rockage/article/details/103776251

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //初始化
	_ "hash/fnv"
	"math/rand"
	_ "math/rand"
	"strings"
	"time"
)

//写到服务层里面，所有字段用于和数据库对接（服务层调用数据库和链码）
type MedicalRecord struct {
	Groups         string `json:"groups"`         //分组
	SubjectMark    string `json:"subjectMark"`    //样本标识符
	Name           string `json:"name"`           //姓名
	NameInitials   string `json:"nameInitials"`   //姓名缩写
	CaseNumber     string `json:"caseNumber"`     //病例号，并上链
	Sex            string `json:"sex"`            //性别
	Nation         string `json:"nation"`         //民族
	Diseases       string `json:"diseases"`       //疾病种类
	MedicalHistory string `json:"medicalHistory"` //现病史
	NativePlace    string `json:"nativePlace"`    //现住地
	Diagnose       string `json:"diagnose"`       //诊断
	Researcher     string `json:"researcher"`     //研究者，并上链
	Organization   string `json:"organization"`   //机构，并上链
	Addition1      string `json:"addition1"`      //预留信息
	Addition2      string `json:"addition2"`      //预留信息
	Addition3      string `json:"addition3"`      //预留信息
	Status         string `json:"status"`         //状态
	EntryTime      string `json:"entryTime"`      //入组时间
	BaseTime       string `json:"baseTime"`       //基准时间
	GatherTime     string `json:"gatherTime"`     //采集时间-1，并上链
}

// 数据本身分为3类，第一类是无关数据，第二类是属性数据，第三类是数据索引——上链，
// 无关数据可以都为nil，但属性数据必须给，并且应该要有一个表
// 存数据索引与案例——这里我没涉及到json解析，字段可以先从你那个结构体中直接获取
// 属性数据
type Date2DB struct {
	Groups       string
	SubjectMark  string
	Diseases     string
	Researcher   string
	Organization string
}

// 无关数据且不被直接引用
type DateInfo struct {
	name           string
	nameInitials   string
	sex            string
	nation         string
	medicalHistory string
	nativePlace    string
	diagnose       string
	addition1      string
	addition2      string
	addition3      string
	status         string
}

// 时间数据
type DateTimeInfo struct {
	EntryTime string
	BaseTime  string
}

// 最终医疗数据结构体设计
type MedicalDate struct {
	Attr       Date2DB
	info       DateInfo
	DateTime   DateTimeInfo
	CaseNumber string
	GatherTime string
}

//数据库配置
const (
	userName = "root"
	password = "root"
	ip       = "localhost"
	port     = "3306"
	dbName   = "itbtsql"
)

// 全局变量 sqlmap存表内字段
var SqlMap map[int]string
var SqlMapResearcher map[int]string
var SqlMapOrganization map[int]string
var SqlMapDiseases map[int]string

func InitSql() {
	// 配置sql关键字
	SqlMap = make(map[int]string)
	SqlMap[0] = "base_info" //table name
	SqlMap[1] = "_SubjectMark"
	SqlMap[2] = "_Researcher"
	SqlMap[3] = "_Organization"
	SqlMap[4] = "_Diseases"
	SqlMap[5] = "_CaseNumber"
	SqlMap[6] = "_GatherTime"
	SqlMap[7] = "_Groups"
	// role
	SqlMapResearcher = make(map[int]string)
	SqlMapResearcher[0] = "admin"
	SqlMapResearcher[1] = "u1"
	SqlMapResearcher[2] = "u2"
	SqlMapResearcher[3] = "u3"
	// organization
	SqlMapOrganization = make(map[int]string)
	SqlMapOrganization[0] = "卫生监督管理局"
	SqlMapOrganization[1] = "北大六院"
	SqlMapOrganization[2] = "北京天坛医院"
	SqlMapOrganization[3] = "北京大学第一医院"
	SqlMapOrganization[4] = "天津医科大学附属医院"
	SqlMapOrganization[5] = "山东大学齐鲁医院"
	// organization
	SqlMapDiseases = make(map[int]string)
	SqlMapDiseases[0] = "未知"
	SqlMapDiseases[1] = "神经系统疾病"
	SqlMapDiseases[2] = "内分泌、营养及代谢疾病"
	SqlMapDiseases[3] = "眼和附器疾病"
	SqlMapDiseases[4] = "鼻咽喉疾病"
	SqlMapDiseases[5] = "呼吸系统疾病"
	SqlMapDiseases[6] = "循环系统疾病"
	SqlMapDiseases[7] = "血液、造血器官、免疫系统疾病"
	SqlMapDiseases[8] = "消化系统疾病"
	SqlMapDiseases[9] = "肾脏和泌尿道疾病"
	SqlMapDiseases[10] = "肌肉骨骼系统和结缔组织疾病"
	SqlMapDiseases[11] = "皮肤、皮下组织、乳腺疾病和烧伤"
}

//Db数据库连接池
var DB *sql.DB

//注意方法名大写，就是public
func InitDB() bool {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := userName + ":" + password + "@tcp(" + ip + ":" + port + ")/" + dbName
	fmt.Println(path)
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)
	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		fmt.Println("opon database fail")
		return false
	}
	fmt.Println("connnect success")
	return true
}

// 用于生成sql查询语句
func SetSQL(sqlmap interface{}, queryString string) string {
	temp := sqlmap.(map[int]string)
	for k := range temp {
		fmt.Println(k, temp[k])
	}
	// "SELECT _SubjectMark From baseinfo WHERE _Groups='喹硫平'"
	// 可自行更改字段
	SQLString := "SELECT " + temp[1] + " From " + temp[0] + " WHERE " + temp[2] + "='" + queryString + "'"
	fmt.Println(SQLString)
	return SQLString
}

// 查询并返回数据
func queryDB(SQL string) map[int]string {
	var subId string
	rows, err := DB.Query(SQL)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(subId)
	data := make(map[int]string)
	index := 0
	for rows.Next() {
		rows.Scan(&subId)
		data[index] = subId
		index = index + 1
		//fmt.Println(subId)
	}
	defer rows.Close()
	return data
}

// 定义返回的时间格式
var returnTimeType string

func ReturnTime(t int32) string {
	if t == 1 {
		returnTimeType = "2006-01-02"
	} else if t == 2 {
		returnTimeType = "2006-01-02 15:03:04"
	}
	now := time.Now()
	seconds := time.Unix(now.Unix(), 0)
	timeString := seconds.Format(returnTimeType)
	// 确定输出格式为string
	// fmt.Printf("%T\n", timeString)
	// 输出 时间戳，完整时间，最后返回内容
	// fmt.Println(now.Unix(), seconds, timeString)
	return timeString
}

// 参考https://blog.csdn.net/chenxing1230/article/details/83784063
func HashSHA256(a Date2DB, t string) string {
	baseString := []string{a.Groups, a.SubjectMark, a.Diseases, a.Researcher, a.Organization, t}
	//fmt.Println(baseString)
	str := strings.Join(baseString, "")
	//创建一个基于SHA256算法的hash.Hash接口的对象
	hash := sha256.New()
	//输入数据
	hash.Write([]byte(str))
	//计算哈希值
	bytes := hash.Sum(nil)
	//将字符串编码为16进制格式,返回字符串
	hashCode := hex.EncodeToString(bytes)
	//返回哈希值
	//fmt.Println(hashCode)
	return hashCode
}

func GetUserLogin(DB *sql.DB, SQLString string) string {
	var str string
	//SQLString2 := "select username from login where state='1'"
	err := DB.QueryRow(SQLString).Scan(&str)
	if err == sql.ErrNoRows { //没有结果
		fmt.Println("err")
		return ""
	} else {
		return str
	}
}

func GetCaseNumber(data []string) string {
	attr := Date2DB{
		Groups:       data[0],
		SubjectMark:  data[1],
		Diseases:     data[7],
		Researcher:   data[11],
		Organization: data[12],
	}
	t := ReturnTime(2) //含时分秒
	//t := "2022-12-18 14:02:35"      //8244d8d6ccc1fac9e9333b3466571efef3678a691e69f3f2a0b972dfaf1091b2
	id := HashSHA256(attr, t)[0:13] //按 六院 数据元 长度为13位
	return id
}

func GenerateDate(data []string) MedicalDate {
	attr := Date2DB{
		Groups:       data[0],
		SubjectMark:  data[1],
		Diseases:     data[7],
		Researcher:   data[11],
		Organization: data[12],
	}
	t := ReturnTime(2) //含时分秒
	//t := "2022-12-18 14:02:35"      //8244d8d6ccc1fac9e9333b3466571efef3678a691e69f3f2a0b972dfaf1091b2
	id := HashSHA256(attr, t)[0:13] //按 六院 数据元 长度为13位
	fmt.Printf("%T\n", id)
	temp := MedicalDate{
		Attr: attr,
		info: DateInfo{
			name: "name", nameInitials: "nameInitials", sex: "sex",
			nation: "nation", nativePlace: "nativePlace",
			status: "status",
		},
		DateTime: DateTimeInfo{
			EntryTime: ReturnTime(1), BaseTime: ReturnTime(1),
		},
		CaseNumber: id,
		GatherTime: t,
	}
	fmt.Println(temp)
	return temp
}

// 插入数据——与数据上传对应
func InsertDB2BaseInfo(data MedicalDate) bool {
	InitSql()
	temp := data
	SQLString := "insert into base_info(_Groups,_SubjectMark,_CaseNumber,_GatherTime)values(?, ?, ?,?)"
	r, err := DB.Exec(SQLString, temp.Attr.Groups, temp.Attr.SubjectMark, temp.CaseNumber, temp.GatherTime)
	if err != nil {
		fmt.Println("插入失败", err)
		return false
	}
	id, err := r.LastInsertId()
	if err != nil {
		fmt.Println("exec failed", err)
		return false
	}
	fmt.Println("插入成功", id)
	return true
}

// 插入数据——填充 机构与案例的 表 insti_coop
func InsertDB2Insti() bool {
	InitSql()
	SQLString := "select _CaseNumber from base_info"
	caseNumber := make(map[int]string)
	caseNumber = queryDB(SQLString)
	for i := 0; i < len(caseNumber); i++ {
		fmt.Println(caseNumber[i])
		//SQLStringFromCase := "select _Organization from base_info where _CaseNumber='" + caseNumber[i] + "'"
		//var organization string
		//rows := DB.QueryRow(SQLStringFromCase)
		//rows.Scan(&organization)
		//fmt.Println(organization)
		result := make(map[int]int)
		result = ChooseRand()
		fmt.Println(result[0], result[1])
		organization := SqlMapOrganization[result[1]]
		fmt.Println(caseNumber[i], organization)
		var str string
		SQLString2 := "select * from insti_coop where _CaseNumber='" + caseNumber[i] + "' and insti_name='" + organization + "'"
		err := DB.QueryRow(SQLString2).Scan(&str)
		if err == sql.ErrNoRows { //没有结果
			SQLString3 := "insert into insti_coop(_CaseNumber,insti_name)values(?,?)"
			_, err := DB.Exec(SQLString3, caseNumber[i], organization)
			if err != nil {
				fmt.Println("err")
			}
		}
	}

	return true
}

func ChooseRand() map[int]int {
	result := make(map[int]int)
	result1 := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(2)
	if result1 == 0 {
		result1 = 1
	}
	//fmt.Println(result1)
	result2 := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(6)
	if result2 == 0 {
		result2 = 1
	}
	//if result1 == 1 {
	//	result2 = 1
	//}
	//fmt.Println(result2)
	result[0] = result1
	result[1] = result2
	return result
}
func UpdateDB(data map[int]string) {
	InitSql()
	result := make(map[int]int)
	gatherTIme := make(map[int]string)
	for i := 0; i < len(data); i++ {
		fmt.Println(i, data[i])
		result = ChooseRand()
		gatherTIme = queryDB("select _GatherTime,_Groups,_Researcher FROM base_info where _SubjectMark='" + data[i] + "'")
		Organization := SqlMapOrganization[result[0]]
		Diseases := SqlMapDiseases[result[1]]
		attr := Date2DB{
			Groups:       gatherTIme[1],
			SubjectMark:  data[i],
			Diseases:     Diseases,
			Researcher:   gatherTIme[2],
			Organization: Organization,
		}
		CaseNumber := HashSHA256(attr, gatherTIme[0])[0:13]
		sql := "UPDATE base_info SET _Organization= '" + Organization + "', _Diseases='" + Diseases + "', _CaseNumber='" + CaseNumber + "' WHERE _SubjectMark='" + data[i] + "'"
		fmt.Println(sql)
		_, err := DB.Exec(sql)
		if err != nil {
			fmt.Println("err")
		}
	}
}

//func main() {
//
//	InitDB()
//	//InsertDB2Insti()
//	SQLString := "select _CaseNumber from base_info"
//	caseNumber := make(map[int]string)
//	caseNumber = queryDB(SQLString)
//	GeneratePolicy(caseNumber)
//	//InsertDB2Insti()
//}

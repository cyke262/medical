package service

import (
	"database/sql"
	"fmt"
	"medical/sdkInit"
	"time"

	_ "github.com/go-sql-driver/mysql" //初始化
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// 写到服务层里面，所有字段用于和数据库对接（服务层调用数据库和链码）
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

// 操作记录，记录在区块链上不可被修改
type OperationRecord struct {
	OperationRecordID string `json:"operationRecordId"` //操作记录ID
	UserID            string `json:"userId"`            //操作者用户
	OrganizationID    string `json:"organizationID"`    //用户组织ID
	DataType          string `json:"dataType"`          //记录类型：操作记录operation
	ActionType        string `json:"actionType"`        //操作类型：访问、修改、删除（Access、Update、Delete）
	MedicalRecordID   string `json:"medicalRecordId"`   //涉及到的医疗记录ID
	OriginalAuthorID  string `json:"originalAuthorId"`  //涉及到的医疗记录创建者
	PatientID         string `json:"patientId"`         //医疗记录涉及到的患者ID
	EntryMethod       string `json:"entryMethod"`       //操作记录的输入方式：自动输入
	IsSuccess         bool   `json:"isSuccess"`         //*新增：操作是否成功，成功为True
	Time              string `json:"time"`              //时间
}

type OperationRecordArr struct {
	OperationRecord []OperationRecord
}

// 审计记录，记录在区块链上不可被修改
type AuditRecord struct {
	AuditRecordID  string `json:"auditRecordId"`  //审计记录ID
	UserID         string `json:"userId"`         //操作者用户（此处为管理员）
	OrganizationID string `json:"organizationID"` //用户组织ID
	DataType       string `json:"dataType"`       //记录类型：审计Audit
	EntryMethod    string `json:"entryMethod"`    //审计记录的输入方式：自动输入
	Time           string `json:"time"`           //时间
}

// *新增：审计报告，记录组织的失败率信息
type AuditReport struct {
	TargetOrg       string     `json:"targetOrg"`       //被审计组织ID
	CurrentCredit   float64    `json:"currentCredit"`   //组织现在的信誉值
	PreviousCredit  float64    `json:"previousCredit"`  //组织之前的信誉值
	CreditChange    string     `json:"creditChange"`    //组织信誉值变动情况：上升、下降、不变
	ReferenceRange  [2]float64 `json:"referenceRange"`  //参照区间
	TotalOperations int64      `json:"totalOperations"` //组织总操作次数
	FailOperations  int64      `json:"failOperations"`  //组织失败操作次数
	FailRate        float64    `json:"failRate"`        //组织失败操作率
	MaxFailRateUser string     `json:"maxFailRateUser"` //组织中失败率最高的用户ID
	MaxFailRate     float64    `json:"maxFailRate"`     //组织所有用户中最高的失败率
}

type ServiceSetup struct {
	ChaincodeID string
	Client      *channel.Client
}

// 数据库配置
const (
	userName = "root"
	password = "2001"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "itbtsql"
)

// Db数据库连接池
var DB *sql.DB

func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}

func InitService(chaincodeID, channelID string, org *sdkInit.OrgInfo, sdk *fabsdk.FabricSDK) (*ServiceSetup, error) {
	handler := &ServiceSetup{
		ChaincodeID: chaincodeID,
	}
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org.OrgUser), fabsdk.WithOrg(org.OrgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new channel client: %s", err)
	}
	handler.Client = client
	return handler, nil
}

// 注意方法名大写，就是public
func InitDB() bool {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := userName + ":" + password + "@tcp(" + ip + ":" + port + ")/" + dbName + "?allowNativePasswords=true"
	fmt.Println(path)
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)
	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail")
		fmt.Println(err)
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
	sql := "SELECT " + temp[1] + " From " + temp[0] + " WHERE " + temp[2] + "='" + queryString + "'"
	fmt.Println(sql)
	return sql
}

// 查询并返回数据
func QueryDB(Sql string) map[int]string {
	var subId string
	rows, err := DB.Query(Sql)
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
		fmt.Println(subId)
	}
	defer rows.Close()
	return data
}

func InsertDB(data []string) bool {
	sqlString := "insert into base_info(_Groups,_SubjectMark,_Name,_NameInitials,_CaseNumber,_Sex,_Nation,_Diseases,_MedicalHistory,_NativePlace,_Diagnose,_Researcher,_Organization,_Addition1,_Addition2,_Addition3,_Status,_EntryTime,_BaseTime,_GatherTime) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	r, err0 := DB.Exec(sqlString, data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7], data[8], data[9], data[10], data[11], data[12], data[13], data[14], data[15], data[16], time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02 15:04:05"))
	if err0 != nil {
		fmt.Println("插入失败：", err0)
		return false
	}
	id, err1 := r.LastInsertId()
	if err1 != nil {
		fmt.Println("操作失败：", err1)
		return false
	}
	fmt.Println("插入成功：", id)
	//受影响的行数
	row_affect, err2 := r.RowsAffected()
	if err2 != nil {
		fmt.Println("受影响行数获取失败:", err2)
		return false
	}
	fmt.Println("受影响的行数：", row_affect)
	return true
}

func SelectDBSingle(data []string) *MedicalRecord {
	sqlString := "select * from base_info where _CaseNumber = ?"
	stmt, _ := DB.Prepare(sqlString)
	row := stmt.QueryRow(data[3])
	if row == nil {
		fmt.Println("未获取到该记录:", data[3])
		return nil
	}
	m := &MedicalRecord{}
	err0 := row.Scan(&m.Groups, &m.SubjectMark, &m.Name, &m.NameInitials, &m.CaseNumber, &m.Sex, &m.Nation, &m.Diseases, &m.MedicalHistory, &m.NativePlace, &m.Diagnose, &m.Researcher, &m.Organization, &m.Addition1, &m.Addition2, &m.Addition3, &m.Status, &m.EntryTime, &m.BaseTime, &m.GatherTime)
	if err0 != nil {
		return nil
	}
	return m
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

// 返回hash——这个先放弃
/*
func hash(a Date2DB, t string) uint32 {
	baseString := []string{a.Groups, a.SubjectMark, a.Diseases, a.Researcher, a.Organization, ReturnTime(2)}
	s := strings.Join(baseString, "")
	h := fnv.New32a()
	h.Write([]byte(s))
	fmt.Println(h.Sum32())
	return h.Sum32()
}
*/

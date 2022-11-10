package service

import (
	"fmt"
	"medical/sdkInit"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// 医疗记录
type MedicalRecord struct {
	MedicalRecordID string `json:"medicalRecordId"` //医疗记录ID
	UserID          string `json:"userId"`          //医疗记录创建者用户
	PatientID       string `json:"patientId"`       //医疗记录涉及到的患者ID
	OrganizationID  string `json:"organizationID"`  //用户组织ID
	DataType        string `json:"dataType"`        //记录类型：医疗记录medical
	DataField       string `json:"dataField"`       //医疗记录的数据：比如患者的健康状况
	Data            string `json:"data"`            //医疗记录的数据值，可能被修改
	EntryMethod     string `json:"entryMethod"`     //医疗记录的输入方式：手动输入或自动输入
	Time            string `json:"time"`            //时间

	Historys []HistoryItem // 当前的历史记录
}

type HistoryItem struct {
	TxId          string
	MedicalRecord MedicalRecord
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

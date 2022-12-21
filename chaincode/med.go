package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

const (
	MedicalRecordString   = "medical-record-key"
	OperationRecordString = "operation-record-key"
	AuditRecordString     = "audit-record-key"
)

//写到链码里，只保留上链字段
type MedicalRecord struct {
	CaseNumber      string `json:"caseNumber"`      //病例号，上链
	Researcher      string `json:"researcher"`      //研究者，上链
	Organization    string `json:"organization"`    //机构，上链
	OperationType   string `json:"operationType"`   //操作类型
	OperationResult string `json:"operationResult"` //操作结果
	GatherTime      string `json:"gatherTime"`      //采集时间-1，上链
}
type SmartContract struct {
}

func (t *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println(" ==== Init ====")

	return shim.Success(nil)
}

func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// 获取用户意图
	fun, args := stub.GetFunctionAndParameters()

	if fun == "UploadMedicalRecord" {
		return t.UploadMedicalRecord(stub, args)
	} else if fun == "AccessMedicalRecord" {
		return t.AccessMedicalRecord(stub, args)
	} else if fun == "DeleteMedicalRecord" {
		return t.DeleteMedicalRecord(stub, args)
	} else if fun == "UpdateMedicalRecord" {
		return t.UpdateMedicalRecord(stub, args)
	}

	return shim.Error("指定的函数名称错误")

}

// 上传医疗记录
func (t *SmartContract) UploadMedicalRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	r, _ := stub.GetState(args[0])
	if r != nil {
		return shim.Error("医疗记录ID已经存在！")
	}
	m := MedicalRecord{
		CaseNumber:      args[0],
		Researcher:      args[1],
		Organization:    args[2],
		OperationType:   "Medical",
		OperationResult: args[3],
		GatherTime:      time.Now().Format("2006-01-02 15:04:05"),
	}
	b, err1 := json.Marshal(m)
	fmt.Println(err1)
	if err1 != nil {
		return shim.Error(err1.Error())
	}
	err2 := stub.PutState(m.CaseNumber, b)
	fmt.Println(err2)
	if err2 != nil {
		return shim.Error(err2.Error())
	}
	err3 := stub.SetEvent(args[4], []byte{})
	fmt.Println(err3)
	if err3 != nil {
		return shim.Error(err3.Error())
	}
	return shim.Success([]byte("医疗记录添加成功！"))
}

//查看医疗记录
func (t *SmartContract) AccessMedicalRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	result0, err0 := stub.GetState(args[3])
	if err0 != nil {
		return shim.Error("医疗记录获取失败！")
	}
	if result0 == nil {
		return shim.Error("医疗记录ID获取失败！")
	}
	err1 := stub.SetEvent(args[4], []byte{})
	if err1 != nil {
		return shim.Error(err1.Error())
	}
	return shim.Success(result0)
}

//查看医疗记录
func (t *SmartContract) DeleteMedicalRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	err0 := stub.DelState(args[3])
	if err0 != nil {
		return shim.Error("医疗记录删除失败！")
	}
	err1 := stub.SetEvent(args[4], []byte{})
	if err1 != nil {
		return shim.Error(err1.Error())
	}
	return shim.Success([]byte("医疗记录删除成功！"))
}

// 对医疗记录进行修改
func (t *SmartContract) UpdateMedicalRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	result0, err0 := stub.GetState(args[0])
	if err0 != nil {
		return shim.Error("医疗记录获取失败！")
	}
	if result0 == nil {
		return shim.Error("医疗记录ID获取失败！")
	}
	var m MedicalRecord
	err1 := json.Unmarshal(result0, &m)
	if err1 != nil {
		return shim.Error("反序列化医疗记录失败！")
	}
	m.Researcher = args[1]
	m.Organization = args[2]
	m.OperationResult = args[3]
	newMedicalRecordAsBytes, err2 := json.Marshal(m)
	if err2 != nil {
		return shim.Error("医疗记录加工失败！")
	}
	err3 := stub.PutState(m.CaseNumber, newMedicalRecordAsBytes)
	if err3 != nil {
		return shim.Error("医疗记录修改失败！")
	}
	err4 := stub.SetEvent(args[4], []byte{})
	if err4 != nil {
		return shim.Error(err4.Error())
	}
	return shim.Success([]byte("医疗记录修改成功！"))
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("启动链码时发生错误: %s", err)
	}
}

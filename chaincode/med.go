package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

const (
	MedicalRecordString = "medical-record-key"
	OperationRecordString = "operation-record-key"
	AuditRecordString = "audit-record-key"
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
	IsSuccess         bool   `json:"isSuccess"`         //[fxb] 操作是否成功
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

type SmartContract struct {
}

func (t *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println(" ==== Init ====")

	return shim.Success(nil)
}

func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// 获取用户意图
	fun, args := stub.GetFunctionAndParameters()

	if fun == "uploadMedicalRecord" {
		return t.uploadMedicalRecord(stub, args)
	} else if fun == "operateMedicalRecord" {
		return t.operateMedicalRecord(stub, args)
	} else if fun == "updateMedicalRecord" {
		return t.updateMedicalRecord(stub, args)
	} else if fun == "deleteMedicalRecord" {
		return t.deleteMedicalRecord(stub, args)
	} else if fun == "auditForAllLogs" {
		return t.auditForAllLogs(stub, args)
	} else if fun == "auditForTimeRange" {
		return t.auditForTimeRange(stub, args)
	} else if fun == "auditForUser" {
		return t.auditForUser(stub, args)
	} else if fun == "auditForOrganisation" {
		return t.auditForOrganisation(stub, args)
	} else if fun == "auditForMedicalRecord" {
		return t.auditForMedicalRecord(stub, args)
	} else if fun == "auditForOriginalAuthor" {
		return t.auditForOriginalAuthor(stub, args)
	} else if fun == "auditForPatient" {
		return t.auditForPatient(stub, args)
	} else if fun == "getMedicalRecordHistory" {
		return t.getMedicalRecordHistory(stub, args)
	}

	return shim.Error("指定的函数名称错误")

}

// 上传医疗记录
func (t *SmartContract) uploadMedicalRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 8 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	r, _ := stub.GetState(args[0])
	if r != nil {
		return shim.Error("医疗记录ID已经存在！")
	}
	// fmt.Println(err0.Error())
	m := MedicalRecord{
		MedicalRecordID: args[0],
		UserID:          args[1],
		PatientID:       args[2],
		OrganizationID:  args[3],
		DataType:        "Medical",
		DataField:       args[4],
		Data:            args[5],
		EntryMethod:     args[6],
		Time:            time.Now().Format("2006-01-02 15:04:05"),
	}
	b, err1 := json.Marshal(m)
	if err1 != nil {
		return shim.Error(err1.Error())
	}
	err2 := stub.PutState(m.MedicalRecordID, b)
	if err2 != nil {
		return shim.Error(err2.Error())
	}
	err3 := stub.SetEvent(args[7], []byte{})
	if err3 != nil {
		return shim.Error(err3.Error())
	}
	return shim.Success([]byte("医疗记录添加成功！"))
}

// 操作医疗记录
func (t *SmartContract) operateMedicalRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
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
	result1, _ := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{args[0]})
	defer result1.Close()
	if result1.HasNext() {
		return shim.Error("操作记录ID存在！")
	}
	var m MedicalRecord
	err2 := json.Unmarshal(result0, &m)
	if err2 != nil {
		return shim.Error("反序列化医疗记录失败！")
	}
	key, errK := stub.CreateCompositeKey(OperationRecordString, []string{args[0], args[1], args[2], args[3], m.UserID, m.PatientID})
	if errK != nil {
		return shim.Error("组合键创建失败！")
	}
	o := OperationRecord{
		OperationRecordID: args[0],
		UserID:            args[1],
		OrganizationID:    args[2],
		DataType:          "Operation",
		ActionType:        "Access",
		MedicalRecordID:   args[3],
		OriginalAuthorID:  m.UserID,
		PatientID:         m.PatientID,
		EntryMethod:       "Auto",
		IsSuccess:         args[2] == m.OrganizationID, //如果组织不一致则判断为False
		Time:              time.Now().Format("2006-01-02 15:04:05"),
	}
	//ActionType为Access
	op, err3 := json.Marshal(o)
	if err3 != nil {
		return shim.Error("操作记录加工失败！")
	}
	err4 := stub.PutState(key, op)
	if err4 != nil {
		return shim.Error("操作失败！")
	}
	err5 := stub.SetEvent(args[4], []byte{})
	if err5 != nil {
		return shim.Error(err5.Error())
	}

	//如果组织不一致则报错
	if !o.IsSuccess {
		return shim.Error("错误操作！")
	}
	return shim.Success(result0)
}

func (t *SmartContract) getMedicalRecordHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
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
	result1, _ := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{args[0]})
	defer result1.Close()
	if result1.HasNext() {
		return shim.Error("操作记录ID存在！")
	}
	var m MedicalRecord
	err2 := json.Unmarshal(result0, &m)
	if err2 != nil {
		return shim.Error("反序列化医疗记录失败！")
	}
	key, errK := stub.CreateCompositeKey(OperationRecordString, []string{args[0], args[1], args[2], args[3], m.UserID, m.PatientID})
	if errK != nil {
		return shim.Error("组合键创建失败！")
	}
	var historys []HistoryItem
	var hisMed MedicalRecord
	iterator, err3 := stub.GetHistoryForKey(args[3])
	if err3 != nil {
		return shim.Error("根据医疗记录ID查询对应的历史变更数据失败！")
	}
	defer iterator.Close()
	for iterator.HasNext() {
		hisData, err4 := iterator.Next()
		if err4 != nil {
			return shim.Error("获取历史变更数据失败！")
		}

		var historyItem HistoryItem
		historyItem.TxId = hisData.TxId
		json.Unmarshal(hisData.Value, &hisMed)

		if hisData.Value == nil {
			var empty MedicalRecord
			historyItem.MedicalRecord = empty
		} else {
			historyItem.MedicalRecord = hisMed
		}

		historys = append(historys, historyItem)

	}
	o := OperationRecord{
		OperationRecordID: args[0],
		UserID:            args[1],
		OrganizationID:    args[2],
		DataType:          "Operation",
		ActionType:        "Access",
		MedicalRecordID:   args[3],
		OriginalAuthorID:  m.UserID,
		PatientID:         m.PatientID,
		EntryMethod:       "Auto",
		IsSuccess:         args[2] == m.OrganizationID,
		Time:              time.Now().Format("2006-01-02 15:04:05"),
	}
	//ActionType为Access
	op, err5 := json.Marshal(o)
	if err5 != nil {
		return shim.Error("操作记录加工失败！")
	}
	err6 := stub.PutState(key, op)
	if err6 != nil {
		return shim.Error("操作失败！")
	}
	err7 := stub.SetEvent(args[4], []byte{})
	if err7 != nil {
		return shim.Error(err5.Error())
	}
	m.Historys = historys
	// 返回
	result2, err8 := json.Marshal(m)
	if err8 != nil {
		return shim.Error("序列化时发生错误！")
	}

	if !o.IsSuccess {
		return shim.Error("失败操作！")
	}
	return shim.Success(result2)
}

// 对医疗记录进行修改，此次修改的为医疗记录的Data字段
func (t *SmartContract) updateMedicalRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 6 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	result0, err0 := stub.GetState(args[3])
	if err0 != nil {
		return shim.Error("医疗记录获取失败！")
	}
	if result0 == nil {
		return shim.Error("医疗记录ID获取失败！")
	}
	result1, _ := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{args[0]})
	defer result1.Close()
	if result1.HasNext() {
		return shim.Error("操作记录ID存在！")
	}
	var m MedicalRecord
	err2 := json.Unmarshal(result0, &m)
	if err2 != nil {
		return shim.Error("反序列化医疗记录失败！")
	}
	key, errK := stub.CreateCompositeKey(OperationRecordString, []string{args[0], args[1], args[2], args[3], m.UserID, m.PatientID})
	if errK != nil {
		return shim.Error("组合键创建失败！")
	}
	m.Data = args[4]
	newMedicalRecordAsBytes, err3 := json.Marshal(m)
	if err3 != nil {
		return shim.Error("医疗记录加工失败！")
	}
	err4 := stub.PutState(m.MedicalRecordID, newMedicalRecordAsBytes)
	if err4 != nil {
		return shim.Error("医疗记录修改失败！")
	}
	o := OperationRecord{
		OperationRecordID: args[0],
		UserID:            args[1],
		OrganizationID:    args[2],
		DataType:          "Operation",
		ActionType:        "Update",
		MedicalRecordID:   args[3],
		OriginalAuthorID:  m.UserID,
		PatientID:         m.PatientID,
		EntryMethod:       "Auto",
		IsSuccess:         args[2] == m.OrganizationID,
		Time:              time.Now().Format("2006-01-02 15:04:05"),
	}
	//ActionType为Update
	op, err5 := json.Marshal(o)
	if err5 != nil {
		return shim.Error("操作记录加工失败！")
	}
	err6 := stub.PutState(key, op)
	if err6 != nil {
		return shim.Error("操作失败！")
	}
	err7 := stub.SetEvent(args[5], []byte{})
	if err7 != nil {
		return shim.Error(err7.Error())
	}

	if !o.IsSuccess {
		return shim.Error("失败操作！")
	}

	return shim.Success([]byte("医疗记录修改成功！"))
}

// 删除医疗记录
func (t *SmartContract) deleteMedicalRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
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
	result1, _ := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{args[0]})
	defer result1.Close()
	if result1.HasNext() {
		return shim.Error("操作记录ID存在！")
	}
	var m MedicalRecord
	err2 := json.Unmarshal(result0, &m)
	if err2 != nil {
		return shim.Error("反序列化医疗记录失败！")
	}
	key, errK := stub.CreateCompositeKey(OperationRecordString, []string{args[0], args[1], args[2], args[3], m.UserID, m.PatientID})
	if errK != nil {
		return shim.Error("组合键创建失败！")
	}
	o := OperationRecord{
		OperationRecordID: args[0],
		UserID:            args[1],
		OrganizationID:    args[2],
		DataType:          "Operation",
		ActionType:        "Delete",
		MedicalRecordID:   args[3],
		OriginalAuthorID:  m.UserID,
		PatientID:         m.PatientID,
		EntryMethod:       "Auto",
		IsSuccess:         args[2] == m.OrganizationID,
		Time:              time.Now().Format("2006-01-02 15:04:05"),
	}
	err3 := stub.DelState(args[3])
	if err3 != nil {
		return shim.Error("删除医疗记录失败！")
	}
	op, err4 := json.Marshal(o)
	if err4 != nil {
		return shim.Error("操作记录加工失败！")
	}
	err5 := stub.PutState(key, op)
	if err5 != nil {
		return shim.Error("操作失败！")
	}
	err6 := stub.SetEvent(args[4], []byte{})
	if err6 != nil {
		return shim.Error(err6.Error())
	}

	if !o.IsSuccess {
		return shim.Error("")
	}
	return shim.Success([]byte("操作成功！"))

}

// 获取所有时间的操作记录
func (t *SmartContract) auditForAllLogs(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 4 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	resultsIterator, err0 := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{})
	if err0 != nil {
		return shim.Error("操作记录获取失败！")
	}
	result1, err1 := stub.GetState(args[0])
	if err1 == nil && result1 != nil {
		return shim.Error("审计记录存在！")
	}
	defer resultsIterator.Close()
	results := OperationRecordArr{}
	for resultsIterator.HasNext() {
		queryResult, err2 := resultsIterator.Next()
		if err2 != nil {
			return shim.Error("迭代失败！")
		}
		var o OperationRecord
		err3 := json.Unmarshal(queryResult.Value, &o)
		if err3 != nil {
			return shim.Error("反序列化操作记录失败！")
		}
		if reflect.TypeOf(o).Name() == "OperationRecord" && o.DataType == "Operation" && (o.ActionType == "Access" || o.ActionType == "Change" || o.ActionType == "Delete") {
			results.OperationRecord = append(results.OperationRecord, o)
		}
	}
	a := AuditRecord{
		AuditRecordID:  args[0],
		UserID:         args[1],
		OrganizationID: args[2],
		DataType:       "Audit",
		EntryMethod:    "Auto",
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	op, err4 := json.Marshal(a)
	if err4 != nil {
		return shim.Error("审计记录加工失败！")
	}
	err5 := stub.PutState(args[0], op)
	if err5 != nil {
		return shim.Error("审计失败！")
	}
	err6 := stub.SetEvent(args[3], []byte{})
	if err6 != nil {
		return shim.Error(err6.Error())
	}
	r, err7 := json.Marshal(results)
	if err7 != nil {
		return shim.Error("结果加工失败！")
	}
	return shim.Success(r)
}

// 获取指定开始结束时间的操作记录
func (t *SmartContract) auditForTimeRange(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 6 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	resultsIterator, err0 := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{})
	if err0 != nil {
		return shim.Error("操作记录获取失败！")
	}
	result1, err1 := stub.GetState(args[0])
	if err1 == nil && result1 != nil {
		return shim.Error("审计记录存在！")
	}
	defer resultsIterator.Close()
	results := OperationRecordArr{}
	for resultsIterator.HasNext() {
		queryResult, err2 := resultsIterator.Next()
		if err2 != nil {
			return shim.Error("迭代失败！")
		}
		var o OperationRecord
		err3 := json.Unmarshal(queryResult.Value, &o)
		if err3 != nil {
			return shim.Error("反序列化医疗记录失败！")
		}
		if reflect.TypeOf(o).Name() == "OperationRecord" && o.DataType == "Operation" && (o.ActionType == "Access" || o.ActionType == "Change" || o.ActionType == "Delete") {
			formatTime, errT := time.Parse("2006-01-02 15:04:05", o.Time)
			if errT != nil {
				return shim.Error("记录时间转换失败！")
			}
			startTime, errST := time.Parse("2006-01-02 15:04:05", args[3])
			if errST != nil {
				return shim.Error("开始时间转换失败！")
			}
			endTime, errET := time.Parse("2006-01-02 15:04:05", args[4])
			if errET != nil {
				return shim.Error("结束时间转换失败！")
			}
			if formatTime.Unix() >= startTime.Unix() && formatTime.Unix() <= endTime.Unix() {
				// arg[3]：开始时间
				// arg[4]：结束时间
				results.OperationRecord = append(results.OperationRecord, o)
			}
		}
	}
	a := AuditRecord{
		AuditRecordID:  args[0],
		UserID:         args[1],
		OrganizationID: args[2],
		DataType:       "Audit",
		EntryMethod:    "Auto",
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	op, err4 := json.Marshal(a)
	if err4 != nil {
		return shim.Error("审计记录加工失败！")
	}
	err5 := stub.PutState(args[0], op)
	if err5 != nil {
		return shim.Error("审计失败！")
	}
	err6 := stub.SetEvent(args[5], []byte{})
	if err6 != nil {
		return shim.Error(err6.Error())
	}
	r, err7 := json.Marshal(results)
	if err7 != nil {
		return shim.Error("结果加工失败！")
	}
	return shim.Success(r)
}

// 获取指定用户的操作记录
func (t *SmartContract) auditForUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	resultsIterator, err0 := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{})
	if err0 != nil {
		return shim.Error("操作记录获取失败！")
	}
	result1, err1 := stub.GetState(args[0])
	if err1 == nil && result1 != nil {
		return shim.Error("审计记录存在！")
	}
	defer resultsIterator.Close()
	results := OperationRecordArr{}
	for resultsIterator.HasNext() {
		queryResult, err2 := resultsIterator.Next()
		if err2 != nil {
			return shim.Error("迭代失败！")
		}
		var o OperationRecord
		err3 := json.Unmarshal(queryResult.Value, &o)
		if err3 != nil {
			return shim.Error("反序列化医疗记录失败！")
		}
		if reflect.TypeOf(o).Name() == "OperationRecord" && o.DataType == "Operation" && (o.ActionType == "Access" || o.ActionType == "Change" || o.ActionType == "Delete") {
			if o.UserID == args[3] {
				results.OperationRecord = append(results.OperationRecord, o)
			}
		}
	}
	a := AuditRecord{
		AuditRecordID:  args[0],
		UserID:         args[1],
		OrganizationID: args[2],
		DataType:       "Audit",
		EntryMethod:    "Auto",
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	op, err4 := json.Marshal(a)
	if err4 != nil {
		return shim.Error("审计记录加工失败！")
	}
	err5 := stub.PutState(args[0], op)
	if err5 != nil {
		return shim.Error("审计失败！")
	}
	err6 := stub.SetEvent(args[4], []byte{})
	if err6 != nil {
		return shim.Error(err6.Error())
	}
	r, err7 := json.Marshal(results)
	if err7 != nil {
		return shim.Error("结果加工失败！")
	}
	return shim.Success(r)
}

// 获取关于指定组织的操作记录
func (t *SmartContract) auditForOrganisation(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	resultsIterator, err0 := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{})
	if err0 != nil {
		return shim.Error("操作记录获取失败！")
	}
	result1, err1 := stub.GetState(args[0])
	if err1 == nil && result1 != nil {
		return shim.Error("审计记录存在！")
	}
	defer resultsIterator.Close()
	results := OperationRecordArr{}
	for resultsIterator.HasNext() {
		queryResult, err2 := resultsIterator.Next()
		if err2 != nil {
			return shim.Error("迭代失败！")
		}
		var o OperationRecord
		err3 := json.Unmarshal(queryResult.Value, &o)
		if err3 != nil {
			return shim.Error("反序列化医疗记录失败！")
		}
		if reflect.TypeOf(o).Name() == "OperationRecord" && o.DataType == "Operation" && (o.ActionType == "Access" || o.ActionType == "Change" || o.ActionType == "Delete") {
			if o.OrganizationID == args[3] {
				results.OperationRecord = append(results.OperationRecord, o)
			}
		}
	}
	a := AuditRecord{
		AuditRecordID:  args[0],
		UserID:         args[1],
		OrganizationID: args[2],
		DataType:       "Audit",
		EntryMethod:    "Auto",
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	op, err4 := json.Marshal(a)
	if err4 != nil {
		return shim.Error("审计记录加工失败！")
	}
	err5 := stub.PutState(args[0], op)
	if err5 != nil {
		return shim.Error("审计失败！")
	}
	err6 := stub.SetEvent(args[4], []byte{})
	if err6 != nil {
		return shim.Error(err6.Error())
	}
	r, err7 := json.Marshal(results)
	if err7 != nil {
		return shim.Error("结果加工失败！")
	}
	return shim.Success(r)
}

// 获取特定医疗记录的操作记录
func (t *SmartContract) auditForMedicalRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	resultsIterator, err0 := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{})
	if err0 != nil {
		return shim.Error("操作记录获取失败！")
	}
	result1, err1 := stub.GetState(args[0])
	if err1 == nil && result1 != nil {
		return shim.Error("审计记录存在！")
	}
	defer resultsIterator.Close()
	results := OperationRecordArr{}
	for resultsIterator.HasNext() {
		queryResult, err2 := resultsIterator.Next()
		if err2 != nil {
			return shim.Error("迭代失败！")
		}
		var o OperationRecord
		err3 := json.Unmarshal(queryResult.Value, &o)
		if err3 != nil {
			return shim.Error("反序列化医疗记录失败！")
		}
		if reflect.TypeOf(o).Name() == "OperationRecord" && o.DataType == "Operation" && (o.ActionType == "Access" || o.ActionType == "Change" || o.ActionType == "Delete") {
			if o.MedicalRecordID == args[3] {
				results.OperationRecord = append(results.OperationRecord, o)
			}
		}
	}
	a := AuditRecord{
		AuditRecordID:  args[0],
		UserID:         args[1],
		OrganizationID: args[2],
		DataType:       "Audit",
		EntryMethod:    "Auto",
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	op, err4 := json.Marshal(a)
	if err4 != nil {
		return shim.Error("审计记录加工失败！")
	}
	err5 := stub.PutState(args[0], op)
	if err5 != nil {
		return shim.Error("审计失败！")
	}
	err6 := stub.SetEvent(args[4], []byte{})
	if err6 != nil {
		return shim.Error(err6.Error())
	}
	r, err7 := json.Marshal(results)
	if err7 != nil {
		return shim.Error("结果加工失败！")
	}
	return shim.Success(r)
}

// 获取特定医疗记录创建者的操作记录
func (t *SmartContract) auditForOriginalAuthor(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	resultsIterator, err0 := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{})
	if err0 != nil {
		return shim.Error("操作记录获取失败！")
	}
	result1, err1 := stub.GetState(args[0])
	if err1 == nil && result1 != nil {
		return shim.Error("审计记录存在！")
	}
	defer resultsIterator.Close()
	results := OperationRecordArr{}
	for resultsIterator.HasNext() {
		queryResult, err2 := resultsIterator.Next()
		if err2 != nil {
			return shim.Error("迭代失败！")
		}
		var o OperationRecord
		err3 := json.Unmarshal(queryResult.Value, &o)
		if err3 != nil {
			return shim.Error("反序列化医疗记录失败！")
		}
		if reflect.TypeOf(o).Name() == "OperationRecord" && o.DataType == "Operation" && (o.ActionType == "Access" || o.ActionType == "Change" || o.ActionType == "Delete") {
			if o.OriginalAuthorID == args[3] {
				results.OperationRecord = append(results.OperationRecord, o)
			}
		}
	}
	a := AuditRecord{
		AuditRecordID:  args[0],
		UserID:         args[1],
		OrganizationID: args[2],
		DataType:       "Audit",
		EntryMethod:    "Auto",
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	op, err4 := json.Marshal(a)
	if err4 != nil {
		return shim.Error("审计记录加工失败！")
	}
	err5 := stub.PutState(args[0], op)
	if err5 != nil {
		return shim.Error("审计失败！")
	}
	err6 := stub.SetEvent(args[4], []byte{})
	if err6 != nil {
		return shim.Error(err6.Error())
	}
	r, err7 := json.Marshal(results)
	if err7 != nil {
		return shim.Error("结果加工失败！")
	}
	return shim.Success(r)
}

// 获取指定病人医疗记录的操作记录
func (t *SmartContract) auditForPatient(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("给定的参数个数不符合要求！")
	}
	resultsIterator, err0 := stub.GetStateByPartialCompositeKey(OperationRecordString, []string{})
	if err0 != nil {
		return shim.Error("操作记录获取失败！")
	}
	result1, err1 := stub.GetState(args[0])
	if err1 == nil && result1 != nil {
		return shim.Error("审计记录存在！")
	}
	defer resultsIterator.Close()
	results := OperationRecordArr{}
	for resultsIterator.HasNext() {
		queryResult, err2 := resultsIterator.Next()
		if err2 != nil {
			return shim.Error("迭代失败！")
		}
		var o OperationRecord
		err3 := json.Unmarshal(queryResult.Value, &o)
		if err3 != nil {
			return shim.Error("反序列化医疗记录失败！")
		}
		if reflect.TypeOf(o).Name() == "OperationRecord" && o.DataType == "Operation" && (o.ActionType == "Access" || o.ActionType == "Change" || o.ActionType == "Delete") {
			if o.PatientID == args[3] {
				results.OperationRecord = append(results.OperationRecord, o)
			}
		}
	}
	a := AuditRecord{
		AuditRecordID:  args[0],
		UserID:         args[1],
		OrganizationID: args[2],
		DataType:       "Audit",
		EntryMethod:    "Auto",
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	op, err4 := json.Marshal(a)
	if err4 != nil {
		return shim.Error("审计记录加工失败！")
	}
	err5 := stub.PutState(args[0], op)
	if err5 != nil {
		return shim.Error("审计失败！")
	}
	err6 := stub.SetEvent(args[4], []byte{})
	if err6 != nil {
		return shim.Error(err6.Error())
	}
	r, err7 := json.Marshal(results)
	if err7 != nil {
		return shim.Error("结果加工失败！")
	}
	return shim.Success(r)
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("启动链码时发生错误: %s", err)
	}
}

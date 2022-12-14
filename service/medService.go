package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"medical_testdemo/sqlaction"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *ServiceSetup) UploadMed(args []string) (string, error) {
	DB := InitDB()
	args[11] = sqlaction.GetUserLogin(DB, "select username from login where state='1'") //确定用户名
	var casenumer, policy, mess string
	policy = ""
	casenumer = sqlaction.GetCaseNumber(args)
	if !InsertDB(DB, args, casenumer) {
		return "", fmt.Errorf("数据库插入不成功！")
	} else {
		if InsertDB2Insti(DB, casenumer) {
			policy = GeneratePolicy(DB, casenumer)
			//fmt.Println(policy)
		}
	}
	eventID := "eventUploadMed"
	resultStr := "Success"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "UploadMedicalRecord", Args: [][]byte{[]byte(casenumer), []byte(args[11]), []byte(args[12]), []byte(resultStr), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return "", err0
	}

	err1 := eventResult(notifier, eventID)
	fmt.Println(err1)
	if err1 != nil {
		return "", err1
	}

	mess = string(respone.TransactionID)[0:6] + " : " + casenumer + "-policy = " + policy
	//return policy, nil
	//return string(respone.TransactionID), nil
	return mess, nil
}

func (t *ServiceSetup) OperateMed(args []string) ([]byte, error) {
	if len(args) != 4 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	DB := InitDB()
	casenumer := args[0]
	if !CheckAction(DB, casenumer, "r") {
		return nil, fmt.Errorf("权限不足，无法操作")
	}
	eventID := "eventAccessMed"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "AccessMedicalRecord", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return respone.Payload, err0
	}

	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return []byte{0x00}, err1
	}

	mp := SelectDBSingle(DB, args)
	if mp == nil {
		return []byte{0x00}, fmt.Errorf("数据库查询不成功！")
	}
	m := *mp
	b, err2 := json.Marshal(m)
	if err2 != nil {
		return []byte{0x00}, err2
	}
	return b, nil
}

func (t *ServiceSetup) DeleteMed(args []string) (string, error) {
	if len(args) != 4 {
		return "", fmt.Errorf("给定的参数个数不符合要求！")
	}
	DB := InitDB()
	// DB := InitDB()
	casenumer := args[0]
	if !CheckAction(DB, casenumer, "d") {
		return "", fmt.Errorf("权限不足，无法操作")
	}
	if !DeleteDB(DB, args) {
		return "", fmt.Errorf("删除数据不成功！")
	}
	eventID := "eventDeleteMed"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "DeleteMedicalRecord", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return "", err0
	}

	err1 := eventResult(notifier, eventID)
	fmt.Println(err1)
	if err1 != nil {
		return "", err1
	}
	return string(respone.TransactionID), nil
}

func (t *ServiceSetup) UpdateMed(args []string) (string, error) {
	if len(args) != 17 {
		return "", fmt.Errorf("给定的参数个数不符合要求！")
	}
	DB := InitDB()
	casenumer := args[0]
	if !CheckAction(DB, casenumer, "w") {
		return "", fmt.Errorf("权限不足，无法操作")
	}
	if !UpdateDB(DB, args) {
		return "", fmt.Errorf("数据库修改不成功！")
	}
	eventID := "eventUpdateMed"
	resultStr := "Success"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "UpdateMedicalRecord", Args: [][]byte{[]byte(args[4]), []byte(args[11]), []byte(args[12]), []byte(resultStr), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return "", err0
	}

	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return "", err1
	}

	return string(respone.TransactionID), nil
}

func (t *ServiceSetup) UserLogin(username string, password string) (bool, error) {
	DB := InitDB()
	SQLString := "select username from login"
	Username := make(map[int]string)
	Username = queryDB(DB, SQLString)

	for _, user := range Username {
		if user == username {
			var str string
			SQLString2 := "select password from login where username='" + user + "'"
			err := DB.QueryRow(SQLString2).Scan(&str)
			if err != sql.ErrNoRows && str == password {
				SQLString3 := "UPDATE login SET state= '1' where username='" + user + "'"
				_, err := DB.Exec(SQLString3)
				if err != nil {
					return false, err
				}
				return true, nil
			}
		}
	}
	return false, nil
}

func (t *ServiceSetup) UserLoginInfo() (map[int]string, error) {
	Result := make(map[int]string)
	DB := InitDB()
	var str string
	SQLString1 := "select usertype from login where state ='1'"
	err := DB.QueryRow(SQLString1).Scan(&str)
	if err != sql.ErrNoRows {
		SQLString2 := "select * from user_type where user_id ='" + str + "'"
		Result = queryDB(DB, SQLString2)
	}
	return Result, nil
}

func (t *ServiceSetup) UserLoginOut() (bool, error) {
	DB := InitDB()
	Userinfo := sqlaction.GetUserLogin(DB, "select username from login where state='1'") //确定用户名
	SQLString3 := "UPDATE login SET state= '0' where username='" + Userinfo + "'"
	_, err := DB.Exec(SQLString3)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (t *ServiceSetup) AuditAll(args []string) ([]byte, error) {
	if len(args) != 3 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventAuditAll"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "auditForAllLogs", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return []byte{0x00}, err0
	}
	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return []byte{0x00}, err1
	}
	return respone.Payload, nil
}
func (t *ServiceSetup) AuditTimeRange(args []string) ([]byte, error) {
	if len(args) != 5 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventAuditTimeRange"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "auditForTimeRange", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return []byte{0x00}, err0
	}

	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return []byte{0x00}, err1
	}

	return respone.Payload, nil
}
func (t *ServiceSetup) AuditUser(args []string) ([]byte, error) {
	if len(args) != 4 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventAuditUser"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "auditForUser", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return []byte{0x00}, err0
	}

	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return []byte{0x00}, err1
	}

	return respone.Payload, nil
}
func (t *ServiceSetup) AuditOrganisation(args []string) ([]byte, error) {
	if len(args) != 4 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventAuditOrganisation"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "auditForOrganisation", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return []byte{0x00}, err0
	}

	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return []byte{0x00}, err1
	}

	return respone.Payload, nil
}
func (t *ServiceSetup) AuditMedicalRecord(args []string) ([]byte, error) {
	if len(args) != 4 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventAuditMedicalRecord"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "auditForMedicalRecord", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return []byte{0x00}, err0
	}

	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return []byte{0x00}, err1
	}

	return respone.Payload, nil
}
func (t *ServiceSetup) AuditOriginalAuthor(args []string) ([]byte, error) {
	if len(args) != 4 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventAuditOriginalAuthor"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "auditForOriginalAuthor", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return []byte{0x00}, err0
	}

	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return []byte{0x00}, err1
	}

	return respone.Payload, nil
}
func (t *ServiceSetup) AuditPatient(args []string) ([]byte, error) {
	if len(args) != 4 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventAuditPatient"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "auditForPatient", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(eventID)}}
	respone, err0 := t.Client.Execute(req)
	if err0 != nil {
		return []byte{0x00}, err0
	}

	err1 := eventResult(notifier, eventID)
	if err1 != nil {
		return []byte{0x00}, err1
	}

	return respone.Payload, nil
}

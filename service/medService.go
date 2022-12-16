package service

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *ServiceSetup) UploadMed(args []string) (string, error) {
	if len(args) != 17 {
		return "", fmt.Errorf("给定的参数个数不符合要求！")
	}
	if !InitDB() {
		return "", fmt.Errorf("数据库连接不成功！")
	}
	if !InsertDB(args) {
		return "", fmt.Errorf("数据库插入不成功！")
	}
	eventID := "eventUploadMed"
	resultStr := "Success"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "uploadMedicalRecord", Args: [][]byte{[]byte(args[4]), []byte(args[11]), []byte(args[12]), []byte(resultStr), []byte(eventID)}}
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
func (t *ServiceSetup) OperateMed(args []string) ([]byte, error) {
	if len(args) != 4 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
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
	if !InitDB() {
		return []byte{0x00}, fmt.Errorf("数据库连接不成功！")
	}
	mp := SelectDBSingle(args)
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
func (t *ServiceSetup) GetMedHistory(args []string) ([]byte, error) {
	if len(args) != 4 {
		return []byte{0x00}, fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventGetMedicalRecordHistory"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "getMedicalRecordHistory", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(eventID)}}
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
func (t *ServiceSetup) UpdateMed(args []string) (string, error) {
	if len(args) != 5 {
		return "", fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventUpdateMed"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "updateMedicalRecord", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4]), []byte(eventID)}}
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
func (t *ServiceSetup) DeleteMed(args []string) (string, error) {
	if len(args) != 4 {
		return "", fmt.Errorf("给定的参数个数不符合要求！")
	}
	eventID := "eventDeleteMed"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "deleteMedicalRecord", Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(eventID)}}
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

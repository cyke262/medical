package controller

import (
	"encoding/json"
	"medical/service"
	"net/http"
	"reflect"
)

var cuser User
var data Data

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "login.html", nil)
}

func (app *Application) Index(w http.ResponseWriter, r *http.Request) {
	ShowView(w, r, "index.html", nil)
}

// 用户登录
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	loginName := r.FormValue("loginName")
	password := r.FormValue("password")

	var flag bool
	for _, user := range users {
		if user.LoginName == loginName && user.Password == password {
			cuser = user
			flag = true
			break
		}
	}
	data.CurrentUser = cuser
	data.Flag = false

	if flag {
		// 登录成功
		ShowView(w, r, "index.html", data)
	} else {
		// 登录失败
		data.Flag = true
		data.CurrentUser.LoginName = loginName
		ShowView(w, r, "login.html", data)
	}
}

// 用户登出
func (app *Application) LoginOut(w http.ResponseWriter, r *http.Request) {
	cuser = User{}
	ShowView(w, r, "login.html", nil)
}

func (app *Application) UploadMed(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Flag = true
	data.Msg = ""

	arr := [7]string{r.FormValue("medicalRecordID"), cuser.LoginName, r.FormValue("patientID"), r.FormValue("organizationID"), r.FormValue("dataField"), r.FormValue("data"), r.FormValue("entryMethod")}
	transactionID, err := app.Setup.UploadMed(arr[:])

	if err != nil {
		data.Msg = err.Error()
	} else {
		data.Msg = "信息添加成功:" + transactionID
	}
	ShowView(w, r, "uploadMed.html", data)
}

func (app *Application) OperateMed(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	ShowView(w, r, "operateMed.html", data)
}

func (app *Application) AuditMed(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	ShowView(w, r, "auditMed.html", data)
}

func (app *Application) AccessMed(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	ShowView(w, r, "accessMed.html", data)
}

func (app *Application) AccessMedHistory(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = true
	ShowView(w, r, "accessMedHistory.html", data)
}

func (app *Application) AccessMedResult(w http.ResponseWriter, r *http.Request) {
	arr := [4]string{r.FormValue("operationRecordID"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("medicalRecordID")}
	var result []byte
	var err error
	if data.History {
		result, err = app.Setup.GetMedHistory(arr[:])
	} else {
		result, err = app.Setup.OperateMed(arr[:])
	}
	var med = service.MedicalRecord{}
	json.Unmarshal(result, &med)
	data.Med = med
	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}
	if reflect.DeepEqual(med, service.MedicalRecord{}) {
		if data.History {
			ShowView(w, r, "accessMedHistory.html", data)
		} else {
			ShowView(w, r, "accessMed.html", data)
		}
	} else {
		ShowView(w, r, "accessMedResult.html", data)
	}
}

func (app *Application) UpdateMed(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Flag = true
	data.Msg = ""
	arr := [5]string{r.FormValue("operationRecordID"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("medicalRecordID"), r.FormValue("medicalRecordData")}
	transactionID, err := app.Setup.UpdateMed(arr[:])

	if err != nil {
		data.Msg = err.Error()
	} else {
		data.Msg = "医疗记录删除成功:" + transactionID
	}
	ShowView(w, r, "updateMed.html", data)
}

func (app *Application) DeleteMed(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Flag = true
	data.Msg = ""
	arr := [4]string{r.FormValue("operationRecordID"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("medicalRecordID")}
	transactionID, err := app.Setup.DeleteMed(arr[:])

	if err != nil {
		data.Msg = err.Error()
	} else {
		data.Msg = "医疗记录删除成功:" + transactionID
	}
	ShowView(w, r, "deleteMed.html", data)
}

func (app *Application) AuditAllRecords(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditAll"
	ShowView(w, r, "auditAllRecords.html", data)
}

func (app *Application) AuditTimeRangeStartEnd(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditTimeRangeStartEnd"
	ShowView(w, r, "auditTimeRangeStartEnd.html", data)
}

func (app *Application) AuditByUser(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditByUser"
	ShowView(w, r, "auditByUser.html", data)
}

func (app *Application) AuditByPatient(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditByPatient"
	ShowView(w, r, "auditByPatient.html", data)
}

func (app *Application) AuditByOrganisation(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditByOrganisation"
	ShowView(w, r, "auditByOrganisation.html", data)
}

func (app *Application) AuditByMedicalRecord(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditByMedicalRecord"
	ShowView(w, r, "auditByMedicalRecord.html", data)
}

func (app *Application) AuditByOriginalAuthor(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditByOriginalAuthor"
	ShowView(w, r, "auditByOriginalAuthor.html", data)
}

func (app *Application) AuditResult(w http.ResponseWriter, r *http.Request) {
	if data.AuditString == "AuditAll" {
		arr := [3]string{r.FormValue("auditRecordID"), cuser.LoginName, r.FormValue("organisationID")}
		result, err := app.Setup.AuditAll(arr[:])
		var ops = service.OperationRecordArr{}
		json.Unmarshal(result, &ops)
		data.Ops = ops
		data.CurrentUser = cuser
		data.Msg = ""
		data.Flag = false
		if err != nil {
			data.Msg = err.Error()
			data.Flag = true
		}
		if reflect.DeepEqual(data.Ops, service.OperationRecordArr{}) {
			ShowView(w, r, "auditAllRecords.html", nil)
		} else {
			ShowView(w, r, "auditResult.html", data)
		}
	} else if data.AuditString == "AuditTimeRangeStartEnd" {
		arr := [5]string{r.FormValue("auditRecordID"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("startTimePoint"), r.FormValue("endTimePoint")}
		result, err := app.Setup.AuditTimeRange(arr[:])
		var ops = service.OperationRecordArr{}
		json.Unmarshal(result, &ops)
		data.Ops = ops
		data.CurrentUser = cuser
		data.Msg = ""
		data.Flag = false
		if err != nil {
			data.Msg = err.Error()
			data.Flag = true
		}
		if reflect.DeepEqual(data.Ops, service.OperationRecordArr{}) {
			ShowView(w, r, "auditTimeRangeStartEnd.html", nil)
		} else {
			ShowView(w, r, "auditResult.html", data)
		}
	} else if data.AuditString == "AuditByMedicalRecord" {
		arr := [4]string{r.FormValue("auditRecordID"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("str")}
		result, err := app.Setup.AuditMedicalRecord(arr[:])
		var ops = service.OperationRecordArr{}
		json.Unmarshal(result, &ops)
		data.Ops = ops
		data.CurrentUser = cuser
		data.Msg = ""
		data.Flag = false
		if err != nil {
			data.Msg = err.Error()
			data.Flag = true
		}
		if reflect.DeepEqual(data.Ops, service.OperationRecordArr{}) {
			ShowView(w, r, "auditByMedicalRecord.html", nil)
		} else {
			ShowView(w, r, "auditResult.html", data)
		}
	} else if data.AuditString == "AuditByUser" {
		arr := [4]string{r.FormValue("auditRecordID"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("str")}
		result, err := app.Setup.AuditUser(arr[:])
		var ops = service.OperationRecordArr{}
		json.Unmarshal(result, &ops)
		data.Ops = ops
		data.CurrentUser = cuser
		data.Msg = ""
		data.Flag = false
		if err != nil {
			data.Msg = err.Error()
			data.Flag = true
		}
		if reflect.DeepEqual(data.Ops, service.OperationRecordArr{}) {
			ShowView(w, r, "auditByUser.html", nil)
		} else {
			ShowView(w, r, "auditResult.html", data)
		}
	} else if data.AuditString == "AuditByPatient" {
		arr := [4]string{r.FormValue("auditRecordID"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("str")}
		result, err := app.Setup.AuditPatient(arr[:])
		var ops = service.OperationRecordArr{}
		json.Unmarshal(result, &ops)
		data.Ops = ops
		data.CurrentUser = cuser
		data.Msg = ""
		data.Flag = false
		if err != nil {
			data.Msg = err.Error()
			data.Flag = true
		}
		if reflect.DeepEqual(data.Ops, service.OperationRecordArr{}) {
			ShowView(w, r, "auditByPatient.html", nil)
		} else {
			ShowView(w, r, "auditResult.html", data)
		}
	} else if data.AuditString == "AuditByOrganisation" {
		arr := [4]string{r.FormValue("auditRecordID"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("str")}
		result, err := app.Setup.AuditOrganisation(arr[:])
		var ops = service.OperationRecordArr{}
		json.Unmarshal(result, &ops)
		data.Ops = ops
		data.CurrentUser = cuser
		data.Msg = ""
		data.Flag = false
		if err != nil {
			data.Msg = err.Error()
			data.Flag = true
		}
		if reflect.DeepEqual(data.Ops, service.OperationRecordArr{}) {
			ShowView(w, r, "auditByOrganisation.html", nil)
		} else {
			ShowView(w, r, "auditResult.html", data)
		}
	} else if data.AuditString == "AuditByOriginalAuthor" {
		arr := [4]string{r.FormValue("auditRecordID"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("str")}
		result, err := app.Setup.AuditOriginalAuthor(arr[:])
		var ops = service.OperationRecordArr{}
		json.Unmarshal(result, &ops)
		data.Ops = ops
		data.CurrentUser = cuser
		data.Msg = ""
		data.Flag = false
		if err != nil {
			data.Msg = err.Error()
			data.Flag = true
		}
		if reflect.DeepEqual(data.Ops, service.OperationRecordArr{}) {
			ShowView(w, r, "auditByOriginalAuthor.html", nil)
		} else {
			ShowView(w, r, "auditResult.html", data)
		}
	} else {
		ShowView(w, r, "auditMed.html", nil)
	}
}

// *新增：调用两次查询函数，返回指定时间段、指定组织的审计报告
func (app *Application) AuditReportResult(w http.ResponseWriter, r *http.Request) {
	if data.AuditString == "AuditReport" {
		arr := []string{r.FormValue("auditRecordID0"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("startTime"), r.FormValue("endTime")}
		result0, err0 := app.Setup.AuditTimeRange(arr[:])
		var ops0 = service.OperationRecordArr{}
		json.Unmarshal(result0, &ops0)
		arr = []string{r.FormValue("auditRecordID1"), cuser.LoginName, r.FormValue("organisationID"), r.FormValue("auditOrg")}
		result1, err1 := app.Setup.AuditOrganisation(arr[:])
		var ops1 = service.OperationRecordArr{}
		json.Unmarshal(result1, &ops1)

		ops := intersection(ops0, ops1)
		data.Ops = ops
		data.CurrentUser = cuser
		data.Msg = ""
		data.Flag = false
		if err0 != nil {
			data.Msg = err0.Error()
			data.Flag = true
		}
		if err1 != nil {
			data.Msg = err1.Error()
			data.Flag = true
		}

		var repo = service.AuditReport{}
		repo.TargetOrg = r.FormValue("auditOrg")

		// 组织操作信息
		total := 0
		fail := 0

		// m0：每个用户的成功操作数
		// m1：每个用户的失败操作数
		m0 := make(map[string]int, 0)
		m1 := make(map[string]int, 0)

		for _, op := range ops.OperationRecord {
			total++
			m0[op.UserID]++
			if !op.IsSuccess {
				fail++
				m1[op.UserID]++
			}
		}

		curFailRate := float64(0)
		maxFailRate := float64(0)

		for user, v := range m0 {
			curFailRate = float64(m1[user]) / float64(v)
			if curFailRate > maxFailRate {
				maxFailRate = curFailRate
				repo.MaxFailRateUser = user
				repo.MaxFailRate = maxFailRate
			}
		}

		repo.TotalOperations = int64(total)
		repo.FailOperations = int64(fail)
		repo.FailRate = float64(fail) / float64(total)

		data.Repo = repo

		if reflect.DeepEqual(ops, service.OperationRecordArr{}) {
			ShowView(w, r, "auditReportByTimeRangeAndOrg.html", nil)
		} else {
			ShowView(w, r, "auditReportResult.html", data)
		}
	} else {
		ShowView(w, r, "auditMed.html", nil)
	}
}

// *新增：同时依照时间段和组织ID审计
func (app *Application) AuditReportByTimeRangeAndOrg(w http.ResponseWriter, r *http.Request) {
	data.CurrentUser = cuser
	data.Msg = ""
	data.Flag = false
	data.History = false
	data.AuditString = "AuditReport"
	ShowView(w, r, "auditReportByTimeRangeAndOrg.html", data)
}

// *新增：返回两个OperationRecordArr的交集
func intersection(nums1 service.OperationRecordArr, nums2 service.OperationRecordArr) service.OperationRecordArr {

	m := make(map[string]int, 0)

	for _, v := range nums1.OperationRecord {
		m[v.OperationRecordID] += 1
	}

	count := 0 //记录新数组长度
	for _, v := range nums2.OperationRecord {
		if m[v.OperationRecordID] > 0 {
			m[v.OperationRecordID] = 0
			nums1.OperationRecord[count] = v
			count++
		}
	}

	return service.OperationRecordArr{OperationRecord: nums1.OperationRecord[:count]}

}

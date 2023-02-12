package web

import (
	"fmt"
	"medical/web/controller"
	"net/http"
)

func WebStart(app controller.Application) {

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 指定路由信息(匹配请求)
	http.HandleFunc("/", app.LoginView)

	// 登录
	http.HandleFunc("/login", app.Login)
	http.HandleFunc("/loginout", app.LoginOut)
	// 注册
	http.HandleFunc("/register", app.Register)
	// 忘记密码
	http.HandleFunc("/forgotpassword", app.Forgotpassword)

	http.HandleFunc("/index", app.Index)

	http.HandleFunc("/uploadMed", app.UploadMed)
	http.HandleFunc("/manageMed", app.ManageMed)
	// http.HandleFunc("/uploadOldMed", app.UploadOldMed)
	http.HandleFunc("/operateMed", app.OperateMed)
	http.HandleFunc("/auditMed", app.AuditMed)

	http.HandleFunc("/accessMed", app.AccessMed)
	http.HandleFunc("/accessMedHistory", app.AccessMedHistory)
	http.HandleFunc("/accessMedResult", app.AccessMedResult)
	http.HandleFunc("/updateMed", app.UpdateMed)
	http.HandleFunc("/deleteMed", app.DeleteMed)

	http.HandleFunc("/auditResult", app.AuditResult)
	http.HandleFunc("/auditReportResult", app.AuditReportResult)
	http.HandleFunc("/03医疗数据审计", app.AuditReportByTimeRangeAndOrg)
	http.HandleFunc("/auditAllRecords", app.AuditAllRecords)
	http.HandleFunc("/auditTimeRangeStartEnd", app.AuditTimeRangeStartEnd)
	http.HandleFunc("/auditByUser", app.AuditByUser)
	http.HandleFunc("/auditByOrganisation", app.AuditByOrganisation)
	http.HandleFunc("/auditByMedicalRecord", app.AuditByMedicalRecord)
	http.HandleFunc("/auditByOriginalAuthor", app.AuditByOriginalAuthor)
	http.HandleFunc("/auditByPatient", app.AuditByPatient)

	http.HandleFunc("/dataupload", app.DataUpload)

	fmt.Println("启动Web服务, 监听端口号为: 8088")
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		fmt.Printf("Web服务启动失败: %v", err)
	}

}

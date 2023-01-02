package abac

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func AutoPolicyCreate(attr Date2DB, CaseNumber string, GatherTime string) string {
	// policy模板
	template := `{
		"CaseNum": {
			"role": {
				"admin": {
					"action": "all"
				},
				"u1": {
					"researce": "",
					"action": "rw"
				},
				"u2": {
					"researce": "",
					"action": "r"
				},
				"u3": {
					"action": "-"
				}
			},
			"owner": "",
			"allowOrg": "",
			"time": ""
		}

	}`
	var mpJson map[string]interface{}
	json.Unmarshal([]byte(template), &mpJson)

	// TODO: 后期根据具体值传参
	casenumber := CaseNumber // _CaseNumber 数据ID
	owner := attr.Researcher //_Researcher 研究者
	ThisTime := GatherTime   //_GatherTime 时间戳

	// TODO: 参与研究的各个单位，应该是由数据库查询后给我，目前暂时写作参数
	org_list := attr.Organization
	// // 传owner的属性
	// // TODO：这里应该是由数据库调用获得Uargs，由于数据库连接我不是很清楚，所以这里写成了传参，到时候改一下就行
	// uid := Uargs[0]
	// role := Uargs[1]
	// department := Uargs[2]
	// Orgs := Uargs[3]
	Diseases := attr.Diseases //指研究方向or疾病

	policy := make(map[string]interface{})
	policy[casenumber] = mpJson["CaseNum"]
	policy[casenumber].(map[string]interface{})["owner"] = owner
	policy[casenumber].(map[string]interface{})["role"].(map[string]interface{})["u1"].(map[string]interface{})["researce"] = Diseases
	policy[casenumber].(map[string]interface{})["role"].(map[string]interface{})["u2"].(map[string]interface{})["researce"] = Diseases
	policy[casenumber].(map[string]interface{})["allowOrg"] = org_list
	policy[casenumber].(map[string]interface{})["time"] = ThisTime

	sub := policy[casenumber].(map[string]interface{})["role"].(map[string]interface{})
	var subrules []string
	for k, v := range sub {
		value, _ := json.Marshal(v.(map[string]interface{}))
		sv := string(value)
		sv = strings.Replace(sv, "\"", "", -1)
		// fmt.Println(" and v is ", sv)
		subrules = append(subrules, "role:"+k+","+sv[1:len(sv)-1])
	}

	// fmt.Println(len(subrules))
	// 把policy标准化

	p := Policy{}
	p.Obj = casenumber
	p.Owner = owner
	p.Env.AllowOrg = org_list
	p.Env.CreatedTime = ThisTime
	tt, _ := time.ParseInLocation("2006-01-02 15:04:05", ThisTime, time.Local)
	ttemp := tt.AddDate(+1, 0, 0)
	p.Env.EndTime = time.Unix(ttemp.Unix(), 0).Format("2006-01-02 15:03:04")
	p.SubRules = subrules
	// fmt.Println(p.SubRules)
	// TODO： 这里需要吧policy写入数据库
	by, err := json.Marshal(p)
	if err != nil {
		fmt.Printf("Map转化为byte数组失败,异常:%s\n", err)
		// return err
	}
	str := string(by)
	fmt.Println(str)
	return (str)
}

// json ==> abacrequest
func parseReq(arg string) (ABACRequest, error) {
	reqAsBytes := []byte(arg)
	req := ABACRequest{}
	curtime := time.Now()
	// policy := m.Policy{}
	err := json.Unmarshal(reqAsBytes, &req)
	req.CurTime = curtime.Unix()
	return req, err
}
func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

//判断是否可以获得权限
func CheckAccess(req ABACRequest, Parg string, user Sub) bool {
	// 第一个参数是（数据，用户，访问操作）
	// 传的参数参数化为ABACRequest
	// req, err := parseReq(args)
	// eventID := args[1]
	// TODO：这个参数是必须的，不用从数据库中调取
	// req := ABACRequest{}
	// req.Op = Rargs[2]
	req.CurTime = time.Now().Unix()

	// 1. get user
	// TODO：第三个参数Uargs应该是由数据库给出，这里需要调用数据库
	//  参数化User
	// user := Sub{}
	// user.UID = Uargs[0]
	// user.Role = Uargs[1]
	// user.Department = Uargs[2]
	// user.Org = Uargs[3]
	// user.Group = Uargs[4] //指研究方向or疾病
	// if user.UID ==

	//2. get policy, 调用链码查询访问规则（获得到某操作对应的权限）
	// TODO：第二个参数不是必要的，应该是由数据的casenum来进行查询得到，这里只是没有进行数据库查询而给了一个参数
	// 应该是向数据库查询访问策略 req.Obj
	// resp := QueryPolicy(APIstub, qparg)
	p := Parg
	// 参数化Policy
	policy := Policy{}
	err := json.Unmarshal([]byte(p), &policy)
	if err != nil {
		fmt.Println("参数化policy失败")
		return false
	}

	// 3.1 check有效期
	CreateTime, _ := time.ParseInLocation("2006-01-02 15:04:05", policy.Env.CreatedTime, time.Local)
	EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", policy.Env.EndTime, time.Local)
	if req.CurTime > EndTime.Unix() || req.CurTime < CreateTime.Unix() {
		fmt.Println("访问时间失效")
		return false
	}
	var tra []interface{}
	// TODO:这里需要根据数据库传参（即第二个参数，Pargs）的具体格式进行修改
	if user.Role == "admin" {
		return true

	} else {
		tra = []interface{}{"role:" + user.Role, "action:" + req.Op, "researce:" + user.Group}
	}
	// user.Org
	org_list := strings.Split(policy.Env.AllowOrg[1:], " ")
	fmt.Println(org_list)
	if !IsContain(org_list, user.Org) {
		fmt.Println("not in org_list")
		return false
	}

	//3.2 check sub的属性，这里使用树形结构判断
	// subrules := policy.SubRules
	// fmt.Println("policy is :", policy)
	ptree := PolicyToTree(&policy)
	// // stree := m.SubRuleToTree(&user)

	fmt.Println(tra)
	PreorderPrint(ptree.Root)
	fmt.Println("")
	if ptree.Search(tra) {
		fmt.Println("有权限")
	} else {
		fmt.Println("无权限")
		return false
	}
	return true
}

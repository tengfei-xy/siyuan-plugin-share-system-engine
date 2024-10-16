package main

import "encoding/json"

type resStruct struct {
	Err  int    `json:"err"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func (r *resStruct) setOK(str string) *resStruct {
	r.Err = ERR_CODE_OK
	r.Msg = ERR_MSG_OK
	r.Data = str
	return r
}
func (r *resStruct) setErrJson() *resStruct {
	r.Err = 1
	r.Msg = "json解析错误"
	r.Data = ""
	return r
}
func (r *resStruct) setErrSystem() *resStruct {
	r.Err = ERR_CODE_SYSTEM
	r.Msg = ERR_MSG_SYSTEM
	r.Data = ""
	return r
}
func (r *resStruct) setNoShare() *resStruct {
	r.Err = 3
	r.Msg = "此页面没有共享"
	r.Data = ""
	return r
}
func (r *resStruct) setErrParam() *resStruct {
	r.Err = ERR_CODE_PARAM
	r.Msg = ERR_MSG_PARAM
	r.Data = ""
	return r
}
func (r *resStruct) setErrURL() *resStruct {
	r.Err = 5
	r.Msg = " 错误请求地址"
	r.Data = ""
	return r
}

// func resOK() resStruct{
// 	return resStruct{
// 		Err : 0,
// 		Msg : "处理完成",
// 		Data : "",
// 	}
// }

// json to string
func (r *resStruct) toString() string {
	v, err := json.Marshal(r)
	if err != nil {
		r.setErrJson()
	}
	return string(v)
}

// json to msg
func (r *resStruct) toMsg() string {
	return r.Msg
}

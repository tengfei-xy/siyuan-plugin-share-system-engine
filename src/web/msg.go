package web

import (
	"encoding/json"
	"net/http"
	"sys"

	"github.com/gin-gonic/gin"
)

type resquest struct {
	Err  int         `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//	func ok(c *gin.Context) {
//		c.JSON(http.StatusOK, resquest{
//			Err: 0,
//			Msg: "ok",
//		})
//	}
func okData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, resquest{
		Err:  0,
		Msg:  "ok",
		Data: data,
	})
}
func internalSystem(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, resquest{
		Err:  1,
		Msg:  "系统错误",
		Data: nil,
	})
}
func badRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, resquest{
		Err:  2,
		Msg:  "参数错误",
		Data: nil,
	})
}

func noShare(c *gin.Context) {
	c.JSON(http.StatusOK, resquest{
		Err:  3,
		Msg:  "此页面没有共享",
		Data: nil,
	})
}

// 以下是旧版的结构体
type resStruct struct {
	Err  int    `json:"err"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type resArray struct {
	Err  int         `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (r *resStruct) setOK(str string) *resStruct {
	r.Err = sys.ERR_CODE_OK
	r.Msg = sys.ERR_MSG_OK
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
	r.Err = sys.ERR_CODE_SYSTEM
	r.Msg = sys.ERR_MSG_SYSTEM
	r.Data = ""
	return r
}
func (r *resStruct) setErrParam() *resStruct {
	r.Err = sys.ERR_CODE_PARAM
	r.Msg = sys.ERR_MSG_PARAM
	r.Data = ""
	return r
}
func (r *resStruct) setErrURL() *resStruct {
	r.Err = 5
	r.Msg = "错误请求地址"
	r.Data = ""
	return r
}
func (r *resStruct) setNoPage() *resStruct {
	r.Err = 6
	r.Msg = "分享文档不存在或未设置首页或"
	r.Data = ""
	return r
}
func msgOK(d interface{}) resArray {
	return resArray{
		Err:  sys.ERR_CODE_OK,
		Msg:  sys.ERR_MSG_OK,
		Data: d,
	}
}
func msgInternalSystemErr() resArray {
	return resArray{
		Err:  sys.ERR_CODE_SYSTEM,
		Msg:  sys.ERR_MSG_SYSTEM,
		Data: nil,
	}
}

// json to string
func (r *resStruct) toString() string {
	v, err := json.Marshal(r)
	if err != nil {
		r.setErrJson()
	}
	return string(v)
}

// json to msg
// func (r *resStruct) toMsg() string {
// 	return r.Msg
// }

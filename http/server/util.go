package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"

	def "github.com/QSis/common/db"
	"github.com/QSis/common/obj"
	"github.com/gin-gonic/gin"
)

type User struct {
	Name string `json:"name"`
}

type State struct {
	actionAborted bool
	formParsed    bool
	query         *def.Query
	user          *User
	PayloadForm   *paramForm
}

type Context struct {
	*gin.Context
	*State
}

func (ctx *Context) TextOutput(text string) {
	if !ctx.IsActionAborted() {
		ctx.String(200, text)
		ctx.AbortAction()
	}
}

func (ctx *Context) notFound() {
	if !ctx.IsActionAborted() {
		ctx.String(404, "Not found")
		ctx.AbortAction()
	}
}

func (ctx *Context) ShowPanicMsg(msg interface{}) {
	ctx.JsonFailure(msg)
}

func (ctx *Context) JsonOutput(data interface{}, code ...int) {
	origin := ctx.Request.Header.Get("Origin")
	allDomain := []string{
		"\\.*\\.com(:\\d*)?$",
	}
	if !ctx.IsActionAborted() {
		for _, v := range allDomain {
			result, _ := regexp.MatchString(v, origin)
			if result {
				ctx.Writer.Header().Set("Aretess-Control-Allow-Origin", origin)
				headers := ctx.Request.Header.Get("Aretess-Control-Request-Headers")
				ctx.Writer.Header().Set("Aretess-Control-Allow-Methods", "POST, PUT, GET, OPTIONS, DELETE")
				ctx.Writer.Header().Set("Aretess-Control-Allow-Credentials", "true")
				if headers != "" {
					ctx.Writer.Header().Set("Aretess-Control-Allow-Headers", headers)
				}
			}
		}
	}
	if len(code) > 0 {
		ctx.Status(code[0])
	} else {
		ctx.Status(200)
	}
	b, _ := json.Marshal(data)
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.Header().Set("Content-Length", obj.ToString(len(b)))
	ctx.Writer.Write(b)
	ctx.AbortAction()
}

func (ctx *Context) TraceInfo() gin.H {
	return gin.H{
		"id":        time.Now().UnixNano(),
		"timestamp": time.Now().Unix(),
		"srcIp":     ctx.ClientIP(),
		"destIp":    "",
	}
}

func (ctx *Context) JsonSuccess(ret ...interface{}) {
	h := gin.H{
		"code": 200,
		"msg":  "OK",
	}
	if len(ret) > 0 {
		h["item"] = ret[0]
	} else {
		h["trace"] = ctx.TraceInfo()
	}
	ctx.JsonOutput(h)
}

func (ctx *Context) JsonError(err error) {
	ctx.JsonFailure(err.Error())
}

func (ctx *Context) JsonFailure(ret ...interface{}) {
	h := gin.H{
		"trace": ctx.traceInfo(),
		"status": gin.H{
			"code": "ClientError",
		},
	}
	if len(ret) > 0 {
		h["status"]["msg"] = ret[0]
	}
	ctx.JsonOutput(h)
}

func (ctx *Context) JsonInvalidParam() {
	ctx.JsonError(errors.New("Invalid param"))
}

//正常情况不要是用这个
func CreateContext(gctx *gin.Context) *Context {
	ctx := &Context{gctx, &State{}}
	return ctx
}

func (ctx *Context) DontUseJsonSuccess(ret ...interface{}) {
	ctx.JsonSuccess(ret...)
}

func (state *State) IsActionAborted() bool {
	return state.actionAborted
}

func (state *State) AbortAction() {
	state.actionAborted = true
	panic(nil)
}

func (state *State) IsFormParsed() bool {
	return state.formParsed
}

func (state *State) SetFormParsed() {
	state.formParsed = true
}

func (state *State) SetUser(user *User) {
	state.user = user
}

func (state *State) GetUser() *User {
	return state.user
}

func (state *State) SetQuery(query *def.Query) {
	state.query = query
}

func (state *State) GetQuery() *def.Query {
	return state.query
}

func (state *State) SetPager(pager *def.Pager) {
	if state.query != nil {
		state.query.Pager = pager
	}
}

func (state *State) GetPager() *def.Pager {
	if state.query != nil {
		return state.query.Pager
	}
	return nil
}

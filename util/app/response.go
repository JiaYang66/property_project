package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"relocate/util/errors"
)

type Gin struct {
	Ctx *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (g *Gin) Response(httpCode int, errMsg string, data interface{}) {
	g.Ctx.JSON(httpCode, Response{
		Code: httpCode,
		Msg:  errMsg,
		Data: data,
	})
	return
}

func (g *Gin) SuccessResponse(data interface{}) {
	g.Ctx.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  "成功",
		Data: data,
	})
	return
}

func (g *Gin) BadResponse(data interface{}) {
	g.Ctx.JSON(http.StatusBadRequest, Response{
		Code: http.StatusBadRequest,
		Msg:  "参数校验失败",
		Data: data,
	})
	return
}

func (g *Gin) UnauthorizedResponse(data interface{}) {
	g.Ctx.JSON(http.StatusUnauthorized, Response{
		Code: http.StatusUnauthorized,
		Msg:  "请求未授权",
		Data: data,
	})
	return
}

func (g *Gin) VisitorsResponse(data interface{}) {
	g.Ctx.JSON(http.StatusNotFound, Response{
		Code: http.StatusNotFound,
		Msg:  "请求未授权",
		Data: data,
	})
	return
}

func (g *Gin) ErrorResponse(data interface{}) {
	g.Ctx.JSON(http.StatusInternalServerError, Response{
		Code: http.StatusInternalServerError,
		Msg:  "服务器内部错误",
		Data: data,
	})
	return
}

func (g *Gin) HasError(err error) bool {
	if err != nil {
		if code, ok := errors.Code(err); ok {
			//业务错误
			g.Response(code, "失败", err.Error())
		} else {
			//其他错误
			g.ErrorResponse(err.Error())
		}
		return true
	}
	return false
}

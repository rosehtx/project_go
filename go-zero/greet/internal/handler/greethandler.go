package handler

import (
	"net/http"
	"reflect"

	"github.com/zeromicro/go-zero/rest/httpx"
	"goZero/greet/internal/logic"
	"goZero/greet/internal/svc"
	"goZero/greet/internal/types"
)

func GreetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGreetLogic(r.Context(), svcCtx)
		resp, err := l.Greet(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func GreetHandlerByName(svcCtx *svc.ServiceContext,name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l 	   := logic.NewGreetLogic(r.Context(), svcCtx)
		op     := reflect.ValueOf(l)
		method := op.MethodByName(name)
		res    := method.Call([]reflect.Value{
			reflect.ValueOf(&req),
		})
		resOpt := res[0]
		resErr := res[1]
		resp := resOpt.Interface().(*types.Response)
		//这边会直接返回nil，
		err  := resErr.Interface()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err.(error))
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

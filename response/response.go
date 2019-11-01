package response

import (
        "context"
        "net/http"

        libCtx "github.com/ae-gis/volt/context"
)

var CtxResponse = libCtx.Key{Name: "response_context"}

type Pagination struct {
        Page  int `json:"page"`
        Size  int `json:"size"`
        Total int `json:"total"`
}

type MetaData struct {
        Code    string `json:"code,omitempty"`
        Type    string `json:"error_type,omitempty"`
        Message string `json:"error_message,omitempty"`
}

type result struct {
        Meta       interface{} `json:"meta,omitempty"`
        Data       interface{} `json:"data,omitempty"`
        Pagination interface{} `json:"pagination,omitempty"`
}

func New() *result {
        null := make(map[string]interface{})
        return &result{
                Meta:       null,
                Data:       null,
                Pagination: null,
        }
}

func (r *result) Error(err MetaData) {
        r.Meta = err
}

func (r *result) Success(code string) {
        r.Meta = MetaData{Code: code}
}

func (r *result) Body(body interface{}) {
        r.Data = body
}

func (r *result) Page(p Pagination) {
        r.Pagination = p
}

func Status(r *http.Request, status int) {
        *r = *r.WithContext(context.WithValue(r.Context(), CtxResponse, status))
}

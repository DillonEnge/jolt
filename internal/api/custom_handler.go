package api

import "net/http"

type HandlerFuncWithError func(w http.ResponseWriter, r *http.Request) *ApiError

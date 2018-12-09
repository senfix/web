package web

import (
	"encoding/json"
	"fmt"
	"github.com/senfix/php-status-page/pkg/errors"
	"github.com/senfix/php-status-page/src/constants"
	app "github.com/senfix/php-status-page/src/context"
	"github.com/senfix/php-status-page/src/web/response"
	"net/http"

	"io"

	"strconv"

	gorillaContext "github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func GetParamString(r *http.Request, name string) (value string, err error) {
	//parsing params from gorilla/mux url
	vars := mux.Vars(r)
	value, ok := vars[name]
	if ok {
		return
	}

	//fallback load from query string
	keys, ok := r.URL.Query()[name]

	if !ok || len(keys[0]) < 1 {
		err = errors.UrlMissingParam.Val(name)
		return
	}
	value = keys[0]

	return
}

func GetParamInt(r *http.Request, name string) (value int, err error) {
	str, err := GetParamString(r, name)
	if err != nil {
		return
	}
	value, err = strconv.Atoi(str)
	return
}

func Send(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		EmitError(w, http.StatusInternalServerError, err)
		return
	}
	json.NewEncoder(w).Encode(&resp)
}

func Decode(w http.ResponseWriter, body io.Reader, data interface{}) (err error) {
	decoder := json.NewDecoder(body)
	err = decoder.Decode(&data)
	if err != nil {
		EmitError(w, http.StatusUnprocessableEntity, err)
	}
	return
}

func LoadContext(req *http.Request) (err error, ctx app.Context) {
	cUser, ok := gorillaContext.GetOk(req, constants.CTX_USER)

	if !ok {
		err = errors.ApiInvalidUser
		return
	}

	ctx = cUser.(app.Context)
	return
}

func EmitError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&response.Error{code, fmt.Sprintf("%v", err)})
}

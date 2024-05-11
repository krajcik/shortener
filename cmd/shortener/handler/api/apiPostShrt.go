package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func (h *PostShrtHandler) PostShrt(w http.ResponseWriter, r *http.Request) {
	rm, done := h.loadModel(w, r)
	if !done {
		return
	}

	res, err := h.S.ShrtByURL(rm.URL)
	if err != nil {
		h.errServer(w, errors.Join(errors.New("ShrtByURL service"), err))

		return
	}

	var host string
	if h.P == nil || h.P.B == "" {
		host = "http://" + r.Host
	} else {
		host = h.P.B
	}
	res = fmt.Sprintf("%s/%s", host, res)
	w.WriteHeader(http.StatusCreated)

	respModel := PostShrtRespModel{ShrtURL: res}

	result, err := json.Marshal(respModel)
	if err != nil {
		h.errServer(w, err)
		return
	}
	_, err = w.Write(result)
	if err != nil {
		h.errServer(w, err)
		return
	}
}

func (h *PostShrtHandler) loadModel(w http.ResponseWriter, r *http.Request) (*PostShrtReqModel, bool) {
	rm := PostShrtReqModel{}
	err := json.NewDecoder(r.Body).Decode(&rm)
	if err != nil {
		h.L.Error("failed to decode request body", zap.Error(err))
		http.Error(w, "invalid json", http.StatusBadRequest)
		return nil, false
	}
	return &rm, true
}

func (h *PostShrtHandler) errServer(w http.ResponseWriter, err error) {
	h.L.Error("apiPostShrt:", zap.Error(err))
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	panic(err)
}

type PostShrtReqModel struct {
	URL string `json:"url"`
}

type PostShrtRespModel struct {
	ShrtURL string `json:"result"`
}

package web

import (
	"net/http"
	"io/ioutil"
	"errors"
)

type Request struct {
	Url string
	body []byte
	Err error
}

func (r *Request) Get() {

	resp, err := http.Get(r.Url)

	if err != nil {
		r.Err = err
		return
	}

	if resp.StatusCode != http.StatusOK {
		r.Err = errors.New(resp.Status)
		return
	}

	defer resp.Body.Close()

	r.body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		r.Err = err
		return
	}

}

func (r *Request) IsOK() (bool) {
	return r.Err == nil
}

func (r *Request) ToString() (string) {
	return string(r.body)
}

func (r *Request) Reset() {
	r.Url  = ""
	r.body = nil
	r.Err  = nil

}


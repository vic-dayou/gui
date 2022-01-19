package password

import (
	"errors"
	"time"
)

type Password struct {
	K          string
	V          string
	Pk         interface{}
	Ext        string
	ExpireTime int64
}

var pmap map[string]*Password

func init() {
	pmap = make(map[string]*Password)
}

func Get(k string) *Password {
	if s, ok := pmap[k]; ok {
		if s.ExpireTime <= time.Now().Unix() {
			return s
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func Put(p *Password) error {
	if p == nil {
		return errors.New("nil pointer")
	}

	pmap[p.K] = p
	return nil
}

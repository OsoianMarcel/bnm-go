package bnm_test

import (
	"testing"
	"github.com/osoianmarcel/bnm"
	"time"
)

// Test Request() method
func TestInst_Request_Success(t *testing.T) {
	inst := bnm.NewBnm()
	res, err := inst.Request(bnm.NewQuery("ro", time.Now()))
	if err != nil {
		t.Error(err)
		return
	}

	if len(res.Rates) == 0 {
		t.Error("empty rates")
	}
}


func TestInst_Request_Error(t *testing.T) {
	inst := bnm.NewBnm()
	_, err := inst.Request(bnm.Query{})

	if err == nil {
		t.Error(err)
		return
	}
}
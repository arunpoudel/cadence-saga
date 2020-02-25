package cadence_saga

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type sageAction struct {
	Fail bool
}

func (s sageAction) Act(_ context.Context) error {
	if s.Fail == true {
		return errors.New("some random error in action")
	}
	return nil
}

type compensator struct {
	Fail bool
}

func (c compensator) Compensate(_ context.Context) error {
	if c.Fail == true {
		return errors.New("some random error in compensator")
	}
	return nil
}

func Test_CompensationDoesNotRunOnActionSuccess(t *testing.T) {
	ctx := context.Background()
	saga := NewSaga()
	err := saga.Action(sageAction{Fail: false}).WithCompensation(compensator{}).Run(ctx)
	if err != nil {
		t.Errorf("error: %v", err)
	}
}

func Test_CompensationRunsOnActionFailure(t *testing.T) {
	ctx := context.Background()
	saga := NewSaga()
	err := saga.Action(sageAction{Fail: true}).WithCompensation(compensator{}).Run(ctx)
	if err == nil {
		t.Errorf("expected error to be present")
	}
}

func TestSaga_ParallelCompensation(t *testing.T) {
	ctx := context.Background()
	saga := NewSaga()
	err := saga.Action(sageAction{Fail: true}).WithCompensation(compensator{Fail: true}).ParallelCompensation(true).Run(ctx)
	if err == nil {
		t.Errorf("expected error to be present")
	}
}

func TestSaga_ParallelCompensation2(t *testing.T) {
	ctx := context.Background()
	saga := NewSaga()
	err := saga.Action(sageAction{Fail: true}).WithCompensation(compensator{Fail: true}).ParallelCompensation(false).Run(ctx)
	if err == nil {
		t.Errorf("expected error to be present")
	}
}

func TestSaga_ContinueWithCompensationError(t *testing.T) {
	ctx := context.Background()
	saga := NewSaga()
	err := saga.Action(sageAction{Fail: true}).WithCompensation(compensator{Fail: true}).ParallelCompensation(false).ContinueWithCompensationError(true).Run(ctx)
	_, ok := err.(CompensationError)
	if ok {
		t.Errorf("compensation error not supposed to be present")
	}
	if err == nil {
		t.Errorf("expected error to be present")
	}
}

func TestSaga_ContinueWithCompensationError2(t *testing.T) {
	ctx := context.Background()
	saga := NewSaga()
	err := saga.Action(sageAction{Fail: true}).WithCompensation(compensator{Fail: true}).ParallelCompensation(false).ContinueWithCompensationError(false).Run(ctx)
	compErr, ok := err.(CompensationError)
	if !ok {
		t.Errorf("expected error to be CompensationError, got %v, %s", reflect.TypeOf(err), err.Error())
	}
	actErr := compErr.ActionError()
	if actErr == nil {
		t.Errorf("compensation was triggered without action error")
	}
}

package cadence_saga

import (
	"context"
)

type (
	action interface {
		Act(context.Context) error
	}

	compensation interface {
		Compensate(context.Context) error
	}

	saga struct {
		parallelCompensation          bool
		continueWithCompensationError bool
		action                        action
		compensations                 []compensation
	}
)

func (s *saga) Action(act action) *saga {
	s.action = act
	return s
}

func (s *saga) WithCompensation(comp compensation) *saga {
	s.compensations = append(s.compensations, comp)
	return s
}

func (s *saga) Run(ctx context.Context) error {
	err := s.action.Act(ctx)
	if err != nil {
		s.Compensate(ctx)
	}
	return err
}

func (s *saga) Compensate(ctx context.Context) {
	for _, compensation := range s.compensations {
		err := compensation.Compensate(ctx)
		if err != nil && s.parallelCompensation != true && s.continueWithCompensationError == false {
			panic(err)
		}
	}
}

func (s *saga) ContinueWithCompensationError(continueWithCompensationError bool) *saga {
	s.continueWithCompensationError = continueWithCompensationError
	return s
}

func (s *saga) ParallelCompensation(parallelCompesation bool) *saga {
	s.parallelCompensation = parallelCompesation
	return s
}

func NewSaga() saga {
	return saga{}
}

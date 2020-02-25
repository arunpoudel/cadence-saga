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

func compensate(ctx context.Context, c compensation, continueWithCompensationError bool) error {
	err := c.Compensate(ctx)
	if err != nil && continueWithCompensationError == false {
		return err
	}
	return nil
}

func (s *saga) Compensate(ctx context.Context) {
	if s.parallelCompensation == true {
		for _, compensation := range s.compensations {
			func() {
				err := compensate(ctx, compensation, s.continueWithCompensationError)
				if err != nil {
					panic(err)
				}
			}()
		}
		return
	}
	for _, compensation := range s.compensations {
		err := compensate(ctx, compensation, s.continueWithCompensationError)
		if err != nil {
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

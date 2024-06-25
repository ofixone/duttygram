package internal

type Stage int

const (
	StageGreetings Stage = iota
	StageMain
)

type Scenario struct {
	Stages []Stage
}

func (s *Scenario) GetStartStage() Stage {
	return StageGreetings
}

func (s *Scenario) GetMainStage() Stage {
	return StageMain
}

func NewScenario() *Scenario {
	return &Scenario{Stages: []Stage{
		StageGreetings,
		StageMain,
	}}
}

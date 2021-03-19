package main

import (
	"fmt"
)

func run_simulation(
	initial_state Variables,
	parameters Parameters,
	timestep_block TimestepBlock,
	timesteps int,
	samples int) RunResults {

	results := make(RunResults, samples)
	for sample := 0; sample < samples; sample++ {
		history := make(History)
		state := initial_state
		new_state := state
		for timestep := 0; timestep < timesteps; timestep++ {
			substep_history := make(SubstepHistory)
			for substep, substep_block := range timestep_block {
				state.timestep = timestep
				state.substep = substep
				new_state = state
				for _, variable_fn := range substep_block {
					new_state = variable_fn(parameters, uint32(substep), history, state)
				}
				substep_history[substep] = new_state
				state = new_state
			}
			history[timestep] = substep_history
		}
		results[sample] = history
	}
	return results
}

type RunResults = []History

type VariableUpdate = func(
	Parameters,
	uint32,
	History,
	Variables) Variables

type SubstepBlock = map[string]VariableUpdate

type TimestepBlock []SubstepBlock

type Variables struct {
	timestep            int
	substep             int
	prey_population     float32
	predator_population float32
}

type Parameters struct {
	increase_size float32
}

type History = map[int]SubstepHistory

type SubstepHistory = map[int]Variables

func s_prey_population(p Parameters,
	substep uint32,
	history History,
	state Variables) Variables {
	state.prey_population += 1
	return state
}

func s_predator_population(p Parameters,
	substep uint32,
	history History,
	state Variables) Variables {
	state.predator_population += 2
	return state
}

func main() {

	initial_state := Variables{
		prey_population:     10,
		predator_population: 5,
	}

	params := Parameters{
		increase_size: 1,
	}

	substep_block_1 := SubstepBlock{
		"prey_population":     s_prey_population,
		"predator_population": s_predator_population,
	}

	substep_block_2 := SubstepBlock{
		"prey_population": s_prey_population,
	}

	timestep_block := TimestepBlock{
		substep_block_1,
		substep_block_2,
	}

	var TIMESTEPS int = 10
	var SAMPLES int = 2

	result := run_simulation(initial_state, params, timestep_block, TIMESTEPS, SAMPLES)

	for _, sample := range result {
		for _, timestep := range sample {
			fmt.Println(timestep)
		}
		fmt.Println("---")
	}
}

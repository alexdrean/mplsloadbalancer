package state

import "mplsloadbalancer/config"

type State struct {
	Paths []Path
	Capacity uint64
}

type Path struct {
	Links []Link
	Capacity uint64
}

type Link struct {
	Label uint32
	Capacity uint64
}

type Compiler struct {
	paths [][]link
}

func NewCompiler(config config.Config) (error, *Compiler) {
	compiler := new(Compiler)
	err := compiler.loadConfig(config)
	if err != nil {
		return err, nil
	}
	return nil, compiler
}

func (compiler Compiler) GetState() *State {
	// Query all statuses
	for _, path := range compiler.paths {
		for _, link := range path {
			link.radio.StatusChan <- nil
		}
	}

	// Read all statuses into State
	var fabricCapacity uint64 = -1
	paths := make([]Path, len(compiler.paths))
	for i, path := range compiler.paths {
		links := make([]Link, len(path))
		var pathCapacity uint64
		for j, pathLink := range path {
			var linkCapacity uint64
			if res := <- pathLink.radio.StatusChan; res != nil {
				if res.CapacityMatters > 0 {
					linkCapacity = uint64(res.CapacityMatters)
					pathCapacity += linkCapacity
				}
			}
			links[j] = Link{
				Label:    pathLink.label,
				Capacity: linkCapacity,
			}
		}
		paths[i] = Path{
			Links: links,
			Capacity: pathCapacity,
		}
		if fabricCapacity == -1 { // Initialize
			fabricCapacity = pathCapacity
		} else if pathCapacity < fabricCapacity { // Replace if lower. Total capacity = lowest capacity of the fabric
			fabricCapacity = pathCapacity
		}
	}
	state := State{
		Paths: paths,
		Capacity: fabricCapacity,
	}
	return &state
}

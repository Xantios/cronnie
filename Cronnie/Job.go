package Cronnie

type Job func(map[string]string) bool

func (ci *Instance) executeRunnerFunction(name string, args map[string]string) bool {
	for key, value := range ci.jobMap {
		if key == name {
			return value(args)
		}
	}

	return false
}

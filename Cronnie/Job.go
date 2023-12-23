package Cronnie

type Job func(map[string]string) bool

var jobMap map[string]Job

// Runners Set runners
func Runners(jobs map[string]Job) {
	jobMap = jobs
}

func exectueRunnerFunction(name string, args map[string]string) bool {
	for key, value := range jobMap {
		if key == name {
			return value(args)
		}
	}

	return false
}

package numatopo

const (
	policyNone = "none"
)

type numaTopoName string

const (
	memoryNumaTopoName numaTopoName = "memory"
	cpuNumaTopoName    numaTopoName = "cpu"
)

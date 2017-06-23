package checks

type Result struct {
	err      error
	hostPort string
}

func (r Result) String() string {
	if r.err != nil {
		return r.hostPort + " " + r.err.Error()
	}
	return r.hostPort + " All right"
}

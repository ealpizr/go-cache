package types

type Job struct {
	value interface{}
}

func NewJob(v interface{}) *Job {
	return &Job{v}
}

package main

type JobParams struct {
}

type Job struct {
	Action string    `json:"action"`
	Params JobParams `json:"params"`
}

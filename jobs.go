package main

type MasterJob interface {
	GetName() string
}

type GenerateWorldJob struct {
	MasterJob
}

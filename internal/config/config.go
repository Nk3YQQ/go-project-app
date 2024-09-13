package config

type GenderChoice string

const (
	Male   GenderChoice = "male"
	Female GenderChoice = "female"
)

type StatusChoice string

const (
	Created   StatusChoice = "created"
	InProcess StatusChoice = "in_process"
	Completed StatusChoice = "completed"
	Expired   StatusChoice = "expired"
)

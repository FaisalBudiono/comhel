package domain

type Status int

const (
	Paused Status = iota
	Restarting
	Removing
	Running
	Dead
	Created
	Exited
)

func (s Status) String() string {
	switch s {
	case Paused:
		return "paused"
	case Restarting:
		return "restarting"
	case Removing:
		return "removing"
	case Running:
		return "running"
	case Dead:
		return "dead"
	case Created:
		return "created"
	case Exited:
		return "exited"
	default:
		return ""
	}
}

func StatusFrom(s string) Status {
	switch s {
	case "paused":
		return Paused
	case "restarting":
		return Restarting
	case "removing":
		return Removing
	case "running":
		return Running
	case "dead":
		return Dead
	case "created":
		return Created
	case "exited":
		return Exited
	default:
		return -1
	}
}

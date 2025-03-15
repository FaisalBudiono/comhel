package cmdmain

import "github.com/FaisalBudiono/comhel/internal/app/domain"

func offService(serviceName string) renderableService {
	return renderableService{
		name:   serviceName,
		status: "off",
	}
}

func fromDomain(s domain.Service) renderableService {
	return renderableService{
		name:   s.Name,
		status: s.Status.String(),
	}
}

type renderableService struct {
	name   string
	status string
}

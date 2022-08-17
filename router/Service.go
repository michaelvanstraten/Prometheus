package router

import "net/http"

type Service struct {
	MiddleServices 	[]*Service
	ServeHTTP 		func(http.ResponseWriter, *http.Request) bool
}

func (s *Service) AddMiddelware(Middelwares ...*Service) {
	for i := 0; i < len(Middelwares); i++ {
		s.MiddleServices 		= append(s.MiddleServices, Middelwares[i])
		copy(s.MiddleServices[1:], s.MiddleServices)
		s.MiddleServices[0] 	= Middelwares[i]
	}
} 
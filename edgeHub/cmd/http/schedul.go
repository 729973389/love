package http

import "github.com/wuff1996/edgeHub/config"

type Schedule struct {
	*config.Url
	Action chan string
}

func NewSchedule()*Schedule{
	return &Schedule{GetConfig(),make(chan string,10)}
}

func Run(){
	for  {


	}
}
package server

func (s *server) routes() {
	s.router.HandleFunc("/", s.handleSomething())
	s.router.HandleFunc("/info/{repository}/{chart}", s.handleChartInfo())
	s.router.HandleFunc("/values/{repository}/{chart}", s.handleRenderValues())
}

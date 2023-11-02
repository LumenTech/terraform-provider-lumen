package bare_metal

type Server struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (s Server) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":   s.Name,
		"status": s.Status,
	}
}

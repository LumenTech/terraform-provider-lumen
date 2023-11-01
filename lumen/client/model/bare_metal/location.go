package bare_metal

type Location struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Region string `json:"region"`
}

func (l Location) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":     l.ID,
		"name":   l.Name,
		"status": l.Status,
		"region": l.Region,
	}
}

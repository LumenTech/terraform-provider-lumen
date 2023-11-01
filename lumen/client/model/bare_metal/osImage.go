package bare_metal

type OsImage struct {
	Name  string `json:"name"`
	Ready bool   `json:"ready"`
	Price Price  `json:"price"`
}

func (n OsImage) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":  n.Name,
		"ready": n.Ready,
		"price": n.Price.String(),
	}
}

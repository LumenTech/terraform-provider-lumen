package bare_metal

import "fmt"

type OsImage struct {
	Name   string      `json:"name"`
	Ready  bool        `json:"ready"`
	Price  Price       `json:"price"`
	Prices []TierPrice `json:"prices"`
}

func ConvertOSImagesToListMap(images []OsImage) []map[string]interface{} {
	var mapList []map[string]interface{}
	for _, image := range images {
		if len(image.Prices) == 0 {
			mapList = append(mapList, map[string]interface{}{
				"name":  image.Name,
				"tier":  "Default (All)",
				"price": image.Price.String(),
			})
		} else {
			for _, price := range image.Prices {
				tier := "Default (All)"
				if price.Tier != nil {
					tier = fmt.Sprintf("%d", *price.Tier)
				}

				mapList = append(mapList, map[string]interface{}{
					"name":  image.Name,
					"tier":  tier,
					"price": price.String(),
				})
			}
		}
	}
	return mapList
}

package bare_metal

type MappableStruct interface {
	ToMap() map[string]interface{}
}

func ConvertToListMap[T MappableStruct](list []T) []map[string]interface{} {
	mapList := make([]map[string]interface{}, len(list))
	for index, item := range list {
		mapList[index] = item.ToMap()
	}
	return mapList
}

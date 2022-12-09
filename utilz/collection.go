package utilz

func IsInArray(resultStatus int, array []int) bool {
	for _, status := range array {
		if status == resultStatus {
			return true
		}
	}
	return false
}
func MapMerge(mp ...map[string]string) map[string]string {
	m2 := make(map[string]string)
	for _, m1 := range mp {
		for k, v := range m1 {
			m2[k] = v
		}
	}
	return m2
}

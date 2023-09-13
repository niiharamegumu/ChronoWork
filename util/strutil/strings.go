package strutil

func RemoveDuplicates(input []string) []string {
	encountered := map[string]int{}
	result := []string{}
	for _, v := range input {
		encountered[v]++
	}
	for _, v := range input {
		if encountered[v] == 1 {
			result = append(result, v)
		}
	}
	return result
}

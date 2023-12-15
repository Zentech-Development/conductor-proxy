package handlers

func checkForGroupMatch(userGroups []string, allowedGroups []string) bool {
	vals := make(map[string]bool)

	for _, group := range userGroups {
		vals[group] = true

		if group == "admin" {
			return true
		}
	}

	for _, group := range allowedGroups {
		if _, present := vals[group]; present {
			return true
		}
	}

	return false
}

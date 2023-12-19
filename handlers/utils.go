package handlers

import "slices"

func checkForGroupMatch(userGroups []string, allowedGroups []string) bool {
	for _, group := range allowedGroups {
		if slices.Contains(userGroups, group) {
			return true
		}
	}

	return false
}

func isAdmin(userGroups []string) bool {
	for _, group := range userGroups {
		if group == "admin" {
			return true
		}
	}

	return false
}

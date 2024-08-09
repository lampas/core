package models

func localize(lang string, ru string, uk string) string {
	if ru == "" || (lang == "uk" && uk != "") {
		return uk
	}

	return ru
}

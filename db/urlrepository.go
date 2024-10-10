package db

func CreateURL(shortId string, longURL string) error {
	err := storeURL(shortId, longURL)
	return err
}

func LoadFromDB(shortId string) (string, error) {
	longURL, err := getLong(shortId)
	return longURL, err
}

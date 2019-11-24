package db

func Begin() *MyDB {
	tx := DB.Begin()
	return &MyDB{tx}
}

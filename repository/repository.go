package repository

const (
	selectAllQuery  = "select * from %v"
	selectByIDQuery = "select * from %v where id=%v"
	findByQuery     = "select * from %v where %v in (%v)"
	updateQuery     = "update %v set %v=%v where id=%v"
)

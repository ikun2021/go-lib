package main

import "gorm.io/gen"

type Method interface {

	//SELECT * FROM @@table WHERE id = @id
	GetByID(id int) (gen.T, error)
}

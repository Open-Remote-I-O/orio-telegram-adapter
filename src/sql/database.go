package sql

type IOrioDatabase interface{}

type OrioDatabase struct{}

func (odb *OrioDatabase) testQuery(){}

package preload

type Service interface {
	PopulateMongoDB() error
}

type Repository interface {
	PopulateMongoDB() error
}
package preload

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) PopulateMongoDB() error {
	err := s.repository.PopulateMongoDB()
	if err != nil {
		return err
	}

	return nil
}
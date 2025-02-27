package chatrooms

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) FindChatroomByCode(code string) (*Chatroom, error) {
	chatroom, err := s.repository.FindChatroomByCode(code)
	if err != nil {
		return nil, err
	}

	return chatroom, nil
}

func (s *service) GetChatrooms() ([]Chatroom, error) {
	chatrooms, err := s.repository.GetChatrooms()
	if err != nil {
		return nil, err
	}

	return chatrooms, nil
}
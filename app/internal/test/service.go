package test

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) SetValue(key, value string) {
	s.repo.Save(key, value)
}

func (s *Service) GetValue(key string) (string, bool) {
	return s.repo.Get(key)
}

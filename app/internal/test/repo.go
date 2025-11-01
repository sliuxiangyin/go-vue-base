package test

type Repo struct {
	data map[string]string
}

func NewRepo() *Repo {
	return &Repo{
		data: make(map[string]string),
	}
}

func (r *Repo) Save(key, value string) {
	r.data[key] = value
}

func (r *Repo) Get(key string) (string, bool) {
	val, ok := r.data[key]
	return val, ok
}

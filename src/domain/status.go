package domain

type (
	Status struct {
		Tasks Stats
		Urls  Stats
	}

	Stats struct {
		Complete int
		Error    int
		InQueue  int
	}
)

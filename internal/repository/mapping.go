package repository

type dataMapper interface {
	DataTo(p interface{}) error
}

func toPredictions(docs []dataMapper) ([]Prediction, error) {
	all := make([]Prediction, 0, len(docs))
	for _, doc := range docs {
		var p Prediction
		if err := doc.DataTo(&p); err != nil {
			return nil, err
		}
		all = append(all, p)
	}
	return all, nil
}

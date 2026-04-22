package repository

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type dataMapper interface {
	DataTo(p interface{}) error
}

func isNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
}

func toPrediction(doc dataMapper) (*Prediction, error) {
	var p Prediction
	if err := doc.DataTo(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func toPredictions(docs []dataMapper) ([]Prediction, error) {
	all := make([]Prediction, 0, len(docs))
	for _, doc := range docs {
		p, err := toPrediction(doc)
		if err != nil {
			return nil, err
		}
		all = append(all, *p)
	}
	return all, nil
}

func toResult(doc dataMapper) (*Result, error) {
	var r Result
	if err := doc.DataTo(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

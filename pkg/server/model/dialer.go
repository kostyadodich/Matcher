package model

import "github.kostyadodich/demo/pkg/model"

type Dialer struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	PhoneNumber int     `json:"phone_number"`
	ItemsIDs    []int64 `json:"item"`
}

func (d *Dialer) Convert() model.Dialer {
	dialers := model.Dialer{
		ID:          d.ID,
		Name:        d.Name,
		PhoneNumber: d.PhoneNumber,
	}

	for _, v := range d.ItemsIDs {
		dialers.Items = append(dialers.Items, model.Item{ID: v})
	}

	return dialers
}

package usecase

import "github.com/usk81/ashihara/services/holidays/core/domain/entity"

var defaultFields = []string{
	"date",
	"name",
	"difinition_id",
}

func dropFields(h *entity.Holiday, fields []string) *entity.Holiday {
	fs := fields
	if h == nil {
		return h
	}
	if len(fs) == 0 {
		fs = defaultFields
	}
	r := &entity.Holiday{}
	for _, f := range fs {
		switch f {
		case "date":
			r.Date = h.Date
		case "name":
			r.Name = h.Name
		case "difinition_id":
			r.DifinitionID = h.DifinitionID
		case "summary":
			r.Summary = h.Summary
		case "description":
			r.Description = h.Description
		}
	}
	return r
}

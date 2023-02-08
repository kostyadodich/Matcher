package matcher

import (
	"database/sql"
	"github.kostyadodich/demo/pkg/model"
)

type Matcher interface {
	Match(user model.User) (model.Dialer, error)
}

type DefaultMatcher struct {
	db *sql.DB
}

func NewDefaultMatcher(db *sql.DB) *DefaultMatcher {
	return &DefaultMatcher{db: db}
}

func (d *DefaultMatcher) Match(user model.User) (model.Dialer, error) {
	return model.Dialer{}, nil
}

package auth

import app "azarole/internal"

type FinderResult struct {
	UserId int64
}

type UserFinder struct {
	application *app.Application
	identifier  string
}

func NewUserFinder(application *app.Application, identifier string) *UserFinder {
	return &UserFinder{
		application: application,
		identifier:  identifier,
	}
}

func (finder *UserFinder) Execute() (*FinderResult, error) {
	result := FinderResult{
		UserId: 1,
	}
	return &result, nil
}

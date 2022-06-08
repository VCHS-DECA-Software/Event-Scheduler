package users

import "main/components/links"

type Admin struct {
	Name string
}

type Judge struct {
	Name string

	AssignedTeams links.Link
	Events        links.Link
}

type Student struct {
	Name string

	Teams links.Link
}

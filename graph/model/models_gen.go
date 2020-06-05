// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AddRole struct {
	Name        string    `json:"name"`
	Permissions []*string `json:"permissions"`
	Parents     []*string `json:"parents"`
}

type DeletePermission struct {
	Name       string `json:"name"`
	Permission string `json:"permission"`
}

type DeleteRole struct {
	Name string `json:"name"`
}

type Jwt struct {
	User       string      `json:"user"`
	Roles      []string    `json:"roles"`
	Properties []*Property `json:"properties"`
}

type NewJwt struct {
	User  string   `json:"user"`
	Roles []string `json:"roles"`
}

type Property struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Role struct {
	Name        string    `json:"name"`
	Permissions []*string `json:"permissions"`
	Parents     []*string `json:"parents"`
}

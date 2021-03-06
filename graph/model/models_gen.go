// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type AddPhoto struct {
	Newspaper string `json:"newspaper"`
	Caption   string `json:"caption"`
	Filename  string `json:"filename"`
}

type AddRole struct {
	Name        string    `json:"name"`
	Permissions []*string `json:"permissions"`
	Parents     []*string `json:"parents"`
}

type AddStory struct {
	Newspaper string `json:"newspaper"`
	Headline  string `json:"headline"`
	Story     string `json:"story"`
}

type DeleteMedia struct {
	Newspaper string `json:"newspaper"`
	UUID      string `json:"uuid"`
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

type ModStaff struct {
	Newspaper string `json:"newspaper"`
	Name      string `json:"name"`
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

type Domain string

const (
	DomainNewspaper Domain = "newspaper"
)

var AllDomain = []Domain{
	DomainNewspaper,
}

func (e Domain) IsValid() bool {
	switch e {
	case DomainNewspaper:
		return true
	}
	return false
}

func (e Domain) String() string {
	return string(e)
}

func (e *Domain) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Domain(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DOMAIN", str)
	}
	return nil
}

func (e Domain) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Rbac string

const (
	RbacJwtQuery     Rbac = "JWT_QUERY"
	RbacJwtMutate    Rbac = "JWT_MUTATE"
	RbacRbacQuery    Rbac = "RBAC_QUERY"
	RbacRbacMutate   Rbac = "RBAC_MUTATE"
	RbacModNewspaper Rbac = "MOD_NEWSPAPER"
	RbacModStaff     Rbac = "MOD_STAFF"
	RbacModStory     Rbac = "MOD_STORY"
	RbacModPhoto     Rbac = "MOD_PHOTO"
	RbacDelMedia     Rbac = "DEL_MEDIA"
)

var AllRbac = []Rbac{
	RbacJwtQuery,
	RbacJwtMutate,
	RbacRbacQuery,
	RbacRbacMutate,
	RbacModNewspaper,
	RbacModStaff,
	RbacModStory,
	RbacModPhoto,
	RbacDelMedia,
}

func (e Rbac) IsValid() bool {
	switch e {
	case RbacJwtQuery, RbacJwtMutate, RbacRbacQuery, RbacRbacMutate, RbacModNewspaper, RbacModStaff, RbacModStory, RbacModPhoto, RbacDelMedia:
		return true
	}
	return false
}

func (e Rbac) String() string {
	return string(e)
}

func (e *Rbac) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Rbac(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RBAC", str)
	}
	return nil
}

func (e Rbac) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Product struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Price       int      `json:"price"`
	Preview     string   `json:"preview"`
	Images      []string `json:"images"`
	Tags        []string `json:"tags"`
	Description *string  `json:"description"`
}

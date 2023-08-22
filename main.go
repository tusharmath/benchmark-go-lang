package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// Post structure
type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"userId": &graphql.Field{
			Type: graphql.Int,
		},
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"body": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var client = resty.New()

func postsResolver(p graphql.ResolveParams) (interface{}, error) {

	var posts []Post
	resp, err := client.R().SetResult(&posts).Get("http://127.0.0.1:3000/posts")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.Status())
	}
	return posts, nil
}

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"posts": &graphql.Field{
			Type:        graphql.NewList(postType),
			Description: "Get list of posts",
			Resolve:     postsResolver,
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
})

func postResolver(p graphql.ResolveParams) (interface{}, error) {

	client := resty.New()
	var post Post
	resp, err := client.R().SetResult(&post).Get("http://127.0.0.1:3000/posts")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.Status())
	}
	return post, nil
}

func main() {
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

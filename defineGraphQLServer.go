package main

import (
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/bags2on/api-gateway/gclients"
	"github.com/bags2on/api-gateway/graph/generated"
)

type Server struct {
	productsClient *gclients.ProductsClient
}

func NewGraphQLServer() *Server {
	productsClient, err := gclients.GetProductsClient("localhost:50051")

	if err != nil {
		fmt.Println("cannot connect to products microservice")
		productsClient.Close()
	}

	return &Server{productsClient}
}

//srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

func (s *Server) Query() generated.QueryResolver {
	return &queryResolver{server: s}
}

//func (s *Server) Products() generated.QueryResolver {
//	return &queryResolver{server: s}
//}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: s,
	})
}

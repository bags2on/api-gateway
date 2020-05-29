package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/bags2on/api-gateway/gclients"
	"github.com/bags2on/api-gateway/graph/generated"
	"github.com/bags2on/api-gateway/graph/model"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

const (
	port = "8080"
)

type QueryResolver struct {
	server *Server
}

func (r *QueryResolver) Products(c context.Context) ([]*model.Product, error) {

	productsList, err := r.server.productsClient.GetProducts(c)
	if err != nil {
		fmt.Println("cannot get products list")
	}

	var products []*model.Product

	for _, value := range productsList {
		product := &model.Product{
			ID:      value.ID,
			Title:   value.Title,
			Price:   value.Price,
			Preview: value.Preview,
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *QueryResolver) Product(c context.Context, id string) (*model.Product, error) {
	requestedProduct, err := r.server.productsClient.GetProductByID(c, id)
	fmt.Println(id)
	if err != nil {
		fmt.Println("cannot get product by id")
	}

	return requestedProduct, nil

}

type Server struct {
	productsClient *gclients.ProductsClient
}

func Gozs() *Server {
	productsClient, err := gclients.GetProductsClient("localhost:50051")

	if err != nil {
		fmt.Println("cannot connect to products microservice")
		productsClient.Close()
	}

	return &Server{productsClient}
}

func (s *Server) Query() generated.QueryResolver {
	return &QueryResolver{server: s}
}

//func (s *Server) Products() generated.QueryResolver {
//	return &QueryResolver{server: s}
//}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: s,
	})
}

func main() {
	//srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}})
	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", srv)

	s := Gozs()

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	router.Handle("/", handler.NewDefaultServer(s.ToExecutableSchema()))
	router.Handle("/playground", playground.Handler("GraphQL playground", "/query"))

	log.Printf("connect to http://localhost:%s/playground for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

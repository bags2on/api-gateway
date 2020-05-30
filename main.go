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

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: s,
	})
}

func wrapGraphQl() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
}

func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		next.ServeHTTP(w, r)
	})
}

func main() {

	s := Gozs()
	r := http.NewServeMux()

	r.Handle("/", middlewareOne(handler.NewDefaultServer(s.ToExecutableSchema())))

	r.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	log.Printf("connect to http://localhost:%s/playground for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

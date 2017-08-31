/**
@file
Defines and initialize the GraphQL schema.
 */

package interfaces

import (
	"github.com/graphql-go/graphql"
	"strconv"
	"git.humbkr.com/jgalletta/alba-player/business"
)

type graphQLInteractor struct {
	Schema graphql.Schema
	Library *business.CollectionInteractor
}

func NewGraphQLInteractor(ci *business.CollectionInteractor) *graphQLInteractor {
	interactor := &graphQLInteractor{Library:ci}

	// Intanciate types.
	artistType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Artist",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
				Resolve: func (p graphql.ResolveParams) (interface{}, error) {
					return "1", nil
				},
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func (p graphql.ResolveParams) (interface{}, error) {
					return "Stub", nil
				},
			},
		},
	})

	/**
	 * This is the type that will be the root of our query,
	 * and the entry point into our schema.
	 */
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"artist": &graphql.Field{
				Type: artistType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "Artist ID",
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					i := p.Args["id"].(string)
					id, err := strconv.Atoi(i)
					if err != nil {
						return nil, err
					}
					return interactor.Library.ArtistRepository.Find(id)
				},
			},
		},
	})

	/**
	* Finally, we construct our schema (whose starting query type is the query
	* type we defined above) and export it.
	 */
	var err error
	interactor.Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
	if err != nil {
		panic(err)
	}

	return interactor
}

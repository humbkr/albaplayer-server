/**
@file
Defines and initialize the GraphQL schema.
 */

package interfaces

import (
	"github.com/graphql-go/graphql"
	"strconv"
	"git.humbkr.com/jgalletta/alba-player/business"
	"git.humbkr.com/jgalletta/alba-player/domain"
)

type graphQLInteractor struct {
	Schema graphql.Schema
	Library *business.LibraryInteractor
}

var artistType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Artist",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Name: "Artist ID",
			Description: "Artist unique Identifier.",
			Type: graphql.NewNonNull(graphql.ID),
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if artist, ok := p.Source.(domain.Artist); ok == true {
					return artist.Id, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Name: "Artist name",
			Description: "Name of the artist.",
			Type: graphql.NewNonNull(graphql.String),
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if artist, ok := p.Source.(domain.Artist); ok == true {
					return artist.Name, nil
				}
				return nil, nil
			},
		},
		"albums": &graphql.Field{
			Name: "Artist albums",
			Description: "Albums of the artist.",
			Type: graphql.NewList(graphql.NewNonNull(albumType)),
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if artist, ok := p.Source.(domain.Artist); ok == true {
					return artist.Albums, nil
				}
				return nil, nil
			},
		},
	},
})

var albumType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Album",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Name: "Album ID",
			Description: "Album unique identifier.",
			Type: graphql.NewNonNull(graphql.ID),
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if album, ok := p.Source.(domain.Album); ok == true {
					return album.Id, nil
				}
				return nil, nil
			},
		},
		"title": &graphql.Field{
			Name: "Album title",
			Description: "Title of the album.",
			Type: graphql.NewNonNull(graphql.String),
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if album, ok := p.Source.(domain.Album); ok == true {
					return album.Title, nil
				}
				return nil, nil
			},
		},
		"year": &graphql.Field{
			Name: "Album year",
			Description: "Year the album was released in, or the year-span in case of a compilation of tracks from released in different years.",
			Type: graphql.NewNonNull(graphql.String),
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if album, ok := p.Source.(domain.Album); ok == true {
					return album.Year, nil
				}
				return nil, nil
			},
		},
		"tracks": &graphql.Field{
			Name: "Album tracks",
			Description: "Tracks of album.",
			Type: graphql.NewList(graphql.NewNonNull(trackType)),
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if album, ok := p.Source.(domain.Album); ok == true {
					return album.Tracks, nil
				}
				return nil, nil
			},
		},
	},
})

var trackType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Track",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Name: "Track ID",
			Description: "Track unique Identifier.",
			Type: graphql.NewNonNull(graphql.ID),
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if track, ok := p.Source.(domain.Track); ok == true {
					return track.Id, nil
				}
				return nil, nil
			},
		},
		"title": &graphql.Field{
			Name: "Track title",
			Description: "Title of the track.",
			Type: graphql.NewNonNull(graphql.String),
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if track, ok := p.Source.(domain.Track); ok == true {
					return track.Title, nil
				}
				return nil, nil
			},
		},
		"disc": &graphql.Field{
			Name: "Track disc",
			Description: "If the album this track is on has multiple discs, specify the disc on which the track is on.",
			Type: graphql.String,
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if track, ok := p.Source.(domain.Track); ok == true {
					return track.Disc, nil
				}
				return nil, nil
			},
		},
		"number": &graphql.Field{
			Name: "Track number",
			Description: "Position of the track on the album or disc.",
			Type: graphql.Int,
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if track, ok := p.Source.(domain.Track); ok == true {
					return track.Number, nil
				}
				return nil, nil
			},
		},
		"duration": &graphql.Field{
			Name: "Track duration",
			Description: "Track duration in seconds.",
			Type: graphql.Int,
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if track, ok := p.Source.(domain.Track); ok == true {
					return track.Duration, nil
				}
				return nil, nil
			},
		},
		"genre": &graphql.Field{
			Name: "Track genre",
			Description: "Music genre.",
			Type: graphql.String,
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if track, ok := p.Source.(domain.Track); ok == true {
					return track.Genre, nil
				}
				return nil, nil
			},
		},
		"path": &graphql.Field{
			Name: "Track path",
			Description: "Localisation of the media file on the computer.",
			Type: graphql.NewNonNull(graphql.String),
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if track, ok := p.Source.(domain.Track); ok == true {
					return track.Path, nil
				}
				return nil, nil
			},
		},
	},
})

/**
Creates a new GraphQL interactor.

Builds GraphQL Schema, initialise dynamic fields on types.
 */
func NewGraphQLInteractor(ci *business.LibraryInteractor) *graphQLInteractor {
	interactor := &graphQLInteractor{Library:ci}

	/**
	 * This is the type that will be the root of our query,
	 * and the entry point into our schema.
	 */
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"artists": &graphql.Field{
				Type: graphql.NewList(artistType),
				Resolve: func (g graphql.ResolveParams) (interface{}, error) {
					return interactor.Library.ArtistRepository.FindAll(true)
				},
			},
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
			"albums": &graphql.Field{
				Type: graphql.NewList(albumType),
				Resolve: func (g graphql.ResolveParams) (interface{}, error) {
					return interactor.Library.AlbumRepository.FindAll(true)
				},
			},
			"album": &graphql.Field{
				Type: albumType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "Album ID",
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					i := p.Args["id"].(string)
					id, err := strconv.Atoi(i)
					if err != nil {
						return nil, err
					}

					return interactor.Library.AlbumRepository.Find(id)
				},
			},
			"tracks": &graphql.Field{
				Type: graphql.NewList(trackType),
				Resolve: func (g graphql.ResolveParams) (interface{}, error) {
					return interactor.Library.TrackRepository.FindAll()
				},
			},
			"track": &graphql.Field{
				Type: trackType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "Track ID",
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					i := p.Args["id"].(string)
					id, err := strconv.Atoi(i)
					if err != nil {
						return nil, err
					}

					return interactor.Library.TrackRepository.Find(id)
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
		Query: queryType,
	})
	if err != nil {
		panic(err)
	}

	return interactor
}

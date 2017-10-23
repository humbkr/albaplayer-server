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

// Defines static parts of artist type.
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

// Defines static parts of album type.
var albumType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Album",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Name: "Album ID",
			Description: "Album unique identifier.",
			Type: graphql.NewNonNull(graphql.ID),
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if album, ok := p.Source.(business.AlbumView); ok == true {
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
				if album, ok := p.Source.(business.AlbumView); ok == true {
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
				if album, ok := p.Source.(business.AlbumView); ok == true {
					return album.Year, nil
				}
				return nil, nil
			},
		},
		"artistName": &graphql.Field{
			Name: "Artist name",
			Description: "Shorthand property for performance, avoid loading an artist for each album.",
			Type: graphql.NewNonNull(graphql.String),
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if album, ok := p.Source.(business.AlbumView); ok == true {
					return album.ArtistName, nil
				}
				return nil, nil
			},
		},
		"tracks": &graphql.Field{
			Name: "Album tracks",
			Description: "Tracks of album.",
			Type: graphql.NewList(graphql.NewNonNull(trackType)),
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				if album, ok := p.Source.(business.AlbumView); ok == true {
					return album.Tracks, nil
				}
				return nil, nil
			},
		},
	},
})

// Defines static parts of track type.
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

var queueType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Queue",
	Fields: graphql.Fields{
		// TODO see if we can create an "options" field including repeat + random.
		"repeat": &graphql.Field{
			Name: "Queue option repeat",
			Description: "Queue playing option: repeat.",
			Type: graphql.String,
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				// TODO change this for multi users.
				queue := business.GetQueueInstance()
				value := "no"
				switch queue.Options.Repeat {
				case business.QueueOptionRepeatSingle:
					value = "all"
				case business.QueueOptionRepeatAll:
					value = "track"
				}

				return value, nil
			},
		},
		// TODO see if we can create an "options" field including repeat + random.
		"random": &graphql.Field{
			Name: "Queue option random",
			Description: "Queue playing option: random.",
			Type: graphql.Boolean,
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				// TODO change this for multi users.
				queue := business.GetQueueInstance()
				return queue.Options.Random, nil
			},
		},
		"current": &graphql.Field{
			Name: "Current track",
			Description: "The current track playing.",
			Type: graphql.NewNonNull(trackType),
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				// TODO change this for multi users.
				queue := business.GetQueueInstance()
				return queue.Current()
			},
		},
		"tracks": &graphql.Field{
			Name: "Tracklist",
			Description: "Tracks in the queue.",
			Type: graphql.NewList(graphql.NewNonNull(trackType)),
			Resolve: func (p graphql.ResolveParams) (interface{}, error) {
				// TODO change this for multi users.
				queue := business.GetQueueInstance()
				tracks := domain.Tracks{}
				for _, trackId := range queue.PlayingOrder {
					tracks = append(tracks, queue.Tracklist[trackId])
				}

				return tracks, nil
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

	// Define dynamic fields on types.
	albumType.AddFieldConfig("artist", &graphql.Field{
		Type: artistType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if album, ok := p.Source.(business.AlbumView); ok == true && album.ArtistId != 0 {
				return interactor.Library.ArtistRepository.Get(album.ArtistId)
			}

			return nil, nil
		},
	})

	trackType.AddFieldConfig("artist", &graphql.Field{
		Type: artistType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if track, ok := p.Source.(domain.Track); ok == true && track.ArtistId != 0 {
				return interactor.Library.ArtistRepository.Get(track.ArtistId)
			}

			return nil, nil
		},
	})
	trackType.AddFieldConfig("album", &graphql.Field{
		Type: artistType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if track, ok := p.Source.(domain.Track); ok == true && track.AlbumId != 0 {
				return interactor.Library.AlbumRepository.Get(track.AlbumId)
			}

			return nil, nil
		},
	})
	trackType.AddFieldConfig("cover", &graphql.Field{
		Name: "Track cover",
		Description: "Localisation of the cover file on the computer.",
		Type: graphql.String,
		Resolve: func (p graphql.ResolveParams) (interface{}, error) {
			if track, ok := p.Source.(domain.Track); ok == true && track.CoverId != 0 {
					cover, err := interactor.Library.CoverRepository.Get(track.CoverId)
					if err == nil {
						return cover.Path, nil
					}
			}

			return nil, nil
		},
	})

	// This is the type that will be the root of our query,
	// and the entry point into our schema.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"artists": &graphql.Field{
				Type: graphql.NewList(artistType),
				Resolve: func (g graphql.ResolveParams) (interface{}, error) {
					return interactor.Library.GetAllArtists(false)
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

					return interactor.Library.ArtistRepository.Get(id)
				},
			},
			"albums": &graphql.Field{
				Type: graphql.NewList(albumType),
				Args: graphql.FieldConfigArgument{
					"hydrate": &graphql.ArgumentConfig{
						Description: "Enable possibility to get tracks from albums list. Default to false.",
						Type: graphql.Boolean,
					},
				},
				Resolve: func (p graphql.ResolveParams) (interface{}, error) {
					if p.Args["hydrate"] != nil {
						if hydrate, ok := p.Args["hydrate"].(bool); ok {
							return interactor.Library.AlbumRepository.GetAll(hydrate)
						}
					}

					return interactor.Library.AlbumRepository.GetAll(false)
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

					return interactor.Library.AlbumRepository.Get(id)
				},
			},
			"tracks": &graphql.Field{
				Type: graphql.NewList(trackType),
				Resolve: func (g graphql.ResolveParams) (interface{}, error) {
					return interactor.Library.TrackRepository.GetAll()
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

					return interactor.Library.TrackRepository.Get(id)
				},
			},
			"queue": &graphql.Field{
				Type: queueType,
				Resolve: func (p graphql.ResolveParams) (interface{}, error) {
					// TODO change this for multi users.
					return business.GetQueueInstance(), nil
				},
			},
		},
	})

	// Queue operations.
	queueMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "QueueMutation",
		Fields: graphql.Fields{
			"next": &graphql.Field{
				Type: queueType,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					queue := business.GetQueueInstance()
					queue.Next()
					return queue, nil
				},
			},
			"previous": &graphql.Field{
				Type: queueType,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					queue := business.GetQueueInstance()
					queue.Previous()
					return queue, nil
				},
			},
			"appendArtist": &graphql.Field{
				Type: queueType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					queue := business.GetQueueInstance()
					artistId, ok := params.Args["id"].(int)
					if ok {
						queue.AppendArtist(artistId)
					}

					return queue, nil
				},
			},
			"appendAlbum": &graphql.Field{
				Type: queueType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					queue := business.GetQueueInstance()
					albumId, ok := params.Args["id"].(int)
					if ok {
						queue.AppendAlbum(albumId)
					}

					return queue, nil
				},
			},
			"appendTrack": &graphql.Field{
				Type: queueType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					queue := business.GetQueueInstance()
					trackId, ok := params.Args["id"].(int)
					if ok {
						queue.AppendTrack(trackId)
					}

					return queue, nil
				},
			},
			"playArtist": &graphql.Field{
				Type: queueType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					queue := business.GetQueueInstance()
					artistId, ok := params.Args["id"].(int)
					if ok {
						queue.PlayArtist(artistId)
					}

					return queue, nil
				},
			},
			"playAlbum": &graphql.Field{
				Type: queueType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					queue := business.GetQueueInstance()
					albumId, ok := params.Args["id"].(int)
					if ok {
						queue.PlayAlbum(albumId)
					}

					return queue, nil
				},
			},
			"playTrack": &graphql.Field{
				Type: queueType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					queue := business.GetQueueInstance()
					trackId, ok := params.Args["id"].(int)
					if ok {
						queue.PlayTrack(trackId)
					}

					return queue, nil
				},
			},
			"clear": &graphql.Field{
				Type: queueType,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					queue := business.GetQueueInstance()
					queue.Clear()
					return queue, nil
				},
			},
			"setOptions": &graphql.Field{
				Type: queueType,
				Args: graphql.FieldConfigArgument{
					"repeat": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"random": &graphql.ArgumentConfig{
						Type: graphql.Boolean,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					queue := business.GetQueueInstance()
					repeat, ok := params.Args["repeat"].(string)
					if ok {
						switch repeat {
						case "track":
							queue.Options.Repeat = business.QueueOptionRepeatSingle
						case "all":
							queue.Options.Repeat = business.QueueOptionRepeatAll
						}
					}
					random, ok := params.Args["random"].(bool)
					if ok {
						queue.Options.Random = random
					}

					return queue, nil
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
		Mutation: queueMutation,
	})
	if err != nil {
		panic(err)
	}

	return interactor
}

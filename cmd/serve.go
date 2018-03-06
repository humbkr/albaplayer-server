package cmd

import (
	"fmt"

	"git.humbkr.com/jgalletta/alba-player/internal/alba"
	"git.humbkr.com/jgalletta/alba-player/internal/alba/interfaces"
	"github.com/mnmtanish/go-graphiql"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	gqlHandler "github.com/graphql-go/handler"
	"net/http"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve app on the specified port",
	Long:  `Launch all services, create all endpoints, and serve UI web app.`,
	Run: func(cmd *cobra.Command, args []string) {
		libraryInteractor := alba.InitApp()

		// Initialize GraphQL stuff.
		graphQLInteractor := interfaces.NewGraphQLInteractor(&libraryInteractor)

		// Create a graphl-go HTTP handler with our previously defined schema
		// and set it to return pretty JSON output.
		graphQLHandler := gqlHandler.New(&gqlHandler.Config{
			Schema: &graphQLInteractor.Schema,
			Pretty: true,
		})

		// Serve a GraphQL endpoint at `/graphql`.
		// Make the server handle cross-domain requests.
		http.Handle("/graphql", cors.Default().Handler(graphQLHandler))

		// Serve media files streaming endpoint.
		// Makes the server handle cross-domain requests.
		mediaFilesHandler := interfaces.NewMediaStreamHandler(&libraryInteractor)
		http.Handle("/stream/", http.StripPrefix("/stream/", cors.Default().Handler(mediaFilesHandler)))

		// Serve media files streaming endpoint.
		// Makes the server handle cross-domain requests.
		coverFilesHandler := interfaces.NewCoverStreamHandler(&libraryInteractor)
		http.Handle("/covers/", http.StripPrefix("/covers/", cors.Default().Handler(coverFilesHandler)))

		// Serve frontend app.
		http.Handle("/static/", http.StripPrefix("/static/", cors.Default().Handler(http.FileServer(http.Dir("web/static")))))
		http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/favicon.ico")
		})
		http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/robots.txt")
		})
		http.HandleFunc("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/manifest.json")
		})
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/index.html")
		})


		if viper.GetBool("DevMode.Enabled") {
			// Serve graphiql.
			http.HandleFunc("/graphiql", graphiql.ServeGraphiQL)
		}

		// Launch the server.
		fmt.Printf("Server is up: http://localhost:%s/graphql\n", viper.GetString("Server.Port"))
		http.ListenAndServe(":" + viper.GetString("Server.Port"), nil)
	},
}

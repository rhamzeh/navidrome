// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/cloudsonic/sonic-server/api"
	"github.com/cloudsonic/sonic-server/domain"
	"github.com/cloudsonic/sonic-server/engine"
	"github.com/cloudsonic/sonic-server/itunesbridge"
	"github.com/cloudsonic/sonic-server/persistence"
	"github.com/cloudsonic/sonic-server/persistence/db_sql"
	"github.com/cloudsonic/sonic-server/scanner"
	"github.com/google/wire"
)

// Injectors from wire_injectors.go:

func CreateApp(musicFolder string) *App {
	provider := createPersistenceProvider()
	checkSumRepository := provider.CheckSumRepository
	itunesScanner := scanner.NewItunesScanner(checkSumRepository)
	mediaFileRepository := provider.MediaFileRepository
	albumRepository := provider.AlbumRepository
	artistRepository := provider.ArtistRepository
	artistIndexRepository := provider.ArtistIndexRepository
	playlistRepository := provider.PlaylistRepository
	propertyRepository := provider.PropertyRepository
	importer := scanner.NewImporter(musicFolder, itunesScanner, mediaFileRepository, albumRepository, artistRepository, artistIndexRepository, playlistRepository, propertyRepository)
	app := NewApp(importer)
	return app
}

func CreateSubsonicAPIRouter() *api.Router {
	provider := createPersistenceProvider()
	propertyRepository := provider.PropertyRepository
	mediaFolderRepository := provider.MediaFolderRepository
	artistIndexRepository := provider.ArtistIndexRepository
	artistRepository := provider.ArtistRepository
	albumRepository := provider.AlbumRepository
	mediaFileRepository := provider.MediaFileRepository
	browser := engine.NewBrowser(propertyRepository, mediaFolderRepository, artistIndexRepository, artistRepository, albumRepository, mediaFileRepository)
	cover := engine.NewCover(mediaFileRepository, albumRepository)
	nowPlayingRepository := provider.NowPlayingRepository
	listGenerator := engine.NewListGenerator(albumRepository, mediaFileRepository, nowPlayingRepository)
	itunesControl := itunesbridge.NewItunesControl()
	playlistRepository := provider.PlaylistRepository
	playlists := engine.NewPlaylists(itunesControl, playlistRepository, mediaFileRepository)
	ratings := engine.NewRatings(itunesControl, mediaFileRepository, albumRepository, artistRepository)
	scrobbler := engine.NewScrobbler(itunesControl, mediaFileRepository, nowPlayingRepository)
	search := engine.NewSearch(artistRepository, albumRepository, mediaFileRepository)
	router := api.NewRouter(browser, cover, listGenerator, playlists, ratings, scrobbler, search)
	return router
}

func createPersistenceProvider() *Provider {
	albumRepository := db_sql.NewAlbumRepository()
	artistRepository := db_sql.NewArtistRepository()
	checkSumRepository := db_sql.NewCheckSumRepository()
	artistIndexRepository := db_sql.NewArtistIndexRepository()
	mediaFileRepository := db_sql.NewMediaFileRepository()
	mediaFolderRepository := persistence.NewMediaFolderRepository()
	nowPlayingRepository := persistence.NewNowPlayingRepository()
	playlistRepository := db_sql.NewPlaylistRepository()
	propertyRepository := db_sql.NewPropertyRepository()
	provider := &Provider{
		AlbumRepository:       albumRepository,
		ArtistRepository:      artistRepository,
		CheckSumRepository:    checkSumRepository,
		ArtistIndexRepository: artistIndexRepository,
		MediaFileRepository:   mediaFileRepository,
		MediaFolderRepository: mediaFolderRepository,
		NowPlayingRepository:  nowPlayingRepository,
		PlaylistRepository:    playlistRepository,
		PropertyRepository:    propertyRepository,
	}
	return provider
}

// wire_injectors.go:

type Provider struct {
	AlbumRepository       domain.AlbumRepository
	ArtistRepository      domain.ArtistRepository
	CheckSumRepository    scanner.CheckSumRepository
	ArtistIndexRepository domain.ArtistIndexRepository
	MediaFileRepository   domain.MediaFileRepository
	MediaFolderRepository domain.MediaFolderRepository
	NowPlayingRepository  domain.NowPlayingRepository
	PlaylistRepository    domain.PlaylistRepository
	PropertyRepository    domain.PropertyRepository
}

var allProviders = wire.NewSet(itunesbridge.NewItunesControl, engine.Set, scanner.Set, api.NewRouter, wire.FieldsOf(new(*Provider), "AlbumRepository", "ArtistRepository", "CheckSumRepository",
	"ArtistIndexRepository", "MediaFileRepository", "MediaFolderRepository", "NowPlayingRepository",
	"PlaylistRepository", "PropertyRepository"), createPersistenceProvider,
)

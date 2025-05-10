//go:build !ci
// +build !ci

package qbittorrent_test

import (
	"encoding/base64"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/autobrr/go-qbittorrent"
)

const (
	// No magic here. It's a torrent that I have crafted with qBittorrent.
	// Replace it with any valid torrent should make no difference.
	sampleTorrent   = "ZDEwOmNyZWF0ZWQgYnkxODpxQml0dG9ycmVudCB2NS4xLjAxMzpjcmVhdGlvbiBkYXRlaTE3NDY4NjI1MDFlNDppbmZvZDY6bGVuZ3RoaTIxZTQ6bmFtZTEyOnVudGl0bGVkLnR4dDEyOnBpZWNlIGxlbmd0aGkxNjM4NGU2OnBpZWNlczIwOrV8kDHOo9sgQCTOvdOwDtO6wMy9Nzpwcml2YXRlaTFlZWU="
	sampleMagnetURL = "magnet:?xt=urn:btih:8f25668abe58fff8a1e11387e8d6475dc6346669&dn=untitled.txt&xl=21"
	sampleInfoHash  = "8f25668abe58fff8a1e11387e8d6475dc6346669"
)

var (
	qBittorrentBaseURL  string
	qBittorrentUsername string
	qBittorrentPassword string
)

func init() {
	qBittorrentBaseURL = "http://127.0.0.1:8080/"
	if val := os.Getenv("QBIT_BASE_URL"); val != "" {
		qBittorrentBaseURL = val
	}
	qBittorrentUsername = "admin"
	if val := os.Getenv("QBIT_USERNAME"); val != "" {
		qBittorrentUsername = val
	}
	qBittorrentPassword = "admin"
	if val := os.Getenv("QBIT_PASSWORD"); val != "" {
		qBittorrentPassword = val
	}
}

func TestClient_Login(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	err := client.Login()
	assert.NoError(t, err)
}

func TestClient_Login_Anonymous(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{Host: qBittorrentBaseURL})

	err := client.Login()
	assert.NoError(t, err)
}

func TestClient_Login_BadCredentials(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword + "invalid",
	})

	err := client.Login()
	assert.ErrorContains(t, err, "bad credentials")
}

func TestClient_Shutdown(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	// comment following lines to run this test manually
	t.Skip("manually run only")
	return

	err := client.Shutdown()
	assert.NoError(t, err)
}

func TestClient_GetTorrents(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	_, err := client.GetTorrents(qbittorrent.TorrentFilterOptions{})
	assert.NoError(t, err)
}

func TestClient_GetTorrents_Whitelist(t *testing.T) {
	authorizedClient := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	// decide if we need to run this test
	preferences, err := authorizedClient.GetAppPreferences()
	assert.NoError(t, err)
	if !preferences.BypassAuthSubnetWhitelistEnabled {
		t.Skip("IP whitelist not enabled")
	}

	unauthorizedClient := qbittorrent.NewClient(qbittorrent.Config{Host: qBittorrentBaseURL})
	_, err = unauthorizedClient.GetTorrents(qbittorrent.TorrentFilterOptions{})
	assert.NoError(t, err)
}

func TestClient_GetTorrentsActiveDownloads(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	_, err := client.GetTorrentsActiveDownloads()
	assert.NoError(t, err)
}

func TestClient_GetTorrentProperties(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	data, err := client.GetTorrents(qbittorrent.TorrentFilterOptions{})
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	hash := data[0].Hash
	properties, err := client.GetTorrentProperties(hash)
	assert.NoError(t, err)
	assert.Equal(t, hash, properties.Hash)
}

func TestClient_GetTorrentsRaw(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	data, err := client.GetTorrentsRaw()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestClient_GetTorrentTrackers(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	data, err := client.GetTorrents(qbittorrent.TorrentFilterOptions{})
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	hash := data[0].Hash
	trackers, err := client.GetTorrentTrackers(hash)
	assert.NoError(t, err)
	assert.NotEmpty(t, trackers)
}

func TestClient_AddTorrentFromMemory(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	torrent, err := base64.StdEncoding.DecodeString(sampleTorrent)
	assert.NoError(t, err)

	err = client.AddTorrentFromMemory(torrent, nil)
	assert.NoError(t, err)
}

func TestClient_AddTorrentFromFile(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	// prepare torrent file
	file, err := os.CreateTemp("", "")
	assert.NoError(t, err)
	defer func(file *os.File) {
		_ = file.Close()
		_ = os.Remove(file.Name())
	}(file)

	torrent, err := base64.StdEncoding.DecodeString(sampleTorrent)
	assert.NoError(t, err)
	_, err = file.Write(torrent)
	assert.NoError(t, err)
	err = file.Close()
	assert.NoError(t, err)

	err = client.AddTorrentFromFile(file.Name(), nil)
	assert.NoError(t, err)
}

func TestClient_AddTorrentFromURL(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	// Theoretically, magnet links and http links are handled in the same way.
	// Replace magnet links with HTTP links should make no difference.
	err := client.AddTorrentFromUrl(sampleMagnetURL, nil)
	assert.NoError(t, err)
}

func TestClient_DeleteTorrents(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	err := client.DeleteTorrents([]string{sampleInfoHash}, false)
	assert.NoError(t, err)
}

func TestClient_GetDefaultSavePath(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	_, err := client.GetDefaultSavePath()
	assert.NoError(t, err)
}

func TestClient_GetAppCookies(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	_, err := client.GetAppCookies()
	assert.NoError(t, err)
}

func TestClient_SetAppCookies(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	var err error
	var cookies = []qbittorrent.Cookie{
		{
			Name:           "test",
			Domain:         "example.com",
			Path:           "/",
			Value:          "test",
			ExpirationDate: time.Now().Add(time.Hour).Unix(),
		},
	}
	err = client.SetAppCookies(cookies)
	assert.NoError(t, err)

	resp, err := client.GetAppCookies()
	assert.NoError(t, err)
	assert.NotEmpty(t, cookies)
	assert.Equal(t, cookies, resp)
}

func TestClient_BanPeers(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	err := client.BanPeers([]string{"127.0.0.1:80"})
	assert.NoError(t, err)
}

func TestClient_GetBuildInfo(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	bi, err := client.GetBuildInfo()
	assert.NoError(t, err)
	assert.NotEmpty(t, bi.Qt)
	assert.NotEmpty(t, bi.Libtorrent)
	assert.NotEmpty(t, bi.Boost)
	assert.NotEmpty(t, bi.Openssl)
	assert.NotEmpty(t, bi.Bitness)
}

func TestClient_GetTorrentDownloadLimit(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	data, err := client.GetTorrents(qbittorrent.TorrentFilterOptions{})
	assert.NoError(t, err)
	var hashes []string
	for _, torrent := range data {
		hashes = append(hashes, torrent.Hash)
	}

	limits, err := client.GetTorrentDownloadLimit(hashes)
	assert.NoError(t, err)
	assert.Equal(t, len(hashes), len(limits))

	// FIXME: The following assertion will fail.
	// Neither "hashes=all" nor "all" is working.
	// I have no idea. Maybe the document is lying?
	//
	// limits, err = client.GetTorrentDownloadLimit([]string{"all"})
	// assert.NoError(t, err)
	// assert.Equal(t, len(hashes), len(limits))
}

func TestClient_GetTorrentUploadLimit(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	data, err := client.GetTorrents(qbittorrent.TorrentFilterOptions{})
	assert.NoError(t, err)
	var hashes []string
	for _, torrent := range data {
		hashes = append(hashes, torrent.Hash)
	}

	limits, err := client.GetTorrentUploadLimit(hashes)
	assert.NoError(t, err)
	assert.Equal(t, len(hashes), len(limits))

	// FIXME: The following assertion will fail.
	// Neither "hashes=all" nor "all" is working.
	// I have no idea. Maybe the document is lying?
	// Just as same as Client.GetTorrentDownloadLimit.
	//
	// limits, err = client.GetTorrentDownloadLimit([]string{"all"})
	// assert.NoError(t, err)
	// assert.Equal(t, len(hashes), len(limits))
}

func TestClient_ToggleTorrentSequentialDownload(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	var err error

	data, err := client.GetTorrents(qbittorrent.TorrentFilterOptions{})
	assert.NoError(t, err)
	var hashes []string
	for _, torrent := range data {
		hashes = append(hashes, torrent.Hash)
	}

	err = client.ToggleTorrentSequentialDownload(hashes)
	assert.NoError(t, err)

	// No idea why this is working but downloadLimit/uploadLimit are not.
	err = client.ToggleTorrentSequentialDownload([]string{"all"})
	assert.NoError(t, err)
}

func TestClient_SetTorrentSuperSeeding(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	var err error

	data, err := client.GetTorrents(qbittorrent.TorrentFilterOptions{})
	assert.NoError(t, err)
	var hashes []string
	for _, torrent := range data {
		hashes = append(hashes, torrent.Hash)
	}

	err = client.SetTorrentSuperSeeding(hashes, true)
	assert.NoError(t, err)

	// FIXME: following test not fail but has no effect.
	// qBittorrent doesn't return any error but super seeding status is not changed.
	// I tried specify hashes as "all" but it's not working too.
	err = client.SetTorrentSuperSeeding([]string{"all"}, false)
	assert.NoError(t, err)
}

func TestClient_GetTorrentPieceStates(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	data, err := client.GetTorrents(qbittorrent.TorrentFilterOptions{})
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	hash := data[0].Hash
	states, err := client.GetTorrentPieceStates(hash)
	assert.NoError(t, err)
	assert.NotEmpty(t, states)
}

func TestClient_GetTorrentPieceHashes(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	data, err := client.GetTorrents(qbittorrent.TorrentFilterOptions{})
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	hash := data[0].Hash
	states, err := client.GetTorrentPieceHashes(hash)
	assert.NoError(t, err)
	assert.NotEmpty(t, states)
}

func TestClient_AddPeersForTorrents(t *testing.T) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     qBittorrentBaseURL,
		Username: qBittorrentUsername,
		Password: qBittorrentPassword,
	})

	data, err := client.GetTorrents(qbittorrent.TorrentFilterOptions{})
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	hashes := []string{data[0].Hash}
	peers := []string{"127.0.0.1:12345"}
	err = client.AddPeersForTorrents(hashes, peers)
	// It seems qBittorrent doesn't actually check whether given peers are available.
	assert.NoError(t, err)
}

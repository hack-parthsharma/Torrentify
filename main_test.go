package main

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestApp_Run(t *testing.T) {
	t.Run("creates torrents", func(t *testing.T) {
		var parsedTorrent torrent
		buf := &bytes.Buffer{}

		app := &App{
			torrentCreator: func(t *torrent, w io.Writer) error {
				parsedTorrent = *t
				_, err := w.Write([]byte("torrent"))
				return err
			},
			fileCreator: func(path string) (io.Writer, error) {
				return buf, nil
			},
		}

		_ = os.Setenv("CREATED_BY", "tester")
		args := []string{"", "-announce", "http://tracker1.com,http://tracker2.com", "-output", "test.torrent", "-name", "test torrent", "/path/to/root"}
		err := app.Run(args)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected := torrent{
			AnnounceUrls: []string{"http://tracker1.com", "http://tracker2.com"},
			Comment:      "",
			CreatedBy:    "tester",
			Name:         "test torrent",
			Private:      false,
			Root:         "/path/to/root",
			PieceLength:  1048576,
		}
		if !reflect.DeepEqual(parsedTorrent, expected) {
			t.Errorf("expected torrent: %+v, got: %+v", expected, parsedTorrent)
		}

		if buf.String() != "torrent" {
			t.Errorf("expected torrent data %v, got: %v", "torrent", buf.String())
		}
	})
}

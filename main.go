package main

import (
	_ "embed"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	app := &App{
		torrentCreator: makeTorrent,
		fileCreator:    createFile,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

//go:embed version.txt
var version string

type torrent struct {
	AnnounceUrls []string
	Comment      string
	CreatedBy    string
	Name         string
	Private      bool
	Root         string
	PieceLength  uint64
}

type TorrentCreatorFunc func(t *torrent, w io.Writer) error
type FileCreatorFunc func(path string) (io.Writer, error)

type App struct {
	torrentCreator TorrentCreatorFunc
	fileCreator    FileCreatorFunc
}

func (a *App) Run(args []string) error {
	cliApp := &cli.App{
		Name:        "torrentify",
		Usage:       "torrent creator",
		ArgsUsage:   "<torrent root>",
		Version:     version,
		Description: "torrentify creates torrent files from given root directory.",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "announce",
				Aliases:  []string{"a"},
				Usage:    "tracker announce URLs, separated with commas.",
				EnvVars:  []string{"ANNOUNCE_URL"},
				Required: true,
			},
			&cli.PathFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "output path",
				DefaultText: "stdout",
				Required:    true,
				TakesFile:   true,
				Value:       "-",
			},
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "torrent name",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "comment",
				Aliases: []string{"c"},
				Usage:   "torrent comment",
			},
			&cli.StringFlag{
				Name:    "createdby",
				Usage:   "torrent creator name",
				EnvVars: []string{"CREATED_BY"},
			},
			&cli.BoolFlag{
				Name:    "private",
				Usage:   "set torrent as private (useful for private trackers)",
				EnvVars: []string{"PRIVATE"},
			},
			&cli.Uint64Flag{
				Name:    "piecelength",
				Usage:   "torrent piece length in bytes.",
				EnvVars: []string{"PIECE_LENGTH"},
				Value:   1024 * 1024,
			},
		},
		HideHelpCommand: true,
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 1 {
				log.Printf("root directory is not specified")
				return cli.ShowAppHelp(ctx)
			}

			rootDir := ctx.Args().Get(0)

			t := &torrent{
				AnnounceUrls: ctx.StringSlice("announce"),
				Comment:      ctx.String("comment"),
				CreatedBy:    ctx.String("createdby"),
				Name:         ctx.String("name"),
				Private:      ctx.Bool("private"),
				Root:         rootDir,
				PieceLength:  ctx.Uint64("piecelength"),
			}

			outputPath := ctx.Path("output")
			w, err := a.fileCreator(outputPath)
			if err != nil {
				return err
			}

			return a.torrentCreator(t, w)
		},
	}
	return cliApp.Run(args)
}

func createFile(path string) (io.Writer, error) {
	if path == "-" {
		return os.Stdout, nil
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, errors.Wrap(err, "create file")
	}
	return f, nil
}

func makeTorrent(t *torrent, w io.Writer) error {
	mi := metainfo.MetaInfo{
		AnnounceList: make([][]string, 0),
	}
	for _, a := range t.AnnounceUrls {
		mi.AnnounceList = append(mi.AnnounceList, []string{a})
	}

	mi.CreationDate = time.Now().Unix()
	if len(t.Comment) > 0 {
		mi.Comment = t.Comment
	}

	if len(t.CreatedBy) > 0 {
		mi.CreatedBy = t.CreatedBy
	} else {
		mi.CreatedBy = "torrentify"
	}

	info := metainfo.Info{
		PieceLength: int64(t.PieceLength),
		Private:     &t.Private,
	}
	err := info.BuildFromFilePath(t.Root)
	if err != nil {
		return errors.Wrap(err, "hash files")
	}

	if t.Name != "" {
		info.Name = t.Name
	}
	mi.InfoBytes, err = bencode.Marshal(info)
	if err != nil {
		return errors.Wrap(err, "marshall torrent")
	}

	return mi.Write(w)
}

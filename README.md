# torrentify

A mini CLI application to create torrent files from a directory.

## Usage

**torrentify** requires the root directory of the files to be specified as the only argument.
The remaining parameters are specified as flags,
some of which can be set using environment variables.

```shell
> ./torrentify --help

NAME:
   torrentify - torrent creator

USAGE:
   torrentify [global options] <torrent root>

VERSION:
   v0.0.1-fa5d041

DESCRIPTION:
   torrentify creates torrent files from given root directory.

GLOBAL OPTIONS:
   --announce value, -a value  tracker announce URLs, separated with commas.  (accepts multiple inputs) [$ANNOUNCE_URL]
   --comment value, -c value   torrent comment
   --createdby value           torrent creator name [$CREATED_BY]
   --help, -h                  show help (default: false)
   --name value, -n value      torrent name
   --output value, -o value    output path, defaults to stdout (default: "-")
   --piecelength value         torrent piece length in bytes. (default: 1048576) [$PIECE_LENGTH]
   --private                   set torrent as private (useful for private trackers) (default: false) [$PRIVATE]
   --version, -v               print the version (default: false)
```

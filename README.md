# find-dupe-music

A simple program to scan a supplied set of folders for music files, and detect duplicates based on Album/Artist metadate.

## Building

`find-dupe-music` is a fairly simple go based tool.

```
$ go get
$ go install
```


## Running

```bash`
./find-dupe-music scan /folder/music/*
```

will run `find-dupe-music` with each subdirectory of `/folder/music/` - you could specify each folder
separately as well. For each directory, `find-dupe-music` starts a thread/go-routine scanning it's tree
for all `alac`, `flac`,`mp3`, `m4p`, and `m4a` files, and for every path that has not been visited,
checks the audiofiles metadata to group by AlbumArtist/Artist:Album and simply reporting duplicates into
the file `dupes.txt`:

## Configuration

### Path Configuration

For the convienence of reuse, `find-dupe-music` can be configured with a predefined set of paths to scan,
and any command line specified paths will be appended to the configued list.

By default, the `~/.find-dupes.yaml` or `$HOME/.find-dupes.yaml` file is checked, but can be specified with `--config`.

```yaml
path:
- /folder/musica
- /folder/musicb
```

### Album Skipping

In some cases, albums with different edtions will be reported as false positive duplicates, this
is because they simply report the same `artist:album` pair, these albums can be skipped from being
reported by defining a `skip` list in the config file:

```yaml
skip:
- ArtistA:AlbumA
- ArtistB:AlbumA
```

If you wish to temporally skip, the album skipping - simply use the `--ignore-skips` or `-s` flag to the `scan` command:

```bash
./find-dupe-music scan --ignore-skips
```

Duplicate Music Report

Found duplicates for House of Mythology:Watch & Pray - Five Years of Studious Decrepitude
   - /Volumes/Media Content/Plex Content/Music/Plex/House of Mythology/2020 - House of Mythology - Watch & Pray - Five Years of Studious Decrepitude
   - /Volumes/Media Content/Plex Content/Music/Plex/House of Mythology/2011 - House of Mythology - Watch & Pray - Five Years of Studious Decrepitude
   - /Volumes/Media Content/Plex Content/Music/Plex/House of Mythology/2018 - House of Mythology - Watch & Pray - Five Years of Studious Decrepitude
Found duplicates for Dead By Monday:Dead by Monday
   - /Volumes/Media Content/Plex Content/Music/Plex/Dead By Monday/2020 - Dead By Monday - Dead by Monday
   - /Volumes/Media Content/Plex Content/Music/iTunes Music/Music/Dead By Monday/Dead by Monday
```

The report is not _exactly_ reporting duplicate files, more duplicate copies of albums - as defined as disparate
paths. This includes traping 'singles' and their albums tracks sharing metadata, or multi-year compilations being separated into paths.

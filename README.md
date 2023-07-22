# find-dupe-music

A simple program to scan a supplied set of folders for music files, and detect duplicates based on Album/Artist metadate.

## Building

`find-dupe-music` is a fairly simple go based tool.

```
$ go get
$ go install
```


## Running

```
$ find-dupe-music /folder/music/*
```

will run `find-dupe-music` with each subdirectory of `/folder/music/` - you could specify each folder
separately as well. For each directory, `find-dupe-music` starts a thread/go-routine scanning it's tree
for all `alac`, `flac`,`mp3`, `m4p`, and `m4a` files, and for every path that has not been visited,
checks the audiofiles metadata to group by AlbumArtist/Artist:Album and simply reporting duplicates into
the file `dupes.txt`:

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

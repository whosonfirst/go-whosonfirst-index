package index

import (
	"bufio"
	"errors"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"io"
	_ "log"
	"os"
	"path/filepath"
	"strings"
)

type IndexerFunc func(path string, info os.FileInfo, args ...interface{}) error

type Indexer struct {
	Mode string
	Func IndexerFunc
}

func NewIndexer(mode string, f IndexerFunc) (*Indexer, error) {

	i := Indexer{
		Mode: mode,
		Func: f,
	}

	return &i, nil
}

func (i *Indexer) IndexPaths(paths []string, args ...interface{}) error {

	for _, path := range paths {

		err := i.IndexPath(path, args...)

		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Indexer) IndexPath(path string, args ...interface{}) error {

	var abs_path string
	var info os.FileInfo

	if i.Mode != "meta" {

		_path, err := filepath.Abs(path)

		if err != nil {
			return err
		}

		_info, err := os.Stat(_path)

		if err != nil {
			return err
		}

		abs_path = _path
		info = _info
	}

	if i.Mode == "directory" {

		return i.IndexDirectory(abs_path, args...)

	} else if i.Mode == "repo" {

		data := filepath.Join(abs_path, "data")

		_, err := os.Stat(data)

		if err != nil {
			return err
		}

		return i.IndexDirectory(abs_path, args...)

	} else if i.Mode == "filelist" {

		return i.IndexFileList(abs_path, args...)

	} else if i.Mode == "meta" {

		parts := strings.Split(path, ":")

		if len(parts) != 2 {
			return errors.New("Invalid path declaration for a meta file")
		}

		for _, p := range parts {

			_, err := os.Stat(p)

			if os.IsNotExist(err) {
				return errors.New("Path does not exist")
			}
		}

		meta_file := parts[0]
		data_root := parts[1]

		return i.IndexMetaFile(meta_file, data_root, args...)
		return nil

	} else if i.Mode == "files" {

		return i.Func(abs_path, info, args...)

	} else {

		return errors.New("Invalid indexer")
	}

}

func (i *Indexer) IndexDirectory(path string, args ...interface{}) error {

	cb := func(path string, info os.FileInfo) error {
		return i.Func(path, info, args...)
	}

	c := crawl.NewCrawler(path)
	return c.Crawl(cb)
}

func (i *Indexer) IndexMetaFile(meta_file string, data_root string, args ...interface{}) error {

	reader, err := csv.NewDictReaderFromPath(meta_file)

	if err != nil {
		return err
	}

	for {
		row, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		rel_path, ok := row["path"]

		if !ok {
			return errors.New("Missing path key")
		}

		// TO DO: make this work with a row["repo"] key
		// (20170809/thisisaaronland)

		file_path := filepath.Join(data_root, rel_path)
		file_info, err := os.Stat(file_path)

		if err != nil {
			return err
		}

		err = i.Func(file_path, file_info, args...)

		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Indexer) IndexFileList(path string, args ...interface{}) error {

	fh, err := os.Open(path)

	if err != nil {
		return nil
	}

	defer fh.Close()

	scanner := bufio.NewScanner(fh)

	for scanner.Scan() {

		file_path := scanner.Text()

		file_info, err := os.Stat(file_path)

		if err != nil {
			return err
		}

		err = i.Func(file_path, file_info, args...)

		if err != nil {
			return err
		}
	}

	return nil

}

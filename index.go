package index

import (
	"bufio"
	"errors"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-timer"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

type IndexerFunc func(path string, info os.FileInfo, args ...interface{}) error

type Indexer struct {
	Mode     string
	Func     IndexerFunc
	Logger   *log.WOFLogger
	Indexed  int64
	Indexing chan bool
	count    int64
}

func NewIndexer(mode string, f IndexerFunc) (*Indexer, error) {

	logger := log.SimpleWOFLogger("index")
	indexing := make(chan bool)

	i := Indexer{
		Mode:     mode,
		Func:     f,
		Logger:   logger,
		Indexing: indexing,
		Indexed:  0,
		count:    0,
	}

	return &i, nil
}

func (i *Indexer) NewTimer(mode string, path string) (*timer.Timer, error) {

	cb := func(t timer.Timing) {
		i.Logger.Status("%s %s %v", mode, path, t.Duration())
	}

	tm, err := timer.NewDefaultTimer()

	if err != nil {
		return nil, err
	}

	tm.Callback = cb
	return tm, nil
}

func (i *Indexer) IndexPaths(paths []string, args ...interface{}) error {

	tm, err := i.NewTimer("paths", "...")

	if err != nil {
		return err
	}

	defer tm.Stop()

	i.increment()
	defer i.decrement()

	for _, path := range paths {

		err := i.IndexPath(path, args...)

		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Indexer) IndexPath(path string, args ...interface{}) error {

	i.increment()
	defer i.decrement()

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

		return i.IndexDirectory(data, args...)

	} else if i.Mode == "filelist" {

		return i.IndexFileList(abs_path, args...)

	} else if i.Mode == "meta" {

		parts := strings.Split(path, ":")

		if len(parts) == 1 {

			meta_root := filepath.Dir(parts[0])
			repo_root := filepath.Dir(meta_root)
			data_root := filepath.Join(repo_root, "data")

			parts = append(parts, data_root)
		}

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

		return i.process(abs_path, info, args...)

	} else {

		return errors.New("Invalid indexer")
	}

}

func (i *Indexer) IndexDirectory(path string, args ...interface{}) error {

	tm, err := i.NewTimer("directory", path)

	if err != nil {
		return err
	}

	defer tm.Stop()

	i.increment()
	defer i.decrement()

	cb := func(path string, info os.FileInfo) error {

		return i.process(path, info, args...)
	}

	c := crawl.NewCrawler(path)
	return c.Crawl(cb)
}

func (i *Indexer) IndexMetaFile(path string, data_root string, args ...interface{}) error {

	tm, err := i.NewTimer("metafile", path)

	if err != nil {
		return err
	}

	defer tm.Stop()

	i.increment()
	defer i.decrement()

	reader, err := csv.NewDictReaderFromPath(path)

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

		err = i.process(file_path, file_info, args...)

		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Indexer) IndexFileList(path string, args ...interface{}) error {

	tm, err := i.NewTimer("filelist", path)

	if err != nil {
		return err
	}

	defer tm.Stop()

	i.increment()
	defer i.decrement()

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

		i.Logger.Debug("process %s", file_path)

		err = i.process(file_path, file_info, args...)

		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Indexer) process(abs_path string, info os.FileInfo, args ...interface{}) error {

	i.increment()
	defer i.decrement()

	i.Logger.Debug("process %s", abs_path)
	err := i.Func(abs_path, info, args...)

	if err != nil {
		return err
	}

	atomic.AddInt64(&i.Indexed, 1)
	return nil
}

func (i *Indexer) increment() {

	go func(){
	   i.Indexing <- true
	}()
}

func (i *Indexer) decrement() {

	count := atomic.AddInt64(&i.count, -1)

	if count <= 0 {

		go func(){
			i.Indexing <- false
		}()
	}
}

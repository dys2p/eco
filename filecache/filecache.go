// Package filecache caches remote http files to local disk.
package filecache

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var ErrStatus = errors.New("status is not ok")

var client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

// A Cache caches upstream HTTP files to the filesystem.
//
// It ignores the HTTP Cache-Control header. This makes it suitable for use cases where upstream cache invalidation is done by changing the upstream URL but localpath does not change.
type Cache struct {
	adding       sync.Map
	maxAge       time.Duration
	modCache     map[string]int64 // uri => unix timestamp
	modCacheLock sync.Mutex
}

func NewCache(maxAge time.Duration) *Cache {
	return &Cache{
		maxAge:   maxAge,
		modCache: make(map[string]int64),
	}
}

func (cache *Cache) Add(uri, localpath string) error {
	threshold := time.Now().Add(-cache.maxAge).Unix()

	// check in-memory modCache first (before accessing the disk)
	cache.modCacheLock.Lock()
	mod, ok := cache.modCache[uri]
	cache.modCacheLock.Unlock()
	if ok && mod > threshold {
		return nil
	}

	// mod time is not in modCache, so get it from disk
	stat, err := os.Stat(localpath)
	if err == nil {
		mod := stat.ModTime().Unix()
		if mod > threshold {
			cache.modCacheLock.Lock()
			cache.modCache[uri] = mod
			cache.modCacheLock.Unlock()
			return nil
		}
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	// sync by uri
	_, ok = cache.adding.LoadOrStore(uri, struct{}{})
	if ok {
		// check back every second for a while
		for i := 0; i < 15; i++ { // a bit longer than http client timeout; TODO go1.22 range over int
			time.Sleep(time.Second)
			if _, ok := cache.adding.Load(uri); !ok {
				return nil
			}
		}
		return nil
	}
	defer cache.adding.Delete(uri)

	// set mod time now, preventing load on upstream in case of errors
	cache.modCacheLock.Lock()
	cache.modCache[uri] = time.Now().Unix()
	cache.modCacheLock.Unlock()

	// if file exists on disk, get Last-Modified header and compare
	if stat != nil {
		resp, err := client.Head(uri)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return ErrStatus
		}
		lastModified, err := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", resp.Header.Get("Last-Modified"))
		if err == nil && lastModified.Before(stat.ModTime()) {
			// store time of last upstream check (= now) in filesystem (modCache has already been updated above)
			err := os.Chtimes(localpath, time.Time{}, time.Now())
			if err == nil {
				return nil
			} // else maybe file on disk has been deleted in the meantime, so let's continue
		}
	}

	// download file
	err = os.MkdirAll(filepath.Dir(localpath), 0755)
	if err != nil {
		return err
	}
	resp, err := client.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	outfile, err := os.Create(localpath) // should set modification time to "now" (modCache has already been updated above)
	if err != nil {
		return err
	}
	defer outfile.Close()
	_, err = io.Copy(outfile, resp.Body)
	return err
}

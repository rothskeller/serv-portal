package ui

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"sunnyvaleserv.org/portal/util/request"
)

type assetInfo struct {
	name       string
	ctype      string
	data       []byte
	hash       uint32
	compressed bool
}

var assets []*assetInfo

var hashRE = regexp.MustCompile(`--[0-9a-f]{8}\.`)

// RegisterAsset is called by pages at init time to register their assets.
func RegisterAsset(name, ctype string, data []byte, hash uint32, compressed bool) {
	assets = append(assets, &assetInfo{name, ctype, data, hash, compressed})
}

// ServeAsset sends the requested asset.
func ServeAsset(r *request.Request, asset string) {
	// NOTE:  In production, ServeAsset only gets called when the requested
	// asset doesn't exist in the filesystem; otherwise, the asset is served
	// by Apache directly and the request never gets to this code.  If the
	// asset doesn't exist in the filesystem, it's probably a request for an
	// outdated version of a real asset (i.e., right name, wrong hash).
	// This code will respond with the current version of that asset.  Of
	// course, if it's truly a nonexistent asset, this code will return an
	// error.
	//
	// During development, when the portal is running as a full web server,
	// this code handles all asset requests.
	var info *assetInfo

	asset = hashRE.ReplaceAllLiteralString(asset, ".") // remove hash
	for _, ai := range assets {
		if ai.name == asset {
			info = ai
			break
		}
	}
	if info == nil {
		http.Error(r, "404 Not Found", http.StatusNotFound)
		return
	}
	// Assets are set to expire basically never (actually after one year),
	// because if we change them, they'll have a new hash and a new URL.
	r.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	r.Header().Set("Content-Type", info.ctype)
	r.DisableCompression() // either it's already compressed or it's not compressible
	if strings.Contains(r.Request.Header.Get("Accept-Encoding"), "gzip") && info.compressed {
		r.Header().Set("Content-Encoding", "gzip")
		r.Write(info.data)
	} else if info.compressed {
		// We have it compressed, but caller wants uncompressed, so uncompress it.
		// This is an unusual case; most browsers ask for compressed.
		if gz, err := gzip.NewReader(bytes.NewReader(info.data)); err == nil {
			io.Copy(r, gz)
			gz.Close()
		} else {
			panic(err)
		}
	} else { // it's not compressible; send it uncompressed
		r.Write(info.data)
	}
}

// AssetURL returns the URL for the asset with the specified name.  It panics if
// there is no such asset.
func AssetURL(asset string) (url string) {
	if strings.HasPrefix(asset, "https://") {
		return asset
	}
	for _, ai := range assets {
		if ai.name == asset {
			ext := filepath.Ext(asset)
			return fmt.Sprintf("/assets/%s--%08x%s", asset[:len(asset)-len(ext)], ai.hash, ext)
		}
	}
	panic(fmt.Sprintf("no such asset %q", asset))
}

// WriteAssetFiles writes all defined asset files to ../assets, removing any
// other files already there.
func WriteAssetFiles() (err error) {
	const path = "../assets"
	var (
		dir   *os.File
		names []string
		nmap  = make(map[string]bool)
	)

	if err = os.MkdirAll(path, 0777); err != nil {
		return fmt.Errorf("mkdir ../assets: %s", err)
	}
	if dir, err = os.Open(path); err != nil {
		return fmt.Errorf("open ../assets: %s", err)
	}
	if names, err = dir.Readdirnames(0); err != nil {
		return fmt.Errorf("read ../assets: %s", err)
	}
	for _, n := range names {
		if n != ".htaccess" {
			nmap[n] = true
		}
	}
	for _, ai := range assets {
		ext := filepath.Ext(ai.name)
		data := ai.data
		if ai.compressed {
			fname := fmt.Sprintf("%s--%08x%s.gz", ai.name[:len(ai.name)-len(ext)], ai.hash, ext)
			if nmap[fname] {
				delete(nmap, fname)
			} else {
				if err = os.WriteFile(filepath.Join(path, fname), data, 0666); err != nil {
					return fmt.Errorf("write %s: %s", fname, err)
				}
				fmt.Fprintf(os.Stderr, "Wrote %s\n", fname)
			}
			var rdr *gzip.Reader
			if rdr, err = gzip.NewReader(bytes.NewReader(data)); err != nil {
				return fmt.Errorf("uncompress %s: %s", fname, err)
			}
			if data, err = io.ReadAll(rdr); err != nil {
				rdr.Close()
				return fmt.Errorf("uncompress %s: %s", fname, err)
			}
			rdr.Close()
		}
		fname := fmt.Sprintf("%s--%08x%s", ai.name[:len(ai.name)-len(ext)], ai.hash, ext)
		if nmap[fname] {
			delete(nmap, fname)
		} else {
			if err = os.WriteFile(filepath.Join(path, fname), data, 0666); err != nil {
				return fmt.Errorf("write %s: %s", fname, err)
			}
			fmt.Fprintf(os.Stderr, "Wrote %s\n", fname)
		}
	}
	for n := range nmap {
		if err = os.Remove(filepath.Join(path, n)); err != nil {
			return fmt.Errorf("remove %s: %s", n, err)
		}
		fmt.Fprintf(os.Stderr, "Removed %s\n", n)
	}
	return nil
}

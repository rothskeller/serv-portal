package phys

import (
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"sunnyvaleserv.org/portal/util/config"
)

// IndexEntry is a single entry in the search index.
type IndexEntry struct {
	// Key is the unique identifier of the index entry (and of the item
	// indexed).
	Key string `json:"objectID"`
	// Type is the object's type, as shown in search results.
	Type string `json:"type"`
	// Label is the object's label, as shown in search results.
	Label string `json:"label"`
	// Context is an optional context for the object in search results
	// (e.g., the event containing a task).
	Context string `json:"context,omitempty"`
	// Name is the searchable name of the object.
	Name string `json:"name"`
	// Date is the optional searchable date corresponding to this object.
	Date string `json:"date,omitempty"`
	// CallSign is the optional searchable callsign of the object (person).
	CallSign string `json:"callsign,omitempty"`
}

// Indexer is an interface implemented by any object that can be indexed.
type Indexer interface {
	// IndexKey returns the index key for the object.
	IndexKey(storer Storer) string
	// IndexEntry returns a complete index entry for the object, or nil if
	// the object should not be indexed.  If the object has a parent object,
	// parent *may* give a pointer to it; otherwise, storer can be used to
	// fetch the parent.
	IndexEntry(storer Storer) *IndexEntry
}

// IndexParent is an interface implemented by an object that has dependent
// objects that can be indexed.
type IndexParent interface {
	// IndexKey returns the index key for the object.
	IndexKey(storer Storer) string
	// IndexContext returns a string describing the object, to be used in
	// the context strings of index entries for dependent objects.
	IndexContext(storer Storer) string
}

// client is the search client.
var client *search.Client

// openSearch creates the client interface for making search index updates.
func openSearch() {
	client = search.NewClient(config.Get("algoliaApplicationID"), config.Get("algoliaUpdateKey"))
}

// Search runs a search for the specified query string.  Only the Key and Type
// fields of the returned results are provided.
func Search(storer Storer, query string) (results []*IndexEntry, err error) {
	index := client.InitIndex(config.Get("algoliaIndex"))
	apires, err := index.Search(query, opt.HitsPerPage(50), opt.AttributesToRetrieve("objectID", "type"))
	if err != nil {
		return nil, err
	}
	for _, hit := range apires.Hits {
		var ie = IndexEntry{
			Key:  hit["objectID"].(string),
			Type: hit["type"].(string),
		}
		results = append(results, &ie)
	}
	return results, nil
}

// Index queues an object to be added or updated in the search index.
func Index(storer Storer, object Indexer) {
	var (
		entry *IndexEntry
		store = storer.AsStore()
	)
	if entry = object.IndexEntry(storer); entry == nil {
		return
	}
	store.tx.searchOps = append(store.tx.searchOps, search.BatchOperationIndexed{
		IndexName: config.Get("algoliaIndex"),
		BatchOperation: search.BatchOperation{
			Action: search.UpdateObject,
			Body:   entry,
		},
	})
}

// Unindex queues an object to be removed from the search index.
func Unindex(storer Storer, object Indexer) {
	key := object.IndexKey(storer)
	if key == "" {
		return
	}
	store := storer.AsStore()
	store.tx.searchOps = append(store.tx.searchOps, search.BatchOperationIndexed{
		IndexName: config.Get("algoliaIndex"),
		BatchOperation: search.BatchOperation{
			Action: search.DeleteObject,
			Body:   deleteObjectID{key},
		},
	})
}

type deleteObjectID struct {
	ObjectID string `json:"objectID"`
}

// applySearchOps sends the accumulated search index operations to the server.
func (store *Store) applySearchOps() (err error) {
	if len(store.tx.searchOps) == 0 {
		return nil
	}
	_, err = client.MultipleBatch(store.tx.searchOps)
	return err
}

// EmptyEntireIndex deletes all index entries from the index.  Unlike other
// index methods, it waits for remote completion before returning.
func EmptyEntireIndex() {
	opened.Do(open)
	index := client.InitIndex(config.Get("algoliaIndex"))
	res, err := index.ClearObjects()
	if err != nil {
		panic(err)
	}
	if err = res.Wait(); err != nil {
		panic(err)
	}
}

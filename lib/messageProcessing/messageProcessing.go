package messageProcessing

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"regexp"

	"github.com/EVE-Tools/emdr-to-nsq/lib/emds"
	"github.com/boltdb/bolt"
	"github.com/buger/jsonparser"
	"github.com/sirupsen/logrus"
	"github.com/spaolacci/murmur3"
)

var db *bolt.DB

// Initialize creates the hash cache
func Initialize(database *bolt.DB) {
	db = database

	// Initialize buckets
	err := db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("rowsetHashes"))
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// FilterMessage filters and parses EMDR messages into structs.
func FilterMessage(nameRegex *regexp.Regexp, versionRegex *regexp.Regexp, message []byte) ([]emds.Rowset, error) {
	// First, check message type and generator info
	discard, err := discardMessage(message, nameRegex, versionRegex)
	if err != nil {
		logrus.WithError(err).Warn("Got UUDIF uncompliant message.")
		return nil, err
	}
	if discard {
		return nil, err
	}

	// Second, parse region/type info and check if we've already seen these orders
	rawRowsets, err := emds.ExtractRawRowsets(message)
	if err != nil {
		logrus.WithError(err).Warn("Got UUDIF uncompliant message.")
		return nil, err
	}

	filteredRawRowsets, err := filterRowsets(rawRowsets)
	if err != nil {
		logrus.WithError(err).Warn("Error filtering rowsets.")
		return nil, err
	}

	if len(filteredRawRowsets) == 0 {
		return nil, nil
	}

	// Finally, get column indices and parse orders
	indices, err := emds.GetColumnIndices(message)
	if err != nil {
		logrus.WithError(err).Error("Error extracting columns.")
		return nil, err
	}

	return emds.ParseRawRowsets(filteredRawRowsets, indices)
}

// Filters raw rowsets which we've already seen. Only do this async to leverage BoltDB batch processing.
func filterRowsets(rowsets []emds.RawRowset) ([]emds.RawRowset, error) {
	var acceptedRowsets []emds.RawRowset
	success := make(chan *emds.RawRowset, len(rowsets))
	failure := make(chan error, len(rowsets))

	for index := range rowsets {
		go discardRowsetAsync(success, failure, &rowsets[index])
	}

	for range rowsets {
		select {
		case rowset := <-success:
			if rowset != nil {
				acceptedRowsets = append(acceptedRowsets, *rowset)
			}
		case err := <-failure:
			logrus.WithError(err).Warn("Error filtering rowset.")
		}
	}

	if len(rowsets) > 0 {
		logrus.WithFields(logrus.Fields{
			"regionID":   rowsets[0].RegionID,
			"rowsets":    len(rowsets),
			"newRowsets": len(acceptedRowsets),
		}).Info("Received rowsets from EMDR.")
	}

	return acceptedRowsets, nil
}

// Tests if message should be discarded (e.g. it contains history instead of orders).
func discardMessage(message []byte, nameRegex *regexp.Regexp, versionRegex *regexp.Regexp) (bool, error) {
	// Check if order message
	resultType, err := jsonparser.GetString(message, "resultType")
	if err != nil {
		return true, err
	}

	if resultType != "orders" {
		return true, nil
	}

	// Check NameRegex if needed
	if nameRegex != nil {
		generatorName, err := jsonparser.GetString(message, "generator", "name")
		if err != nil {
			return true, err
		}

		if !nameRegex.MatchString(generatorName) {
			return true, err
		}
	}

	// Check VersionRegex if needed
	if versionRegex != nil {
		generatorVersion, err := jsonparser.GetString(message, "generator", "version")
		if err != nil {
			logrus.WithError(err).Warn("Got UUDIF uncompliant message.")
			return true, err
		}

		if !versionRegex.MatchString(generatorVersion) {
			return true, err
		}
	}

	return false, nil
}

// Checks a rowset asnychronously and replies via channels, returns nil if discarded
func discardRowsetAsync(success chan<- *emds.RawRowset, failure chan<- error, rowset *emds.RawRowset) {
	discard, err := discardRowset(rowset)
	if err != nil {
		logrus.WithError(err).Warn("Error parsing rowset.")
		failure <- err
		return
	}

	if !discard {
		success <- rowset
	} else {
		success <- nil
	}
}

// Tests if rowset should be discarded (duplicate content for a region/type or no changes since last time).
func discardRowset(rowset *emds.RawRowset) (bool, error) {
	var hash []byte
	key := []byte(fmt.Sprintf("%v-%v", rowset.RegionID, rowset.TypeID))

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("rowsetHashes"))
		if bucket == nil {
			panic("Bucket not found! This should never happen!")
		}

		hash = bucket.Get(key)
		return nil
	})

	// Calculate hash and compare, store new hash if needed
	newHash := make([]byte, 4) // 32 bits
	binary.LittleEndian.PutUint32(newHash, murmur3.Sum32(rowset.Rows))

	// Equal hash -> filter
	if bytes.Compare(hash, newHash) == 0 {
		return true, nil
	}

	// New hash -> store
	err := db.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("rowsetHashes"))
		if bucket == nil {
			panic("Bucket not found! This should never happen!")
		}
		return bucket.Put(key, newHash)
	})

	if err != nil {
		return false, err
	}

	return false, nil
}

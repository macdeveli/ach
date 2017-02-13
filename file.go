// Copyright 2017 The ACH Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

// Package ach reads and writes (ACH) Automated Clearing House files. ACH is the
// primary method of electronic money movemenet through the United States.
//
// https://en.wikipedia.org/wiki/Automated_Clearing_House
package ach

import "errors"

// First position of all Record Types. These codes are uniquily assigned to
// the first byte of each row in a file.
const (
	headerPos       = "1"
	batchPos        = "5"
	entryDetailPos  = "6"
	entryAddendaPos = "7"
	batchControlPos = "8"
	fileControlPos  = "9"
)

// Errors specific to parsing a Batch container
var (
	ErrFileBatchCount = errors.New("Total Number of Batches in file is out-of-balance with File Control")
)

// File contains the structures of a parsed ACH File.
type File struct {
	Header  FileHeader
	Batches []Batch
	Control FileControl

	// TODO: remove
	Addenda
}

// addEntryDetail appends an EntryDetail to the Batch
func (f *File) addBatch(batch Batch) []Batch {
	f.Batches = append(f.Batches, batch)
	return f.Batches
}

// Validate NACHA rules on the entire batch before being added to a File
func (f *File) Validate() error {
	// The value of the Batch Count Field is equal to the number of Company/Batch/Header Records in the file.
	if f.Control.BatchCount != len(f.Batches) {
		return ErrFileBatchCount
	}
	return nil
}

// TODO: isEntryHashMismatch
// This field is prepared by hashing the RDFI’s 8-digit Routing Number in each entry.
//The Entry Hash provides a check against inadvertent alteration of data

// TODO: isFileAmountMismatch
// The Total Debit and Credit Entry Dollar Amounts Fields contain accumulated
// Entry Detail debit and credit totals within the file

// TODO:

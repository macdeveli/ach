// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"strings"
	"testing"
)

func mockBatchADVControl() *BatchADVControl {
	bc := NewBatchADVControl()
	bc.ServiceClassCode = 220
	bc.ACHOperatorData = "T-BANK"
	bc.ODFIIdentification = "12104288"
	return bc
}

// testMockBatchADVControl tests mock batch control
func testMockBatchADVControl(t testing.TB) {
	bc := mockBatchADVControl()
	if err := bc.Validate(); err != nil {
		t.Error("mockBatchADVControl does not validate and will break other tests")
	}
	if bc.ServiceClassCode != 220 {
		t.Error("ServiceClassCode depedendent default value has changed")
	}
	if bc.ACHOperatorData != "T-BANK" {
		t.Error("ACHOperatorData depedendent default value has changed")
	}
	if bc.ODFIIdentification != "12104288" {
		t.Error("ODFIIdentification depedendent default value has changed")
	}
}

// TestMockBatchADVControl test mock batch control
func TestMockBatchADVControl(t *testing.T) {
	testMockBatchADVControl(t)
}

// BenchmarkMockBatchADVControl benchmarks mock batch control
func BenchmarkMockBatchADVControl(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testMockBatchADVControl(b)
	}
}

// TestParseBatchADVControl parses a known Batch ControlRecord string.
func testParseBatchADVControl(t testing.TB) {
	var line = "822500000100053200010000000000000001050000000000000000000000T-BANK             076401250000001"
	r := NewReader(strings.NewReader(line))
	r.line = line
	bh := BatchHeader{BatchNumber: 1,
		StandardEntryClassCode: "ADV",
		ServiceClassCode:       225,
		CompanyIdentification:  "origid",
		ODFIIdentification:     "7640125"}
	r.addCurrentBatch(NewBatchADV(&bh))

	r.currentBatch.AddEntry(&EntryDetail{TransactionCode: 27, Amount: 10500, RDFIIdentification: "5320001", TraceNumber: 76401255655291})
	if err := r.parseBatchControl(); err != nil {
		t.Errorf("%T: %s", err, err)
	}
	record := r.currentBatch.GetADVControl()

	if record.recordType != "8" {
		t.Errorf("RecordType Expected '8' got: %v", record.recordType)
	}
	if record.ServiceClassCode != 225 {
		t.Errorf("ServiceClassCode Expected '225' got: %v", record.ServiceClassCode)
	}
	if record.EntryAddendaCountField() != "000001" {
		t.Errorf("EntryAddendaCount Expected '000001' got: %v", record.EntryAddendaCountField())
	}
	if record.EntryHashField() != "0005320001" {
		t.Errorf("EntryHash Expected '0005320001' got: %v", record.EntryHashField())
	}
	if record.TotalDebitEntryDollarAmountField() != "00000000000000010500" {
		t.Errorf("TotalDebitEntryDollarAmount Expected '00000000000000010500' got: %v", record.TotalDebitEntryDollarAmountField())
	}
	if record.TotalCreditEntryDollarAmountField() != "00000000000000000000" {
		t.Errorf("TotalCreditEntryDollarAmount Expected '00000000000000000000' got: %v", record.TotalCreditEntryDollarAmountField())
	}
	if record.ACHOperatorDataField() != "T-BANK             " {
		t.Errorf("ACHOperatorData Expected 'T-BANK             ' got: %v", record.ACHOperatorDataField())
	}
	if record.ODFIIdentificationField() != "07640125" {
		t.Errorf("OdfiIdentification Expected '07640125' got: %v", record.ODFIIdentificationField())
	}
	if record.BatchNumberField() != "0000001" {
		t.Errorf("BatchNumber Expected '0000001' got: %v", record.BatchNumberField())
	}
}

// TestParseBatchADVControl tests parsing a known Batch ControlRecord string.
func TestParseBatchADVControl(t *testing.T) {
	testParseBatchADVControl(t)
}

// BenchmarkParseBatchADVControl benchmarks parsing a known Batch ControlRecord string.
func BenchmarkParseBatchADVControl(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testParseBatchADVControl(b)
	}
}

// testADVBCString validates that a known parsed file can be return to a string of the same value
func testADVBCString(t testing.TB) {
	var line = "822500000100053200010000000000000001050000000000000000000000T-BANK             076401250000001"
	r := NewReader(strings.NewReader(line))
	r.line = line
	bh := BatchHeader{BatchNumber: 1,
		StandardEntryClassCode: "ADV",
		ServiceClassCode:       225,
		CompanyIdentification:  "origid",
		ODFIIdentification:     "7640125"}
	r.addCurrentBatch(NewBatchADV(&bh))

	r.currentBatch.AddEntry(&EntryDetail{TransactionCode: 27, Amount: 10500, RDFIIdentification: "5320001", TraceNumber: 76401255655291})
	if err := r.parseBatchControl(); err != nil {
		t.Errorf("%T: %s", err, err)
	}
	record := r.currentBatch.GetADVControl()

	if record.String() != line {
		t.Errorf("Strings do not match")
	}
}

// TestADVBCString tests validating that a known parsed file can be return to a string of the same value
func TestADVBCString(t *testing.T) {
	testADVBCString(t)
}

// BenchmarkADVBCString benchmarks validating that a known parsed file can be return to a string of the same value
func BenchmarkADVBCString(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testADVBCString(b)
	}
}

// testValidateADVBCRecordType ensure error if recordType is not 8
func testValidateADVBCRecordType(t testing.TB) {
	bc := mockBatchADVControl()
	bc.recordType = "2"
	if err := bc.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "recordType" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestValidateADVBCRecordType tests ensuring an error if recordType is not 8
func TestValidateADVBCRecordType(t *testing.T) {
	testValidateADVBCRecordType(t)
}

// BenchmarkValidateADVBCRecordType benchmarks ensuring an error if recordType is not 8
func BenchmarkValidateADVBCRecordType(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testValidateBCRecordType(b)
	}
}

// testADVisServiceClassErr verifies service class code
func testADVBCisServiceClassErr(t testing.TB) {
	bc := mockBatchADVControl()
	bc.ServiceClassCode = 123
	if err := bc.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "ServiceClassCode" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestADVBCisServiceClassErr tests verifying service class code
func TestADVBCisServiceClassErr(t *testing.T) {
	testADVBCisServiceClassErr(t)
}

// BenchmarkADVBCisServiceClassErr benchmarks verifying service class code
func BenchmarkADVBCisServiceClassErr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testADVBCisServiceClassErr(b)
	}
}

// testADVBCBatchNumber verifies batch number
func testADVBCBatchNumber(t testing.TB) {
	bc := mockBatchADVControl()
	bc.BatchNumber = 0
	if err := bc.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "BatchNumber" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestADVBCBatchNumber tests verifying batch number
func TestADVBCBatchNumber(t *testing.T) {
	testADVBCBatchNumber(t)
}

// BenchmarkADVBCBatchNumber benchmarks verifying batch number
func BenchmarkADVBCBatchNumber(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testADVBCBatchNumber(b)
	}
}

// testADVBCACHOperatorDataAlphaNumeric verifies Company Identification is AlphaNumeric
func testADVBCACHOperatorDataAlphaNumeric(t testing.TB) {
	bc := mockBatchADVControl()
	bc.ACHOperatorData = "®"
	if err := bc.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "ACHOperatorData" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestADVBCACHOperatorDataAlphaNumeric tests verifying Company Identification is AlphaNumeric
func TestADVBCACHOperatorDataAlphaNumeric(t *testing.T) {
	testADVBCACHOperatorDataAlphaNumeric(t)
}

// BenchmarkADVACHOperatorDataAlphaNumeric benchmarks verifying Company Identification is AlphaNumeric
func BenchmarkADVACHOperatorDataAlphaNumeric(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testADVBCACHOperatorDataAlphaNumeric(b)
	}
}

// testADVBCFieldInclusionRecordType verifies Record Type is included
func testADVBCFieldInclusionRecordType(t testing.TB) {
	bc := mockBatchADVControl()
	bc.recordType = ""
	if err := bc.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if !strings.Contains(e.Msg, msgFieldInclusion) {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestADVBCFieldInclusionRecordType tests verifying Record Type is included
func TestADVBCFieldInclusionRecordType(t *testing.T) {
	testADVBCFieldInclusionRecordType(t)
}

// BenchmarkADVBCFieldInclusionRecordType benchmarks verifying Record Type is included
func BenchmarkADVBCFieldInclusionRecordType(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testADVBCFieldInclusionRecordType(b)
	}
}

// testADVBCFieldInclusionServiceClassCode verifies Service Class Code is included
func testADVBCFieldInclusionServiceClassCode(t testing.TB) {
	bc := mockBatchADVControl()
	bc.ServiceClassCode = 0
	if err := bc.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if !strings.Contains(e.Msg, msgFieldInclusion) {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestADVBCFieldInclusionServiceClassCode tests verifying Service Class Code is included
func TestADVBCFieldInclusionServiceClassCode(t *testing.T) {
	testADVBCFieldInclusionServiceClassCode(t)
}

// BenchmarkADVBCFieldInclusionServiceClassCod benchmarks verifying Service Class Code is included
func BenchmarkADVBCFieldInclusionServiceClassCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testADVBCFieldInclusionServiceClassCode(b)
	}
}

// testADVBCFieldInclusionODFIIdentification verifies batch control ODFIIdentification
func testADVBCFieldInclusionODFIIdentification(t testing.TB) {
	bc := mockBatchADVControl()
	bc.ODFIIdentification = "000000000"
	if err := bc.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if !strings.Contains(e.Msg, msgFieldInclusion) {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestADVBCFieldInclusionODFIIdentification tests verifying batch control ODFIIdentification
func TestADVBCFieldInclusionODFIIdentification(t *testing.T) {
	testADVBCFieldInclusionODFIIdentification(t)
}

// BenchmarkADVBCFieldInclusionODFIIdentification benchmarks verifying batch control ODFIIdentification
func BenchmarkADVBCFieldInclusionODFIIdentification(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testADVBCFieldInclusionODFIIdentification(b)
	}
}

// testBatchADVControlLength verifies batch control length
func testBatchADVControlLength(t testing.TB) {
	bc := NewBatchADVControl()
	recordLength := len(bc.String())
	if recordLength != 94 {
		t.Errorf("Instantiated length of Batch Control string is not 94 but %v", recordLength)
	}
}

// TestBatchADVControlLength tests verifying batch control length
func TestBatchADVControlLength(t *testing.T) {
	testBatchADVControlLength(t)
}

// BenchmarkBatchADVControlLength benchmarks verifying batch control length
func BenchmarkBatchADVControlLength(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchADVControlLength(b)
	}
}

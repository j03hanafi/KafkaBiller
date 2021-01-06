package main

import (
	"github.com/go-yaml/yaml"
	"github.com/mofax/iso8583"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

// Any helper to process ISO data
// converter, formatter, etc

// Convert JSON data to ISO8583 format
func convertIso(transaction Transaction) iso8583.IsoStruct {

	log.Println("Converting JSON to ISO8583")

	cardAcceptorTerminalId := transaction.CardAcceptorData.CardAcceptorTerminalId
	cardAcceptorName := transaction.CardAcceptorData.CardAcceptorName
	cardAcceptorCity := transaction.CardAcceptorData.CardAcceptorCity
	cardAcceptorCountryCode := transaction.CardAcceptorData.CardAcceptorCountryCode

	if len(transaction.CardAcceptorData.CardAcceptorTerminalId) < 16 {
		cardAcceptorTerminalId = rightPad(transaction.CardAcceptorData.CardAcceptorTerminalId, 16, " ")
	}
	if len(transaction.CardAcceptorData.CardAcceptorName) < 25 {
		cardAcceptorName = rightPad(transaction.CardAcceptorData.CardAcceptorName, 25, " ")
	}
	if len(transaction.CardAcceptorData.CardAcceptorCity) < 13 {
		cardAcceptorCity = rightPad(transaction.CardAcceptorData.CardAcceptorCity, 13, " ")
	}
	if len(transaction.CardAcceptorData.CardAcceptorCountryCode) < 2 {
		cardAcceptorCountryCode = rightPad(transaction.CardAcceptorData.CardAcceptorCountryCode, 2, " ")
	}
	cardAcceptor := cardAcceptorName + cardAcceptorCity + cardAcceptorCountryCode

	trans := map[int64]string{
		2:  transaction.Pan,
		3:  transaction.ProcessingCode,
		4:  strconv.Itoa(transaction.TotalAmount),
		5:  transaction.SettlementAmount,
		6:  transaction.CardholderBillingAmount,
		7:  transaction.TransmissionDateTime,
		9:  transaction.SettlementConversionRate,
		10: transaction.CardHolderBillingConvRate,
		11: transaction.Stan,
		12: transaction.LocalTransactionTime,
		13: transaction.LocalTransactionDate,
		17: transaction.CaptureDate,
		18: transaction.CategoryCode,
		22: transaction.PointOfServiceEntryMode,
		37: transaction.Refnum,
		41: cardAcceptorTerminalId,
		43: cardAcceptor,
		48: transaction.AdditionalData,
		49: transaction.Currency,
		50: transaction.SettlementCurrencyCode,
		51: transaction.CardHolderBillingCurrencyCode,
		57: transaction.AdditionalDataNational,
	}

	one := iso8583.NewISOStruct("spec1987.yml", false)
	spec, _ := specFromFile("spec1987.yml")

	if one.Mti.String() != "" {
		log.Printf("Empty generates invalid MTI")
	}

	for field, data := range trans {

		fieldSpec := spec.fields[int(field)]

		if fieldSpec.LenType == "fixed" {
			lengthValidate, _ := iso8583.FixedLengthIntegerValidator(int(field), fieldSpec.MaxLen, data)

			if lengthValidate == false {
				if fieldSpec.ContentType == "n" {
					data = leftPad(data, fieldSpec.MaxLen, "0")
				} else {
					data = rightPad(data, fieldSpec.MaxLen, " ")
				}
			}
		}

		one.AddField(field, data)

		log.Printf("Field[%s]: %s (%v)", strconv.Itoa(int(field)), data, fieldSpec.Label)

	}

	log.Println("Convert Success")
	return one
}

// Spec contains a structured description of an iso8583 spec
// properly defined by a spec file
type Spec struct {
	fields map[int]fieldDescription
}

// readFromFile reads a yaml specfile and loads
// and iso8583 spec from it
func (s *Spec) readFromFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	yaml.Unmarshal(content, &s.fields) // expecting content to be valid yaml
	return nil
}

// SpecFromFile returns a brand new empty spec
func specFromFile(filename string) (Spec, error) {
	s := Spec{}
	err := s.readFromFile(filename)
	if err != nil {
		return s, err
	}
	return s, nil
}

// Add pad on left of data,
// Used to format number by adding "0" in front of number data
func leftPad(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(pad, length-len(s))
	return padding + s
}

// Add pad on right of data,
// Used to format string by adding " " at the end of string data
func rightPad(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(pad, length-len(s))
	return s + padding
}

package sign

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"math/big"

	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
)

func MarshalToEtreeDocument(object any) (*etree.Document, error) {
	bytes, err := xml.Marshal(object)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(bytes); err != nil {
		return nil, err
	}
	return doc, nil
}

func calculateDigest(template *Signature, path string, c14n dsig.Canonicalizer) ([]byte, error) {
	document, err := MarshalToEtreeDocument(template)
	if err != nil {
		return nil, err
	}

	element := document.Root().FindElement(path)
	if element == nil {
		return nil, fmt.Errorf("calculateDigest: element with path %s not found in root", path)
	}

	canonical, err := c14n.Canonicalize(element)
	if err != nil {
		return nil, err
	}

	digest := sha256.Sum256(canonical)

	return digest[:], nil
}

func calculateSignature(template *Signature, path string, key *ecdsa.PrivateKey, c14n dsig.Canonicalizer) (*big.Int, *big.Int, error) {
	digest, err := calculateDigest(template, path, c14n)
	if err != nil {
		return nil, nil, err
	}
	return ecdsa.Sign(rand.Reader, key, digest)
}

func CanonicalSerialize(el *etree.Element) ([]byte, error) {
	doc := etree.NewDocument()
	doc.SetRoot(el.Copy())

	doc.WriteSettings = etree.WriteSettings{
		CanonicalAttrVal: true,
		CanonicalEndTags: false,
		CanonicalText:    true,
	}

	return doc.WriteToBytes()
}

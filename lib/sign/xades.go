package sign

import (
	"crypto/sha256"
	"fmt"

	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
)

func buildQualifiedSignedProperties(signedProperties *etree.Element) *etree.Element {
	qualifiedSignedProperties := signedProperties.Copy()
	qualifiedSignedProperties.Attr = append(
		signedProperties.Attr,
		etree.Attr{Space: "xmlns", Key: "", Value: dsig.Namespace},
		etree.Attr{Space: "xmlns", Key: "xades", Value: "http://uri.etsi.org/01903/v1.3.2#"},
	)
	return qualifiedSignedProperties
}

func calculateSignedPropertiesDigest(signature *Signature, c14n dsig.Canonicalizer) ([]byte, error) {
	doc, err := MarshalToEtreeDocument(signature)
	if err != nil {
		return nil, err
	}

	element := doc.Root().FindElement(".//[@Id='SignedProperties']")
	if element == nil {
		return nil, fmt.Errorf("SignedProperties element not found")
	}

	qualified := buildQualifiedSignedProperties(element)

	canonical, err := c14n.Canonicalize(qualified)
	if err != nil {
		return nil, err
	}

	digest := sha256.Sum256(canonical)
	return digest[:], nil
}

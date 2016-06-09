package xml

import (
	"encoding/xml"
)

type Bandwidth struct {
	XMLName xml.Name `xml:"config"`
	Comment string   `xml:"config,omitempty"`

	Identity       *BandwidthIdentity `xml:"identity"`
	BandwidthLimit *BandwidthLimit    `xml:"limit"`
}

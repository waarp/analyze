package rest

import (
	"fmt"
	"strings"
	"time"
)

// TransferServices handles communication with the transfers related
// methods.
type TransferServices struct {
	client *RestClient
}

// Transfer represents a R66 transfer
type Transfer struct {
	ti *transferInfo

	// R66 internal transfer ID
	Id    int64  `json:"idInt"`
	IdStr string `json:"id"`

	// the original basename of the transfered file.
	Filename string            `json:"filename"`
	FileInfo map[string]string `json:"fileinfo"`

	// The size is only available on the origin partner. On the
	// destination partner, it is -1.
	Size int64 `json:"size"`

	// The name of the rule used for the transfer
	Rule         string         `json:"rulename"`
	TransferMode waarp.RuleMode `json:"transferMode"`

	// When the transfer started
	StartTime time.Time `json:"starttime"`
	// When the transfer ended. If the transfer is running, it is the last update time
	EndTime time.Time `json:"endtime"`
	// Duration of the transfer
	Duration time.Duration `json:"duration"`
	// Speed of the transfer
	Speed float64 `json:"speed"`

	// Name of the origin partner
	Origin string `json:"originHostid"`
	// Name of the destination partner
	Destination string `json:"destinationHostid"`

	// The state of the transfer
	Status waarp.TransferStatus `json:"status"`
	// Detailed information about the status
	StatusInfo string `json:"statusinfo"`
}

//----------------------------------------------------------------------
// Methods
//----------------------------------------------------------------------

// Get a list of transfers from a Waarp R66 server. data is a
// ListTransferRequest object containing query filters.
func (t *TransferServices) List(data ListTransferRequest) ([]*Transfer, error) {
	req, err := NewRequest("GET", t.client.url("/transfers"), data)
	if err != nil {
		return nil, fmt.Errorf(`An error occured while preparing the request: %s`, err.Error())
	}

	results := []transferInfo{}

	if err := t.client.exec(req, &results); err != nil {
		return nil, err
	}

	rv := make([]*Transfer, len(results))
	for i, ti := range results {
		rv[i] = ti.ToTransfer()
	}

	return rv, nil
}

func (t *TransferServices) ListAll(data ListTransferRequest) ([]*Transfer, error) {
	max := time.Now()
	// min := max.Add(-12 * 30 * 24 * time.Hour)
	min := maxDate(*data.StartAfter, max.Add(-12*30*24*time.Hour))
	var allTransfers []*Transfer

	for {
		data.StartAfter = &min
		data.StartBefore = &max
		data.Limit = 50

		transfers, err := t.List(data)
		if err != nil {
			return nil, err
		}

		for _, transfer := range transfers {
			allTransfers = append(allTransfers, transfer)
			max = minDate(max, transfer.StartTime).Add(-1 * time.Millisecond)
		}
		if len(transfers) < data.Limit {
			break
		}
	}

	return allTransfers, nil
}

func (t *TransferServices) Restart(data TransferControlRequest) error {
	data.Class = "org.waarp.openr66.protocol.localhandler.packet.json.RestartTransferJsonPacket"
	data.RequestUserPacket = 4
	return t.doControlRequest(data)
}

func (t *TransferServices) Stop(data TransferControlRequest) error {
	data.Class = "org.waarp.openr66.protocol.localhandler.packet.json.StopOrCancelJsonPacket"
	data.RequestUserPacket = 9
	return t.doControlRequest(data)
}

func (t *TransferServices) Cancel(data TransferControlRequest) error {
	data.Class = "org.waarp.openr66.protocol.localhandler.packet.json.StopOrCancelJsonPacket"
	data.RequestUserPacket = 10
	return t.doControlRequest(data)
}

func (t *TransferServices) doControlRequest(data TransferControlRequest) error {
	req, err := NewRequest("PUT", t.client.url("/control"), data)
	if err != nil {
		return fmt.Errorf(`An error occured while preparing the request: %s`, err.Error())
	}

	results := []TransferControlResponse{}

	if err := t.client.exec(req, &results); err != nil {
		return err
	}

	return nil
}

//----------------------------------------------------------------------
// Messages
//----------------------------------------------------------------------

// A transfer information as returned by Waarp R66 REST server
type transferInfo struct {
	Model          string             `json:"@model"` //    "@model":"DbTaskRunner"
	GlobalStep     waarp.Step           //    "GLOBALSTEP":4
	GlobalLastStep waarp.Step           //    "GLOBALLASTSTEP":4
	Step           int                  //    "STEP":0
	Rank           int                  //    "RANK":1
	StepStatus     waarp.ErrorCode      //    "STEPSTATUS":"O  "
	RetrieveMode   bool                 //    "RETRIEVEMODE":false
	Filename       string               //    "FILENAME":"/data/in/-9223372036854775801_8996631326837891415_client.xml"
	IsMoved        bool                 //    "ISMOVED":false
	Rule           string               `json:"IDRULE"`  //    "IDRULE":"push"
	BlockSize      int                  `json:"BLOCKSZ"` //    "BLOCKSZ":65536
	OriginalName   string               //    "ORIGINALNAME":"/out/client.xml"
	FileInfo       string               //    "FILEINFO":"noinfo"
	TransferInfo   string               //    "TRANSFERINFO":"{}"
	ModeTrans      waarp.RuleMode       //    "MODETRANS":1
	StartTrans     int64                //    "STARTTRANS":1413386036548
	StopTrans      int64                //    "STOPTRANS":1413386036722
	InfoStatus     waarp.ErrorCode      //    "INFOSTATUS":"O  "
	OwnerReq       string               //    "OWNERREQ":"server1"
	Requester      string               //    "REQUESTER":"client11"
	Requested      string               //    "REQUESTED":"server1"
	Id             int64                `json:"SPECIALID"` //    "SPECIALID":-9223372036854775801
	OriginalSize   int64                //    "ORIGINALSIZE":-1
	Status         waarp.TransferStatus `json:"UPDATEDINFO"`
}

func (t *transferInfo) ToTransfer() *Transfer {
	transfer := &Transfer{
		ti:        t,
		Id:        t.Id,
		IdStr:     fmt.Sprintf("%d|%s|%s", t.Id, t.Requester, t.Requested),
		Filename:  basename(t.OriginalName),
		Rule:      t.Rule,
		StartTime: tsToTime(t.StartTrans),
		EndTime:   tsToTime(t.StopTrans),
		Status:    waarp.TransferStatus(t.Status),
	}
	transfer.Duration = transfer.EndTime.Sub(transfer.StartTime)
	if t.OriginalSize != -1 {
		transfer.Size = t.OriginalSize
	} else {
		transfer.Size = int64(t.BlockSize * t.Rank)
	}
	transfer.Speed = float64(t.BlockSize*t.Rank) / transfer.Duration.Seconds()

	transfer.TransferMode = t.ModeTrans

	transfer.FileInfo = makeFileinfo(t.FileInfo)

	switch t.ModeTrans {
	case waarp.SEND, waarp.SENDMD5, waarp.SENDTHROUGH, waarp.SENDMD5THROUGH:
		transfer.Origin = t.Requester
		transfer.Destination = t.Requested
	case waarp.RECV, waarp.RECVMD5, waarp.RECVTHROUGH, waarp.RECVMD5THROUGH:
		transfer.Origin = t.Requested
		transfer.Destination = t.Requester
	}

	if transfer.Status == waarp.TS_INERROR {
		switch t.InfoStatus {
		case waarp.EC_CANCELED_TRANSFER:
			transfer.Status = waarp.TS_CANCELED
			transfer.StatusInfo = ""

		case waarp.EC_STOPPED_TRANSFER:
			transfer.Status = waarp.TS_STOPPED
			transfer.StatusInfo = ""

		default:
			transfer.StatusInfo = fmt.Sprintf(
				"L'erreur suivante s'est produite durant l'étape %s: %s",
				t.GlobalLastStep, t.InfoStatus.FullString())
		}
	}

	return transfer
}

type TransferControlResponse struct {
	Class             string `json:"@class"`            //"org.waarp.openr66.protocol.localhandler.packet.json.RestartTransferJsonPacket",
	RequestUserPacket int    `json:"requestUserPacket"` // 4
	Requester         string `json:"requester"`
	Requested         string `json:"requested"`
	TransferId        int64  `json:"specialid"`
	restarttime       int64  `json:"restarttime,omitempty"`
}

// The request parameters for the metod "list transfer"
type ListTransferRequest struct {
	// Number of results to fetch (max: 100)
	Limit int `json:"LIMIT"`

	// Filter on statuses
	AllStatus  bool `json:"ALLSTATUS,omitempty"`
	InError    bool `json:"INERROR,omitempty"`
	InTransfer bool `json:"INTRANSFER,omitempty"`
	Done       bool `json:"DONE,omitempty"`
	Pending    bool `json:"PENDING,omitempty"` // "boolean"

	// Get only transfers that started before StartBefore and after StartAfter
	StartAfter  *time.Time `json:"STARTTRANS,omitempty"`
	StartBefore *time.Time `json:"STOPTRANS,omitempty"`

	// Order the transfers by id. By default, they are sortet by start
	// date descending
	OrderById bool `json:"ORDERBYID,omitempty"` // "boolean"

	// Filter on a range of transfer ids
	StartId int64 `json:"STARTID,omitempty"` // "transfer id"
	StopId  int64 `json:"STOPID,omitempty"`  // "transfer id",

	// Get only transfers for the rule Rule
	Rule string `json:"IDRULE,omitempty"` // "rule name"

	// Get only transfers with partner Partner
	Partner string `json:"PARTNER,omitempty"` // "partner (requester or requested) name"

	// R66 instance that owns the request
	Owner string `json:"OWNERREQ,omitempty"`
}

type TransferControlRequest struct {
	Class             string `json:"@class"`            //"org.waarp.openr66.protocol.localhandler.packet.json.RestartTransferJsonPacket",
	RequestUserPacket int    `json:"requestUserPacket"` // 4
	Requester         string `json:"requester"`
	Requested         string `json:"requested"`
	TransferId        int64  `json:"specialid"`
	RestartTime       int64  `json:"restarttime,omitempty"`
}

//----------------------------------------------------------------------
// Utilities
//----------------------------------------------------------------------

func basename(p string) string {
	var sep string
	if strings.Contains(p, "/") {
		sep = "/"
	} else {
		sep = "\\"

	}
	parts := strings.Split(p, sep)
	return parts[len(parts)-1]
}

func tsToTime(ts int64) time.Time {
	sec := ts / 1000
	nano := (ts % 1000) * 1e6
	return time.Unix(sec, nano)
}

func minDate(first time.Time, others ...time.Time) time.Time {
	min := first
	for _, d := range others {
		if d.Before(min) {
			min = d
		}
	}
	return min
}

func maxDate(first time.Time, others ...time.Time) time.Time {
	max := first
	for _, d := range others {
		if d.After(max) {
			max = d
		}
	}
	return max
}

func makeFileinfo(fi string) map[string]string {
	rv := map[string]string{}
	for _, info := range strings.Fields(fi) {
		if !strings.Contains(info, "=") {
			continue
		}

		kv := strings.SplitN(info, "=", 2)
		rv[kv[0]] = kv[1]
	}
	return rv
}

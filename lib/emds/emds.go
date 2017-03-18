//
// EVE Market Data Structures in go
// Provides functions for dealing with UUDIF-formatted data.
// See: http://dev.eve-central.com/unifieduploader/start
//

package emds

import (
	log "github.com/Sirupsen/logrus"
	"github.com/buger/jsonparser"
)

// Order stores information about a single order on the market including info from the rowset.
type Order struct {
	OrderID       int64   `json:"orderID"`
	RegionID      int64   `json:"regionID"`
	TypeID        int64   `json:"typeID"`
	GeneratedAt   string  `json:"generatedAt"`
	Price         float64 `json:"price"`
	VolRemaining  int64   `json:"volRemaining"`
	OrderRange    int64   `json:"range"`
	VolEntered    int64   `json:"volEntered"`
	MinVolume     int64   `json:"minVolume"`
	Bid           bool    `json:"bid"`
	IssueDate     string  `json:"issueDate"`
	Duration      int64   `json:"duration"`
	StationID     int64   `json:"stationID"`
	SolarSystemID int64   `json:"solarSystemID"`
}

// Rowset contains all orders for a given region/type combination at generatedAt (market snapshot).
type Rowset struct {
	GeneratedAt string  `json:"generatedAt"`
	RegionID    int64   `json:"regionID"`
	TypeID      int64   `json:"typeID"`
	Rows        []Order `json:"orders"`
}

// RawRowset contains unparsed orders for a given region/type combination at generatedAt (market snapshot). Can be used for deduplication of orders without parsing.
type RawRowset struct {
	GeneratedAt string `json:"generatedAt"`
	RegionID    int64  `json:"regionID"`
	TypeID      int64  `json:"typeID"`
	Rows        []byte `json:"-"`
}

// ColumnIndices keeps index values of the rowsets from UUDIF.
type ColumnIndices struct {
	price, volRemaining, orderRange, orderID, volEntered, minVolume, bid, issueDate, duration, stationID, solarSystemID int
}

// ParseUUDIF message into structs. You probably want to use this function if all you want is to get all data.
func ParseUUDIF(message []byte) ([]Rowset, error) {
	var rowsets []Rowset
	var rawRowsets []RawRowset

	indices, err := GetColumnIndices(message)
	if err != nil {
		log.Errorf("Error extracting columns: %s", err.Error())
		return rowsets, err
	}

	rawRowsets, err = ExtractRawRowsets(message)
	if err != nil {
		log.Errorf("Error extracting raw rowsets: %s", err.Error())
		return rowsets, err
	}

	rowsets, err = ParseRawRowsets(rawRowsets, indices)
	if err != nil {
		log.Errorf("Error parsing rowsets: %s", err.Error())
		return rowsets, err
	}

	return rowsets, nil
}

// GetColumnIndices gets column indices for mapping keys correctly.
func GetColumnIndices(message []byte) (ColumnIndices, error) {

	var indices ColumnIndices
	var i int

	_, err := jsonparser.ArrayEach(message, func(column []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			return
		}

		switch string(column) {
		case "price":
			indices.price = i
		case "volRemaining":
			indices.volRemaining = i
		case "orderRange":
			indices.orderRange = i
		case "orderID":
			indices.orderID = i
		case "volEntered":
			indices.volEntered = i
		case "minVolume":
			indices.minVolume = i
		case "bid":
			indices.bid = i
		case "issueDate":
			indices.issueDate = i
		case "duration":
			indices.duration = i
		case "stationID":
			indices.stationID = i
		case "solarSystemID":
			indices.solarSystemID = i
		}

		i++
	}, "columns")

	if err != nil {
		return indices, err
	}

	return indices, nil
}

// ExtractRawRowsets extracts raw rowsets (orders are not being parsed).
func ExtractRawRowsets(message []byte) ([]RawRowset, error) {
	var rowsets []RawRowset

	_, err := jsonparser.ArrayEach(message, func(rowset []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			log.Warnf("Error splitting rowsets: %s", err.Error())
			return
		}

		regionID, err := jsonparser.GetInt(rowset, "regionID")
		if err != nil {
			log.Warnf("Error parsing regionID: %s", err.Error())
			return
		}

		typeID, err := jsonparser.GetInt(rowset, "typeID")
		if err != nil {
			log.Warnf("Error parsing typeID: %s", err.Error())
			return
		}

		generatedAt, err := jsonparser.GetString(rowset, "generatedAt")
		if err != nil {
			log.Warnf("Error parsing generatedAt: %s", err.Error())
			return
		}

		rawOrders, _, _, err := jsonparser.Get(rowset, "rows")
		if err != nil {
			log.Warnf("Error extracting rawOrders: %s", err.Error())
			return
		}

		rowsets = append(rowsets, RawRowset{
			RegionID:    regionID,
			TypeID:      typeID,
			GeneratedAt: generatedAt,
			Rows:        rawOrders,
		})
	}, "rowsets")

	return rowsets, err
}

// ParseRawRowsets parses orders in raw rowset.
func ParseRawRowsets(rawRowsets []RawRowset, indices ColumnIndices) ([]Rowset, error) {
	var parsedRowsets []Rowset

	for _, rowset := range rawRowsets {
		orders, err := ParseOrders(rowset.Rows, indices, rowset.RegionID, rowset.TypeID, rowset.GeneratedAt)

		if err != nil {
			log.Warnf("Error parsing orders: %s", err.Error())
			return nil, err
		}

		parsedRowsets = append(parsedRowsets, Rowset{
			RegionID:    rowset.RegionID,
			TypeID:      rowset.TypeID,
			GeneratedAt: rowset.GeneratedAt,
			Rows:        orders,
		})
	}

	return parsedRowsets, nil
}

// ParseOrders parses orders from rows.
func ParseOrders(rows []byte, indices ColumnIndices, regionID int64, typeID int64, generatedAt string) ([]Order, error) {

	var orders []Order

	_, err := jsonparser.ArrayEach(rows, func(row []byte, dataType jsonparser.ValueType, offset int, err error) {
		var columnIndex int
		order := Order{RegionID: regionID, TypeID: typeID, GeneratedAt: generatedAt}

		jsonparser.ArrayEach(row, func(column []byte, dataType jsonparser.ValueType, offset int, err error) {
			switch columnIndex {
			case indices.price:
				price, err := jsonparser.GetFloat(column)
				if err != nil {
					log.Warnf("Unable to parse price: %s", err.Error())
				}
				order.Price = price
			case indices.volRemaining:
				volRemaining, err := jsonparser.GetInt(column)
				if err != nil {
					log.Warnf("Unable to parse volRemaining: %s", err.Error())
				}
				order.VolRemaining = volRemaining
			case indices.orderRange:
				orderRange, err := jsonparser.GetInt(column)
				if err != nil {
					log.Warnf("Unable to parse orderRange: %s", err.Error())
				}
				order.OrderRange = orderRange
			case indices.orderID:
				orderID, err := jsonparser.GetInt(column)
				if err != nil {
					log.Warnf("Unable to parse orderID: %s", err.Error())
				}
				order.OrderID = orderID
			case indices.volEntered:
				volEntered, err := jsonparser.GetInt(column)
				if err != nil {
					log.Warnf("Unable to parse volEntered: %s", err.Error())
				}
				order.VolEntered = volEntered
			case indices.minVolume:
				minVolume, err := jsonparser.GetInt(column)
				if err != nil {
					log.Warnf("Unable to parse minVolume: %s", err.Error())
				}
				order.MinVolume = minVolume
			case indices.bid:
				bid, err := jsonparser.GetBoolean(column)
				if err != nil {
					log.Warnf("Unable to parse bid: %s", err.Error())
				}
				order.Bid = bid
			case indices.issueDate:
				order.IssueDate = string(column)
			case indices.duration:
				duration, err := jsonparser.GetInt(column)
				if err != nil {
					log.Warnf("Unable to parse duration: %s", err.Error())
				}
				order.Duration = duration
			case indices.stationID:
				stationID, err := jsonparser.GetInt(column)
				if err != nil {
					log.Warnf("Unable to parse stationID: %s", err.Error())
				}
				order.StationID = stationID
			case indices.solarSystemID:
				solarSystemID, err := jsonparser.GetInt(column)
				if err != nil {
					// log.Warnf("Unable to parse solarSystemID: %s", err.Error())
				}
				order.SolarSystemID = solarSystemID
			}
			columnIndex++
		})

		columnIndex = 0
		orders = append(orders, order)
	})

	if err != nil {
		return orders, err
	}

	return orders, nil
}

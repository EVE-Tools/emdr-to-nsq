//
// EVE Market Data Structures in go
// Provides functions for dealing with UUDIF-formatted data.
// See: http://dev.eve-central.com/unifieduploader/start
//

package emds

import (
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/buger/jsonparser"
	"github.com/mailru/easyjson/jwriter"
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
		logrus.WithError(err).Error("Error extracting columns.")
		return rowsets, err
	}

	rawRowsets, err = ExtractRawRowsets(message)
	if err != nil {
		logrus.WithError(err).Error("Error extracting raw rowsets.")
		return rowsets, err
	}

	rowsets, err = ParseRawRowsets(rawRowsets, indices)
	if err != nil {
		logrus.WithError(err).Error("Error parsing rowsets.")
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
			logrus.WithError(err).Warn("Error splitting rowsets.")
			return
		}

		regionID, err := jsonparser.GetInt(rowset, "regionID")
		if err != nil {
			logrus.WithError(err).Warn("Error parsing regionID.")
			return
		}

		typeID, err := jsonparser.GetInt(rowset, "typeID")
		if err != nil {
			logrus.WithError(err).Warn("Error parsing typeID.")
			return
		}

		generatedAt, err := jsonparser.GetString(rowset, "generatedAt")
		if err != nil {
			logrus.WithError(err).Warn("Error parsing generatedAt.")
			return
		}

		rawOrders, _, _, err := jsonparser.Get(rowset, "rows")
		if err != nil {
			logrus.WithError(err).Warn("Error extracting rawOrders.")
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
			logrus.WithError(err).Warn("Error parsing orders.")
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
					logrus.WithError(err).Warn("Unable to parse price.")
				}
				order.Price = price
			case indices.volRemaining:
				volRemaining, err := jsonparser.GetInt(column)
				if err != nil {
					logrus.WithError(err).Warn("Unable to parse volRemaining.")
				}
				order.VolRemaining = volRemaining
			case indices.orderRange:
				orderRange, err := jsonparser.GetInt(column)
				if err != nil {
					logrus.WithError(err).Warn("Unable to parse orderRange.")
				}
				order.OrderRange = orderRange
			case indices.orderID:
				orderID, err := jsonparser.GetInt(column)
				if err != nil {
					logrus.WithError(err).Warn("Unable to parse orderID.")
				}
				order.OrderID = orderID
			case indices.volEntered:
				volEntered, err := jsonparser.GetInt(column)
				if err != nil {
					logrus.WithError(err).Warn("Unable to parse volEntered.")
				}
				order.VolEntered = volEntered
			case indices.minVolume:
				minVolume, err := jsonparser.GetInt(column)
				if err != nil {
					logrus.WithError(err).Warn("Unable to parse minVolume.")
				}
				order.MinVolume = minVolume
			case indices.bid:
				bid, err := jsonparser.GetBoolean(column)
				if err != nil {
					logrus.WithError(err).Warn("Unable to parse bid.")
				}
				order.Bid = bid
			case indices.issueDate:
				order.IssueDate = string(column)
			case indices.duration:
				duration, err := jsonparser.GetInt(column)
				if err != nil {
					logrus.WithError(err).Warn("Unable to parse duration.")
				}
				order.Duration = duration
			case indices.stationID:
				stationID, err := jsonparser.GetInt(column)
				if err != nil {
					logrus.WithError(err).Warn("Unable to parse stationID.")
				}
				order.StationID = stationID
			case indices.solarSystemID:
				solarSystemID, err := jsonparser.GetInt(column)
				if err != nil {
					//.WithError(err) logrusWarnf("Unable to parse solarSystemID.")
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

// ConvertRange converts a string range returned by the API into an int for use in UUDIF
func ConvertRange(rangeString string) (int64, error) {
	switch rangeString {
	case "station":
		return -1, nil
	case "solarsystem":
		return 0, nil
	case "region":
		return 32767, nil
	}

	intRange, err := strconv.ParseInt(rangeString, 10, 64)
	if err != nil {
		return 0, err
	}

	return intRange, nil
}

// RowsetsToUUDIF converts a list of rowsets to a UUDIF message
func RowsetsToUUDIF(rowsets []Rowset, generatorName string, generatorVersion string) ([]byte, error) {
	now := time.Now().Format(time.RFC3339)
	out := jwriter.Writer{}

	out.RawString("{\"resultType\":\"orders\",\"version\":\"0.1\",\"uploadKeys\":[{\"key\":\"e43\",\"name\":\"EMDR\"}],\"generator\":{\"name\":")
	out.String(generatorName)
	out.RawString(",\"version\":")
	out.String(generatorVersion)
	out.RawString("},\"columns\":[\"price\",\"volRemaining\",\"range\",\"orderID\",\"volEntered\",\"minVolume\",\"bid\",\"issueDate\",\"duration\",\"stationID\",\"solarSystemID\"],\"currentTime\":")
	out.String(now)
	out.RawString(",\"rowsets\":[")

	for rowsetIndex, rowset := range rowsets {

		out.RawString("{\"generatedAt\":")
		out.String(rowset.GeneratedAt)
		out.RawString(",\"regionID\":")
		out.Int64(rowset.RegionID)
		out.RawString(",\"typeID\":")
		out.Int64(rowset.TypeID)
		out.RawString(",\"rows\":[")

		for rowIndex, row := range rowset.Rows {
			out.RawByte('[')
			out.Float32(float32(row.Price))
			out.RawByte(',')
			out.Int64(row.VolRemaining)
			out.RawByte(',')
			out.Int64(row.OrderRange)
			out.RawByte(',')
			out.Int64(row.OrderID)
			out.RawByte(',')
			out.Int64(row.VolEntered)
			out.RawByte(',')
			out.Int64(row.MinVolume)
			out.RawByte(',')
			out.Bool(row.Bid)
			out.RawByte(',')
			out.String(row.IssueDate)
			out.RawByte(',')
			out.Int64(row.Duration)
			out.RawByte(',')
			out.Int64(row.StationID)
			out.RawByte(',')
			out.Int64(row.SolarSystemID)
			out.RawByte(']')

			if rowIndex != (len(rowset.Rows) - 1) {
				out.RawByte(',')
			}
		}

		out.RawString("]}")

		if rowsetIndex != (len(rowsets) - 1) {
			out.RawByte(',')
		}
	}

	out.RawString("]}")

	return out.Buffer.BuildBytes(), out.Error
}

// ByOrderID sorts orders by order ID
type ByOrderID []Order

func (o ByOrderID) Len() int           { return len(o) }
func (o ByOrderID) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o ByOrderID) Less(i, j int) bool { return o[i].OrderID < o[j].OrderID }

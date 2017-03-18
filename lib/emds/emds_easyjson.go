// AUTOGENERATED FILE: easyjson marshaller/unmarshallers.

package emds

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson7c6c15edDecodeGithubComEVEToolsEmdrToNsqLibEmds(in *jlexer.Lexer, out *ColumnIndices) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson7c6c15edEncodeGithubComEVEToolsEmdrToNsqLibEmds(out *jwriter.Writer, in ColumnIndices) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ColumnIndices) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7c6c15edEncodeGithubComEVEToolsEmdrToNsqLibEmds(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ColumnIndices) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7c6c15edEncodeGithubComEVEToolsEmdrToNsqLibEmds(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ColumnIndices) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7c6c15edDecodeGithubComEVEToolsEmdrToNsqLibEmds(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ColumnIndices) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7c6c15edDecodeGithubComEVEToolsEmdrToNsqLibEmds(l, v)
}
func easyjson7c6c15edDecodeGithubComEVEToolsEmdrToNsqLibEmds1(in *jlexer.Lexer, out *RawRowset) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "generatedAt":
			out.GeneratedAt = string(in.String())
		case "regionID":
			out.RegionID = int64(in.Int64())
		case "typeID":
			out.TypeID = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson7c6c15edEncodeGithubComEVEToolsEmdrToNsqLibEmds1(out *jwriter.Writer, in RawRowset) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"generatedAt\":")
	out.String(string(in.GeneratedAt))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"regionID\":")
	out.Int64(int64(in.RegionID))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"typeID\":")
	out.Int64(int64(in.TypeID))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RawRowset) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7c6c15edEncodeGithubComEVEToolsEmdrToNsqLibEmds1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RawRowset) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7c6c15edEncodeGithubComEVEToolsEmdrToNsqLibEmds1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RawRowset) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7c6c15edDecodeGithubComEVEToolsEmdrToNsqLibEmds1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RawRowset) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7c6c15edDecodeGithubComEVEToolsEmdrToNsqLibEmds1(l, v)
}
func easyjson7c6c15edDecodeGithubComEVEToolsEmdrToNsqLibEmds2(in *jlexer.Lexer, out *Rowset) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "generatedAt":
			out.GeneratedAt = string(in.String())
		case "regionID":
			out.RegionID = int64(in.Int64())
		case "typeID":
			out.TypeID = int64(in.Int64())
		case "orders":
			if in.IsNull() {
				in.Skip()
				out.Rows = nil
			} else {
				in.Delim('[')
				if !in.IsDelim(']') {
					out.Rows = make([]Order, 0, 1)
				} else {
					out.Rows = []Order{}
				}
				for !in.IsDelim(']') {
					var v1 Order
					(v1).UnmarshalEasyJSON(in)
					out.Rows = append(out.Rows, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson7c6c15edEncodeGithubComEVEToolsEmdrToNsqLibEmds2(out *jwriter.Writer, in Rowset) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"generatedAt\":")
	out.String(string(in.GeneratedAt))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"regionID\":")
	out.Int64(int64(in.RegionID))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"typeID\":")
	out.Int64(int64(in.TypeID))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"orders\":")
	if in.Rows == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in.Rows {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Rowset) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7c6c15edEncodeGithubComEVEToolsEmdrToNsqLibEmds2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Rowset) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7c6c15edEncodeGithubComEVEToolsEmdrToNsqLibEmds2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Rowset) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7c6c15edDecodeGithubComEVEToolsEmdrToNsqLibEmds2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Rowset) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7c6c15edDecodeGithubComEVEToolsEmdrToNsqLibEmds2(l, v)
}
func easyjson7c6c15edDecodeGithubComEVEToolsEmdrToNsqLibEmds3(in *jlexer.Lexer, out *Order) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "orderID":
			out.OrderID = int64(in.Int64())
		case "regionID":
			out.RegionID = int64(in.Int64())
		case "typeID":
			out.TypeID = int64(in.Int64())
		case "generatedAt":
			out.GeneratedAt = string(in.String())
		case "price":
			out.Price = float64(in.Float64())
		case "volRemaining":
			out.VolRemaining = int64(in.Int64())
		case "range":
			out.OrderRange = int64(in.Int64())
		case "volEntered":
			out.VolEntered = int64(in.Int64())
		case "minVolume":
			out.MinVolume = int64(in.Int64())
		case "bid":
			out.Bid = bool(in.Bool())
		case "issueDate":
			out.IssueDate = string(in.String())
		case "duration":
			out.Duration = int64(in.Int64())
		case "stationID":
			out.StationID = int64(in.Int64())
		case "solarSystemID":
			out.SolarSystemID = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson7c6c15edEncodeGithubComEVEToolsEmdrToNsqLibEmds3(out *jwriter.Writer, in Order) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"orderID\":")
	out.Int64(int64(in.OrderID))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"regionID\":")
	out.Int64(int64(in.RegionID))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"typeID\":")
	out.Int64(int64(in.TypeID))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"generatedAt\":")
	out.String(string(in.GeneratedAt))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"price\":")
	out.Float64(float64(in.Price))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"volRemaining\":")
	out.Int64(int64(in.VolRemaining))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"range\":")
	out.Int64(int64(in.OrderRange))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"volEntered\":")
	out.Int64(int64(in.VolEntered))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"minVolume\":")
	out.Int64(int64(in.MinVolume))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"bid\":")
	out.Bool(bool(in.Bid))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"issueDate\":")
	out.String(string(in.IssueDate))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"duration\":")
	out.Int64(int64(in.Duration))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"stationID\":")
	out.Int64(int64(in.StationID))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"solarSystemID\":")
	out.Int64(int64(in.SolarSystemID))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Order) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7c6c15edEncodeGithubComEVEToolsEmdrToNsqLibEmds3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Order) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7c6c15edEncodeGithubComEVEToolsEmdrToNsqLibEmds3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Order) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7c6c15edDecodeGithubComEVEToolsEmdrToNsqLibEmds3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Order) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7c6c15edDecodeGithubComEVEToolsEmdrToNsqLibEmds3(l, v)
}

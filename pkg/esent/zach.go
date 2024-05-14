package esent

import (
	"encoding/binary"

	"github.com/charmbracelet/log"
)

// Get all of the columns for a record
func (e *Esent_record) GetColumns() []string {
	keys := make([]string, len(e.column))

	i := 0

	for k := range e.column {
		keys[i] = k
		i++
	}

	return keys
}

// Get the type of a column
func (e *Esent_record) GetColumnType(column string) recordTyp {
	return e.GetRecord(column).GetType()
}

func (e esent_recordVal) ValueAsUint16() uint16 {
	return uint16(binary.LittleEndian.Uint16(e.val))
}

func (e esent_recordVal) ValueAsUint32() uint32 {
	return uint32(binary.LittleEndian.Uint32(e.val))
}

func (e esent_recordVal) ValueAsUint64() uint64 {
	return uint64(binary.LittleEndian.Uint64(e.val))
}

func (e esent_recordVal) ValueAsInt16() int16 {
	return int16(binary.LittleEndian.Uint16(e.val))
}

func (e esent_recordVal) ValueAsInt32() int32 {
	return int32(binary.LittleEndian.Uint32(e.val))
}

func (e esent_recordVal) ValueAsInt64() int64 {
	return int64(binary.LittleEndian.Uint64(e.val))
}

func (e esent_recordVal) ValueAsFloat32() float32 {
	return Float32frombytes(e.val)
}

func (e esent_recordVal) ValueAsFloat64() float64 {
	return Float64frombytes(e.val)
}

func (e *Esent_record) GetShortVal(column string) (int16, bool) {
	v, ok := e.column[column]
	if v != nil && ok {
		return v.ValueAsInt16(), ok
	}
	return 0, ok
}

func (e *Esent_record) GetCurrencyVal(column string) (uint64, bool) {
	v, ok := e.column[column]
	if v != nil && ok {
		return v.ValueAsUint64(), ok
	}
	return 0, ok
}

func (e *Esent_record) GetIEEESinglVal(column string) (float32, bool) {
	v, ok := e.column[column]
	if v != nil && ok {
		return v.ValueAsFloat32(), ok
	}
	return 0, ok
}

func (e *Esent_record) GetIEEEDoublVal(column string) (float64, bool) {
	v, ok := e.column[column]
	if v != nil && ok {
		return v.ValueAsFloat64(), ok
	}
	return 0, ok
}

func (e *Esent_record) GetDateTimeVal(column string) (uint64, bool) {
	v, ok := e.column[column]
	if v != nil && ok {
		return v.ValueAsUint64(), ok
	}
	return 0, ok
}

func (e *Esent_record) GetUnsLngVal(column string) (uint32, bool) {
	v, ok := e.column[column]
	if v != nil && ok {
		return v.ValueAsUint32(), ok
	}
	return 0, ok
}

func (e *Esent_record) GetLngLngVal(column string) (uint64, bool) {
	v, ok := e.column[column]
	if v != nil && ok {
		return v.ValueAsUint64(), ok
	}
	return 0, ok
}

func (e *Esent_record) GetGuidVal(column string) ([]byte, bool) {
	v, ok := e.column[column]
	if v != nil && ok {
		return v.Bytes(), ok
	}
	return nil, ok
}

func (e *Esent_record) GetUnsShrtVal(column string) (uint16, bool) {
	v, ok := e.column[column]
	if v != nil && ok {
		return v.ValueAsUint16(), ok
	}
	return 0, ok
}

// Try to convert a column's value to a type other than bytes
//
// Reference from elsewhere in the code:
//
// bit        bool
// unsByt     byte
// short      int16
// long       int32
// curr       uint64
// iEEESingl  float32
// iEEEDoubl  float64
// dateTim    uint64
// unsLng     uint32
// lngLng     uint64
// guid       [16]byte
// unsShrt    uint16
// nils for binary, text, longbin, longtext and slv?
func (e *Esent_record) ConvertValue(column string) interface{} {
	switch typ := e.GetColumnType(column); typ {
	case Byt:
		// log.Infof("Type of column %s: %d (Byt)", column, typ)

		value, ok := e.GetBytVal(column)
		if !ok {
			log.Fatalf("Failed to GetBytVal for %s", column)
		}

		// log.Infof("Value of column %s: %v", column, value)
		return value

	// TODO: May not want bytes
	// NOTE: Haven't encountered yet
	case Tup:
		log.Infof("Type of column %s: %d (Tup)", column, typ)

		value, ok := e.GetBytVal(column)
		if !ok {
			log.Fatalf("Failed to GetBytVal for %s", column)
		}

		log.Infof("Value of column %s: %v", column, value)
		return value

	case Str:
		// log.Infof("Type of column %s: %d (Str)", column, typ)

		value, err := e.StrVal(column)
		if err != nil {
			panic(err)
		}

		// log.Infof("Value of column %s: %v", column, value)
		return value

	case Nil:
		// log.Infof("Type of column %s: %d (Nil)", column, typ)

	// TODO: May not want bytes
	// NOTE: Haven't encountered yet
	case Bit:
		log.Infof("Type of column %s: %d (Bit)", column, typ)

		value, ok := e.GetBytVal(column)
		if !ok {
			log.Fatalf("Failed to GetBytVal for %s", column)
		}

		log.Infof("Value of column %s: %v", column, value)
		return value

	// TODO: Difference w/ Byt?
	case UnsByt:
		// log.Infof("Type of column %s: %d (UnsByt)", column, typ)

		value, ok := e.GetBytVal(column)
		if !ok {
			log.Fatalf("Failed to GetBytVal for %s", column)
		}

		// log.Infof("Value of column %s: %v", column, value)
		return value

	case Short:
		// log.Infof("Type of column %s: %d (Short)", column, typ)

		value, ok := e.GetShortVal(column)
		if !ok {
			log.Fatalf("Failed to GetShortVal for %s", column)
		}

		// log.Infof("Value of column %s: %v", column, value)
		return value

	case Long:
		// log.Infof("Type of column %s: %d (Long)", column, typ)

		value, ok := e.GetLongVal(column)
		if !ok {
			log.Fatalf("Failed to GetLongVal for %s", column)
		}

		// log.Infof("Value of column %s: %v", column, value)
		return value

	case Curr:
		// log.Infof("Type of column %s: %d (Curr)", column, typ)

		value, ok := e.GetCurrencyVal(column)
		if !ok {
			log.Fatalf("Failed to GetCurrencyVal for %s", column)
		}

		// log.Infof("Value of column %s: %v", column, value)
		return value

	case IEEESingl:
		// log.Infof("Type of column %s: %d (IEEESingl)", column, typ)

		value, ok := e.GetIEEESinglVal(column)
		if !ok {
			log.Fatalf("Failed to GetIEEESinglVal for %s", column)
		}

		// log.Infof("Value of column %s: %v", column, value)
		return value

	case IEEEDoub:
		// log.Infof("Type of column %s: %d (IEEEDoub)", column, typ)

		value, ok := e.GetIEEEDoublVal(column)
		if !ok {
			log.Fatalf("Failed to GetIEEEDoublVal for %s", column)
		}

		// log.Infof("Value of column %s: %v", column, value)
		return value

	// NOTE: Haven't encountered yet
	case DateTim:
		log.Infof("Type of column %s: %d (Datetim)", column, typ)

		value, ok := e.GetDateTimeVal(column)
		if !ok {
			log.Fatalf("Failed to GetDateTimeVal for %s", column)
		}

		// log.Infof("Value of column %s: %v", column, value)
		return value

	// TODO: May not want bytes
	// NOTE: Haven't encountered yet
	case Bin:
		log.Infof("Type of column %s: %d (Bin)", column, typ)

		value, ok := e.GetBytVal(column)
		if !ok {
			log.Fatalf("Failed to GetBytVal for %s", column)
		}

		log.Infof("Value of column %s: %v", column, value)
		return value

	// TODO: May not want bytes (especially here!)
	// NOTE: Haven't encountered yet
	case Txt:
		log.Infof("Type of column %s: %d (Txt)", column, typ)

		value, ok := e.GetBytVal(column)
		if !ok {
			log.Fatalf("Failed to GetBytVal for %s", column)
		}

		log.Infof("Value of column %s: %v", column, value)
		return value

	// TODO: Does this need to be decoded at all? Has stuff like objectSid and objectGUID
	case LongBin:
		// log.Infof("Type of column %s: %d (LongBin)", column, typ)

		value, ok := e.GetBytVal(column)
		if !ok {
			log.Fatalf("Failed to GetBytVal for %s", column)
		}

		// log.Infof("Value of column %s: %v", column, value)
		return value

	// TODO: May not want bytes (especially here!)
	// NOTE: Haven't encountered yet
	case LongTxt:
		log.Infof("Type of column %s: %d (LongTxt)", column, typ)

		value, ok := e.GetBytVal(column)
		if !ok {
			log.Fatalf("Failed to GetBytVal for %s", column)
		}

		log.Infof("Value of column %s: %v", column, value)
		return value

	// TODO: May not want bytes
	// NOTE: Haven't encountered yet
	case SLV:
		log.Infof("Type of column %s: %d (SLV)", column, typ)

		value, ok := e.GetBytVal(column)
		if !ok {
			log.Fatalf("Failed to GetBytVal for %s", column)
		}

		log.Infof("Value of column %s: %v", column, value)
		return value

	case UnsLng:
		// log.Infof("Type of column %s: %d (UnsLng)", column, typ)

		value, ok := e.GetUnsLngVal(column)
		if !ok {
			log.Fatalf("Failed to GetUnsLngVal for %s", column)
		}

		// log.Infof("Value of column %s: %v", column, value)
		return value

	case LngLng:
		// log.Infof("Type of column %s: %d (LngLng)", column, typ)

		value, ok := e.GetLngLngVal(column)
		if !ok {
			log.Fatalf("Failed to GetLngLngVal for %s", column)
		}

		// log.Infof("Value of column %s: %v", column, value)
		return value

	// NOTE: Haven't encountered yet
	case Guid:
		log.Infof("Type of column %s: %d (Guid)", column, typ)

		value, ok := e.GetGuidVal(column)
		if !ok {
			log.Fatalf("Failed to GetGuidVal for %s", column)
		}

		// log.Infof("Value of column %s: %v", column, value)
		return value

	case UnsShrt:
		// log.Infof("Type of column %s: %d (UnsShrt)", column, typ)

		value, ok := e.GetUnsShrtVal(column)
		if !ok {
			log.Fatalf("Failed to GetUnsShrtVal for %s", column)
		}

		// log.Infof("Value of column %s: %v", column, value)
		return value

	// TODO: May not want bytes
	// NOTE: Haven't encountered yet
	case Max:
		log.Infof("Type of column %s: %d (Max)", column, typ)

		value, ok := e.GetBytVal(column)
		if !ok {
			log.Fatalf("Failed to GetBytVal for %s", column)
		}

		// log.Infof("Value of column %s: %v", column, value)
		return value

	default:
		log.Fatalf("Unknown type for column %s: %d", column, typ)
		return nil // for static analysis
	}

	log.Fatalf("Failed to enter switch statement for column %s...?", column)
	return nil // for static analysis
}

// Main function for converting DB values to something usable
func (e *Esent_record) ZachsRecordParse() map[string]interface{} {

	log := GetLogger()

	ditDump := make(map[string]interface{})

	for _, column := range e.GetColumns() {
		// oMSyntax (type 7: Long)
		// if column == "ATTj131303" {
		// 	log.Fatal(e.GetColumnType(column))
		// }

		// objectGUID (type 14: LongBin)
		// if column == "ATTk589826" {
		// 	log.Fatal(e.GetColumnType(column))
		// }

		// pwdLastSet (type 8: Curr)
		// if column == "ATTq589920" {
		// 	log.Fatal(e.GetColumnType(column))
		// }

		value := e.ConvertValue(column)
		log.Infof("Type of column %s: %d (Byt)", column, e.GetColumnType(column))
		log.Infof("Value of column %s: %v", column, value)
		ditDump[column] = value
	}

	// log.Infof("Full record:%v", ditDump)

	// jsonString, err := json.Marshal(ditDump)
	// if err != nil {
	// 	panic(err)
	// }

	// log.Infof("JSON:%s", jsonString)

	return ditDump
}

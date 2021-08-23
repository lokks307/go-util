package djson

import (
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	V_TYPE_NULL int = iota
	V_TYPE_INT
	V_TYPE_FLOAT
	V_TYPE_NUMBER
	V_TYPE_STRING
	V_TYPE_OBJECT
	V_TYPE_ARRAY
	V_TYPE_MULTI
)

var CountryCodes = []string{
	"GH", "GA", "GY", "GM", "GG", "GP", "GT", "GU", "GD", "GR", "GL", "GW", "GN", "NA", "NR", "NG", "AQ", "SS", "ZA", "AN", "NL",
	"NP", "NO", "NF", "NC", "NZ", "NU", "NE", "NI", "KR", "DK", "DO", "DM", "DE", "TL", "LA", "LR", "LV", "RU", "LB", "LS", "RE",
	"RO", "LU", "RW", "LY", "LT", "LI", "MG", "MQ", "MH", "YT", "MO", "MW", "MY", "ML", "IM", "MX", "MC", "MA", "MU", "MR", "MZ",
	"ME", "MS", "MD", "MV", "MT", "MN", "UM", "VI", "US", "MM", "FM", "VU", "BH", "BB", "VA", "BS", "BD", "BM", "BJ", "VE", "VN",
	"BE", "BY", "BZ", "BA", "BW", "BO", "BI", "BF", "BV", "BT", "MP", "MK", "BG", "BR", "BN", "WS", "SA", "GS", "SM", "ST", "PM",
	"EH", "SN", "RS", "SC", "LC", "VC", "KN", "SH", "SO", "SB", "SD", "SR", "LK", "SJ", "SE", "CH", "ES", "SK", "SI", "SY", "SL",
	"SX", "SG", "AE", "AW", "AM", "AR", "AS", "IS", "HT", "IE", "AZ", "AF", "AD", "AL", "DZ", "AO", "AG", "AI", "ER", "SZ", "EE",
	"EC", "ET", "SV", "VG", "IO", "GB", "YE", "OM", "AU", "AT", "HN", "AX", "WF", "JO", "UG", "UY", "UZ", "UA", "IQ", "IR", "IL",
	"EG", "IT", "ID", "IN", "JP", "JM", "ZM", "JE", "GQ", "KP", "GE", "CN", "CF", "DJ", "GI", "ZW", "TD", "CZ", "CL", "CM", "CV",
	"KZ", "QA", "KH", "CA", "KE", "KY", "KM", "CR", "CC", "CI", "CO", "CG", "CD", "CU", "KW", "CK", "HR", "CX", "KG", "KI", "CY",
	"TW", "TJ", "TZ", "TH", "TC", "TR", "TG", "TK", "TO", "TM", "TV", "TN", "TT", "PA", "PY", "PK", "PG", "PW", "PS", "FO", "PE",
	"PT", "FK", "PL", "PR", "GF", "TF", "PF", "FR", "FJ", "FI", "PH", "PN", "HM", "HU", "HK",
}

var HexRegExp *regexp.Regexp
var TimestampRegExp *regexp.Regexp
var YYYYMMDDRegExp *regexp.Regexp
var YYMMDDRegExp *regexp.Regexp
var HHMMSSRegExp *regexp.Regexp
var HHMMRegExp *regexp.Regexp
var EmailRegExp *regexp.Regexp
var UUIDRegExp *regexp.Regexp
var TelRegExp *regexp.Regexp

func CheckFuncHex(ts string) bool {
	return HexRegExp.Match([]byte(ts))
}

func CheckFuncTimestamp(ts string) bool {
	return TimestampRegExp.Match([]byte(ts))
}

func CheckFuncYYYYMMDD(ts string) bool {
	return YYYYMMDDRegExp.Match([]byte(ts))
}

func CheckFuncYYMMDD(ts string) bool {
	return YYMMDDRegExp.Match([]byte(ts))
}

func CheckFuncHHMMSS(ts string) bool {
	return HHMMSSRegExp.Match([]byte(ts))
}

func CheckFuncHHMM(ts string) bool {
	return HHMMRegExp.Match([]byte(ts))
}

func CheckFuncEmail(ts string) bool {
	return EmailRegExp.Match([]byte(ts))
}

func CheckFuncIntString(ts string) bool {
	_, err := strconv.Atoi(ts)
	return err == nil
}

func CheckFuncFloatString(ts string) bool {
	_, err := strconv.ParseFloat(ts, 64)
	return err == nil
}

func CheckFuncUUID(ts string) bool {
	return UUIDRegExp.Match([]byte(ts))
}

func CheckISO31661A2(val string) bool {
	if len(val) != 2 {
		return false
	}

	valUpper := strings.ToUpper(val)

	for idx := range CountryCodes {
		if valUpper == CountryCodes[idx] {
			return true
		}
	}

	return false
}

func CheckBase64(ts string) bool {
	_, err := base64.StdEncoding.DecodeString(ts)
	return err == nil
}

func CheckTelephone(ts string) bool {
	return TelRegExp.Match([]byte(ts))
}

// ISO 3166-2 : KR-XX, GH-XX, ...
func CheckISO31662(val string) bool {
	if len(val) < 4 {
		return false
	}

	if val[2:3] != "-" {
		return false
	}

	return CheckISO31661A2(val[0:2])
}

func init() {
	HexRegExp = regexp.MustCompile(`^([A-Fa-f0-9]{2})*$`)
	TimestampRegExp = regexp.MustCompile(`^[0-9]{9,11}$`)
	YYYYMMDDRegExp = regexp.MustCompile(`^[1-2][0-9]{3}-{0,1}(0[1-9]|1[0-2])-{0,1}(0[1-9]|[1-2][0-9]|3[0-1])$`)
	YYMMDDRegExp = regexp.MustCompile(`^[0-9]{2}-{0,1}(0[1-9]|1[0-2])-{0,1}(0[1-9]|[1-2][0-9]|3[0-1])$`)
	HHMMSSRegExp = regexp.MustCompile(`^([0-1][0-9]|2[0-3])\:{0,1}([0-5][0-9])\:{0,1}([0-5][0-9])$`)
	HHMMRegExp = regexp.MustCompile(`^([0-1][0-9]|2[0-3])\:{0,1}([0-5][0-9])$`)
	EmailRegExp = regexp.MustCompile(`^([\w\.\_\-])*[a-zA-Z0-9]+([\w\.\_\-])*([a-zA-Z0-9])+([\w\.\_\-])+@([a-zA-Z0-9]+\.)+[a-zA-Z0-9]{2,8}$`)
	UUIDRegExp = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89AB][0-9a-f]{3}-[0-9a-f]{12}$`)
	TelRegExp = regexp.MustCompile(`^(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$`)
}

type VItem struct {
	Type      int
	Name      string
	Max       int64
	Min       int64
	MaxFloat  float64
	MinFloat  float64
	IsRequred bool
	SubItems  []*VItem
	CheckFunc func(string) bool
	RegExp    *regexp.Regexp
}

type Validator struct {
	Syntax    *DJSON
	RootItems []*VItem
}

func NewValidator() *Validator {
	return &Validator{
		Syntax: NewDJSON(),
	}
}

func (m *Validator) Compile(syntax string) bool {
	m.Syntax.Parse(syntax)

	if !m.Syntax.IsObject() && !m.Syntax.IsString() && !m.Syntax.IsArray() {
		m.Syntax = NewDJSON()
		return false
	}

	m.RootItems = make([]*VItem, 0)

	if m.Syntax.IsObject() || m.Syntax.IsString() {
		vItem := GetVItem("", m.Syntax)
		if vItem != nil {
			m.RootItems = append(m.RootItems, vItem)
		}

	} else if m.Syntax.IsArray() {
		m.Syntax.Seek()
		for es := m.Syntax.Next(); es != nil; es = m.Syntax.Next() {
			vi := GetVItem("", es)
			if vi != nil {
				m.RootItems = append(m.RootItems, vi)
			}
		}
	}

	return true
}

func GetVItem(name string, ejson *DJSON) *VItem {
	eitem := new(VItem)
	eitem.Name = name
	etype := ""

	if ejson.IsString() {
		etype = ejson.GetAsString()

		switch etype {
		case "INT":
			eitem.Type = V_TYPE_INT
			eitem.Min = -9007199254740991
			eitem.Max = 9007199254740991
		case "FLOAT":
			eitem.Type = V_TYPE_FLOAT
			eitem.MinFloat = -1.7976931348623157e+308
			eitem.MaxFloat = 1.7976931348623157e+308
		case "NUMBER":
			eitem.Type = V_TYPE_NUMBER
			eitem.MinFloat = -1.7976931348623157e+308
			eitem.MaxFloat = 1.7976931348623157e+308
		case "STRING":
			eitem.Type = V_TYPE_STRING
			eitem.Min = 0
			eitem.Max = 8192
		case "OBJECT":
			eitem.Type = V_TYPE_OBJECT
		case "ARRAY":
			eitem.Type = V_TYPE_ARRAY
			eitem.Min = 0
			eitem.Max = 9007199254740991
		case "NONEMPTY.STRING":
			eitem.Type = V_TYPE_STRING
			eitem.Min = 1
			eitem.Max = 8192
		case "NONEMPTY.ARRAY":
			eitem.Type = V_TYPE_ARRAY
			eitem.Min = 1
			eitem.Max = 9007199254740991
		case "HEX":
			eitem.Type = V_TYPE_STRING
			eitem.Min = 0
			eitem.Max = 8192
			eitem.CheckFunc = CheckFuncHex
		}

	} else if ejson.IsArray() {

		eitem.Type = V_TYPE_MULTI
		ejson.Seek()
		for es := ejson.Next(); es != nil; es = ejson.Next() {
			vi := GetVItem(name, es)
			if vi != nil {
				eitem.SubItems = append(eitem.SubItems, vi)
			}
		}

	} else if ejson.IsObject() {

		etype = ejson.GetAsString("type")
		eitem.IsRequred = ejson.GetAsBool("required")
		if ejson.GetAsString("regexp") != "" {
			eitem.RegExp, _ = regexp.Compile(ejson.GetAsString("regexp"))
		}

		switch etype {
		case "INT":
			eitem.Type = V_TYPE_INT
			eitem.Min = ejson.GetAsInt("min", -9007199254740991)
			eitem.Max = ejson.GetAsInt("max", 9007199254740991)
		case "FLOAT":
			eitem.Type = V_TYPE_FLOAT
			eitem.MinFloat = ejson.GetAsFloat("min", -1.7976931348623157e+308)
			eitem.MaxFloat = ejson.GetAsFloat("max", 1.7976931348623157e+308)
		case "NUMBER":
			eitem.Type = V_TYPE_NUMBER
			eitem.MinFloat = ejson.GetAsFloat("min", -1.7976931348623157e+308)
			eitem.MaxFloat = ejson.GetAsFloat("max", 1.7976931348623157e+308)
		case "STRING":
			eitem.Type = V_TYPE_STRING
			eitem.Min = ejson.GetAsInt("min", 0)
			eitem.Max = ejson.GetAsInt("max", 8192)
		case "OBJECT":
			subJson, ok := ejson.GetAsObject("object")
			if ok {
				eitem.Type = V_TYPE_OBJECT
				eitem.IsRequred = ejson.GetAsBool("required")

				ks := subJson.GetKeys()
				for _, ek := range ks {
					ejson, ok := subJson.Get(ek)
					if ok {
						vItem := GetVItem(ek, ejson)
						if vItem != nil {
							eitem.SubItems = append(eitem.SubItems, vItem)
						}
					}

				}
			}
		case "NONEMPTY.STRING":
			eitem.Type = V_TYPE_STRING
			eitem.Min = ejson.GetAsInt("min", 1)
			if eitem.Min < 1 {
				eitem.Min = 1
			}
			eitem.Max = ejson.GetAsInt("max", 8192)
		case "ARRAY":
			eitem.Type = V_TYPE_ARRAY
			eitem.Min = ejson.GetAsInt("min", 0)
		case "NONEMPTY.ARRAY":
			eitem.Type = V_TYPE_ARRAY
			eitem.Min = ejson.GetAsInt("min", 1)
			if eitem.Min < 1 {
				eitem.Min = 1
			}
		case "HEX":
			eitem.Type = V_TYPE_STRING
			eitem.Min = ejson.GetAsInt("min", 0)
			eitem.Max = ejson.GetAsInt("max", 8192)
			eitem.CheckFunc = CheckFuncHex
		}

		if etype == "ARRAY" || etype == "NONEMPTY.ARRAY" {
			eitem.Max = ejson.GetAsInt("max", 9007199254740991)
			oa, ok := ejson.Get("array") // type of element
			if ok {
				eitem.SubItems = make([]*VItem, 0)
				if oa.IsArray() {
					oa.Seek()
					for es := oa.Next(); es != nil; es = oa.Next() {
						vi := GetVItem("", es)
						if vi != nil {
							eitem.SubItems = append(eitem.SubItems, vi)
						}
					}
				} else if oa.IsString() || oa.IsObject() {
					vi := GetVItem("", oa)
					eitem.SubItems = append(eitem.SubItems, vi)
				}
			}
		}

	}

	switch etype {
	case "TIMESTAMP":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 0
		eitem.Max = 10
		eitem.CheckFunc = CheckFuncTimestamp
	case "YYYYMMDD":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 8
		eitem.Max = 10
		eitem.CheckFunc = CheckFuncYYYYMMDD
	case "YYMMDD":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 6
		eitem.Max = 8
		eitem.CheckFunc = CheckFuncYYMMDD
	case "HHMMSS":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 6
		eitem.Max = 8
		eitem.CheckFunc = CheckFuncHHMMSS
	case "HHMM":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 4
		eitem.Max = 5
		eitem.CheckFunc = CheckFuncHHMM
	case "EMAIL":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 3
		eitem.Max = 255
		eitem.CheckFunc = CheckFuncEmail
	case "INT_STRING":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 1
		eitem.Max = 17
		eitem.CheckFunc = CheckFuncIntString
	case "FLOAT_STRING":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 1
		eitem.Max = 24
		eitem.CheckFunc = CheckFuncFloatString
	case "UUID":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 36
		eitem.Max = 36
		eitem.CheckFunc = CheckFuncUUID
	case "ISO31661A2":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 2
		eitem.Max = 2
		eitem.CheckFunc = CheckISO31661A2
	case "ISO31662":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 5
		eitem.Max = 5
		eitem.CheckFunc = CheckISO31662
	case "BASE64":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 0
		eitem.Max = 8192
		eitem.CheckFunc = CheckBase64
	case "TELEPHONE":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 4
		eitem.Max = 20
		eitem.CheckFunc = CheckTelephone
	}

	return eitem
}

func (m *Validator) IsValid(tjson *DJSON) bool {
	if tjson == nil {
		return len(m.RootItems) == 0
	}

	if m.Syntax.IsObject() { // json must be valid one
		log.Println("---------------------------------------------")
		for _, vitem := range m.RootItems {
			fmt.Println("---------------------------------------------")
			fmt.Println(vitem.Type, vitem.Name)
			fmt.Println(tjson.ToString())
			fmt.Println("---------------------------------------------")
			return CheckVItem(vitem, tjson)
		}

	} else if m.Syntax.IsArray() || m.Syntax.IsString() {
		// each element must be valid for one of vitems

		if tjson.IsArray() {

			if len(m.RootItems) == 0 {
				return true
			}

			tjson.Seek()
			for et := tjson.Next(); et != nil; et = tjson.Next() {
				eachValid := false
				for _, vitem := range m.RootItems {
					if CheckVItem(vitem, et) {
						eachValid = true
						break
					}
				}

				if !eachValid {
					return false
				}
			}

			return true
		}
		fmt.Println("m.Syntax.IsString")

		for _, vitem := range m.RootItems {
			fmt.Println("---------------------------------------------")
			fmt.Println(vitem.Type, vitem.Name)
			fmt.Println(tjson.ToString())
			fmt.Println("---------------------------------------------")
			if CheckVItem(vitem, tjson) {
				return true
			}
		}

		return false
	}

	return true

}

func CheckVItem(vi *VItem, tjson *DJSON) bool {

	fmt.Println("=============================== B")
	fmt.Println(vi.Type, vi.Name)
	fmt.Println(tjson.ToString())
	fmt.Println("=============================== E")

	if vi.IsRequred && vi.Name != "" && !tjson.HasKey(vi.Name) {
		return false
	}

	vtype := tjson.GetType(vi.Name)

	switch vi.Type {
	case V_TYPE_INT:
		if vtype != "int" {
			return false
		}

		si := tjson.GetAsInt(vi.Name)
		if vi.Max < si || vi.Min > si {
			return false
		}

	case V_TYPE_NUMBER:
		if vtype != "float" && vtype != "int" {
			return false
		}

		sf := tjson.GetAsFloat(vi.Name)
		if vi.MaxFloat < sf || vi.MinFloat > sf {
			return false
		}
	case V_TYPE_FLOAT:
		if vtype != "float" {
			return false
		}

		sf := tjson.GetAsFloat(vi.Name)
		if vi.MaxFloat < sf || vi.MinFloat > sf {
			return false
		}
	case V_TYPE_STRING:
		if vtype != "string" {
			return false
		}

		ss := tjson.GetAsString(vi.Name)
		lenv := int64(len(ss))

		if lenv > vi.Max || lenv < vi.Min {
			return false
		}

		if vi.RegExp != nil {
			return vi.RegExp.Match([]byte(ss))
		}

		if vi.CheckFunc != nil {
			return vi.CheckFunc(ss)
		}

	case V_TYPE_OBJECT:
		so, ok := tjson.GetAsObject(vi.Name)
		if vi.IsRequred && (!ok || so == nil || !so.IsObject()) {
			return false
		}

		if vi.Name != "" && so == nil {
			return true
		}

		if so == nil {
			return false
		}

		for _, svi := range vi.SubItems {
			if !CheckVItem(svi, so) {
				return false
			}
		}

	case V_TYPE_ARRAY:
		sa, ok := tjson.GetAsArray(vi.Name)
		if vi.IsRequred && (!ok || sa == nil || !sa.IsArray()) {
			return false
		}

		if vi.Name != "" && sa == nil {
			return true
		}

		if sa == nil {
			return false
		}

		lenv := int64(sa.Length())
		if lenv > vi.Max || lenv < vi.Min {
			return false
		}

		if ok {
			if len(vi.SubItems) == 0 {
				return true
			}

			sa.Seek()
			for ssa := sa.Next(); ssa != nil; ssa = sa.Next() {
				isValid := false
				for _, svi := range vi.SubItems {
					if CheckVItem(svi, ssa) {
						isValid = true
						break
					}
				}
				if !isValid {
					return false
				}
			}
		}

	case V_TYPE_MULTI:

		isValid := false

		for _, svi := range vi.SubItems {
			if CheckVItem(svi, tjson) {
				isValid = true
				break
			}
		}

		if !isValid {
			return false
		}

	}

	return true
}

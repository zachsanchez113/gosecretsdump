package ditreader

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"unicode"

	"github.com/C-Sto/gosecretsdump/pkg/esent"
	"github.com/C-Sto/gosecretsdump/pkg/systemreader"
)

type M map[string]interface{}

// TODO: Map column names to human-readable
func (d DitReader) DumpJSON() error {
	//if local (always local for now)
	if d.systemHiveLocation != "" {
		ls, err := systemreader.New(d.systemHiveLocation)
		if err != nil {
			return err
		}
		d.bootKey = ls.BootKey()
		if d.ntdsFileLocation != "" {
			d.noLMHash = ls.HasNoLMHashPolicy()
		}
	} else {
		return fmt.Errorf("System hive empty")
	}

	d.getPek()
	if len(d.pek) < 1 {
		return fmt.Errorf("NO PEK FOUND THIS IS VERY BAD")
	}

	var records []M

	for {
		//read each record from the db
		record, err := d.db.GetNextRow(d.cursor)
		if err != nil {
			if err.Error() == "ignore" {
				break //we will get an 'ignore' error when there are no more records
			}
			fmt.Println("Couldn't get row due to error: ", err.Error())
			continue
		}

		samAccountName, err := record.StrVal(nsAMAccountName)
		validUsername := true
		// no value if not nil
		if err == nil {
			for _, char := range samAccountName {
				if !unicode.IsPrint(char) {
					validUsername = false
					break
				}
			}
		}

		if validUsername {
			// RecordToJSON?
			parsedRecord := record.ZachsRecordParse()

			if lmHash, err := d.GetLMHash(record); err == nil {
				parsedRecord["lmHash"] = lmHash
			}

			if ntlmHash, err := d.GetNTLMHash(record); err == nil {
				parsedRecord["ntlmHash"] = ntlmHash
			}

			// Convert bytes to string
			for k, v := range parsedRecord {
				if v2, ok := v.([]byte); ok {
					parsedRecord[k] = hex.EncodeToString(v2)
				}
			}

			records = append(records, parsedRecord)
		}
	}

	fmt.Fprintf(os.Stderr, "Number of records: %d\n", len(records))

	var records2 []M
	for _, record := range records {
		_, ok := record["ntlmHash"]
		if ok {
			records2 = append(records2, record)
		}
	}

	fmt.Fprintf(os.Stderr, "Number of user records: %d\n", len(records2))

	jsonString, err := json.Marshal(records2)
	// jsonString, err := json.Marshal(records2[:100])
	if err != nil {
		panic(err)
	}

	d.userData <- DumpedHash{JsonString: string(jsonString)}

	close(d.userData)
	return nil
}

func (d *DitReader) RecordToJSON(record esent.Esent_record) (map[string]interface{}, error) {
	dh := DumpedHash{}
	v, _ := record.GetBytVal(nobjectSid)
	sid, err := NewSAMRRPCSID(v) //record.Column[z].BytVal)
	if err != nil {
		return nil, err
	}
	//dh.Rid = sid.FormatCanonical()[strings.LastIndex(sid.FormatCanonical(), "-")+1:]
	dh.Rid = sid.Rid()

	//lm hash
	if b, err := record.GetBytVal(ndBCSPwd); err && len(b) > 0 {
		//if record.Column[ndBCSPwd"]].StrVal != "" {
		var tmpLM []byte
		encryptedLM, err := NewCryptedHash(b)
		if err != nil {
			return nil, err
		}
		if bytes.Equal(encryptedLM.Header[:4], []byte("\x13\x00\x00\x00")) {
			encryptedLMW := NewCryptedHashW16(b)
			pekIndex := encryptedLMW.Header
			tmpLM, err = DecryptAES(d.pek[pekIndex[4]], encryptedLMW.EncryptedHash[:16], encryptedLMW.KeyMaterial[:])
			if err != nil {
				return nil, err
			}
		} else {
			tmpLM, err = d.removeRC4(encryptedLM)
			if err != nil {
				return nil, err
			}
		}
		dh.LMHash, err = RemoveDES(tmpLM, dh.Rid)
		if err != nil {
			return nil, err
		}
	} else {
		//hard coded empty lm hash
		dh.LMHash = EmptyLM //, _ = hex.DecodeString("aad3b435b51404eeaad3b435b51404ee")
	}

	//nt hash
	if v, _ := record.GetBytVal(nunicodePwd); len(v) > 0 { //  record.Column[nunicodePwd"]].BytVal; len(v) > 0 {
		var tmpNT []byte
		encryptedNT, err := NewCryptedHash(v)
		if err != nil {
			return nil, err
		}
		if bytes.Equal(encryptedNT.Header[:4], []byte("\x13\x00\x00\x00")) {
			encryptedNTW := NewCryptedHashW16(v)
			pekIndex := encryptedNTW.Header
			tmpNT, err = DecryptAES(d.pek[pekIndex[4]], encryptedNTW.EncryptedHash[:16], encryptedNTW.KeyMaterial[:])
			if err != nil {
				return nil, err
			}
		} else {
			tmpNT, err = d.removeRC4(encryptedNT)
			if err != nil {
				return nil, err
			}
		}
		dh.NTHash, err = RemoveDES(tmpNT, dh.Rid)
		if err != nil {
			return nil, err
		}
	} else {
		//hard coded empty NTLM hash
		dh.NTHash = EmptyNT //, _ = hex.DecodeString("31D6CFE0D16AE931B73C59D7E0C089C0")
	}

	ditDump := record.ZachsRecordParse()
	ditDump["lmHash"] = hex.EncodeToString(dh.LMHash)
	ditDump["ntlmHash"] = hex.EncodeToString(dh.NTHash)

	return ditDump, nil
}

func (d *DitReader) GetLMHash(record esent.Esent_record) (string, error) {
	dh := DumpedHash{}
	v, _ := record.GetBytVal(nobjectSid)
	sid, err := NewSAMRRPCSID(v) //record.Column[z].BytVal)
	if err != nil {
		return "", err
	}
	//dh.Rid = sid.FormatCanonical()[strings.LastIndex(sid.FormatCanonical(), "-")+1:]
	dh.Rid = sid.Rid()

	//lm hash
	if b, err := record.GetBytVal(ndBCSPwd); err && len(b) > 0 {
		//if record.Column[ndBCSPwd"]].StrVal != "" {
		var tmpLM []byte
		encryptedLM, err := NewCryptedHash(b)
		if err != nil {
			return "", err
		}
		if bytes.Equal(encryptedLM.Header[:4], []byte("\x13\x00\x00\x00")) {
			encryptedLMW := NewCryptedHashW16(b)
			pekIndex := encryptedLMW.Header
			tmpLM, err = DecryptAES(d.pek[pekIndex[4]], encryptedLMW.EncryptedHash[:16], encryptedLMW.KeyMaterial[:])
			if err != nil {
				return "", err
			}
		} else {
			tmpLM, err = d.removeRC4(encryptedLM)
			if err != nil {
				return "", err
			}
		}
		dh.LMHash, err = RemoveDES(tmpLM, dh.Rid)
		if err != nil {
			return "", err
		}
	} else {
		//hard coded empty lm hash
		dh.LMHash = EmptyLM //, _ = hex.DecodeString("aad3b435b51404eeaad3b435b51404ee")
	}

	return hex.EncodeToString(dh.LMHash), nil
}

func (d *DitReader) GetNTLMHash(record esent.Esent_record) (string, error) {
	dh := DumpedHash{}
	v, _ := record.GetBytVal(nobjectSid)
	sid, err := NewSAMRRPCSID(v) //record.Column[z].BytVal)
	if err != nil {
		return "", err
	}
	//dh.Rid = sid.FormatCanonical()[strings.LastIndex(sid.FormatCanonical(), "-")+1:]
	dh.Rid = sid.Rid()

	//nt hash
	if v, _ := record.GetBytVal(nunicodePwd); len(v) > 0 { //  record.Column[nunicodePwd"]].BytVal; len(v) > 0 {
		var tmpNT []byte
		encryptedNT, err := NewCryptedHash(v)
		if err != nil {
			return "", err
		}
		if bytes.Equal(encryptedNT.Header[:4], []byte("\x13\x00\x00\x00")) {
			encryptedNTW := NewCryptedHashW16(v)
			pekIndex := encryptedNTW.Header
			tmpNT, err = DecryptAES(d.pek[pekIndex[4]], encryptedNTW.EncryptedHash[:16], encryptedNTW.KeyMaterial[:])
			if err != nil {
				return "", err
			}
		} else {
			tmpNT, err = d.removeRC4(encryptedNT)
			if err != nil {
				return "", err
			}
		}
		dh.NTHash, err = RemoveDES(tmpNT, dh.Rid)
		if err != nil {
			return "", err
		}
	} else {
		//hard coded empty NTLM hash
		dh.NTHash = EmptyNT //, _ = hex.DecodeString("31D6CFE0D16AE931B73C59D7E0C089C0")
	}

	return hex.EncodeToString(dh.NTHash), nil
}

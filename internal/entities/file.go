package entities

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type InputEvent struct {
}

type FileStatus string
type FileType string

const (
	FileStatusInProgress FileStatus = "inProgress"
	FileStatusReady      FileStatus = "ready"
	FileStatusError      FileStatus = "error"
)

const (
	FileCSV  FileType = "csv"
	FileXLSX FileType = "xlsx"
)

type File struct {
	ID          uuid.UUID  `json:"id"`
	Status      FileStatus `json:"status"`
	Type        FileType   `json:"type"`
	Name        string     `json:"name"`
	Data        []byte     `json:"data"`
	Array       []any      `json:"array"`
	ReflectType string     `json:"reflectType"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type FileResponse struct {
	ID        uuid.UUID  `json:"id"`
	Status    FileStatus `json:"status"`
	Type      FileType   `json:"type"`
	Name      string     `json:"name"`
	Error     string     `json:"error"`
	CreatedAt time.Time  `json:"createdAt"`
}

func (s *File) Response() FileResponse {
	res := FileResponse{
		ID:        s.ID,
		Status:    s.Status,
		Type:      s.Type,
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
	}

	if res.Status == FileStatusError {
		res.Error = string(s.Data)
	}

	return res
}

func (s *File) MarshalBinary() (data []byte, err error) {
	return json.Marshal(s)
}

func (s *File) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	switch s.ReflectType {
	case reflect.TypeOf(&FinancialReport{}).Name():
		data, _ = json.Marshal(s.Array)
		s.Array = []any{}
		realData := []*FinancialReport{}
		err = json.Unmarshal(data, &realData)
		if err != nil {
			return err
		}

		for _, elem := range realData {
			s.Array = append(s.Array, elem)
		}
	case reflect.TypeOf(&Spin{}).Name():
		data, _ = json.Marshal(s.Array)
		s.Array = []any{}
		realData := []*Spin{}
		err = json.Unmarshal(data, &realData)
		if err != nil {
			return err
		}

		for _, elem := range realData {
			s.Array = append(s.Array, elem)
		}
	case reflect.TypeOf(&GamingSession{}).Name():
		data, _ = json.Marshal(s.Array)
		s.Array = []any{}
		realData := []*GamingSession{}
		err = json.Unmarshal(data, &realData)
		if err != nil {
			return err
		}

		for _, elem := range realData {
			s.Array = append(s.Array, elem)
		}
	case reflect.TypeOf(&AggregatedReportByGame{}).Name():
		data, _ = json.Marshal(s.Array)
		s.Array = []any{}
		realData := []*GamingSession{}
		err = json.Unmarshal(data, &realData)
		if err != nil {
			return err
		}

		for _, elem := range realData {
			s.Array = append(s.Array, elem)
		}
	case reflect.TypeOf(&AggregatedReportByCountry{}).Name():
		data, _ = json.Marshal(s.Array)
		s.Array = []any{}
		realData := []*GamingSession{}
		err = json.Unmarshal(data, &realData)
		if err != nil {
			return err
		}

		for _, elem := range realData {
			s.Array = append(s.Array, elem)
		}
	}

	return nil
}

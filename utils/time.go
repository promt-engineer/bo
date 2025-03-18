package utils

import (
	"backoffice/internal/constants"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func ParseTimestampPB(ts string) (*timestamppb.Timestamp, error) {
	var res *timestamppb.Timestamp

	if ts != "" {
		t, err := time.Parse(constants.TimeLayout, ts)
		if err != nil {
			return nil, err
		}

		res = timestamppb.New(t)
	}

	return res, nil
}

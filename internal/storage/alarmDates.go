package storage

func CreateAlarmDates( req *als.CreateAlarmDatesRequest) (*als.CreateAlarmDatesResponse, error) {
	db := s.DB()

	var returnDates []*als.AlarmDateTime

	if len(req.ReqFilter) > 0 {
		for _, date := range req.ReqFilter {
			var returnID int64

			err := db.QueryRowx(`insert into 
			alarm_date_time(alarm_id, alarm_day, start_time, end_time) values ($1, $2, $3, $4) returning id`,
				date.AlarmId, date.AlarmDay, date.AlarmStartTime, date.AlarmEndTime).Scan(&returnID)

			if err != nil {
				return &als.CreateAlarmDatesResponse{RespDateTime: returnDates}, s.HandlePSQLError(s.Insert, err, "insert error")
			}
			createdDate := als.AlarmDateTime{
				Id:             returnID,
				AlarmId:        date.AlarmId,
				AlarmDay:       date.AlarmDay,
				AlarmStartTime: date.AlarmStartTime,
				AlarmEndTime:   date.AlarmEndTime,
			}
			returnDates = append(returnDates, &createdDate)
		}
	}
	return &als.CreateAlarmDatesResponse{RespDateTime: returnDates}, nil
}
package backup

import "time"

type Metadata struct {
	ID		    string
	Name	    string
	StartTime   time.Time
	EndTime		time.Time
	Duration    float64
	Size	    int64
	Status	    string
	Error    string
}
package api

import (
	"time"

	"github.com/nvellon/hal"
)

type Report struct {
	Id            int64     `db:"id,omitempty"`
	Origin        string    `db:"origin"`
	Method        string    `db:"method"`
	Status        int       `db:"status"`
	ContentType   string    `db:"content_type"`
	ContentLength uint      `db:"content_length"`
	Host          string    `db:"host"`
	URL           string    `db:"url"`
	Scheme        string    `db:"scheme"`
	Path          string    `db:"path"`
	Body          string    `db:"body"`
	RequestBody   string    `db:"request_body"`
	DateStart     time.Time `db:"date_start"`
	DateEnd       time.Time `db:"date_end"`
	TimeTaken     time.Time `db:"time_taken"`
	Ts            time.Time `db:"ts,omitempty"`
	// Header        Header    `db:"header"`
	// RequestHeader Header    `db:"request_header"`
}

func (c Report) GetMap() hal.Entry {
	return hal.Entry{
		"id":             c.Id,
		"origin":         c.Origin,
		"method":         c.Method,
		"status":         c.Status,
		"content_type":   c.ContentType,
		"content_length": c.ContentLength,
		"host":           c.Host,
		"url":            c.URL,
		"scheme":         c.Scheme,
		"path":           c.Path,
		"body":           c.Body,
		"request_body":   c.RequestBody,
		"date_start":     c.DateStart,
		"date_end":       c.DateEnd,
		"time_taken":     c.TimeTaken,
		"created":        c.Ts,
	}
}

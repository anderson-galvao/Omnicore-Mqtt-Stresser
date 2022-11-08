package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=postgres dbname=iot_core_demo password=iotcore22!@bx#m host=10.56.48.2 sslmode=disable"

	subscription := "epsi-production"
	device := "Stresser"
	registry := "Stresser"
	entries := 100

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		panic(err)
	}

	sqlStatement := `INSERT INTO public."Device" (
		devicename, subscription, registry, parent, numid, fullname, credentials, loglevel, blocked, metadata, createdon, updatedon, capresent, isgateway, gateway, heartbeat, lasttelemetryreceived, laststatereceived, lastconfigsent, clientonline, subscriptions) VALUES (
		$3::character varying, $1::character varying, $6::character varying, $2::character varying, $5::bigint, $4::character varying, '[{"expirationTime":"2023-10-02T15:01:23.045123456Z","id":"1","publicKey":{"format":"RSA_X509_PEM","key":"-----BEGIN CERTIFICATE-----\nMIIDujCCAqKgAwIBAgITaOszYpBme+SRHZUkFWLgDs7EMDANBgkqhkiG9w0BAQsF\nADAeMQ0wCwYDVQQKEwRrb3JlMQ0wCwYDVQQDEwRrb3JlMB4XDTIyMDgwNDEwMzQy\nMloXDTMyMDgwMTEwMzIwOFowADCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC\nggEBAMFvTBHPdgH7+5wVlUnEdIS/0a4p9fkVzMdEMdDVr5s62VoGO7nZWxMCaxxU\nXqQiGuX3N7SINyD7h8LI8CxQsn5zyDda3QVNGU7I96iWjwzOYJmNHAN1nRI2hRDY\n8fJoQgTZI+IiRDBmgkmL9yjTY04qY7UP8zpofuMnKRuTwP6Ey1eFEMBqFfvgwrVl\niLNcq9At0bd/vlQ0VUnKV6oKqSTq9ZDPB6Cxu5amhejVwTeE6p5GGmiKw5vskmtB\ndGNgsom1K/pJdOMes8lODVp00tIVnsplL3jLgrWfbfCPALRnGz/C5XlKW8fNKEuW\nqFw2Lhnk51dtobw/oBo7vJcx2w0CAwEAAaOCAQ0wggEJMA4GA1UdDwEB/wQEAwIF\noDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBQUiRS8X3OWpJNpjQoJ+22xGb3xXTAf\nBgNVHSMEGDAWgBTKO7S10CConGVgZZli7NVAcim/AzCBjQYIKwYBBQUHAQEEgYAw\nfjB8BggrBgEFBQcwAoZwaHR0cDovL3ByaXZhdGVjYS1jb250ZW50LTYyZTM5YmRh\nLTAwMDAtMjI0My1hNjFhLTNjMjg2ZDRlZWUwYS5zdG9yYWdlLmdvb2dsZWFwaXMu\nY29tL2MwZGQxZjg3ZDcwZGZhMDEwNGEwL2NhLmNydDAZBgNVHREBAf8EDzANggtn\nYWRnZW9uLmNvbTANBgkqhkiG9w0BAQsFAAOCAQEAK82b/xGn8B6Nfogw0myKjy3O\nWg53YPXuct3E04qRmD3JJtFzpSkjI2WyRmIkRX1b5SKF+ImOmGzvENZDkjT/Y2I/\nnsBL639OlXnz/+GYSq4rL6fVxXistP4LGA+khoBYSfHFZb7EYoVOYJFzZjnvJbtz\n7XG0jTMeHo8KhCBPxrNWkOERrcc7OWqREldQ36yg7zdbRLjDOjeD6FByTrpRhbDC\n0AeozF9ug9W/gPYtnkI++ksUqjJcV06uGd+9XLJPGcjH0Bai1alxROh+dkWx6TcB\nHC94el4KR6EJijMvylmnOyKHedmYaDvb52+B6zXTW9rQkh1UycmONlmAlA3OeQ==\n-----END CERTIFICATE-----\n"}}]'::json, 'INFO'::character varying, false::boolean, '{}'::json, '2022-10-28 17:17:14.286019+05:30'::timestamp with time zone, '2022-11-05 14:44:46.464838+05:30'::timestamp with time zone, true::boolean, false::boolean, '{}'::character varying[], '2022-10-30 12:36:59.616122+05:30'::timestamp with time zone, '2022-10-30 12:36:48.543632+05:30'::timestamp with time zone, '2022-10-30 12:36:59.616122+05:30'::timestamp with time zone, '2022-10-30 12:36:59.616122+05:30'::timestamp with time zone, true::boolean, '{}'::character varying[])
		ON CONFLICT DO NOTHING;`

	for i := 0; i < entries; i++ {
		d3 := fmt.Sprintf("subscriptions/%s/registries/%s/devices", subscription, registry)
		devId := fmt.Sprintf("%s%d", device, i)
		d5 := fmt.Sprintf("%s/%s", d3, devId)
		_, err = db.Exec(sqlStatement, subscription, d3, devId, d5, i+500, registry)
		if err != nil {
			panic(err)
		}
	}

}

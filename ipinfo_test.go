package ipinfo

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"errors"
	"net/http"
	"testing"
	"time"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const TEST_PORT = "56127"

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

type IPInfoSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&IPInfoSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *IPInfoSuite) SetUpSuite(c *C) {
	API = "http://127.0.0.1:" + TEST_PORT

	mux := http.NewServeMux()
	server := &http.Server{Addr: ":" + TEST_PORT, Handler: mux}

	mux.HandleFunc("GET /1.1.1.1", handlerGetBasic)
	mux.HandleFunc("GET /1.1.1.2", handlerGetError)
	mux.HandleFunc("GET /1.1.1.3", handlerGetInvalidResponse)
	mux.HandleFunc("GET /1.1.1.4", handlerGetInvalidError)

	go server.ListenAndServe()

	time.Sleep(time.Second)
}

func (s *IPInfoSuite) TestBasicErrors(c *C) {
	_, err := Get("1.1.1.1")
	c.Assert(err, Equals, ErrEmptyToken)

	Token = "ABCD"

	_, err = Get("")
	c.Assert(err, Equals, ErrEmptyIP)

	_, err = Get("1.1.1.A")
	c.Assert(err, DeepEquals, errors.New(`IP "1.1.1.A" is not valid`))

	_, err = Get("1234.1234")
	c.Assert(err, DeepEquals, errors.New(`IP "1234.1234" is not valid`))

	API = "http://127.0.0.1:12346"

	_, err = Get("1.1.1.1")
	c.Assert(err, NotNil)

	API = "http://127.0.0.1:" + TEST_PORT
}

func (s *IPInfoSuite) TestResponses(c *C) {
	info, err := Get("1.1.1.1")
	c.Assert(info, NotNil)
	c.Assert(err, IsNil)
	c.Assert(info.IP, Equals, "1.1.1.1")
	c.Assert(info.ASN, Equals, "AS13335")
	c.Assert(info.ASName, Equals, "Cloudflare, Inc.")
	c.Assert(info.ASDomain, Equals, "cloudflare.com")
	c.Assert(info.Country, Equals, "Australia")
	c.Assert(info.CountryCode, Equals, "AU")
	c.Assert(info.Continent, Equals, "Oceania")
	c.Assert(info.ContinentCode, Equals, "OC")

	info, err = Get("1.1.1.2")
	c.Assert(info, IsNil)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, `IPinfo API returned error: Please provide a valid IP address`)

	info, err = Get("1.1.1.3")
	c.Assert(info, IsNil)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, `Can't decode IPinfo API response: invalid character 'X' looking for beginning of object key string`)

	info, err = Get("1.1.1.4")
	c.Assert(info, IsNil)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, `IPinfo API returned non-ok status code (400)`)
}

// ////////////////////////////////////////////////////////////////////////////////// //

func handlerGetBasic(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
	rw.Write([]byte(`{
  "ip": "1.1.1.1",
  "asn": "AS13335",
  "as_name": "Cloudflare, Inc.",
  "as_domain": "cloudflare.com",
  "country_code": "AU",
  "country": "Australia",
  "continent_code": "OC",
  "continent": "Oceania"
}`))
}

func handlerGetError(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(400)
	rw.Write([]byte(`{
  "error": "Please provide a valid IP address"
}`))
}

func handlerGetInvalidResponse(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
	rw.Write([]byte(`{X`))
}

func handlerGetInvalidError(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(400)
	rw.Write([]byte(`{X`))
}

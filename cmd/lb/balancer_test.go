package main

import (
	"net/http"
	"slices"
	"testing"

	. "gopkg.in/check.v1"
)

func TestBalancer(t *testing.T) {
	// TODO: Реалізуйте юніт-тест для балансувальникка.
	// Загалом, в Test_getServer імплементований загальний юніт-тест для балансувальника.
	// Я навіть зловив на ньому помилку, коли сервер брався не з списку здорових серверів, а із загального списку
	TestingT(t)
}

type BalancerSuite struct{}

var _ = Suite(&BalancerSuite{})

func (s *BalancerSuite) Test_getHealthyServers(c *C) {
	serversStates = map[string]bool{}
	tests := []struct {
		name              string
		serversStatesTest map[string]bool
		wantResult        []string
	}{
		{
			"no servers",
			map[string]bool{},
			[]string{},
		},
		{
			"no available servers",
			map[string]bool{
				"server1": false,
			},
			[]string{},
		},
		{
			"single available server",
			map[string]bool{
				"server1": true,
			},
			[]string{"server1"},
		},
		{
			"1 available server, 1 unavailable",
			map[string]bool{
				"server1": true,
				"server2": false,
			},
			[]string{"server1"},
		},
		{
			"2 unavailable servers",
			map[string]bool{
				"server1": false,
				"server2": false,
			},
			[]string{},
		},
	}
	for _, tt := range tests {
		serversStates = tt.serversStatesTest
		gotResult := getHealthyServers()
		if len(gotResult) == len(tt.wantResult) && len(tt.wantResult) == 0 {
			continue
		}
		c.Assert(gotResult, DeepEquals, tt.wantResult)
	}
}

func (s *BalancerSuite) Test_getServer(c *C) {
	serversPool = []string{
		"server1:8080",
		"server2:8080",
		"server3:8080",
	}

	tests := []struct {
		name          string
		r             *http.Request
		serversStates map[string]bool
		wantedServers []string
		wantErr       bool
	}{
		{
			"test no available servers",
			&http.Request{RemoteAddr: "127.0.0.1:6051"},
			map[string]bool{
				"server1:8080": false,
				"server2:8080": false,
				"server3:8080": false,
			},
			[]string{},
			true,
		},
		{
			"test one available server",
			&http.Request{RemoteAddr: "127.0.0.1:6053"},
			map[string]bool{
				"server1:8080": false,
				"server2:8080": true,
				"server3:8080": false,
			},
			[]string{"server2:8080"},
			false,
		},
		{
			"test two available servers",
			&http.Request{RemoteAddr: "127.0.0.1:6054"},
			map[string]bool{
				"server1:8080": false,
				"server2:8080": true,
				"server3:8080": true,
			},
			[]string{"server2:8080", "server3:8080"},
			false,
		},
	}
	for _, tt := range tests {
		serversStates = tt.serversStates
		gotServer, err := getServer(tt.r)

		if tt.wantErr {
			c.Assert(err, NotNil)
		} else {
			c.Assert(err, IsNil)
			c.Assert(slices.Contains(tt.wantedServers, gotServer), Equals, true)
		}
	}
}

func (s *BalancerSuite) Test_getHash(c *C) {
	randomAddresses := []string{
		"220.216.32.147:22166",
		"136.59.94.139:52910",
		"76.241.221.109:58354",
		"231.131.206.235:57372",
		"153.97.93.1:23218",
		"197.29.20.129:38109",
		"120.65.92.149:21862",
		"183.158.32.207:65448",
		"107.216.47.215:43163",
		"24.59.32.47:62307",
		"133.65.232.138:50517",
		"210.230.55.249:28408",
		"238.51.66.183:24121",
		"19.191.55.215:33057",
		"148.90.88.192:1262",
		"187.7.137.127:43996",
		"240.176.63.40:16514",
		"226.230.4.40:65444",
		"87.76.11.182:42912",
		"9.105.58.169:12826",
	}

	duplicatedAddresses := []string{
		"220.216.32.147:22166",
		"220.216.32.147:22166",
		"220.216.32.147:22166",
		"220.216.32.147:22166",
		"220.216.32.147:22166",
	}
	resultHashes := []uint64{}
	for _, address := range randomAddresses {
		hash := getHash(address)
		c.Assert(slices.Contains(resultHashes, hash), Equals, false)
		resultHashes = append(resultHashes, hash)
	}
	dupHashes := []uint64{getHash(duplicatedAddresses[0])}
	for _, address := range duplicatedAddresses {
		hash := getHash(address)
		c.Assert(slices.Contains(dupHashes, hash), Equals, true)
	}
}

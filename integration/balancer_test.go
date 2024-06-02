package integration

import (
	"net/http"
	"slices"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

const baseAddress = "http://balancer:8090"

var client = http.Client{
	Timeout: 3 * time.Second,
}

func TestBalancer(t *testing.T) {
	// DONE: Реалізуйте юніт-тест для балансувальникка.
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

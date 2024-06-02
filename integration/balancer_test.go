package integration

import (
	"net/http"
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

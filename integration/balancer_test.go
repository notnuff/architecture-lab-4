package integration

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

const baseAddress = "http://balancer:8090"

var client = http.Client{
	Timeout: 3 * time.Second,
}

func Test(t *testing.T) {
	if _, exists := os.LookupEnv("INTEGRATION_TEST"); !exists {
		t.Skip("Integration test is not enabled")
	}
	TestingT(t)
}

type TestBalancerSuite struct{}

var _ = Suite(&TestBalancerSuite{})

func (s *TestBalancerSuite) TestBalancer(c *C) {
	// DONE: Реалізуйте інтеграційний тест для балансувальникка.
	// В цьому тесті, на мою думку, варто погратися із математичною статистикою, і, скажімо, подивитися на
	// кінцеві використання кожного із серверів. Якщо використання кожного із серверів більш-менш рівномірне,
	// то і балансувальник працює правильно. Хоч це і не ідеально точний тест, проте він достатньо точний
	// для наших потреб. А ще, це прикольне застосування теорії ймовірностей у теорії програмування
	serversUses := map[string]int{}

	numOfRequests := 10000
	for i := 0; i < numOfRequests; i++ {
		resp, err := client.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))

		if err != nil {
			c.Error(err)
		}

		usedServer := resp.Header.Get("lb-from")
		serversUses[usedServer] += 1
	}

	c.Logf("servers usages: %v", serversUses)
	numOfServers := len(serversUses)
	targetServerUsage := 1.0 / float64(numOfServers)
	acceptableStatisticalError := -1 * targetServerUsage
	acceptableUsagePerServer := targetServerUsage - acceptableStatisticalError
	c.Logf("acceptable usage per server: %f", acceptableUsagePerServer)

	for s := range serversUses {
		serverPercentageLoad := float64(serversUses[s]) / float64(numOfRequests)
		c.Logf("[%s] server percentage usage: %f", s, serverPercentageLoad)

		if serverPercentageLoad <= acceptableUsagePerServer {
			c.Errorf("server percentage is too small: [%s] - %f, target usage: %f", s, serverPercentageLoad, acceptableUsagePerServer)
		}
	}

}

func BenchmarkBalancer(b *testing.B) {
	// TODO: Реалізуйте інтеграційний бенчмарк для балансувальникка.
}

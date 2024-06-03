package integration

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

const baseAddress = "http://localhost:8090"

var client = http.Client{
	Timeout: 3 * time.Second,
}

func Test(t *testing.T) { TestingT(t) }

type TestBalancerSuite struct{}

var _ = Suite(&TestBalancerSuite{})

func (s *TestBalancerSuite) TestBalancer(c *C) {
	//if _, exists := os.LookupEnv("INTEGRATION_TEST"); !exists {
	//	t.Skip("Integration test is not enabled")
	//}

	// DONE: Реалізуйте інтеграційний тест для балансувальникка.
	// В цьому тесті, на мою думку, варто погратися із математичною статистикою, і, скажімо, подивитися на
	// кінцеві використання кожного із серверів. Якщо використання кожного із серверів більш-менш рівномірне,
	// то і балансувальник працює правильно. Хоч це і не ідеально точний тест, проте він достатньо точний
	// для наших потреб. А ще, це прикольне застосування теорії ймовірностей у теорії програмування)

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
	acceptableStatisticalError := 0.1 * targetServerUsage
	acceptableUsagePerServer := targetServerUsage - acceptableStatisticalError
	c.Logf("acceptable usage per server: %f", acceptableUsagePerServer)

	for s := range serversUses {
		serverPercentageLoad := float64(serversUses[s]) / float64(numOfRequests)
		c.Logf("[%s] server percentage usage: %f", s, serverPercentageLoad)

		if serverPercentageLoad <= acceptableUsagePerServer {
			c.Error("server percentage is too small", s, serverPercentageLoad)
		}
	}

}

func BenchmarkBalancer(b *testing.B) {
	// TODO: Реалізуйте інтеграційний бенчмарк для балансувальникка.
}

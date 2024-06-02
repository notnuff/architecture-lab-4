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

func BenchmarkBalancer(b *testing.B) {
	// TODO: Реалізуйте інтеграційний бенчмарк для балансувальникка.
}

package runner

import (
	"runtime"
	"time"
)

// Run executes the test
//
//	report, err := runner.Run(
//		"helloworld.Greeter.SayHello",
//		"localhost:50051",
//		WithProtoFile("greeter.proto", []string{}),
//		WithDataFromFile("data.json"),
//		WithInsecure(true),
//	)
func Run(call, host string, options ...Option) (*Report, error) {
	c, err := NewConfig(call, host, options...)

	if err != nil {
		return nil, err
	}

	oldCPUs := runtime.NumCPU()

	runtime.GOMAXPROCS(c.cpus)
	defer runtime.GOMAXPROCS(oldCPUs)

	reqr, err := NewRequester(c)

	if err != nil {
		return nil, err
	}
	if c.z > 0 {
		go func() {
			time.Sleep(c.z)
			reqr.Stop(ReasonTimeout)
		}()
	}

	rep, err := reqr.Run()

	return rep, err
}

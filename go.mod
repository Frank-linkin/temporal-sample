module github.com/temporalio/samples-go

go 1.16

require (
	github.com/golang/mock v1.6.0
	github.com/golang/snappy v0.0.4
	github.com/hashicorp/go-plugin v1.4.3
	github.com/mitchellh/mapstructure v1.1.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pborman/uuid v1.2.1
	github.com/prometheus/client_golang v1.11.0
	github.com/stretchr/testify v1.7.0
	github.com/uber-go/tally/v4 v4.1.1
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	go.temporal.io/api v1.7.1-0.20220223032354-6e6fe738916a
	go.temporal.io/sdk v1.13.0
	go.temporal.io/sdk/contrib/opentracing v0.1.0
	go.temporal.io/sdk/contrib/tally v0.1.0
	go.temporal.io/server v1.14.1
	go.uber.org/goleak v1.1.12 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.21.0
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
	golang.org/x/time v0.0.0-20211116232009-f0f3c7e86c11 // indirect
	google.golang.org/genproto v0.0.0-20220228195345-15d65a4533f7 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/square/go-jose.v2 v2.6.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace go.temporal.io/sdk => ../sdk-go

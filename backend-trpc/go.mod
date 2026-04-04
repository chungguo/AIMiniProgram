module aiminiprogram/backend-trpc

go 1.21

require (
	trpc.group/trpc-go/trpc-go v1.0.0
	google.golang.org/protobuf v1.32.0
	github.com/lib/pq v1.10.9
)

replace aiminiprogram/proto => ../proto

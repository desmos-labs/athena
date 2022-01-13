module github.com/desmos-labs/djuno/v2

go 1.13

require (
	cloud.google.com/go/iam v0.1.0 // indirect
	firebase.google.com/go v3.13.0+incompatible
	github.com/cosmos/cosmos-sdk v0.44.5
	github.com/cosmos/ibc-go v1.2.5
	github.com/desmos-labs/desmos/v2 v2.3.1
	github.com/forbole/juno/v2 v2.0.0-20220113082840-5b6bf27ac741
	github.com/go-co-op/gocron v1.11.0
	github.com/gogo/protobuf v1.3.3
	github.com/jmoiron/sqlx v1.3.4
	github.com/proullon/ramsql v0.0.0-20181213202341-817cee58a244
	github.com/rs/zerolog v1.26.1
	github.com/spf13/cobra v1.3.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.14
	github.com/ziutek/mymysql v1.5.4 // indirect
	google.golang.org/api v0.65.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/cosmos/cosmos-sdk => github.com/desmos-labs/cosmos-sdk v0.43.0-alpha1.0.20211102084520-683147efd235

replace google.golang.org/grpc => google.golang.org/grpc v1.42.0

replace github.com/cosmos/ledger-cosmos-go => github.com/desmos-labs/ledger-desmos-go v0.11.2-0.20210814121638-5d87e392e8a9

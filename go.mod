module github.com/desmos-labs/djuno

go 1.13

require (
	firebase.google.com/go v3.12.0+incompatible
	github.com/cosmos/cosmos-sdk v0.44.2
	github.com/cosmos/ibc-go v1.2.2
	github.com/desmos-labs/desmos/v2 v2.1.0-testnet
	github.com/forbole/juno/v2 v2.0.0-20211020184842-e358a33007ff
	github.com/go-co-op/gocron v0.3.3
	github.com/gogo/protobuf v1.3.3
	github.com/jmoiron/sqlx v1.2.0
	github.com/proullon/ramsql v0.0.0-20181213202341-817cee58a244
	github.com/rs/zerolog v1.26.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.35.0
	github.com/ziutek/mymysql v1.5.4 // indirect
	google.golang.org/api v0.56.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

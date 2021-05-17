module github.com/desmos-labs/djuno

go 1.13

require (
	cloud.google.com/go v0.55.0 // indirect
	cloud.google.com/go/firestore v1.1.1 // indirect
	firebase.google.com/go v3.12.0+incompatible
	github.com/BurntSushi/toml v0.3.1
	github.com/cosmos/cosmos-sdk v0.42.4
	github.com/desmos-labs/desmos v0.16.3-0.20210517103958-e847d726e914
	github.com/desmos-labs/juno v0.0.0-20210513082948-fad7f160e2cd
	github.com/jmoiron/sqlx v1.2.0
	github.com/proullon/ramsql v0.0.0-20181213202341-817cee58a244
	github.com/rs/zerolog v1.20.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.9
	github.com/ziutek/mymysql v1.5.4 // indirect
	google.golang.org/api v0.20.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/cosmos/cosmos-sdk => github.com/desmos-labs/cosmos-sdk v0.42.5-0.20210510082024-f2b3abdd6ea4

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

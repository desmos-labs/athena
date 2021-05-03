module github.com/desmos-labs/djuno

go 1.13

require (
	cloud.google.com/go v0.55.0 // indirect
	cloud.google.com/go/firestore v1.1.1 // indirect
	firebase.google.com/go v3.12.0+incompatible
	github.com/cosmos/cosmos-sdk v0.42.4
	github.com/desmos-labs/desmos v0.16.0
	github.com/desmos-labs/juno v0.0.0-20210503094329-1df81455bd10
	github.com/jmoiron/sqlx v1.2.0
	github.com/proullon/ramsql v0.0.0-20181213202341-817cee58a244
	github.com/rs/zerolog v1.20.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.9
	github.com/ziutek/mymysql v1.5.4 // indirect
	google.golang.org/api v0.20.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

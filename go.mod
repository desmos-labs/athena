module github.com/desmos-labs/djuno

go 1.13

require (
	cloud.google.com/go v0.55.0 // indirect
	cloud.google.com/go/firestore v1.1.1 // indirect
	firebase.google.com/go v3.12.0+incompatible
	github.com/cosmos/cosmos-sdk v0.40.0-rc4
	github.com/desmos-labs/desmos v0.14.1-0.20201209131257-74ab1ef6ca7d
	github.com/desmos-labs/juno v0.0.0-20201209082915-17b5f3e771be
	github.com/go-co-op/gocron v0.3.3
	github.com/jmoiron/sqlx v1.2.0
	github.com/proullon/ramsql v0.0.0-20181213202341-817cee58a244
	github.com/rs/zerolog v1.18.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/tendermint v0.34.0-rc6
	github.com/ziutek/mymysql v1.5.4 // indirect
	golang.org/x/tools v0.0.0-20200321224714-0d839f3cf2ed // indirect
	google.golang.org/api v0.20.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4

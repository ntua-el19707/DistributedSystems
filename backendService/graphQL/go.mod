module github.com/ntua-el19707/DistributedSystems/backendService/graphQL

go 1.21.3

replace Service v0.0.0 => ../services/service

require Service v0.0.0

replace balanceResolver v0.0.0 => ./balanceResolver

require balanceResolver v0.0.0

replace BlockCoinChainResolver v0.0.0 => ./BlockCoinChainResolver

require BlockCoinChainResolver v0.0.0

replace blockChainMsgResolver v0.0.0 => ./blockChainMsgResolver

require blockChainMsgResolver v0.0.0

replace TransactionCoinsResolver v0.0.0 => ./transactionsCoinsResolver

require TransactionCoinsResolver v0.0.0

replace Logger v0.0.0 => ../services/Logger

require Logger v0.0.0

replace NodeDetails v0.0.0 => ./NodeDeatails

require NodeDetails v0.0.0

replace TransactionMsgResolver v0.0.0 => ./TransactionMsgResolver

require TransactionMsgResolver v0.0.0

require github.com/graphql-go/graphql v0.8.1

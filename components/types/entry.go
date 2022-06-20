package types

import "main/components/proto"

func Initialize() {
	Register([][]any{
		{
			proto.AccountType_STUDENT,
			proto.AccountType_JUDGE,
			proto.AccountType_ADMIN,
		},
		{
			proto.GroupingType_EVENT,
			proto.GroupingType_TEAM,
		},
	})
}

//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/sqlite"
)

var JudgeEventRel = newJudgeEventRelTable("", "judge_event_rel", "")

type judgeEventRelTable struct {
	sqlite.Table

	// Columns
	JudgeID sqlite.ColumnInteger
	EventID sqlite.ColumnString

	AllColumns     sqlite.ColumnList
	MutableColumns sqlite.ColumnList
}

type JudgeEventRelTable struct {
	judgeEventRelTable

	EXCLUDED judgeEventRelTable
}

// AS creates new JudgeEventRelTable with assigned alias
func (a JudgeEventRelTable) AS(alias string) *JudgeEventRelTable {
	return newJudgeEventRelTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new JudgeEventRelTable with assigned schema name
func (a JudgeEventRelTable) FromSchema(schemaName string) *JudgeEventRelTable {
	return newJudgeEventRelTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new JudgeEventRelTable with assigned table prefix
func (a JudgeEventRelTable) WithPrefix(prefix string) *JudgeEventRelTable {
	return newJudgeEventRelTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new JudgeEventRelTable with assigned table suffix
func (a JudgeEventRelTable) WithSuffix(suffix string) *JudgeEventRelTable {
	return newJudgeEventRelTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newJudgeEventRelTable(schemaName, tableName, alias string) *JudgeEventRelTable {
	return &JudgeEventRelTable{
		judgeEventRelTable: newJudgeEventRelTableImpl(schemaName, tableName, alias),
		EXCLUDED:           newJudgeEventRelTableImpl("", "excluded", ""),
	}
}

func newJudgeEventRelTableImpl(schemaName, tableName, alias string) judgeEventRelTable {
	var (
		JudgeIDColumn  = sqlite.IntegerColumn("judge_id")
		EventIDColumn  = sqlite.StringColumn("event_id")
		allColumns     = sqlite.ColumnList{JudgeIDColumn, EventIDColumn}
		mutableColumns = sqlite.ColumnList{}
	)

	return judgeEventRelTable{
		Table: sqlite.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		JudgeID: JudgeIDColumn,
		EventID: EventIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}

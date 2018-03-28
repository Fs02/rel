package sqlutil

import (
	"testing"

	"github.com/Fs02/grimoire"
	. "github.com/Fs02/grimoire/c"
	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	users := grimoire.Query{
		Collection: "users",
		Fields:     []string{"*"},
	}

	tests := []struct {
		QueryString string
		Args        []interface{}
		Query       grimoire.Query
	}{
		{
			"SELECT * FROM users;",
			nil,
			users,
		},
		{
			"SELECT id, name FROM users;",
			nil,
			users.Select("id", "name"),
		},
		{
			"SELECT * FROM users JOIN transactions ON transactions.id=users.transaction_id;",
			nil,
			users.Join("transactions", Eq(I("transactions.id"), I("users.transaction_id"))),
		},
		{
			"SELECT * FROM users WHERE id=?;",
			[]interface{}{10},
			users.Where(Eq(I("id"), 10)),
		},
		{
			"SELECT DISTINCT * FROM users GROUP BY type;",
			nil,
			users.Distinct().Group("type"),
		},
		{
			"SELECT * FROM users JOIN transactions ON transactions.id=users.transaction_id HAVING price>?;",
			[]interface{}{1000},
			users.Join("transactions", Eq(I("transactions.id"), I("users.transaction_id"))).Having(Gt(I("price"), 1000)),
		},
		{
			"SELECT * FROM users ORDER BY created_at ASC;",
			nil,
			users.Order(Asc("created_at")),
		},
		{
			"SELECT * FROM users OFFSET 10 LIMIT 10;",
			nil,
			users.Offset(10).Limit(10),
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder("?", false).Find(tt.Query)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestFindOrdinal(t *testing.T) {
	users := grimoire.Query{
		Collection: "users",
		Fields:     []string{"*"},
	}

	tests := []struct {
		QueryString string
		Args        []interface{}
		Query       grimoire.Query
	}{
		{
			"SELECT * FROM users;",
			nil,
			users,
		},
		{
			"SELECT id, name FROM users;",
			nil,
			users.Select("id", "name"),
		},
		{
			"SELECT * FROM users JOIN transactions ON transactions.id=users.transaction_id;",
			nil,
			users.Join("transactions", Eq(I("transactions.id"), I("users.transaction_id"))),
		},
		{
			"SELECT * FROM users WHERE id=$1;",
			[]interface{}{10},
			users.Where(Eq(I("id"), 10)),
		},
		{
			"SELECT DISTINCT * FROM users GROUP BY type;",
			nil,
			users.Distinct().Group("type"),
		},
		{
			"SELECT * FROM users JOIN transactions ON transactions.id=users.transaction_id HAVING price>$1;",
			[]interface{}{1000},
			users.Join("transactions", Eq(I("transactions.id"), I("users.transaction_id"))).Having(Gt(I("price"), 1000)),
		},
		{
			"SELECT * FROM users ORDER BY created_at ASC;",
			nil,
			users.Order(Asc("created_at")),
		},
		{
			"SELECT * FROM users OFFSET 10 LIMIT 10;",
			nil,
			users.Offset(10).Limit(10),
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder("$", true).Find(tt.Query)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestSelect(t *testing.T) {
	assert.Equal(t, "SELECT *", NewBuilder("?", false).Select(false, "*"))
	assert.Equal(t, "SELECT id, name", NewBuilder("?", false).Select(false, "id", "name"))

	assert.Equal(t, "SELECT DISTINCT *", NewBuilder("?", false).Select(true, "*"))
	assert.Equal(t, "SELECT DISTINCT id, name", NewBuilder("?", false).Select(true, "id", "name"))
}

func TestFrom(t *testing.T) {
	assert.Equal(t, "FROM users", NewBuilder("?", false).From("users"))
}

func TestJoin(t *testing.T) {
	tests := []struct {
		QueryString string
		Args        []interface{}
		JoinClause  []Join
	}{
		{
			"",
			nil,
			nil,
		},
		{
			"JOIN users ON user.id=trxs.user_id",
			nil,
			grimoire.Query{Collection: "trxs"}.Join("users", Eq(I("user.id"), I("trxs.user_id"))).JoinClause,
		},
		{
			"INNER JOIN users ON user.id=trxs.user_id",
			nil,
			grimoire.Query{Collection: "trxs"}.JoinWith("INNER JOIN", "users", Eq(I("user.id"), I("trxs.user_id"))).JoinClause,
		},
		{
			"JOIN users ON user.id=trxs.user_id JOIN payments ON payments.id=trxs.payment_id",
			nil,
			grimoire.Query{Collection: "trxs"}.Join("users", Eq(I("user.id"), I("trxs.user_id"))).
				Join("payments", Eq(I("payments.id"), I("trxs.payment_id"))).JoinClause,
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder("?", false).Join(tt.JoinClause...)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestWhere(t *testing.T) {
	tests := []struct {
		QueryString string
		Args        []interface{}
		Condition   Condition
	}{
		{
			"",
			nil,
			And(),
		},
		{
			"WHERE field=?",
			[]interface{}{"value"},
			Eq(I("field"), "value"),
		},
		{
			"WHERE (field1=? AND field2=?)",
			[]interface{}{"value1", "value2"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder("?", false).Where(tt.Condition)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestWhereOrdinal(t *testing.T) {
	tests := []struct {
		QueryString string
		Args        []interface{}
		Condition   Condition
	}{
		{
			"",
			nil,
			And(),
		},
		{
			"WHERE field=$1",
			[]interface{}{"value"},
			Eq(I("field"), "value"),
		},
		{
			"WHERE (field1=$1 AND field2=$2)",
			[]interface{}{"value1", "value2"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder("$", true).Where(tt.Condition)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestGroupBy(t *testing.T) {
	assert.Equal(t, "", NewBuilder("?", false).GroupBy())
	assert.Equal(t, "GROUP BY city", NewBuilder("?", false).GroupBy("city"))
	assert.Equal(t, "GROUP BY city, nation", NewBuilder("?", false).GroupBy("city", "nation"))
}

func TestHaving(t *testing.T) {
	tests := []struct {
		QueryString string
		Args        []interface{}
		Condition   Condition
	}{
		{
			"",
			nil,
			And(),
		},
		{
			"HAVING field=?",
			[]interface{}{"value"},
			Eq(I("field"), "value"),
		},
		{
			"HAVING (field1=? AND field2=?)",
			[]interface{}{"value1", "value2"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder("?", false).Having(tt.Condition)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestHavingOrdinal(t *testing.T) {
	tests := []struct {
		QueryString string
		Args        []interface{}
		Condition   Condition
	}{
		{
			"",
			nil,
			And(),
		},
		{
			"HAVING field=$1",
			[]interface{}{"value"},
			Eq(I("field"), "value"),
		},
		{
			"HAVING (field1=$1 AND field2=$2)",
			[]interface{}{"value1", "value2"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder("$", true).Having(tt.Condition)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestOrderBy(t *testing.T) {
	assert.Equal(t, "", NewBuilder("?", false).OrderBy())
	assert.Equal(t, "ORDER BY name ASC", NewBuilder("?", false).OrderBy(Asc("name")))
	assert.Equal(t, "ORDER BY name ASC, created_at DESC", NewBuilder("?", false).OrderBy(Asc("name"), Desc("created_at")))
}

func TestOffset(t *testing.T) {
	assert.Equal(t, "", NewBuilder("?", false).Offset(0))
	assert.Equal(t, "OFFSET 10", NewBuilder("?", false).Offset(10))
}

func TestLimit(t *testing.T) {
	assert.Equal(t, "", NewBuilder("?", false).Limit(0))
	assert.Equal(t, "LIMIT 10", NewBuilder("?", false).Limit(10))
}

func TestCondition(t *testing.T) {
	tests := []struct {
		QueryString string
		Args        []interface{}
		Condition   Condition
	}{
		{
			"",
			nil,
			And(),
		},
		{
			"field=?",
			[]interface{}{"value"},
			Eq(I("field"), "value"),
		},
		{
			"?=field",
			[]interface{}{"value"},
			Eq("value", I("field")),
		},
		{
			"?=?",
			[]interface{}{"value1", "value2"},
			Eq("value1", "value2"),
		},
		{
			"field<>?",
			[]interface{}{"value"},
			Ne(I("field"), "value"),
		},
		{
			"?<>field",
			[]interface{}{"value"},
			Ne("value", I("field")),
		},
		{
			"?<>?",
			[]interface{}{"value1", "value2"},
			Ne("value1", "value2"),
		},
		{
			"field<?",
			[]interface{}{10},
			Lt(I("field"), 10),
		},
		{
			"?<field",
			[]interface{}{"value"},
			Lt("value", I("field")),
		},
		{
			"?<?",
			[]interface{}{"value1", "value2"},
			Lt("value1", "value2"),
		},
		{
			"field<=?",
			[]interface{}{10},
			Lte(I("field"), 10),
		},
		{
			"?<=field",
			[]interface{}{"value"},
			Lte("value", I("field")),
		},
		{
			"?<=?",
			[]interface{}{"value1", "value2"},
			Lte("value1", "value2"),
		},
		{
			"field>?",
			[]interface{}{10},
			Gt(I("field"), 10),
		},
		{
			"?>field",
			[]interface{}{"value"},
			Gt("value", I("field")),
		},
		{
			"?>?",
			[]interface{}{"value1", "value2"},
			Gt("value1", "value2"),
		},
		{
			"field>=?",
			[]interface{}{10},
			Gte(I("field"), 10),
		},
		{
			"?>=field",
			[]interface{}{"value"},
			Gte("value", I("field")),
		},
		{
			"?>=?",
			[]interface{}{"value1", "value2"},
			Gte("value1", "value2"),
		},
		{
			"field IS NULL",
			nil,
			Nil("field"),
		},
		{
			"field IS NOT NULL",
			nil,
			NotNil("field"),
		},
		{
			"field IN (?)",
			[]interface{}{"value1"},
			In("field", "value1"),
		},
		{
			"field IN (?,?)",
			[]interface{}{"value1", "value2"},
			In("field", "value1", "value2"),
		},
		{
			"field IN (?,?,?)",
			[]interface{}{"value1", "value2", "value3"},
			In("field", "value1", "value2", "value3"),
		},
		{
			"field NOT IN (?)",
			[]interface{}{"value1"},
			Nin("field", "value1"),
		},
		{
			"field NOT IN (?,?)",
			[]interface{}{"value1", "value2"},
			Nin("field", "value1", "value2"),
		},
		{
			"field NOT IN (?,?,?)",
			[]interface{}{"value1", "value2", "value3"},
			Nin("field", "value1", "value2", "value3"),
		},
		{
			"field LIKE ?",
			[]interface{}{"%value%"},
			Like("field", "%value%"),
		},
		{
			"field NOT LIKE ?",
			[]interface{}{"%value%"},
			NotLike("field", "%value%"),
		},
		{
			"FRAGMENT",
			nil,
			Fragment("FRAGMENT"),
		},
		{
			"(field1=? AND field2=?)",
			[]interface{}{"value1", "value2"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"(field1=? AND field2=? AND field3=?)",
			[]interface{}{"value1", "value2", "value3"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"(field1=? OR field2=?)",
			[]interface{}{"value1", "value2"},
			Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"(field1=? OR field2=? OR field3=?)",
			[]interface{}{"value1", "value2", "value3"},
			Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"(field1=? XOR field2=?)",
			[]interface{}{"value1", "value2"},
			Xor(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"(field1=? XOR field2=? XOR field3=?)",
			[]interface{}{"value1", "value2", "value3"},
			Xor(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"NOT (field1=? AND field2=?)",
			[]interface{}{"value1", "value2"},
			Not(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"NOT (field1=? AND field2=? AND field3=?)",
			[]interface{}{"value1", "value2", "value3"},
			Not(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"((field1=? OR field2=?) AND field3=?)",
			[]interface{}{"value1", "value2", "value3"},
			And(Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")), Eq(I("field3"), "value3")),
		},
		{
			"((field1=? OR field2=?) AND (field3=? OR field4=?))",
			[]interface{}{"value1", "value2", "value3", "value4"},
			And(Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")), Or(Eq(I("field3"), "value3"), Eq(I("field4"), "value4"))),
		},
		{
			"(NOT (field1=? AND field2=?) AND NOT (field3=? OR field4=?))",
			[]interface{}{"value1", "value2", "value3", "value4"},
			And(Not(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")), Not(Or(Eq(I("field3"), "value3"), Eq(I("field4"), "value4")))),
		},
		{
			"NOT (field1=? AND (field2=? OR field3=?) AND NOT (field4=? OR field5=?))",
			[]interface{}{"value1", "value2", "value3", "value4", "value5"},
			And(Not(Eq(I("field1"), "value1"), Or(Eq(I("field2"), "value2"), Eq(I("field3"), "value3")), Not(Or(Eq(I("field4"), "value4"), Eq(I("field5"), "value5"))))),
		},
		{
			"((field1 IN (?,?) OR field2 NOT IN (?)) AND field3 IN (?,?,?))",
			[]interface{}{"value1", "value2", "value3", "value4", "value5", "value6"},
			And(Or(In("field1", "value1", "value2"), Nin("field2", "value3")), In("field3", "value4", "value5", "value6")),
		},
		{
			"(field1 LIKE ? AND field2 NOT LIKE ?)",
			[]interface{}{"%value1%", "%value2%"},
			And(Like(I("field1"), "%value1%"), NotLike(I(I("field2")), "%value2%")),
		},
		{
			"",
			nil,
			Condition{Type: ConditionType(9999)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder("?", false).Condition(tt.Condition)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestConditionOrdinal(t *testing.T) {
	tests := []struct {
		QueryString string
		Args        []interface{}
		Condition   Condition
	}{
		{
			"",
			nil,
			And(),
		},
		{
			"field=$1",
			[]interface{}{"value"},
			Eq(I("field"), "value"),
		},
		{
			"$1=field",
			[]interface{}{"value"},
			Eq("value", I("field")),
		},
		{
			"$1=$2",
			[]interface{}{"value1", "value2"},
			Eq("value1", "value2"),
		},
		{
			"field<>$1",
			[]interface{}{"value"},
			Ne(I("field"), "value"),
		},
		{
			"$1<>field",
			[]interface{}{"value"},
			Ne("value", I("field")),
		},
		{
			"$1<>$2",
			[]interface{}{"value1", "value2"},
			Ne("value1", "value2"),
		},
		{
			"field<$1",
			[]interface{}{10},
			Lt(I("field"), 10),
		},
		{
			"$1<field",
			[]interface{}{"value"},
			Lt("value", I("field")),
		},
		{
			"$1<$2",
			[]interface{}{"value1", "value2"},
			Lt("value1", "value2"),
		},
		{
			"field<=$1",
			[]interface{}{10},
			Lte(I("field"), 10),
		},
		{
			"$1<=field",
			[]interface{}{"value"},
			Lte("value", I("field")),
		},
		{
			"$1<=$2",
			[]interface{}{"value1", "value2"},
			Lte("value1", "value2"),
		},
		{
			"field>$1",
			[]interface{}{10},
			Gt(I("field"), 10),
		},
		{
			"$1>field",
			[]interface{}{"value"},
			Gt("value", I("field")),
		},
		{
			"$1>$2",
			[]interface{}{"value1", "value2"},
			Gt("value1", "value2"),
		},
		{
			"field>=$1",
			[]interface{}{10},
			Gte(I("field"), 10),
		},
		{
			"$1>=field",
			[]interface{}{"value"},
			Gte("value", I("field")),
		},
		{
			"$1>=$2",
			[]interface{}{"value1", "value2"},
			Gte("value1", "value2"),
		},
		{
			"field IS NULL",
			nil,
			Nil("field"),
		},
		{
			"field IS NOT NULL",
			nil,
			NotNil("field"),
		},
		{
			"field IN ($1)",
			[]interface{}{"value1"},
			In("field", "value1"),
		},
		{
			"field IN ($1,$2)",
			[]interface{}{"value1", "value2"},
			In("field", "value1", "value2"),
		},
		{
			"field IN ($1,$2,$3)",
			[]interface{}{"value1", "value2", "value3"},
			In("field", "value1", "value2", "value3"),
		},
		{
			"field NOT IN ($1)",
			[]interface{}{"value1"},
			Nin("field", "value1"),
		},
		{
			"field NOT IN ($1,$2)",
			[]interface{}{"value1", "value2"},
			Nin("field", "value1", "value2"),
		},
		{
			"field NOT IN ($1,$2,$3)",
			[]interface{}{"value1", "value2", "value3"},
			Nin("field", "value1", "value2", "value3"),
		},
		{
			"field LIKE $1",
			[]interface{}{"%value%"},
			Like("field", "%value%"),
		},
		{
			"field NOT LIKE $1",
			[]interface{}{"%value%"},
			NotLike("field", "%value%"),
		},
		{
			"FRAGMENT",
			nil,
			Fragment("FRAGMENT"),
		},
		{
			"(field1=$1 AND field2=$2)",
			[]interface{}{"value1", "value2"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"(field1=$1 AND field2=$2 AND field3=$3)",
			[]interface{}{"value1", "value2", "value3"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"(field1=$1 OR field2=$2)",
			[]interface{}{"value1", "value2"},
			Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"(field1=$1 OR field2=$2 OR field3=$3)",
			[]interface{}{"value1", "value2", "value3"},
			Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"(field1=$1 XOR field2=$2)",
			[]interface{}{"value1", "value2"},
			Xor(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"(field1=$1 XOR field2=$2 XOR field3=$3)",
			[]interface{}{"value1", "value2", "value3"},
			Xor(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"NOT (field1=$1 AND field2=$2)",
			[]interface{}{"value1", "value2"},
			Not(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"NOT (field1=$1 AND field2=$2 AND field3=$3)",
			[]interface{}{"value1", "value2", "value3"},
			Not(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"((field1=$1 OR field2=$2) AND field3=$3)",
			[]interface{}{"value1", "value2", "value3"},
			And(Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")), Eq(I("field3"), "value3")),
		},
		{
			"((field1=$1 OR field2=$2) AND (field3=$3 OR field4=$4))",
			[]interface{}{"value1", "value2", "value3", "value4"},
			And(Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")), Or(Eq(I("field3"), "value3"), Eq(I("field4"), "value4"))),
		},
		{
			"(NOT (field1=$1 AND field2=$2) AND NOT (field3=$3 OR field4=$4))",
			[]interface{}{"value1", "value2", "value3", "value4"},
			And(Not(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")), Not(Or(Eq(I("field3"), "value3"), Eq(I("field4"), "value4")))),
		},
		{
			"NOT (field1=$1 AND (field2=$2 OR field3=$3) AND NOT (field4=$4 OR field5=$5))",
			[]interface{}{"value1", "value2", "value3", "value4", "value5"},
			And(Not(Eq(I("field1"), "value1"), Or(Eq(I("field2"), "value2"), Eq(I("field3"), "value3")), Not(Or(Eq(I("field4"), "value4"), Eq(I("field5"), "value5"))))),
		},
		{
			"((field1 IN ($1,$2) OR field2 NOT IN ($3)) AND field3 IN ($4,$5,$6))",
			[]interface{}{"value1", "value2", "value3", "value4", "value5", "value6"},
			And(Or(In("field1", "value1", "value2"), Nin("field2", "value3")), In("field3", "value4", "value5", "value6")),
		},
		{
			"(field1 LIKE $1 AND field2 NOT LIKE $2)",
			[]interface{}{"%value1%", "%value2%"},
			And(Like(I("field1"), "%value1%"), NotLike(I(I("field2")), "%value2%")),
		},
		{
			"",
			nil,
			Condition{Type: ConditionType(9999)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder("$", true).Condition(tt.Condition)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}
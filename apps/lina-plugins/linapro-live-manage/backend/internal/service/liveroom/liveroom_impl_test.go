// This file verifies buildOptionsCondition query composition so candidate
// queries never leak rows outside the current tenant regardless of keyword
// content.

package liveroom

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"

	"lina-core/pkg/dbdriver"
	"lina-plugin-linapro-live-manage/backend/internal/dao"
)

// orPattern matches the SQL OR operator as a standalone token.
var orPattern = regexp.MustCompile(`\bOR\b`)

// TestOptionsConditionGroupsKeywordInsideParens verifies keyword alternatives
// are composed as one parenthesized group so the OR clauses can never escape
// the tenant predicate when the condition is attached to the model.
func TestOptionsConditionGroupsKeywordInsideParens(t *testing.T) {
	sql, args := buildOptionsCondition(newRoomBuilderForTest(t), OptionsInput{Keyword: "boom"}).Build()

	orIdx := orPattern.FindStringIndex(sql)
	if orIdx == nil {
		t.Fatalf("expected one grouped OR clause, got %q", sql)
	}
	parenOpen := strings.LastIndex(sql[:orIdx[0]], "(")
	parenClose := strings.Index(sql[orIdx[0]:], ")")
	if parenOpen < 0 || parenClose < 0 {
		t.Fatalf("expected OR clause enclosed in parentheses, got %q", sql)
	}
	if !strings.Contains(strings.ToUpper(sql), "NOT IN") {
		t.Fatalf("expected disabled rooms excluded from the same condition, got %q", sql)
	}
	if len(args) != 3 || args[0] != "%boom%" || args[1] != "%boom%" {
		t.Fatalf("expected keyword args %%boom%% twice plus status arg, got %v", args)
	}
}

// TestOptionsConditionWithoutKeywordExcludesDisabledOnly verifies a blank
// keyword query carries no OR clause and keeps the status exclusion.
func TestOptionsConditionWithoutKeywordExcludesDisabledOnly(t *testing.T) {
	sql, _ := buildOptionsCondition(newRoomBuilderForTest(t), OptionsInput{}).Build()

	if orPattern.MatchString(sql) {
		t.Fatalf("expected no OR clause without keyword, got %q", sql)
	}
	if !strings.Contains(strings.ToUpper(sql), "NOT IN") {
		t.Fatalf("expected disabled rooms excluded, got %q", sql)
	}
}

// TestOptionsConditionIncludeDisabledDropsStatusExclusion verifies
// IncludeDisabled drops the status exclusion while keeping keyword grouping.
func TestOptionsConditionIncludeDisabledDropsStatusExclusion(t *testing.T) {
	sql, _ := buildOptionsCondition(newRoomBuilderForTest(t), OptionsInput{
		Keyword:         "boom",
		IncludeDisabled: true,
	}).Build()

	if orPattern.FindStringIndex(sql) == nil {
		t.Fatalf("expected grouped keyword clause, got %q", sql)
	}
	if strings.Contains(strings.ToUpper(sql), "NOT IN") {
		t.Fatalf("expected status exclusion dropped, got %q", sql)
	}
}

// TestOptionsConditionWithoutFilterReturnsNil verifies no condition is
// attached when all statuses are included and no keyword is given.
func TestOptionsConditionWithoutFilterReturnsNil(t *testing.T) {
	if condition := buildOptionsCondition(newRoomBuilderForTest(t), OptionsInput{IncludeDisabled: true}); condition != nil {
		t.Fatal("expected nil condition without keyword and status exclusion")
	}
}

// newRoomBuilderForTest registers a non-connecting gdb config and returns a
// fresh condition builder seeded from the room model.
func newRoomBuilderForTest(t *testing.T) *gdb.WhereBuilder {
	t.Helper()

	previous := gdb.GetConfig("default")
	t.Cleanup(func() {
		if err := gdb.SetConfigGroup("default", previous); err != nil {
			t.Logf("restore gdb config failed: %v", err)
		}
	})
	if err := gdb.SetConfigGroup("default", gdb.ConfigGroup{
		{
			Type: dbdriver.TypePostgreSQL,
			Link: "pgsql:user:pass@tcp(127.0.0.1:5432)/test",
		},
	}); err != nil {
		t.Fatalf("register test gdb config failed: %v", err)
	}

	return dao.Room.Ctx(context.Background()).Builder()
}

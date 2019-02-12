package parser

import (
	"fmt"
	"testing"

	"github.com/pingcap/parser"
	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	p := parser.New()

	// 2. Parse a text SQL into AST([]ast.StmtNode).
	stmtNodes, _, err := p.Parse("select * from tbl where id = 1", "", "")

	// 3. Use AST to do cool things.
	fmt.Println(stmtNodes[0], err)

	assert.Equal(t, 0, 1)
}

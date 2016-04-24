package compiler

import (
	"testing"

	"github.com/twtiger/gosecco/asm"

	"github.com/twtiger/gosecco/tree"

	. "gopkg.in/check.v1"
)

func IncludeTest(t *testing.T) { TestingT(t) }

type IncludeCompilerSuite struct{}

var _ = Suite(&IncludeCompilerSuite{})

func (s *IncludeCompilerSuite) Test_compliationOfIncludeOperation(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Inclusion{
					Positive: true,
					Left:     tree.Argument{Index: 0},
					Rights:   []tree.Numeric{tree.NumericLiteral{1}, tree.NumericLiteral{2}},
				},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+ // syscallNameIndex
		"jeq_k	00	06	1\n"+ // syscall.SYS_WRITE
		"ld_abs	14\n"+ //argumentindex[0][upper]
		"jeq_k	00	04	0\n"+
		"ld_abs	10\n"+ //argumentindex[0][lower]
		"jeq_k	01	00	1\n"+ // compare to first number in list
		"jeq_k	00	01	2\n"+ // compare to second number in list
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}

func (s *IncludeCompilerSuite) Test_compliationOfNotIncludeOperation(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Inclusion{
					Positive: false,
					Left:     tree.Argument{Index: 0},
					Rights:   []tree.Numeric{tree.NumericLiteral{1}, tree.NumericLiteral{2}},
				},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+ // syscallNameIndex
		"jeq_k	00	06	1\n"+ // syscall.SYS_WRITE
		"ld_abs	14\n"+ //argumentindex[0][upper]
		"jeq_k	00	04	0\n"+
		"ld_abs	10\n"+ //argumentindex[0][lower]
		"jeq_k	02	00	1\n"+ // compare to first number in list
		"jeq_k	01	00	2\n"+ // compare to second number in list
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}

func (s *IncludeCompilerSuite) Test_compliationOfArgumentsInIncludeList(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Inclusion{
					Positive: true,
					Left:     tree.NumericLiteral{1},
					Rights:   []tree.Numeric{tree.Argument{Index: 1}, tree.Argument{Index: 0}},
				},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs\t0\n"+ // syscallNameIndex
		"jeq_k\t00\t0B\t1\n"+ // syscall.SYS_WRITE
		"ld_imm\t1\n"+ // load K into A
		"tax\n"+ // move A to X
		"ld_abs\t1C\n"+ // load first half of argument 1
		"jeq_k\t00\t07\t0\n"+ // compare it to 0
		"ld_abs\t18\n"+ //load second half of argument 1
		"jeq_x\t04\t00\n"+ // compare it to X
		"ld_abs\t14\n"+ // load first half of argument 0
		"jeq_k\t00\t03\t0\n"+ // compare it to 0
		"ld_abs\t10\n"+ // load second half of argument 1
		"jeq_x\t00\t01\n"+ // compare it to X
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}

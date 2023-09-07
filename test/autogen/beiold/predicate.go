// Package beiold
// @author tabuyos
// @since 2023/9/5
// @description beiold
package beiold

import (
	"fmt"
	"strings"
)

type Sym func(bool) string

func (s Sym) String() string {
	return s.Self()
}

func (s Sym) Ph() string {
	return s(true)
}

func (s Sym) Self() string {
	return s(false)
}

func andSym(bool) string {
	return "AND"
}

func orSym(bool) string {
	return "OR"
}

func eqSym(has bool) string {
	if has {
		return "= ?"
	}
	return "="
}

func nqSym(has bool) string {
	if has {
		return "<> ?"
	}
	return "<>"
}

var (
	AND Sym = andSym
	OR  Sym = orSym
	EQ  Sym = eqSym
	NQ  Sym = nqSym
)

type Mode int

const (
	DftMode Mode = iota + 1
	AndMode
	OrMode
)

type Predicate struct {
	mod Mode
	fid string
	lft *Predicate
	sym Sym
	rht *Predicate
	ars []any
}

func Once(f string, s Sym, ars ...any) *Predicate {
	return &Predicate{
		mod: DftMode,
		fid: f,
		sym: s,
		ars: ars,
	}
}

func (p *Predicate) And(preds ...*Predicate) *Predicate {
	var fp = p
	for _, pred := range preds {
		fp = &Predicate{
			mod: AndMode,
			lft: fp,
			sym: AND,
			rht: pred,
		}
	}
	return fp
}

func (p *Predicate) Or(preds ...*Predicate) *Predicate {
	var fp = p
	for _, pred := range preds {
		fp = &Predicate{
			mod: OrMode,
			lft: fp,
			sym: OR,
			rht: pred,
		}
	}
	return fp
}

func (p *Predicate) String() string {
	sql, values := p.SQL()
	var index = 0
	var buf strings.Builder
	for _, r := range sql {
		if r == '?' {
			buf.WriteString(fmt.Sprintf("%v", values[index]))
			index++
			continue
		}
		buf.WriteRune(r)
	}
	return buf.String()
}

func (p *Predicate) SQL() (sql string, values []any) {
	if p == nil {
		return "", nil
	}
	if p.lft == nil {
		sql = fmt.Sprintf("%s %s", p.fid, p.sym.Ph())
		values = append(values, p.ars...)
		return
	}
	lftSQL, lftValues := p.lft.SQL()
	rhtSQL, rhtValues := p.rht.SQL()
	if p.lft.mod == DftMode {
		p.lft.mod = p.mod
	}
	if p.rht.mod == DftMode {
		p.rht.mod = p.mod
	}
	if p.mod == p.lft.mod {
		if p.mod == p.rht.mod {
			sql = fmt.Sprintf("%s %s %s", lftSQL, p.sym.Ph(), rhtSQL)
			values = append(values, lftValues...)
			values = append(values, rhtValues...)
		} else {
			sql = fmt.Sprintf("%s %s (%s)", lftSQL, p.sym.Ph(), rhtSQL)
			values = append(values, lftValues...)
			values = append(values, rhtValues...)
		}
	} else {
		if p.mod == p.rht.mod {
			sql = fmt.Sprintf("(%s) %s %s", lftSQL, p.sym.Ph(), rhtSQL)
			values = append(values, lftValues...)
			values = append(values, rhtValues...)
		} else {
			sql = fmt.Sprintf("(%s) %s (%s)", lftSQL, p.sym.Ph(), rhtSQL)
			values = append(values, lftValues...)
			values = append(values, rhtValues...)
		}
	}
	return
}

package provider

import (
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
)

type classicalStrategy struct {
	rules           []C.Rule
	count           int
	shouldResolveIP bool
}

func (c *classicalStrategy) Match(metadata *C.Metadata) bool {
	for _, rule := range c.rules {
		if rule.Match(metadata) {
			return true
		}
	}

	return false
}

func (c *classicalStrategy) Count() int {
	return c.count
}

func (c *classicalStrategy) ShouldResolveIP() bool {
	return c.shouldResolveIP
}

func (c *classicalStrategy) OnUpdate(rules []string) {
	for _, rawRule := range rules {
		ruleType, rule, params := ruleParse(rawRule)
		r, err := parseRule(ruleType, rule, "", params)
		if err != nil {
			log.Warnln("parse rule error:[%s]", err.Error())
		} else {
			if !c.shouldResolveIP {
				c.shouldResolveIP = r.ShouldResolveIP()
			}

			c.rules = append(c.rules, r)
			c.count++
		}
	}
}

func NewClassicalStrategy() *classicalStrategy {
	return &classicalStrategy{rules: []C.Rule{}}
}

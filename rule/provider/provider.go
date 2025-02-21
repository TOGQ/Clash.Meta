package provider

import (
	"encoding/json"
	C "github.com/Dreamacro/clash/constant"
	P "github.com/Dreamacro/clash/constant/provider"
	"gopkg.in/yaml.v2"
	"runtime"
	"time"
)

var (
	ruleProviders = map[string]P.RuleProvider{}
)

type ruleSetProvider struct {
	*fetcher
	behavior P.RuleType
	strategy ruleStrategy
}

type RuleSetProvider struct {
	*ruleSetProvider
}

type RulePayload struct {
	/**
	key: Domain or IP Cidr
	value: Rule type or is empty
	*/
	Rules  []string `yaml:"payload"`
	Rules2 []string `yaml:"rules"`
}

type ruleStrategy interface {
	Match(metadata *C.Metadata) bool
	Count() int
	ShouldResolveIP() bool
	OnUpdate(rules []string)
}

func RuleProviders() map[string]P.RuleProvider {
	return ruleProviders
}

func SetRuleProvider(ruleProvider P.RuleProvider) {
	if ruleProvider != nil {
		ruleProviders[(ruleProvider).Name()] = ruleProvider
	}
}

func (rp *ruleSetProvider) Type() P.ProviderType {
	return P.Rule
}

func (rp *ruleSetProvider) Initial() error {
	elm, err := rp.fetcher.Initial()
	if err != nil {
		return err
	}

	return rp.fetcher.onUpdate(elm)
}

func (rp *ruleSetProvider) Update() error {
	elm, same, err := rp.fetcher.Update()
	if err == nil && !same {
		return rp.fetcher.onUpdate(elm)
	}

	return err
}

func (rp *ruleSetProvider) Behavior() P.RuleType {
	return rp.behavior
}

func (rp *ruleSetProvider) Match(metadata *C.Metadata) bool {
	return rp.strategy != nil && rp.strategy.Match(metadata)
}

func (rp *ruleSetProvider) ShouldResolveIP() bool {
	return rp.strategy.ShouldResolveIP()
}

func (rp *ruleSetProvider) AsRule(adaptor string) C.Rule {
	panic("implement me")
}

func (rp *ruleSetProvider) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		map[string]interface{}{
			"behavior":    rp.behavior.String(),
			"name":        rp.Name(),
			"ruleCount":   rp.strategy.Count(),
			"type":        rp.Type().String(),
			"updatedAt":   rp.updatedAt,
			"vehicleType": rp.VehicleType().String(),
		})
}

func NewRuleSetProvider(name string, behavior P.RuleType, interval time.Duration, vehicle P.Vehicle) P.RuleProvider {
	rp := &ruleSetProvider{
		behavior: behavior,
	}

	onUpdate := func(elm interface{}) error {
		rulesRaw := elm.([]string)
		rp.strategy.OnUpdate(rulesRaw)
		return nil
	}

	fetcher := newFetcher(name, interval, vehicle, rulesParse, onUpdate)
	rp.fetcher = fetcher
	rp.strategy = newStrategy(behavior)

	wrapper := &RuleSetProvider{
		rp,
	}

	final := func(provider *RuleSetProvider) { rp.fetcher.Destroy() }
	runtime.SetFinalizer(wrapper, final)
	return wrapper
}

func newStrategy(behavior P.RuleType) ruleStrategy {
	switch behavior {
	case P.Domain:
		strategy := NewDomainStrategy()
		return strategy
	case P.IPCIDR:
		strategy := NewIPCidrStrategy()
		return strategy
	case P.Classical:
		strategy := NewClassicalStrategy()
		return strategy
	default:
		return nil
	}
}

func rulesParse(buf []byte) (interface{}, error) {
	rulePayload := RulePayload{}
	err := yaml.Unmarshal(buf, &rulePayload)
	if err != nil {
		return nil, err
	}

	return append(rulePayload.Rules, rulePayload.Rules2...), nil
}

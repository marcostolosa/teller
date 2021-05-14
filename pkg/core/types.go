package core

type Severity string

const (
	High   Severity = "high"
	Medium Severity = "medium"
	Low    Severity = "low"
	None   Severity = "none"
)

type KeyPath struct {
	Env        string            `yaml:"env,omitempty"`
	Path       string            `yaml:"path"`
	Field      string            `yaml:"field,omitempty"`
	Remap      map[string]string `yaml:"remap,omitempty"`
	Decrypt    bool              `yaml:"decrypt,omitempty"`
	Optional   bool              `yaml:"optional,omitempty"`
	Severity   Severity          `yaml:"severity,omitempty" default:"high"`
	RedactWith string            `yaml:"redact_with,omitempty" default:"**REDACTED**"`
	Source     string            `yaml:"source,omitempty"`
	Sink       string            `yaml:"sink,omitempty"`
}
type WizardAnswers struct {
	Project      string
	Providers    []string
	ProviderKeys map[string]bool
	Confirm      bool
}

func (k *KeyPath) WithEnv(env string) KeyPath {
	return KeyPath{
		Env:      env,
		Path:     k.Path,
		Field:    k.Field,
		Decrypt:  k.Decrypt,
		Optional: k.Optional,
		Source:   k.Source,
		Sink:     k.Sink,
	}
}
func (k *KeyPath) SwitchPath(path string) KeyPath {
	return KeyPath{
		Path:     path,
		Field:    k.Field,
		Env:      k.Env,
		Decrypt:  k.Decrypt,
		Optional: k.Optional,
		Source:   k.Source,
		Sink:     k.Sink,
	}
}

type DriftedEntriesBySource []DriftedEntry

func (a DriftedEntriesBySource) Len() int           { return len(a) }
func (a DriftedEntriesBySource) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DriftedEntriesBySource) Less(i, j int) bool { return a[i].Source.Source < a[j].Source.Source }

type EntriesByKey []EnvEntry

func (a EntriesByKey) Len() int           { return len(a) }
func (a EntriesByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a EntriesByKey) Less(i, j int) bool { return a[i].Key > a[j].Key }

type EntriesByValueSize []EnvEntry

func (a EntriesByValueSize) Len() int           { return len(a) }
func (a EntriesByValueSize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a EntriesByValueSize) Less(i, j int) bool { return len(a[i].Value) > len(a[j].Value) }

type EnvEntry struct {
	Key          string
	Value        string
	ProviderName string
	Provider     string
	ResolvedPath string
	Severity     Severity
	RedactWith   string
	Source       string
	Sink         string
}
type DriftedEntry struct {
	Diff   string
	Source EnvEntry
	Target EnvEntry
}
type EnvEntryLookup struct {
	Entries []EnvEntry
}

func (ee *EnvEntryLookup) EnvBy(key, provider, path, dflt string) string {
	for i := range ee.Entries {
		e := ee.Entries[i]
		if e.Key == key && e.ProviderName == provider && e.ResolvedPath == path {
			return e.Value
		}
	}
	return dflt
}
func (ee *EnvEntryLookup) EnvByKey(key, dflt string) string {
	for i := range ee.Entries {
		e := ee.Entries[i]
		if e.Key == key {
			return e.Value
		}

	}
	return dflt
}

func (ee *EnvEntryLookup) EnvByKeyAndProvider(key, provider, dflt string) string {
	for i := range ee.Entries {
		e := ee.Entries[i]
		if e.Key == key && e.ProviderName == provider {
			return e.Value
		}

	}
	return dflt
}

type Provider interface {
	Name() string
	// in this case 'env' is empty, but EnvEntries are the value
	GetMapping(p KeyPath) ([]EnvEntry, error)

	// in this case env is filled
	Get(p KeyPath) (*EnvEntry, error)

	Put(p KeyPath, val string) error
}

type Match struct {
	Path       string
	Line       string
	LineNumber int
	MatchIndex int
	Entry      EnvEntry
}

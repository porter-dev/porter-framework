package porter

import (
	"github.com/porterdev/ego/internal/plan"
	v "github.com/porterdev/ego/internal/value"
)

// Object is a Porter object: it resolves to an interface{} type, but internally
// the types are cast and validated
type Object v.Value

// Config is an interface that Porter templates must implement
type Config interface {
	Apply(input Object) (Object, error)
	Data(input Object) (Object, error)
	Generate(input Object) (Object, error)
	Plan(old Object, new Object) (*plan.OpQueue, error)
	Run(op *plan.Operation) error
	Validate(op *plan.Operation) error
	Save(v Object) error
}

// DefaultConfig is the generic implementation of a Porter configuration. Generate()
// and Run() should be overwritten, while Plan() should not typically be overwritten,
// and Validate(), Save() and Data() can be optionally overwritten.
//
// The ID will uniquely identify an instance of a configuration. This is used by the
// default Data storage (as a file on the local fs) to write a certain configuration
// and store configuration backups. It should also be used in implementations of Data.
//
// The DefaultConfig struct logs at two levels: error and info. If LogLevel is 0, only
// errors are logged. If LogLevel is 1, warning and error logs are written. If LogLevel
// is 2, all logs are written.
type DefaultConfig struct {
	ID string

	Logger *Logger
	Store  *LocalStore
}

// CreateDefaultConfig creates a new configuration based on an ID and a logLevel.
func CreateDefaultConfig(id string, stateDir string, backupDir string, logLevel int) *DefaultConfig {
	logger := NewLogger(logLevel)

	store, err := NewLocalStore(id, logger, stateDir, backupDir)
	logger.Check(err, id, "could not initialize store")

	conf := DefaultConfig{
		ID:     id,
		Logger: logger,
		Store:  store,
	}

	return &conf
}

// Apply runs an application loop for a given Config. This should never be overwritten.
// Returns the new configuration.
func (c DefaultConfig) Apply(input Object) (Object, error) {
	old, err := c.Data(input)

	c.Logger.Check(err, c.ID, "data retrieval failed")
	c.Logger.Log(INFO, c.ID, "successfully retrieved data for configuration")

	new, err := c.Generate(input)

	c.Logger.Check(err, c.ID, "config generation failed")
	c.Logger.Log(INFO, c.ID, "successfully generated configuration")

	q, err := c.Plan(old, new)

	c.Logger.Check(err, c.ID, "plan failed")
	c.Logger.Log(INFO, c.ID, "successfully generated plan")

	for !q.IsEmpty() {
		op := q.Dequeue()

		err := c.Run(op)

		c.Logger.Check(err, c.ID, "run", op.ToString(), "failed")
		c.Logger.Log(INFO, c.ID, "successfully ran:", op.ToString())

		err = c.Validate(op)

		c.Logger.Check(err, c.ID, "validate", op.ToString(), "failed")
		c.Logger.Log(INFO, c.ID, "successfully validated:", op.ToString())
	}

	err = c.Save(new)

	c.Logger.Check(err, c.ID, "save failed")
	c.Logger.Log(INFO, c.ID, "successfully saved")

	return new, nil
}

// Data is the default implementation of Config.Data(), and can optionally be overwritten.
// It simply returns the input values as the retrieved data for this configuration.
func (c DefaultConfig) Data(input Object) (Object, error) {
	return c.Store.GetState()
}

// Generate is the default implementation of Config.Generate(), and should be overwritten.
// It simply returns the input values.
func (c DefaultConfig) Generate(input Object) (Object, error) {
	return input, nil
}

// Plan is the default implementation of Config.Plan(), and should **not** be overwritten,
// unless you know what you are doing. It generates an execution plan by comparing the
// passed Value against a previously stored Value.
func (c DefaultConfig) Plan(old Object, new Object) (*plan.OpQueue, error) {
	return plan.CreateOpQueue(old, new), nil
}

// Run is the default implementation of Config.Run(), and should be overwritten. This
// function just returns nil.
func (c DefaultConfig) Run(op *plan.Operation) error {
	return nil
}

// Validate is the default implementation of Config.Validate(), and can optionally be
// overwritten. This function just returns nil.
func (c DefaultConfig) Validate(op *plan.Operation) error {
	return nil
}

// Save is the default implementation of Config.Save(), and can optionally be overwritten.
// By default, this implementation saves the generated configuration to a local file.
func (c DefaultConfig) Save(v Object) error {
	return c.Store.WriteState(v)
}

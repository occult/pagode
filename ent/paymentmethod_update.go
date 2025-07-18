// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/occult/pagode/ent/paymentcustomer"
	"github.com/occult/pagode/ent/paymentmethod"
	"github.com/occult/pagode/ent/predicate"
)

// PaymentMethodUpdate is the builder for updating PaymentMethod entities.
type PaymentMethodUpdate struct {
	config
	hooks    []Hook
	mutation *PaymentMethodMutation
}

// Where appends a list predicates to the PaymentMethodUpdate builder.
func (pmu *PaymentMethodUpdate) Where(ps ...predicate.PaymentMethod) *PaymentMethodUpdate {
	pmu.mutation.Where(ps...)
	return pmu
}

// SetProviderPaymentMethodID sets the "provider_payment_method_id" field.
func (pmu *PaymentMethodUpdate) SetProviderPaymentMethodID(s string) *PaymentMethodUpdate {
	pmu.mutation.SetProviderPaymentMethodID(s)
	return pmu
}

// SetNillableProviderPaymentMethodID sets the "provider_payment_method_id" field if the given value is not nil.
func (pmu *PaymentMethodUpdate) SetNillableProviderPaymentMethodID(s *string) *PaymentMethodUpdate {
	if s != nil {
		pmu.SetProviderPaymentMethodID(*s)
	}
	return pmu
}

// SetProvider sets the "provider" field.
func (pmu *PaymentMethodUpdate) SetProvider(s string) *PaymentMethodUpdate {
	pmu.mutation.SetProvider(s)
	return pmu
}

// SetNillableProvider sets the "provider" field if the given value is not nil.
func (pmu *PaymentMethodUpdate) SetNillableProvider(s *string) *PaymentMethodUpdate {
	if s != nil {
		pmu.SetProvider(*s)
	}
	return pmu
}

// SetType sets the "type" field.
func (pmu *PaymentMethodUpdate) SetType(pa paymentmethod.Type) *PaymentMethodUpdate {
	pmu.mutation.SetType(pa)
	return pmu
}

// SetNillableType sets the "type" field if the given value is not nil.
func (pmu *PaymentMethodUpdate) SetNillableType(pa *paymentmethod.Type) *PaymentMethodUpdate {
	if pa != nil {
		pmu.SetType(*pa)
	}
	return pmu
}

// SetLastFour sets the "last_four" field.
func (pmu *PaymentMethodUpdate) SetLastFour(s string) *PaymentMethodUpdate {
	pmu.mutation.SetLastFour(s)
	return pmu
}

// SetNillableLastFour sets the "last_four" field if the given value is not nil.
func (pmu *PaymentMethodUpdate) SetNillableLastFour(s *string) *PaymentMethodUpdate {
	if s != nil {
		pmu.SetLastFour(*s)
	}
	return pmu
}

// ClearLastFour clears the value of the "last_four" field.
func (pmu *PaymentMethodUpdate) ClearLastFour() *PaymentMethodUpdate {
	pmu.mutation.ClearLastFour()
	return pmu
}

// SetBrand sets the "brand" field.
func (pmu *PaymentMethodUpdate) SetBrand(s string) *PaymentMethodUpdate {
	pmu.mutation.SetBrand(s)
	return pmu
}

// SetNillableBrand sets the "brand" field if the given value is not nil.
func (pmu *PaymentMethodUpdate) SetNillableBrand(s *string) *PaymentMethodUpdate {
	if s != nil {
		pmu.SetBrand(*s)
	}
	return pmu
}

// ClearBrand clears the value of the "brand" field.
func (pmu *PaymentMethodUpdate) ClearBrand() *PaymentMethodUpdate {
	pmu.mutation.ClearBrand()
	return pmu
}

// SetExpMonth sets the "exp_month" field.
func (pmu *PaymentMethodUpdate) SetExpMonth(i int) *PaymentMethodUpdate {
	pmu.mutation.ResetExpMonth()
	pmu.mutation.SetExpMonth(i)
	return pmu
}

// SetNillableExpMonth sets the "exp_month" field if the given value is not nil.
func (pmu *PaymentMethodUpdate) SetNillableExpMonth(i *int) *PaymentMethodUpdate {
	if i != nil {
		pmu.SetExpMonth(*i)
	}
	return pmu
}

// AddExpMonth adds i to the "exp_month" field.
func (pmu *PaymentMethodUpdate) AddExpMonth(i int) *PaymentMethodUpdate {
	pmu.mutation.AddExpMonth(i)
	return pmu
}

// ClearExpMonth clears the value of the "exp_month" field.
func (pmu *PaymentMethodUpdate) ClearExpMonth() *PaymentMethodUpdate {
	pmu.mutation.ClearExpMonth()
	return pmu
}

// SetExpYear sets the "exp_year" field.
func (pmu *PaymentMethodUpdate) SetExpYear(i int) *PaymentMethodUpdate {
	pmu.mutation.ResetExpYear()
	pmu.mutation.SetExpYear(i)
	return pmu
}

// SetNillableExpYear sets the "exp_year" field if the given value is not nil.
func (pmu *PaymentMethodUpdate) SetNillableExpYear(i *int) *PaymentMethodUpdate {
	if i != nil {
		pmu.SetExpYear(*i)
	}
	return pmu
}

// AddExpYear adds i to the "exp_year" field.
func (pmu *PaymentMethodUpdate) AddExpYear(i int) *PaymentMethodUpdate {
	pmu.mutation.AddExpYear(i)
	return pmu
}

// ClearExpYear clears the value of the "exp_year" field.
func (pmu *PaymentMethodUpdate) ClearExpYear() *PaymentMethodUpdate {
	pmu.mutation.ClearExpYear()
	return pmu
}

// SetIsDefault sets the "is_default" field.
func (pmu *PaymentMethodUpdate) SetIsDefault(b bool) *PaymentMethodUpdate {
	pmu.mutation.SetIsDefault(b)
	return pmu
}

// SetNillableIsDefault sets the "is_default" field if the given value is not nil.
func (pmu *PaymentMethodUpdate) SetNillableIsDefault(b *bool) *PaymentMethodUpdate {
	if b != nil {
		pmu.SetIsDefault(*b)
	}
	return pmu
}

// SetMetadata sets the "metadata" field.
func (pmu *PaymentMethodUpdate) SetMetadata(m map[string]interface{}) *PaymentMethodUpdate {
	pmu.mutation.SetMetadata(m)
	return pmu
}

// ClearMetadata clears the value of the "metadata" field.
func (pmu *PaymentMethodUpdate) ClearMetadata() *PaymentMethodUpdate {
	pmu.mutation.ClearMetadata()
	return pmu
}

// SetUpdatedAt sets the "updated_at" field.
func (pmu *PaymentMethodUpdate) SetUpdatedAt(t time.Time) *PaymentMethodUpdate {
	pmu.mutation.SetUpdatedAt(t)
	return pmu
}

// SetCustomerID sets the "customer" edge to the PaymentCustomer entity by ID.
func (pmu *PaymentMethodUpdate) SetCustomerID(id int) *PaymentMethodUpdate {
	pmu.mutation.SetCustomerID(id)
	return pmu
}

// SetCustomer sets the "customer" edge to the PaymentCustomer entity.
func (pmu *PaymentMethodUpdate) SetCustomer(p *PaymentCustomer) *PaymentMethodUpdate {
	return pmu.SetCustomerID(p.ID)
}

// Mutation returns the PaymentMethodMutation object of the builder.
func (pmu *PaymentMethodUpdate) Mutation() *PaymentMethodMutation {
	return pmu.mutation
}

// ClearCustomer clears the "customer" edge to the PaymentCustomer entity.
func (pmu *PaymentMethodUpdate) ClearCustomer() *PaymentMethodUpdate {
	pmu.mutation.ClearCustomer()
	return pmu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pmu *PaymentMethodUpdate) Save(ctx context.Context) (int, error) {
	pmu.defaults()
	return withHooks(ctx, pmu.sqlSave, pmu.mutation, pmu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pmu *PaymentMethodUpdate) SaveX(ctx context.Context) int {
	affected, err := pmu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pmu *PaymentMethodUpdate) Exec(ctx context.Context) error {
	_, err := pmu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pmu *PaymentMethodUpdate) ExecX(ctx context.Context) {
	if err := pmu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pmu *PaymentMethodUpdate) defaults() {
	if _, ok := pmu.mutation.UpdatedAt(); !ok {
		v := paymentmethod.UpdateDefaultUpdatedAt()
		pmu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pmu *PaymentMethodUpdate) check() error {
	if v, ok := pmu.mutation.ProviderPaymentMethodID(); ok {
		if err := paymentmethod.ProviderPaymentMethodIDValidator(v); err != nil {
			return &ValidationError{Name: "provider_payment_method_id", err: fmt.Errorf(`ent: validator failed for field "PaymentMethod.provider_payment_method_id": %w`, err)}
		}
	}
	if v, ok := pmu.mutation.Provider(); ok {
		if err := paymentmethod.ProviderValidator(v); err != nil {
			return &ValidationError{Name: "provider", err: fmt.Errorf(`ent: validator failed for field "PaymentMethod.provider": %w`, err)}
		}
	}
	if v, ok := pmu.mutation.GetType(); ok {
		if err := paymentmethod.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "PaymentMethod.type": %w`, err)}
		}
	}
	if v, ok := pmu.mutation.ExpMonth(); ok {
		if err := paymentmethod.ExpMonthValidator(v); err != nil {
			return &ValidationError{Name: "exp_month", err: fmt.Errorf(`ent: validator failed for field "PaymentMethod.exp_month": %w`, err)}
		}
	}
	if pmu.mutation.CustomerCleared() && len(pmu.mutation.CustomerIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "PaymentMethod.customer"`)
	}
	return nil
}

func (pmu *PaymentMethodUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := pmu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(paymentmethod.Table, paymentmethod.Columns, sqlgraph.NewFieldSpec(paymentmethod.FieldID, field.TypeInt))
	if ps := pmu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pmu.mutation.ProviderPaymentMethodID(); ok {
		_spec.SetField(paymentmethod.FieldProviderPaymentMethodID, field.TypeString, value)
	}
	if value, ok := pmu.mutation.Provider(); ok {
		_spec.SetField(paymentmethod.FieldProvider, field.TypeString, value)
	}
	if value, ok := pmu.mutation.GetType(); ok {
		_spec.SetField(paymentmethod.FieldType, field.TypeEnum, value)
	}
	if value, ok := pmu.mutation.LastFour(); ok {
		_spec.SetField(paymentmethod.FieldLastFour, field.TypeString, value)
	}
	if pmu.mutation.LastFourCleared() {
		_spec.ClearField(paymentmethod.FieldLastFour, field.TypeString)
	}
	if value, ok := pmu.mutation.Brand(); ok {
		_spec.SetField(paymentmethod.FieldBrand, field.TypeString, value)
	}
	if pmu.mutation.BrandCleared() {
		_spec.ClearField(paymentmethod.FieldBrand, field.TypeString)
	}
	if value, ok := pmu.mutation.ExpMonth(); ok {
		_spec.SetField(paymentmethod.FieldExpMonth, field.TypeInt, value)
	}
	if value, ok := pmu.mutation.AddedExpMonth(); ok {
		_spec.AddField(paymentmethod.FieldExpMonth, field.TypeInt, value)
	}
	if pmu.mutation.ExpMonthCleared() {
		_spec.ClearField(paymentmethod.FieldExpMonth, field.TypeInt)
	}
	if value, ok := pmu.mutation.ExpYear(); ok {
		_spec.SetField(paymentmethod.FieldExpYear, field.TypeInt, value)
	}
	if value, ok := pmu.mutation.AddedExpYear(); ok {
		_spec.AddField(paymentmethod.FieldExpYear, field.TypeInt, value)
	}
	if pmu.mutation.ExpYearCleared() {
		_spec.ClearField(paymentmethod.FieldExpYear, field.TypeInt)
	}
	if value, ok := pmu.mutation.IsDefault(); ok {
		_spec.SetField(paymentmethod.FieldIsDefault, field.TypeBool, value)
	}
	if value, ok := pmu.mutation.Metadata(); ok {
		_spec.SetField(paymentmethod.FieldMetadata, field.TypeJSON, value)
	}
	if pmu.mutation.MetadataCleared() {
		_spec.ClearField(paymentmethod.FieldMetadata, field.TypeJSON)
	}
	if value, ok := pmu.mutation.UpdatedAt(); ok {
		_spec.SetField(paymentmethod.FieldUpdatedAt, field.TypeTime, value)
	}
	if pmu.mutation.CustomerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   paymentmethod.CustomerTable,
			Columns: []string{paymentmethod.CustomerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(paymentcustomer.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pmu.mutation.CustomerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   paymentmethod.CustomerTable,
			Columns: []string{paymentmethod.CustomerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(paymentcustomer.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, pmu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{paymentmethod.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	pmu.mutation.done = true
	return n, nil
}

// PaymentMethodUpdateOne is the builder for updating a single PaymentMethod entity.
type PaymentMethodUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *PaymentMethodMutation
}

// SetProviderPaymentMethodID sets the "provider_payment_method_id" field.
func (pmuo *PaymentMethodUpdateOne) SetProviderPaymentMethodID(s string) *PaymentMethodUpdateOne {
	pmuo.mutation.SetProviderPaymentMethodID(s)
	return pmuo
}

// SetNillableProviderPaymentMethodID sets the "provider_payment_method_id" field if the given value is not nil.
func (pmuo *PaymentMethodUpdateOne) SetNillableProviderPaymentMethodID(s *string) *PaymentMethodUpdateOne {
	if s != nil {
		pmuo.SetProviderPaymentMethodID(*s)
	}
	return pmuo
}

// SetProvider sets the "provider" field.
func (pmuo *PaymentMethodUpdateOne) SetProvider(s string) *PaymentMethodUpdateOne {
	pmuo.mutation.SetProvider(s)
	return pmuo
}

// SetNillableProvider sets the "provider" field if the given value is not nil.
func (pmuo *PaymentMethodUpdateOne) SetNillableProvider(s *string) *PaymentMethodUpdateOne {
	if s != nil {
		pmuo.SetProvider(*s)
	}
	return pmuo
}

// SetType sets the "type" field.
func (pmuo *PaymentMethodUpdateOne) SetType(pa paymentmethod.Type) *PaymentMethodUpdateOne {
	pmuo.mutation.SetType(pa)
	return pmuo
}

// SetNillableType sets the "type" field if the given value is not nil.
func (pmuo *PaymentMethodUpdateOne) SetNillableType(pa *paymentmethod.Type) *PaymentMethodUpdateOne {
	if pa != nil {
		pmuo.SetType(*pa)
	}
	return pmuo
}

// SetLastFour sets the "last_four" field.
func (pmuo *PaymentMethodUpdateOne) SetLastFour(s string) *PaymentMethodUpdateOne {
	pmuo.mutation.SetLastFour(s)
	return pmuo
}

// SetNillableLastFour sets the "last_four" field if the given value is not nil.
func (pmuo *PaymentMethodUpdateOne) SetNillableLastFour(s *string) *PaymentMethodUpdateOne {
	if s != nil {
		pmuo.SetLastFour(*s)
	}
	return pmuo
}

// ClearLastFour clears the value of the "last_four" field.
func (pmuo *PaymentMethodUpdateOne) ClearLastFour() *PaymentMethodUpdateOne {
	pmuo.mutation.ClearLastFour()
	return pmuo
}

// SetBrand sets the "brand" field.
func (pmuo *PaymentMethodUpdateOne) SetBrand(s string) *PaymentMethodUpdateOne {
	pmuo.mutation.SetBrand(s)
	return pmuo
}

// SetNillableBrand sets the "brand" field if the given value is not nil.
func (pmuo *PaymentMethodUpdateOne) SetNillableBrand(s *string) *PaymentMethodUpdateOne {
	if s != nil {
		pmuo.SetBrand(*s)
	}
	return pmuo
}

// ClearBrand clears the value of the "brand" field.
func (pmuo *PaymentMethodUpdateOne) ClearBrand() *PaymentMethodUpdateOne {
	pmuo.mutation.ClearBrand()
	return pmuo
}

// SetExpMonth sets the "exp_month" field.
func (pmuo *PaymentMethodUpdateOne) SetExpMonth(i int) *PaymentMethodUpdateOne {
	pmuo.mutation.ResetExpMonth()
	pmuo.mutation.SetExpMonth(i)
	return pmuo
}

// SetNillableExpMonth sets the "exp_month" field if the given value is not nil.
func (pmuo *PaymentMethodUpdateOne) SetNillableExpMonth(i *int) *PaymentMethodUpdateOne {
	if i != nil {
		pmuo.SetExpMonth(*i)
	}
	return pmuo
}

// AddExpMonth adds i to the "exp_month" field.
func (pmuo *PaymentMethodUpdateOne) AddExpMonth(i int) *PaymentMethodUpdateOne {
	pmuo.mutation.AddExpMonth(i)
	return pmuo
}

// ClearExpMonth clears the value of the "exp_month" field.
func (pmuo *PaymentMethodUpdateOne) ClearExpMonth() *PaymentMethodUpdateOne {
	pmuo.mutation.ClearExpMonth()
	return pmuo
}

// SetExpYear sets the "exp_year" field.
func (pmuo *PaymentMethodUpdateOne) SetExpYear(i int) *PaymentMethodUpdateOne {
	pmuo.mutation.ResetExpYear()
	pmuo.mutation.SetExpYear(i)
	return pmuo
}

// SetNillableExpYear sets the "exp_year" field if the given value is not nil.
func (pmuo *PaymentMethodUpdateOne) SetNillableExpYear(i *int) *PaymentMethodUpdateOne {
	if i != nil {
		pmuo.SetExpYear(*i)
	}
	return pmuo
}

// AddExpYear adds i to the "exp_year" field.
func (pmuo *PaymentMethodUpdateOne) AddExpYear(i int) *PaymentMethodUpdateOne {
	pmuo.mutation.AddExpYear(i)
	return pmuo
}

// ClearExpYear clears the value of the "exp_year" field.
func (pmuo *PaymentMethodUpdateOne) ClearExpYear() *PaymentMethodUpdateOne {
	pmuo.mutation.ClearExpYear()
	return pmuo
}

// SetIsDefault sets the "is_default" field.
func (pmuo *PaymentMethodUpdateOne) SetIsDefault(b bool) *PaymentMethodUpdateOne {
	pmuo.mutation.SetIsDefault(b)
	return pmuo
}

// SetNillableIsDefault sets the "is_default" field if the given value is not nil.
func (pmuo *PaymentMethodUpdateOne) SetNillableIsDefault(b *bool) *PaymentMethodUpdateOne {
	if b != nil {
		pmuo.SetIsDefault(*b)
	}
	return pmuo
}

// SetMetadata sets the "metadata" field.
func (pmuo *PaymentMethodUpdateOne) SetMetadata(m map[string]interface{}) *PaymentMethodUpdateOne {
	pmuo.mutation.SetMetadata(m)
	return pmuo
}

// ClearMetadata clears the value of the "metadata" field.
func (pmuo *PaymentMethodUpdateOne) ClearMetadata() *PaymentMethodUpdateOne {
	pmuo.mutation.ClearMetadata()
	return pmuo
}

// SetUpdatedAt sets the "updated_at" field.
func (pmuo *PaymentMethodUpdateOne) SetUpdatedAt(t time.Time) *PaymentMethodUpdateOne {
	pmuo.mutation.SetUpdatedAt(t)
	return pmuo
}

// SetCustomerID sets the "customer" edge to the PaymentCustomer entity by ID.
func (pmuo *PaymentMethodUpdateOne) SetCustomerID(id int) *PaymentMethodUpdateOne {
	pmuo.mutation.SetCustomerID(id)
	return pmuo
}

// SetCustomer sets the "customer" edge to the PaymentCustomer entity.
func (pmuo *PaymentMethodUpdateOne) SetCustomer(p *PaymentCustomer) *PaymentMethodUpdateOne {
	return pmuo.SetCustomerID(p.ID)
}

// Mutation returns the PaymentMethodMutation object of the builder.
func (pmuo *PaymentMethodUpdateOne) Mutation() *PaymentMethodMutation {
	return pmuo.mutation
}

// ClearCustomer clears the "customer" edge to the PaymentCustomer entity.
func (pmuo *PaymentMethodUpdateOne) ClearCustomer() *PaymentMethodUpdateOne {
	pmuo.mutation.ClearCustomer()
	return pmuo
}

// Where appends a list predicates to the PaymentMethodUpdate builder.
func (pmuo *PaymentMethodUpdateOne) Where(ps ...predicate.PaymentMethod) *PaymentMethodUpdateOne {
	pmuo.mutation.Where(ps...)
	return pmuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (pmuo *PaymentMethodUpdateOne) Select(field string, fields ...string) *PaymentMethodUpdateOne {
	pmuo.fields = append([]string{field}, fields...)
	return pmuo
}

// Save executes the query and returns the updated PaymentMethod entity.
func (pmuo *PaymentMethodUpdateOne) Save(ctx context.Context) (*PaymentMethod, error) {
	pmuo.defaults()
	return withHooks(ctx, pmuo.sqlSave, pmuo.mutation, pmuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pmuo *PaymentMethodUpdateOne) SaveX(ctx context.Context) *PaymentMethod {
	node, err := pmuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (pmuo *PaymentMethodUpdateOne) Exec(ctx context.Context) error {
	_, err := pmuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pmuo *PaymentMethodUpdateOne) ExecX(ctx context.Context) {
	if err := pmuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pmuo *PaymentMethodUpdateOne) defaults() {
	if _, ok := pmuo.mutation.UpdatedAt(); !ok {
		v := paymentmethod.UpdateDefaultUpdatedAt()
		pmuo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pmuo *PaymentMethodUpdateOne) check() error {
	if v, ok := pmuo.mutation.ProviderPaymentMethodID(); ok {
		if err := paymentmethod.ProviderPaymentMethodIDValidator(v); err != nil {
			return &ValidationError{Name: "provider_payment_method_id", err: fmt.Errorf(`ent: validator failed for field "PaymentMethod.provider_payment_method_id": %w`, err)}
		}
	}
	if v, ok := pmuo.mutation.Provider(); ok {
		if err := paymentmethod.ProviderValidator(v); err != nil {
			return &ValidationError{Name: "provider", err: fmt.Errorf(`ent: validator failed for field "PaymentMethod.provider": %w`, err)}
		}
	}
	if v, ok := pmuo.mutation.GetType(); ok {
		if err := paymentmethod.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "PaymentMethod.type": %w`, err)}
		}
	}
	if v, ok := pmuo.mutation.ExpMonth(); ok {
		if err := paymentmethod.ExpMonthValidator(v); err != nil {
			return &ValidationError{Name: "exp_month", err: fmt.Errorf(`ent: validator failed for field "PaymentMethod.exp_month": %w`, err)}
		}
	}
	if pmuo.mutation.CustomerCleared() && len(pmuo.mutation.CustomerIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "PaymentMethod.customer"`)
	}
	return nil
}

func (pmuo *PaymentMethodUpdateOne) sqlSave(ctx context.Context) (_node *PaymentMethod, err error) {
	if err := pmuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(paymentmethod.Table, paymentmethod.Columns, sqlgraph.NewFieldSpec(paymentmethod.FieldID, field.TypeInt))
	id, ok := pmuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "PaymentMethod.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := pmuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, paymentmethod.FieldID)
		for _, f := range fields {
			if !paymentmethod.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != paymentmethod.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := pmuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pmuo.mutation.ProviderPaymentMethodID(); ok {
		_spec.SetField(paymentmethod.FieldProviderPaymentMethodID, field.TypeString, value)
	}
	if value, ok := pmuo.mutation.Provider(); ok {
		_spec.SetField(paymentmethod.FieldProvider, field.TypeString, value)
	}
	if value, ok := pmuo.mutation.GetType(); ok {
		_spec.SetField(paymentmethod.FieldType, field.TypeEnum, value)
	}
	if value, ok := pmuo.mutation.LastFour(); ok {
		_spec.SetField(paymentmethod.FieldLastFour, field.TypeString, value)
	}
	if pmuo.mutation.LastFourCleared() {
		_spec.ClearField(paymentmethod.FieldLastFour, field.TypeString)
	}
	if value, ok := pmuo.mutation.Brand(); ok {
		_spec.SetField(paymentmethod.FieldBrand, field.TypeString, value)
	}
	if pmuo.mutation.BrandCleared() {
		_spec.ClearField(paymentmethod.FieldBrand, field.TypeString)
	}
	if value, ok := pmuo.mutation.ExpMonth(); ok {
		_spec.SetField(paymentmethod.FieldExpMonth, field.TypeInt, value)
	}
	if value, ok := pmuo.mutation.AddedExpMonth(); ok {
		_spec.AddField(paymentmethod.FieldExpMonth, field.TypeInt, value)
	}
	if pmuo.mutation.ExpMonthCleared() {
		_spec.ClearField(paymentmethod.FieldExpMonth, field.TypeInt)
	}
	if value, ok := pmuo.mutation.ExpYear(); ok {
		_spec.SetField(paymentmethod.FieldExpYear, field.TypeInt, value)
	}
	if value, ok := pmuo.mutation.AddedExpYear(); ok {
		_spec.AddField(paymentmethod.FieldExpYear, field.TypeInt, value)
	}
	if pmuo.mutation.ExpYearCleared() {
		_spec.ClearField(paymentmethod.FieldExpYear, field.TypeInt)
	}
	if value, ok := pmuo.mutation.IsDefault(); ok {
		_spec.SetField(paymentmethod.FieldIsDefault, field.TypeBool, value)
	}
	if value, ok := pmuo.mutation.Metadata(); ok {
		_spec.SetField(paymentmethod.FieldMetadata, field.TypeJSON, value)
	}
	if pmuo.mutation.MetadataCleared() {
		_spec.ClearField(paymentmethod.FieldMetadata, field.TypeJSON)
	}
	if value, ok := pmuo.mutation.UpdatedAt(); ok {
		_spec.SetField(paymentmethod.FieldUpdatedAt, field.TypeTime, value)
	}
	if pmuo.mutation.CustomerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   paymentmethod.CustomerTable,
			Columns: []string{paymentmethod.CustomerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(paymentcustomer.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pmuo.mutation.CustomerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   paymentmethod.CustomerTable,
			Columns: []string{paymentmethod.CustomerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(paymentcustomer.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &PaymentMethod{config: pmuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, pmuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{paymentmethod.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	pmuo.mutation.done = true
	return _node, nil
}

package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	ApigwRoute struct {
		svc routeService
		ac  apiGwRouteAccessController
	}

	routePayload struct {
		*types.ApigwRoute

		CanGrant            bool `json:"canGrant"`
		CanUpdateApigwRoute bool `json:"canUpdateApigwRoute"`
		CanDeleteApigwRoute bool `json:"canDeleteApigwRoute"`
	}

	routeSetPayload struct {
		Filter types.ApigwRouteFilter `json:"filter"`
		Set    []*routePayload        `json:"set"`
	}

	routeService interface {
		FindByID(ctx context.Context, ID uint64) (*types.ApigwRoute, error)
		Create(ctx context.Context, new *types.ApigwRoute) (*types.ApigwRoute, error)
		Update(ctx context.Context, upd *types.ApigwRoute) (*types.ApigwRoute, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error
		Search(ctx context.Context, filter types.ApigwRouteFilter) (types.ApigwRouteSet, types.ApigwRouteFilter, error)
	}

	apiGwRouteAccessController interface {
		CanGrant(context.Context) bool

		CanCreateApigwRoute(context.Context) bool
		CanUpdateApigwRoute(context.Context, *types.ApigwRoute) bool
		CanDeleteApigwRoute(context.Context, *types.ApigwRoute) bool
	}
)

func (ApigwRoute) New() *ApigwRoute {
	return &ApigwRoute{
		svc: service.DefaultApigwRoute,
		ac:  service.DefaultAccessControl,
	}
}

func (ctrl *ApigwRoute) List(ctx context.Context, r *request.ApigwRouteList) (interface{}, error) {
	var (
		err error
		f   = types.ApigwRouteFilter{
			// todo: this should renamed to r.Endpoint after UI is aligned with this
			Endpoint: r.Query,
			Deleted:  filter.State(r.Deleted),

			// todo: this should dynamic as Delete
			//		but making it default to `1`, until UI is aligned with this
			Disabled: filter.StateInclusive,
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, f, err := ctrl.svc.Search(ctx, f)

	return ctrl.makeFilterPayload(ctx, set, f, err)
}

func (ctrl *ApigwRoute) Create(ctx context.Context, r *request.ApigwRouteCreate) (interface{}, error) {
	var (
		err error
		q   = &types.ApigwRoute{
			Endpoint: r.Endpoint,
			Method:   r.Method,
			Enabled:  r.Enabled,
		}
	)

	q, err = ctrl.svc.Create(ctx, q)

	return ctrl.makePayload(ctx, q, err)
}

func (ctrl *ApigwRoute) Read(ctx context.Context, r *request.ApigwRouteRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.RouteID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ApigwRoute) Update(ctx context.Context, r *request.ApigwRouteUpdate) (interface{}, error) {
	var (
		err error
		q   = &types.ApigwRoute{
			ID:       r.RouteID,
			Endpoint: r.Endpoint,
			Method:   r.Method,
			Group:    r.Group,
			Enabled:  r.Enabled,
		}
	)

	q, err = ctrl.svc.Update(ctx, q)

	return ctrl.makePayload(ctx, q, err)
}

func (ctrl *ApigwRoute) Delete(ctx context.Context, r *request.ApigwRouteDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.RouteID)
}

func (ctrl *ApigwRoute) Undelete(ctx context.Context, r *request.ApigwRouteUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.RouteID)
}

func (ctrl *ApigwRoute) makePayload(ctx context.Context, q *types.ApigwRoute, err error) (*routePayload, error) {
	if err != nil || q == nil {
		return nil, err
	}

	qq := &routePayload{
		ApigwRoute: q,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateApigwRoute: ctrl.ac.CanUpdateApigwRoute(ctx, q),
		CanDeleteApigwRoute: ctrl.ac.CanDeleteApigwRoute(ctx, q),
	}

	return qq, nil
}

func (ctrl *ApigwRoute) makeFilterPayload(ctx context.Context, nn types.ApigwRouteSet, f types.ApigwRouteFilter, err error) (*routeSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &routeSetPayload{Filter: f, Set: make([]*routePayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}

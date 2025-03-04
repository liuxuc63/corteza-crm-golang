package request

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/go-chi/chi/v5"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// dummy vars to prevent
// unused imports complain
var (
	_ = chi.URLParam
	_ = multipart.ErrMessageTooLarge
	_ = payload.ParseUint64s
	_ = strings.ToLower
	_ = io.EOF
	_ = fmt.Errorf
	_ = json.NewEncoder
)

type (
	// Internal API interface
	DalConnectionList struct {
		// ConnectionID GET parameter
		//
		// Filter by connection ID
		ConnectionID []string

		// Handle GET parameter
		//
		// Search handle to match against connections
		Handle string

		// Type GET parameter
		//
		// Search type to match against connections
		Type string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted connections
		Deleted uint

		// IncTotal GET parameter
		//
		// Include total counter
		IncTotal bool
	}

	DalConnectionCreate struct {
		// Handle POST parameter
		//
		//
		Handle string

		// Type POST parameter
		//
		//
		Type string

		// Meta POST parameter
		//
		//
		Meta types.ConnectionMeta

		// Config POST parameter
		//
		//
		Config types.ConnectionConfig
	}

	DalConnectionUpdate struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`

		// Handle POST parameter
		//
		//
		Handle string

		// Type POST parameter
		//
		//
		Type string

		// Meta POST parameter
		//
		//
		Meta types.ConnectionMeta

		// Config POST parameter
		//
		//
		Config types.ConnectionConfig
	}

	DalConnectionRead struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`
	}

	DalConnectionDelete struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`
	}

	DalConnectionUndelete struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`
	}
)

// NewDalConnectionList request
func NewDalConnectionList() *DalConnectionList {
	return &DalConnectionList{}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
		"handle":       r.Handle,
		"type":         r.Type,
		"deleted":      r.Deleted,
		"incTotal":     r.IncTotal,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) GetConnectionID() []string {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) GetType() string {
	return r.Type
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) GetIncTotal() bool {
	return r.IncTotal
}

// Fill processes request and fills internal variables
func (r *DalConnectionList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["connectionID[]"]; ok {
			r.ConnectionID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["connectionID"]; ok {
			r.ConnectionID, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["type"]; ok && len(val) > 0 {
			r.Type, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["deleted"]; ok && len(val) > 0 {
			r.Deleted, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["incTotal"]; ok && len(val) > 0 {
			r.IncTotal, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewDalConnectionCreate request
func NewDalConnectionCreate() *DalConnectionCreate {
	return &DalConnectionCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle": r.Handle,
		"type":   r.Type,
		"meta":   r.Meta,
		"config": r.Config,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetType() string {
	return r.Type
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetMeta() types.ConnectionMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetConfig() types.ConnectionConfig {
	return r.Config
}

// Fill processes request and fills internal variables
func (r *DalConnectionCreate) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["type"]; ok && len(val) > 0 {
				r.Type, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseConnectionMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseConnectionMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["config[]"]; ok {
				r.Config, err = types.ParseConnectionConfig(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["config"]; ok {
				r.Config, err = types.ParseConnectionConfig(val)
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["type"]; ok && len(val) > 0 {
			r.Type, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseConnectionMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseConnectionMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["config[]"]; ok {
			r.Config, err = types.ParseConnectionConfig(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["config"]; ok {
			r.Config, err = types.ParseConnectionConfig(val)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewDalConnectionUpdate request
func NewDalConnectionUpdate() *DalConnectionUpdate {
	return &DalConnectionUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
		"handle":       r.Handle,
		"type":         r.Type,
		"meta":         r.Meta,
		"config":       r.Config,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetType() string {
	return r.Type
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetMeta() types.ConnectionMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetConfig() types.ConnectionConfig {
	return r.Config
}

// Fill processes request and fills internal variables
func (r *DalConnectionUpdate) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["type"]; ok && len(val) > 0 {
				r.Type, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseConnectionMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseConnectionMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["config[]"]; ok {
				r.Config, err = types.ParseConnectionConfig(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["config"]; ok {
				r.Config, err = types.ParseConnectionConfig(val)
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["type"]; ok && len(val) > 0 {
			r.Type, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseConnectionMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseConnectionMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["config[]"]; ok {
			r.Config, err = types.ParseConnectionConfig(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["config"]; ok {
			r.Config, err = types.ParseConnectionConfig(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "connectionID")
		r.ConnectionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDalConnectionRead request
func NewDalConnectionRead() *DalConnectionRead {
	return &DalConnectionRead{}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionRead) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Fill processes request and fills internal variables
func (r *DalConnectionRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "connectionID")
		r.ConnectionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDalConnectionDelete request
func NewDalConnectionDelete() *DalConnectionDelete {
	return &DalConnectionDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionDelete) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Fill processes request and fills internal variables
func (r *DalConnectionDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "connectionID")
		r.ConnectionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDalConnectionUndelete request
func NewDalConnectionUndelete() *DalConnectionUndelete {
	return &DalConnectionUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUndelete) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Fill processes request and fills internal variables
func (r *DalConnectionUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "connectionID")
		r.ConnectionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

package person

import "github.com/omatheuscaetano/planus-api/internal/auth"

const (
    ActionPending  auth.Action = "create"
    ActionUpdating auth.Action = "update"
    ActionDeleting auth.Action = "delete"
)



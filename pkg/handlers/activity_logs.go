package handlers

import (
	"net/http"

	"github.com/AshimKoirala/load-balancer-admin/pkg/db"
	"github.com/AshimKoirala/load-balancer-admin/utils"
)

func GetActivityLogs(w http.ResponseWriter, r *http.Request) {
    logs, err := db.FetchActivityLogs(r.Context()) // Fetch logs from the database
    if err != nil {
        utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Failed to fetch activity logs"})
        return
    }

    if len(logs) == 0 {
        utils.NewSuccessResponse(w, "No activity logs found")
        return
    }

    utils.NewSuccessResponseWithData(w, logs)
}

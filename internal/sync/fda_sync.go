package sync

import (
	"recall-app/cmd/api"
)

func RunSync(app *api.Application) error {
	lastDate, err := app.Handlers.Fda.GetLastSyncedDate("food")
	if err != nil {
		return err
	}

	latestSeen, err := app.Handlers.Fda.FetchAndStoreRecallsSince(lastDate)
	if err != nil {
		return err
	}
	print(latestSeen.GoString())

	return nil
}

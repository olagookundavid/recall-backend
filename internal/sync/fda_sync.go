package sync

import (
	"recall-app/cmd/api"
)

func RunSync(app *api.Application) error {
	///Get last date
	// lastDate, err := app.Handlers.Fda.GetLastSyncedDate("food")
	// if err != nil {
	// 	return err
	// }

	///Idea here is to grab each tracked product, then incrementally call fda function, it handles the rest, the function uses goroutines btw
	// latestSeen, err := app.Handlers.Sync.FetchAndStoreRecallsSince(lastDate)
	// if err != nil {
	// 	return err
	// }

	///update last date
	// print(latestSeen.GoString())

	return nil
}

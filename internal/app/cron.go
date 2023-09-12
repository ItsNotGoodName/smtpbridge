package app

// import (
// 	"context"
// 	"time"
// )

// type Cron struct {
// 	db        database.Querier
// 	fileStore FileStore
// 	bus       core.Bus
// 	config    *models.Config
// }
//
// func (c Cron) Clean(ctx context.Context) {
// 	// trimmerDeleteByAge(cc, cc.Config.RetentionPolicy)
// 	// trimmerDeleteOrphanAttachments(cc)
//
// 	// storage, err := StorageGet(cc)
// 	// if err != nil {
// 	// 	log.Err(err).Msg("Failed to get storage")
// 	// 	return
// 	// }
// 	//
// 	// trimmerDeleteByEnvelopeCount(cc, cc.Config.RetentionPolicy, storage)
// 	// trimmerDeleteByAttachmentSize(cc, cc.Config.RetentionPolicy, storage)
// }
//
// func (c Cron) Serve(ctx context.Context) error {
// 	_ = time.NewTicker(30 * time.Minute)
//
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return nil
// 			// case <-envCreatedC:
// 			// 	storage, err := StorageGet(cc)
// 			// 	if err != nil {
// 			// 		log.Err(err).Msg("Failed to get storage")
// 			// 		continue
// 			// 	}
// 			//
// 			// 	trimmerDeleteByEnvelopeCount(cc, cc.Config.RetentionPolicy, storage)
// 			// 	trimmerDeleteByAttachmentSize(cc, cc.Config.RetentionPolicy, storage)
// 			// case <-envDeletedC:
// 			// 	trimmerDeleteOrphanAttachments(cc)
// 			// case <-ticker.C:
// 			// 	clean()
// 			// case evt := <-evtTrimStart:
// 			// 	clean()
// 			//
// 			// 	select {
// 			// 	case <-ctx.Done():
// 			// 		return
// 			// 	case evt.Response <- true:
// 			// 	}
// 		}
// 	}
// }

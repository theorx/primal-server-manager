package Scheduler

type TriggerLog struct {
	log []WipeTrigger
}

func (t *TriggerLog) Log(wt *WipeTrigger) {
	//Todo: Implement real persistent storage
	t.log = append(t.log, *wt)
}

func (t *TriggerLog) Get(start int64, end int64, limit int64) []*WipeTrigger {
	//query the data
	return nil
}

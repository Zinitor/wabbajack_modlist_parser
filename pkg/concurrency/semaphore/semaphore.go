package semaphore

type New chan struct{}

func (s New) Acquire() {
	s <- struct{}{}
}
func (s New) Release() {
	<-s
}

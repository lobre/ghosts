package main

type listener struct {
	docker     docker
	processors []processor
}

func newListener(docker docker) listener {
	return listener{docker, nil}
}

func (l *listener) addProcessor(p processor) {
	l.processors = append(l.processors, p)
}

func (l listener) start(stop <-chan int) error {
	// Init phase
	for _, p := range l.processors {
		if err := p.init(); err != nil {
			return err
		}
	}

	c, cerr := l.docker.listenContainers()
	for {
		select {
		case <-stop:
			return nil
		case err := <-cerr:
			return err
		case event := <-c:
			switch event.Action {
			case "start":
				for _, p := range l.processors {
					if err := p.startEvent(event.ID); err != nil {
						return err
					}
				}
			case "die":
				for _, p := range l.processors {
					if err := p.dieEvent(event.ID); err != nil {
						return err
					}
				}
			}
		}
	}
}

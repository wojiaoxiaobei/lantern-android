package client

import (
	"github.com/getlantern/balancer"
	"log"
)

func (client *Client) getBalancer() *balancer.Balancer {
	bal := <-client.balCh
	client.balCh <- bal
	return bal
}

func (client *Client) initBalancer() *balancer.Balancer {
	dialers := make([]*balancer.Dialer, 0, len(client.frontedServers))

	for _, s := range client.frontedServers {
		dialer := s.dialer()
		dialers = append(dialers, dialer)
	}

	bal := balancer.New(dialers...)

	if client.balInitialized {
		log.Printf("Draining balancer channel.")
		old := <-client.balCh
		// Close old balancer on a goroutine to avoid blocking here
		go func() {
			old.Close()
			log.Printf("Closed old balancer.")
		}()
	} else {
		log.Printf("Creating balancer channel.")
		client.balCh = make(chan *balancer.Balancer, 1)
	}

	log.Printf("Publishing balancer.")

	client.balInitialized = true
	client.balCh <- bal

	return bal
}

package cmd

import (
	"fmt"

	"github.com/Shopify/sarama"
)

type (
	// DeleteConsumerGroups obtains the kafka admin clients from provided broker strings
	DeleteConsumerGroups struct {
		BrokerList []string `help:"Kafka Brokers Sasl list for the cluster"`
		DryRun     bool     `help:"Run in dummy mode to show what would happen"`
	}
)

// Run use the source and target broker lists to migrate topic definitions from one cluster to another
func (g *DeleteConsumerGroups) Run(globals *Globals) error {

	config := getSaramaConfig(globals.SourceAPIKey)
	admin, err := sarama.NewClusterAdmin(g.BrokerList, config)
	if err != nil {
		return err
	}

	// List consumergroups on the source cluster
	consumergroups, err := admin.ListConsumerGroups()
	if err != nil {
		return err
	}

	// Delete topics on target cluster
	for name := range consumergroups {
		if !g.DryRun {
			err := admin.DeleteConsumerGroup(name)
			if err != nil {
				fmt.Printf("Unable to delete consumer group %s with err %s\n", name, err.Error())
				continue
			}
		}
		fmt.Printf("Deleted consumer group %s\n", name)
	}

	return nil
}

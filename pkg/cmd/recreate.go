package cmd

import (
	"fmt"

	"github.com/Shopify/sarama"
)

type (
	// Recreatae obtains the kafka admin clients from provided broker strings
	Recreate struct {
		BrokerList []string `help:"Kafka Brokers Sasl list for the cluster"`
		DryRun     bool     `help:"Run in dummy mode to show what would happen"`
	}
)

// Run use the source and target broker lists to migrate topic definitions from one cluster to another
func (g *Recreate) Run(globals *Globals) error {

	config := getSaramaConfig(globals.SourceAPIKey)
	admin, err := sarama.NewClusterAdmin(g.BrokerList, config)
	if err != nil {
		return err
	}

	// ListTopics on the source cluster
	topics, err := admin.ListTopics()
	if err != nil {
		return err
	}

	// Delete topics on target cluster
	for name, _ := range topics {
		if !g.DryRun {
			if name == "__consumer_offsets" {
				continue
			}
			err := admin.DeleteTopic(name)
			if err != nil {
				fmt.Printf("Unable to delete topic %s with err %s\n", name, err.Error())
				continue
			}
		}
		fmt.Printf("Deleted topic %s\n", name)
	}
	// Create topics on target cluster
	for name, topic := range topics {
		if name == "__consumer_offsets" {
			continue
		}
		configEntries := map[string]*string{}
		for k, v := range topic.ConfigEntries {
			switch k {
			case "cleanup.policy", "retention.bytes", "retention.ms":
				configEntries[k] = v
			}
		}
		detail := sarama.TopicDetail{
			NumPartitions:     topic.NumPartitions,
			ReplicationFactor: topic.ReplicationFactor,
			ConfigEntries:     configEntries,
		}
		if !g.DryRun {
			err := admin.CreateTopic(name, &detail, false)
			if err != nil {
				fmt.Printf("Unable to create topic %s with err %s\n", name, err.Error())
				continue
			}
		}
		fmt.Printf("Creating topic %s with %d partitions and config parameters {", name, detail.NumPartitions)
		for k, v := range detail.ConfigEntries {
			fmt.Printf(" %s=%s", k, *v)
		}
		fmt.Println(" }")
	}

	return nil
}

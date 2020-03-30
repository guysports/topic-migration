package cmd

import (
	"crypto/tls"
	"fmt"

	"github.com/Shopify/sarama"
)

type (
	// Migrate obtains the kafka admin clients from provided broker strings
	Migrate struct {
		SourceBrokerList []string `help:"Kafka Brokers Sasl list for the source cluster"`
		TargetBrokerList []string `help:"Kafka Brokers Sasl list for the target cluster"`
		DryRun           bool     `help:"Run in dummy mode to show what would happen"`
	}
)

func getSaramaConfig(key string) *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V1_1_0_0
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = &tls.Config{
		// nolint
		InsecureSkipVerify: true,
	}
	config.Net.SASL.Enable = true
	config.Net.SASL.User = "token"
	config.Net.SASL.Password = key

	return config
}

// Run use the source and target broker lists to migrate topic definitions from one cluster to another
func (g *Migrate) Run(globals *Globals) error {

	srcConfig := getSaramaConfig(globals.SourceAPIKey)
	srcAdmin, err := sarama.NewClusterAdmin(g.SourceBrokerList, srcConfig)
	if err != nil {
		return err
	}

	// ListTopics on the source cluster
	topics, err := srcAdmin.ListTopics()
	if err != nil {
		return err
	}

	tgtConfig := getSaramaConfig(globals.TargetAPIKey)
	tgtAdmin, err := sarama.NewClusterAdmin(g.TargetBrokerList, tgtConfig)
	if err != nil {
		return err
	}

	// Create topics on target cluster
	for name, topic := range topics {
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
			err := tgtAdmin.CreateTopic(name, &detail, false)
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

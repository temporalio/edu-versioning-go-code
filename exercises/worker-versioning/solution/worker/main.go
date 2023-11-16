package main

import (
	"log"
	pizza "temporal-versioning/exercises/worker-versioning/solution"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workerOptions := worker.Options{
		BuildID:                 "revision-yymmdd",
		UseBuildIDForVersioning: true,
	}

	w := worker.New(c, pizza.TaskQueueName, workerOptions)

	w.RegisterWorkflow(pizza.PizzaWorkflow)
	w.RegisterActivity(pizza.GetDistance)
	w.RegisterActivity(pizza.SendBill)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}

}

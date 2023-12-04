package main

import (
	"log"
	pizza "temporal-versioning/exercises/worker-versioning/practice"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// TODO Part A: Change your worker.New() call to use a struct, workerOptions,
	// that contains the options `BuildID` and `UseBuildIDForVersioning: true`.

	w := worker.New(c, pizza.TaskQueueName, worker.Options{})

	w.RegisterWorkflow(pizza.PizzaWorkflow)
	w.RegisterActivity(pizza.GetDistance)
	w.RegisterActivity(pizza.SendBill)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}

}

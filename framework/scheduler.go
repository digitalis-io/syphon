package framework

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elodina/syphon/consumer"
	"github.com/golang/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	"github.com/mesos/mesos-go/scheduler"
	"github.com/stealthly/siesta"
	"io/ioutil"
	"net/http"
	"time"
)

type ElodinaTransportSchedulerConfig struct {
	Topics []string

	// Number of CPUs allocated for each created Mesos task.
	CpuPerTask float64

	// Number of RAM allocated for each created Mesos task.
	MemPerTask float64

	// Artifact server host name. Will be used to fetch the executor.
	ServiceHost string

	// Artifact server port.Will be used to fetch the executor.
	ServicePort int

	// Name of the executor binary file contained in the executor archive.
	ExecutorBinaryName string

	// Maximum retries to kill a task.
	KillTaskRetries int

	// time after partition is considered stale
	StaleDuration time.Duration

	// Mirror configuration
	ConsumerConfig consumer.PartitionConsumerConfig

	// Threads per task
	ThreadsPerTask int

	// Target produce URL
	TargetURL string

	//SSL certificate file path
	SSLCertFilePath string

	//SSL key file path
	SSLKeyFilePath string

	//SSL CA certificate file path
	SSLCACertFilePath string

	//Elodina API key
	ApiKey string

	//Elodina API user
	ApiUser string

	//Disable certificate verification
	Insecure bool
}

func NewElodinaTransportSchedulerConfig() ElodinaTransportSchedulerConfig {
	return ElodinaTransportSchedulerConfig{
		CpuPerTask:      0.2,
		MemPerTask:      256,
		KillTaskRetries: 3,
		ThreadsPerTask:  3,
	}
}

type ElodinaTransportScheduler struct {
	config               *ElodinaTransportSchedulerConfig
	taskIdToTaskState    map[string]*ElodinaTransport
	kafkaClient          siesta.Connector
	TakenTopicPartitions *consumer.TopicAndPartitionSet
}

func NewElodinaTransportScheduler(config ElodinaTransportSchedulerConfig) *ElodinaTransportScheduler {
	connectorConfig := siesta.NewConnectorConfig()
	connectorConfig.BrokerList = config.ConsumerConfig.BrokerList
	connectorConfig.ClientID = config.ConsumerConfig.ClientID
	connectorConfig.CommitOffsetBackoff = config.ConsumerConfig.CommitOffsetBackoff
	connectorConfig.CommitOffsetRetries = config.ConsumerConfig.CommitOffsetRetries
	connectorConfig.ConnectTimeout = config.ConsumerConfig.ConnectTimeout
	connectorConfig.ConsumerMetadataBackoff = config.ConsumerConfig.ConsumerMetadataBackoff
	connectorConfig.ConsumerMetadataRetries = config.ConsumerConfig.ConsumerMetadataRetries
	connectorConfig.FetchMaxWaitTime = config.ConsumerConfig.FetchMaxWaitTime
	connectorConfig.FetchMinBytes = config.ConsumerConfig.FetchMinBytes
	connectorConfig.FetchSize = config.ConsumerConfig.FetchSize
	connectorConfig.KeepAlive = config.ConsumerConfig.KeepAlive
	connectorConfig.KeepAliveTimeout = config.ConsumerConfig.KeepAliveTimeout
	connectorConfig.MaxConnections = config.ConsumerConfig.MaxConnections
	connectorConfig.MaxConnectionsPerBroker = config.ConsumerConfig.MaxConnectionsPerBroker
	connectorConfig.MetadataBackoff = config.ConsumerConfig.MetadataBackoff
	connectorConfig.MetadataRetries = config.ConsumerConfig.MetadataRetries
	connectorConfig.ReadTimeout = config.ConsumerConfig.ReadTimeout
	connectorConfig.WriteTimeout = config.ConsumerConfig.WriteTimeout
	kafkaClient, err := siesta.NewDefaultConnector(connectorConfig)
	if err != nil {
		panic(err)
	}

	scheduler := &ElodinaTransportScheduler{
		config:            &config,
		taskIdToTaskState: make(map[string]*ElodinaTransport),
		kafkaClient:       kafkaClient,
	}

	scheduler.TakenTopicPartitions = consumer.NewTopicAndPartitionSet()

	return scheduler
}

// mesos.Scheduler interface method.
// Invoked when the scheduler successfully registers with a Mesos master.
func (this *ElodinaTransportScheduler) Registered(driver scheduler.SchedulerDriver, frameworkId *mesos.FrameworkID,
	masterInfo *mesos.MasterInfo) {
	fmt.Printf("Framework Registered with Master %s\n", masterInfo)
}

// mesos.Scheduler interface method.
// Invoked when the scheduler re-registers with a newly elected Mesos master.
func (this *ElodinaTransportScheduler) Reregistered(driver scheduler.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	fmt.Printf("Framework Re-Registered with Master %s\n", masterInfo)
}

// mesos.Scheduler interface method.
// Invoked when the scheduler becomes "disconnected" from the master.
func (this *ElodinaTransportScheduler) Disconnected(driver scheduler.SchedulerDriver) {
	fmt.Println("Disconnected")
}

// mesos.Scheduler interface method.
// Invoked when resources have been offered to this framework.
func (this *ElodinaTransportScheduler) ResourceOffers(driver scheduler.SchedulerDriver, offers []*mesos.Offer) {
	fmt.Println("Received offers")
	offersAndTasks := make(map[*mesos.Offer][]*mesos.TaskInfo)
	remainingPartitions, err := this.GetTopicPartitions()
	if err != nil {
		return
	}
    remainingPartitions.RemoveAll(this.TakenTopicPartitions.GetArray())
    fmt.Printf("%v\n", remainingPartitions)
    tps := remainingPartitions.GetArray()

    offersAndResources := this.wrapInOfferAndResources(offers)
	for !remainingPartitions.IsEmpty() {
        fmt.Printf("Iteration %v\n", remainingPartitions)
		if this.hasEnoughInstances() {
			for _, transfer := range this.taskIdToTaskState {
				if len(transfer.assignment) < this.config.ThreadsPerTask {
					transfer.assignment = append(transfer.assignment, tps[0])
					remainingPartitions.Remove(tps[0])
					this.TakenTopicPartitions.Add(tps[0])
					if len(tps) > 1 {
						tps = tps[1:]
					} else {
						tps = []consumer.TopicAndPartition{}
					}
				}
			}
		} else {
			fmt.Println("Trying to launch new task")
			offer, task := this.launchNewTask(offersAndResources)
			if offer != nil && task != nil {
				offersAndTasks[offer] = append(offersAndTasks[offer], task)
			} else {
				for _, offer := range offers {
					if _, exists := offersAndTasks[offer]; !exists {
						offersAndTasks[offer] = make([]*mesos.TaskInfo, 0)
					}
				}
				break
			}
		}
	}

	this.assignPendingPartitions()

	for _, offer := range offers {
		if tasks, ok := offersAndTasks[offer]; ok {
            driver.LaunchTasks([]*mesos.OfferID{offer.Id}, tasks, &mesos.Filters{RefuseSeconds: proto.Float64(1)})
        } else {
            driver.DeclineOffer([]*mesos.OfferID{offer.Id}, &mesos.Filters{RefuseSeconds: proto.Float64(10)})
        }
	}
}

// mesos.Scheduler interface method.
// Invoked when the status of a task has changed.
func (this *ElodinaTransportScheduler) StatusUpdate(driver scheduler.SchedulerDriver, status *mesos.TaskStatus) {
	if *status.GetState().Enum() == mesos.TaskState_TASK_RUNNING {
		this.taskIdToTaskState[*status.TaskId.Value].pending = true
	} else if status.GetState() == mesos.TaskState_TASK_LOST || status.GetState() == mesos.TaskState_TASK_FAILED || status.GetState() == mesos.TaskState_TASK_FINISHED {
		this.TakenTopicPartitions.RemoveAll(this.taskIdToTaskState[*status.TaskId.Value].GetAssignment())
		delete(this.taskIdToTaskState, *status.TaskId.Value)
	}
}

// mesos.Scheduler interface method.
// Invoked when an offer is no longer valid.
func (this *ElodinaTransportScheduler) OfferRescinded(driver scheduler.SchedulerDriver, offerId *mesos.OfferID) {
	fmt.Printf("Offer %s is no longer valid\n", *offerId.Value)
}

// mesos.Scheduler interface method.
// Invoked when an executor sends a message.
func (this *ElodinaTransportScheduler) FrameworkMessage(driver scheduler.SchedulerDriver, executorId *mesos.ExecutorID,
	slaveId *mesos.SlaveID, message string) {
	fmt.Printf("Message from executor %s: %s\n", *executorId.Value, message)
}

// mesos.Scheduler interface method.
// Invoked when a slave has been determined unreachable
func (this *ElodinaTransportScheduler) SlaveLost(driver scheduler.SchedulerDriver, slaveId *mesos.SlaveID) {
	fmt.Printf("Slave %s has been lost.\n", *slaveId.Value)
}

// mesos.Scheduler interface method.
// Invoked when an executor has exited/terminated.
func (this *ElodinaTransportScheduler) ExecutorLost(scheduler scheduler.SchedulerDriver, executorId *mesos.ExecutorID,
	slaveId *mesos.SlaveID, exitCode int) {
	fmt.Printf("Executor %s on slave %s has exited with %d status code\n", *executorId.Value, *slaveId.Value, exitCode)
}

// mesos.Scheduler interface method.
// Invoked when there is an unrecoverable error in the scheduler or scheduler driver.
func (this *ElodinaTransportScheduler) Error(driver scheduler.SchedulerDriver, err string) {
	fmt.Printf("Scheduler received error: %s\n", err)
}

// Gracefully shuts down all running tasks.
func (this *ElodinaTransportScheduler) Shutdown(driver scheduler.SchedulerDriver) {
	fmt.Println("Shutting down the scheduler.")
}

func (this *ElodinaTransportScheduler) launchNewTask(offers []*OfferAndResources) (*mesos.Offer, *mesos.TaskInfo) {
	for _, offer := range offers {
		configBlob, err := json.Marshal(this.config.ConsumerConfig)
		if err != nil {
			break
		}
		fmt.Printf("%v\n", offer)
		if this.hasEnoughResources(offer) {
			port := this.takePort(&offer.RemainingPorts)
			taskPort := &mesos.Value_Range{Begin: port, End: port}
			taskId := &mesos.TaskID{
				Value: proto.String(fmt.Sprintf("elodina-mirror-%s-%d", *offer.Offer.Hostname, *port)),
			}

            cpuTaken := this.config.CpuPerTask * float64(this.config.ThreadsPerTask)
            memoryTaken := this.config.MemPerTask * float64(this.config.ThreadsPerTask)
			task := &mesos.TaskInfo{
				Name:     proto.String(taskId.GetValue()),
				TaskId:   taskId,
				SlaveId:  offer.Offer.SlaveId,
				Executor: this.createExecutor(len(this.taskIdToTaskState), *port),
				Resources: []*mesos.Resource{
					util.NewScalarResource("cpus", cpuTaken),
					util.NewScalarResource("mem", memoryTaken),
					util.NewRangesResource("ports", []*mesos.Value_Range{taskPort}),
				},
				Data: configBlob,
			}
			fmt.Printf("Prepared task: %s with offer %s for launch. Ports: %s\n", task.GetName(), offer.Offer.Id.GetValue(), taskPort)

			transport := NewElodinaTransport(fmt.Sprintf("http://%s:%d/assign", *offer.Offer.Hostname, *port), task, this.config.StaleDuration)
			this.taskIdToTaskState[*taskId.Value] = transport

			fmt.Printf("Prepared task: %s with offer %s for launch. Ports: %s\n", task.GetName(), offer.Offer.Id.GetValue(), taskPort)

			offer.RemainingPorts = offer.RemainingPorts[1:]
			offer.RemainingCpu -= cpuTaken
			offer.RemainingMemory -= memoryTaken

			return offer.Offer, task
		} else {
			fmt.Println("Not enough CPU and memory")
		}
	}

	return nil, nil
}

func (this *ElodinaTransportScheduler) hasEnoughResources(offer *OfferAndResources) bool {
	return this.config.CpuPerTask*float64(this.config.ThreadsPerTask) <= offer.RemainingCpu &&
		this.config.MemPerTask*float64(this.config.ThreadsPerTask) <= offer.RemainingMemory &&
		len(offer.RemainingPorts) > 0
}

func (this *ElodinaTransportScheduler) tryKillTask(driver scheduler.SchedulerDriver, taskId *mesos.TaskID) error {
	fmt.Printf("Trying to kill task %s\n", taskId.GetValue())
	var err error
	for i := 0; i <= this.config.KillTaskRetries; i++ {
		if _, err = driver.KillTask(taskId); err == nil {
			return nil
		}
	}
	return err
}

func (this *ElodinaTransportScheduler) removeTask(id *mesos.TaskID) {
	delete(this.taskIdToTaskState, *id.Value)
}

func (this *ElodinaTransportScheduler) takePort(ports *[]*mesos.Value_Range) *uint64 {
	port := (*ports)[0].Begin
	portRange := (*ports)[0]
	portRange.Begin = proto.Uint64((*portRange.Begin) + 1)

	if *portRange.Begin > *portRange.End {
		*ports = (*ports)[1:]
	} else {
		(*ports)[0] = portRange
	}

	return port
}

func (this *ElodinaTransportScheduler) createExecutor(instanceId int, port uint64) *mesos.ExecutorInfo {
	return &mesos.ExecutorInfo{
		ExecutorId: util.NewExecutorID(fmt.Sprintf("elodina-mirror-%d", instanceId)),
		Name:       proto.String("Elodina Mirror Executor"),
		Source:     proto.String("Elodina"),
		Command: &mesos.CommandInfo{
			Value: proto.String(fmt.Sprintf("./%s --port %d --ssl.cert %s --ssl.key %s --ssl.cacert %s --api.key %s --api.user %s --target.url %s --insecure %v",
				this.config.ExecutorBinaryName, port, this.config.SSLCertFilePath, this.config.SSLKeyFilePath, this.config.SSLCACertFilePath, this.config.ApiKey, this.config.ApiUser, this.config.TargetURL, this.config.Insecure)),
			Uris: []*mesos.CommandInfo_URI{&mesos.CommandInfo_URI{
				Value:      proto.String(fmt.Sprintf("http://%s:%d/resource/%s", this.config.ServiceHost, this.config.ServicePort, this.config.ExecutorBinaryName)),
				Executable: proto.Bool(true),
			},
				&mesos.CommandInfo_URI{
					Value:      proto.String(fmt.Sprintf("http://%s:%d/resource/%s", this.config.ServiceHost, this.config.ServicePort, this.config.SSLCertFilePath)),
					Executable: proto.Bool(false),
					Extract:    proto.Bool(false),
				},
				&mesos.CommandInfo_URI{
					Value:      proto.String(fmt.Sprintf("http://%s:%d/resource/%s", this.config.ServiceHost, this.config.ServicePort, this.config.SSLKeyFilePath)),
					Executable: proto.Bool(false),
					Extract:    proto.Bool(false),
				},
				&mesos.CommandInfo_URI{
					Value:      proto.String(fmt.Sprintf("http://%s:%d/resource/%s", this.config.ServiceHost, this.config.ServicePort, this.config.SSLCACertFilePath)),
					Executable: proto.Bool(false),
					Extract:    proto.Bool(false),
				}},
		},
	}
}

func (this *ElodinaTransportScheduler) GetTopicPartitions() (*consumer.TopicAndPartitionSet, error) {
	topicsMetadata, err := this.kafkaClient.GetTopicMetadata(this.config.Topics)
	if err != nil {
		return nil, err
	}
	topicsAndPartitions := make([]consumer.TopicAndPartition, 0)
	for _, topicMetadata := range topicsMetadata.TopicsMetadata {
		for _, partitionMetadata := range topicMetadata.PartitionsMetadata {
			topicsAndPartitions = append(topicsAndPartitions, consumer.TopicAndPartition{
				Topic:     topicMetadata.Topic,
				Partition: partitionMetadata.PartitionID,
			})
		}
	}

	tpSet := consumer.NewTopicAndPartitionSet()
	tpSet.AddAll(topicsAndPartitions)
	fmt.Printf("%v\n", topicsAndPartitions)
	fmt.Printf("%v\n", tpSet.GetArray())

	return tpSet, nil
}

func (this *ElodinaTransportScheduler) wrapInOfferAndResources(offers []*mesos.Offer) []*OfferAndResources {
	offerStates := make([]*OfferAndResources, len(offers))
	for i, offer := range offers {
		offerStates[i] = NewOfferState(offer)
	}

	return offerStates
}

func (this *ElodinaTransportScheduler) hasEnoughInstances() bool {
	for _, transfer := range this.taskIdToTaskState {
		if len(transfer.assignment) < this.config.ThreadsPerTask {
			return true
		}
	}

	return false
}

func (this *ElodinaTransportScheduler) assignPendingPartitions() {
	for _, transfer := range this.taskIdToTaskState {
		if !transfer.IsPending() {
			continue
		}

		data, err := json.Marshal(transfer.GetAssignment())
		if err != nil {
			fmt.Println(err.Error())
		} else {
			request, err := http.NewRequest("POST", transfer.GetConnectUrl(), bytes.NewReader(data))
			if err != nil {
				fmt.Println(err.Error())
			}
			resp, err := http.DefaultClient.Do(request)
			if err != nil {
				panic(err.Error())
			}
			if resp.StatusCode != 200 {
				func() {
					defer resp.Body.Close()
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Println(err.Error())
					}
					fmt.Println(string(body))
				}()
			} else {
				transfer.pending = false
			}
		}
	}
}

type OfferAndResources struct {
	RemainingCpu    float64
	RemainingMemory float64
	RemainingPorts  []*mesos.Value_Range
	Offer           *mesos.Offer
}

func NewOfferState(offer *mesos.Offer) *OfferAndResources {
	cpus := getScalarResources(offer, "cpus")
	memory := getScalarResources(offer, "mem")
	ports := getRangeResources(offer, "ports")

	return &OfferAndResources{
		RemainingCpu:    cpus,
		RemainingMemory: memory,
		RemainingPorts:  ports,
		Offer:           offer,
	}
}

// Code generated by protoc-gen-gogo.
// source: scheduler.proto
// DO NOT EDIT!

package mesosproto

import proto "github.com/elodina/syphon/framework/Godeps/_workspace/src/github.com/gogo/protobuf/proto"
import math "math"

// discarding unused import gogoproto "github.com/gogo/protobuf/gogoproto/gogo.pb"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

// Possible event types, followed by message definitions if
// applicable.
type Event_Type int32

const (
	Event_REGISTERED   Event_Type = 1
	Event_REREGISTERED Event_Type = 2
	Event_OFFERS       Event_Type = 3
	Event_RESCIND      Event_Type = 4
	Event_UPDATE       Event_Type = 5
	Event_MESSAGE      Event_Type = 6
	Event_FAILURE      Event_Type = 7
	Event_ERROR        Event_Type = 8
)

var Event_Type_name = map[int32]string{
	1: "REGISTERED",
	2: "REREGISTERED",
	3: "OFFERS",
	4: "RESCIND",
	5: "UPDATE",
	6: "MESSAGE",
	7: "FAILURE",
	8: "ERROR",
}
var Event_Type_value = map[string]int32{
	"REGISTERED":   1,
	"REREGISTERED": 2,
	"OFFERS":       3,
	"RESCIND":      4,
	"UPDATE":       5,
	"MESSAGE":      6,
	"FAILURE":      7,
	"ERROR":        8,
}

func (x Event_Type) Enum() *Event_Type {
	p := new(Event_Type)
	*p = x
	return p
}
func (x Event_Type) String() string {
	return proto.EnumName(Event_Type_name, int32(x))
}
func (x *Event_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Event_Type_value, data, "Event_Type")
	if err != nil {
		return err
	}
	*x = Event_Type(value)
	return nil
}

// Possible call types, followed by message definitions if
// applicable.
type Call_Type int32

const (
	Call_REGISTER    Call_Type = 1
	Call_REREGISTER  Call_Type = 2
	Call_UNREGISTER  Call_Type = 3
	Call_REQUEST     Call_Type = 4
	Call_DECLINE     Call_Type = 5
	Call_REVIVE      Call_Type = 6
	Call_LAUNCH      Call_Type = 7
	Call_KILL        Call_Type = 8
	Call_ACKNOWLEDGE Call_Type = 9
	Call_RECONCILE   Call_Type = 10
	Call_MESSAGE     Call_Type = 11
)

var Call_Type_name = map[int32]string{
	1:  "REGISTER",
	2:  "REREGISTER",
	3:  "UNREGISTER",
	4:  "REQUEST",
	5:  "DECLINE",
	6:  "REVIVE",
	7:  "LAUNCH",
	8:  "KILL",
	9:  "ACKNOWLEDGE",
	10: "RECONCILE",
	11: "MESSAGE",
}
var Call_Type_value = map[string]int32{
	"REGISTER":    1,
	"REREGISTER":  2,
	"UNREGISTER":  3,
	"REQUEST":     4,
	"DECLINE":     5,
	"REVIVE":      6,
	"LAUNCH":      7,
	"KILL":        8,
	"ACKNOWLEDGE": 9,
	"RECONCILE":   10,
	"MESSAGE":     11,
}

func (x Call_Type) Enum() *Call_Type {
	p := new(Call_Type)
	*p = x
	return p
}
func (x Call_Type) String() string {
	return proto.EnumName(Call_Type_name, int32(x))
}
func (x *Call_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Call_Type_value, data, "Call_Type")
	if err != nil {
		return err
	}
	*x = Call_Type(value)
	return nil
}

// *
// Low-level scheduler event API.
//
// An event is described using the standard protocol buffer "union"
// trick, see https://developers.google.com/protocol-buffers/docs/techniques#union.
type Event struct {
	// Type of the event, indicates which optional field below should be
	// present if that type has a nested message definition.
	Type             *Event_Type         `protobuf:"varint,1,req,name=type,enum=mesosproto.Event_Type" json:"type,omitempty"`
	Registered       *Event_Registered   `protobuf:"bytes,2,opt,name=registered" json:"registered,omitempty"`
	Reregistered     *Event_Reregistered `protobuf:"bytes,3,opt,name=reregistered" json:"reregistered,omitempty"`
	Offers           *Event_Offers       `protobuf:"bytes,4,opt,name=offers" json:"offers,omitempty"`
	Rescind          *Event_Rescind      `protobuf:"bytes,5,opt,name=rescind" json:"rescind,omitempty"`
	Update           *Event_Update       `protobuf:"bytes,6,opt,name=update" json:"update,omitempty"`
	Message          *Event_Message      `protobuf:"bytes,7,opt,name=message" json:"message,omitempty"`
	Failure          *Event_Failure      `protobuf:"bytes,8,opt,name=failure" json:"failure,omitempty"`
	Error            *Event_Error        `protobuf:"bytes,9,opt,name=error" json:"error,omitempty"`
	XXX_unrecognized []byte              `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}

func (m *Event) GetType() Event_Type {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Event_REGISTERED
}

func (m *Event) GetRegistered() *Event_Registered {
	if m != nil {
		return m.Registered
	}
	return nil
}

func (m *Event) GetReregistered() *Event_Reregistered {
	if m != nil {
		return m.Reregistered
	}
	return nil
}

func (m *Event) GetOffers() *Event_Offers {
	if m != nil {
		return m.Offers
	}
	return nil
}

func (m *Event) GetRescind() *Event_Rescind {
	if m != nil {
		return m.Rescind
	}
	return nil
}

func (m *Event) GetUpdate() *Event_Update {
	if m != nil {
		return m.Update
	}
	return nil
}

func (m *Event) GetMessage() *Event_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *Event) GetFailure() *Event_Failure {
	if m != nil {
		return m.Failure
	}
	return nil
}

func (m *Event) GetError() *Event_Error {
	if m != nil {
		return m.Error
	}
	return nil
}

type Event_Registered struct {
	FrameworkId      *FrameworkID `protobuf:"bytes,1,req,name=framework_id" json:"framework_id,omitempty"`
	MasterInfo       *MasterInfo  `protobuf:"bytes,2,req,name=master_info" json:"master_info,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Event_Registered) Reset()         { *m = Event_Registered{} }
func (m *Event_Registered) String() string { return proto.CompactTextString(m) }
func (*Event_Registered) ProtoMessage()    {}

func (m *Event_Registered) GetFrameworkId() *FrameworkID {
	if m != nil {
		return m.FrameworkId
	}
	return nil
}

func (m *Event_Registered) GetMasterInfo() *MasterInfo {
	if m != nil {
		return m.MasterInfo
	}
	return nil
}

type Event_Reregistered struct {
	FrameworkId      *FrameworkID `protobuf:"bytes,1,req,name=framework_id" json:"framework_id,omitempty"`
	MasterInfo       *MasterInfo  `protobuf:"bytes,2,req,name=master_info" json:"master_info,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Event_Reregistered) Reset()         { *m = Event_Reregistered{} }
func (m *Event_Reregistered) String() string { return proto.CompactTextString(m) }
func (*Event_Reregistered) ProtoMessage()    {}

func (m *Event_Reregistered) GetFrameworkId() *FrameworkID {
	if m != nil {
		return m.FrameworkId
	}
	return nil
}

func (m *Event_Reregistered) GetMasterInfo() *MasterInfo {
	if m != nil {
		return m.MasterInfo
	}
	return nil
}

type Event_Offers struct {
	Offers           []*Offer `protobuf:"bytes,1,rep,name=offers" json:"offers,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *Event_Offers) Reset()         { *m = Event_Offers{} }
func (m *Event_Offers) String() string { return proto.CompactTextString(m) }
func (*Event_Offers) ProtoMessage()    {}

func (m *Event_Offers) GetOffers() []*Offer {
	if m != nil {
		return m.Offers
	}
	return nil
}

type Event_Rescind struct {
	OfferId          *OfferID `protobuf:"bytes,1,req,name=offer_id" json:"offer_id,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *Event_Rescind) Reset()         { *m = Event_Rescind{} }
func (m *Event_Rescind) String() string { return proto.CompactTextString(m) }
func (*Event_Rescind) ProtoMessage()    {}

func (m *Event_Rescind) GetOfferId() *OfferID {
	if m != nil {
		return m.OfferId
	}
	return nil
}

type Event_Update struct {
	Uuid             []byte      `protobuf:"bytes,1,req,name=uuid" json:"uuid,omitempty"`
	Status           *TaskStatus `protobuf:"bytes,2,req,name=status" json:"status,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *Event_Update) Reset()         { *m = Event_Update{} }
func (m *Event_Update) String() string { return proto.CompactTextString(m) }
func (*Event_Update) ProtoMessage()    {}

func (m *Event_Update) GetUuid() []byte {
	if m != nil {
		return m.Uuid
	}
	return nil
}

func (m *Event_Update) GetStatus() *TaskStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

type Event_Message struct {
	SlaveId          *SlaveID    `protobuf:"bytes,1,req,name=slave_id" json:"slave_id,omitempty"`
	ExecutorId       *ExecutorID `protobuf:"bytes,2,req,name=executor_id" json:"executor_id,omitempty"`
	Data             []byte      `protobuf:"bytes,3,req,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *Event_Message) Reset()         { *m = Event_Message{} }
func (m *Event_Message) String() string { return proto.CompactTextString(m) }
func (*Event_Message) ProtoMessage()    {}

func (m *Event_Message) GetSlaveId() *SlaveID {
	if m != nil {
		return m.SlaveId
	}
	return nil
}

func (m *Event_Message) GetExecutorId() *ExecutorID {
	if m != nil {
		return m.ExecutorId
	}
	return nil
}

func (m *Event_Message) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type Event_Failure struct {
	SlaveId *SlaveID `protobuf:"bytes,1,opt,name=slave_id" json:"slave_id,omitempty"`
	// If this was just a failure of an executor on a slave then
	// 'executor_id' will be set and possibly 'status' (if we were
	// able to determine the exit status).
	ExecutorId       *ExecutorID `protobuf:"bytes,2,opt,name=executor_id" json:"executor_id,omitempty"`
	Status           *int32      `protobuf:"varint,3,opt,name=status" json:"status,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *Event_Failure) Reset()         { *m = Event_Failure{} }
func (m *Event_Failure) String() string { return proto.CompactTextString(m) }
func (*Event_Failure) ProtoMessage()    {}

func (m *Event_Failure) GetSlaveId() *SlaveID {
	if m != nil {
		return m.SlaveId
	}
	return nil
}

func (m *Event_Failure) GetExecutorId() *ExecutorID {
	if m != nil {
		return m.ExecutorId
	}
	return nil
}

func (m *Event_Failure) GetStatus() int32 {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return 0
}

type Event_Error struct {
	Message          *string `protobuf:"bytes,1,req,name=message" json:"message,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Event_Error) Reset()         { *m = Event_Error{} }
func (m *Event_Error) String() string { return proto.CompactTextString(m) }
func (*Event_Error) ProtoMessage()    {}

func (m *Event_Error) GetMessage() string {
	if m != nil && m.Message != nil {
		return *m.Message
	}
	return ""
}

// *
// Low-level scheduler call API.
//
// Like Event, a Call is described using the standard protocol buffer
// "union" trick (see above).
type Call struct {
	// Identifies who generated this call. Always necessary, but the
	// only thing that needs to be set for certain calls, e.g.,
	// REGISTER, REREGISTER, and UNREGISTER.
	FrameworkInfo *FrameworkInfo `protobuf:"bytes,1,req,name=framework_info" json:"framework_info,omitempty"`
	// Type of the call, indicates which optional field below should be
	// present if that type has a nested message definition.
	Type             *Call_Type        `protobuf:"varint,2,req,name=type,enum=mesosproto.Call_Type" json:"type,omitempty"`
	Request          *Call_Request     `protobuf:"bytes,3,opt,name=request" json:"request,omitempty"`
	Decline          *Call_Decline     `protobuf:"bytes,4,opt,name=decline" json:"decline,omitempty"`
	Launch           *Call_Launch      `protobuf:"bytes,5,opt,name=launch" json:"launch,omitempty"`
	Kill             *Call_Kill        `protobuf:"bytes,6,opt,name=kill" json:"kill,omitempty"`
	Acknowledge      *Call_Acknowledge `protobuf:"bytes,7,opt,name=acknowledge" json:"acknowledge,omitempty"`
	Reconcile        *Call_Reconcile   `protobuf:"bytes,8,opt,name=reconcile" json:"reconcile,omitempty"`
	Message          *Call_Message     `protobuf:"bytes,9,opt,name=message" json:"message,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *Call) Reset()         { *m = Call{} }
func (m *Call) String() string { return proto.CompactTextString(m) }
func (*Call) ProtoMessage()    {}

func (m *Call) GetFrameworkInfo() *FrameworkInfo {
	if m != nil {
		return m.FrameworkInfo
	}
	return nil
}

func (m *Call) GetType() Call_Type {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Call_REGISTER
}

func (m *Call) GetRequest() *Call_Request {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *Call) GetDecline() *Call_Decline {
	if m != nil {
		return m.Decline
	}
	return nil
}

func (m *Call) GetLaunch() *Call_Launch {
	if m != nil {
		return m.Launch
	}
	return nil
}

func (m *Call) GetKill() *Call_Kill {
	if m != nil {
		return m.Kill
	}
	return nil
}

func (m *Call) GetAcknowledge() *Call_Acknowledge {
	if m != nil {
		return m.Acknowledge
	}
	return nil
}

func (m *Call) GetReconcile() *Call_Reconcile {
	if m != nil {
		return m.Reconcile
	}
	return nil
}

func (m *Call) GetMessage() *Call_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

type Call_Request struct {
	Requests         []*Request `protobuf:"bytes,1,rep,name=requests" json:"requests,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *Call_Request) Reset()         { *m = Call_Request{} }
func (m *Call_Request) String() string { return proto.CompactTextString(m) }
func (*Call_Request) ProtoMessage()    {}

func (m *Call_Request) GetRequests() []*Request {
	if m != nil {
		return m.Requests
	}
	return nil
}

type Call_Decline struct {
	OfferIds         []*OfferID `protobuf:"bytes,1,rep,name=offer_ids" json:"offer_ids,omitempty"`
	Filters          *Filters   `protobuf:"bytes,2,opt,name=filters" json:"filters,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *Call_Decline) Reset()         { *m = Call_Decline{} }
func (m *Call_Decline) String() string { return proto.CompactTextString(m) }
func (*Call_Decline) ProtoMessage()    {}

func (m *Call_Decline) GetOfferIds() []*OfferID {
	if m != nil {
		return m.OfferIds
	}
	return nil
}

func (m *Call_Decline) GetFilters() *Filters {
	if m != nil {
		return m.Filters
	}
	return nil
}

type Call_Launch struct {
	TaskInfos        []*TaskInfo `protobuf:"bytes,1,rep,name=task_infos" json:"task_infos,omitempty"`
	OfferIds         []*OfferID  `protobuf:"bytes,2,rep,name=offer_ids" json:"offer_ids,omitempty"`
	Filters          *Filters    `protobuf:"bytes,3,opt,name=filters" json:"filters,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *Call_Launch) Reset()         { *m = Call_Launch{} }
func (m *Call_Launch) String() string { return proto.CompactTextString(m) }
func (*Call_Launch) ProtoMessage()    {}

func (m *Call_Launch) GetTaskInfos() []*TaskInfo {
	if m != nil {
		return m.TaskInfos
	}
	return nil
}

func (m *Call_Launch) GetOfferIds() []*OfferID {
	if m != nil {
		return m.OfferIds
	}
	return nil
}

func (m *Call_Launch) GetFilters() *Filters {
	if m != nil {
		return m.Filters
	}
	return nil
}

type Call_Kill struct {
	TaskId           *TaskID `protobuf:"bytes,1,req,name=task_id" json:"task_id,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Call_Kill) Reset()         { *m = Call_Kill{} }
func (m *Call_Kill) String() string { return proto.CompactTextString(m) }
func (*Call_Kill) ProtoMessage()    {}

func (m *Call_Kill) GetTaskId() *TaskID {
	if m != nil {
		return m.TaskId
	}
	return nil
}

type Call_Acknowledge struct {
	SlaveId          *SlaveID `protobuf:"bytes,1,req,name=slave_id" json:"slave_id,omitempty"`
	TaskId           *TaskID  `protobuf:"bytes,2,req,name=task_id" json:"task_id,omitempty"`
	Uuid             []byte   `protobuf:"bytes,3,req,name=uuid" json:"uuid,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *Call_Acknowledge) Reset()         { *m = Call_Acknowledge{} }
func (m *Call_Acknowledge) String() string { return proto.CompactTextString(m) }
func (*Call_Acknowledge) ProtoMessage()    {}

func (m *Call_Acknowledge) GetSlaveId() *SlaveID {
	if m != nil {
		return m.SlaveId
	}
	return nil
}

func (m *Call_Acknowledge) GetTaskId() *TaskID {
	if m != nil {
		return m.TaskId
	}
	return nil
}

func (m *Call_Acknowledge) GetUuid() []byte {
	if m != nil {
		return m.Uuid
	}
	return nil
}

// Allows the framework to query the status for non-terminal tasks.
// This causes the master to send back the latest task status for
// each task in 'statuses', if possible. Tasks that are no longer
// known will result in a TASK_LOST update. If statuses is empty,
// then the master will send the latest status for each task
// currently known.
// TODO(bmahler): Add a guiding document for reconciliation or
// document reconciliation in-depth here.
type Call_Reconcile struct {
	Statuses         []*TaskStatus `protobuf:"bytes,1,rep,name=statuses" json:"statuses,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *Call_Reconcile) Reset()         { *m = Call_Reconcile{} }
func (m *Call_Reconcile) String() string { return proto.CompactTextString(m) }
func (*Call_Reconcile) ProtoMessage()    {}

func (m *Call_Reconcile) GetStatuses() []*TaskStatus {
	if m != nil {
		return m.Statuses
	}
	return nil
}

type Call_Message struct {
	SlaveId          *SlaveID    `protobuf:"bytes,1,req,name=slave_id" json:"slave_id,omitempty"`
	ExecutorId       *ExecutorID `protobuf:"bytes,2,req,name=executor_id" json:"executor_id,omitempty"`
	Data             []byte      `protobuf:"bytes,3,req,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *Call_Message) Reset()         { *m = Call_Message{} }
func (m *Call_Message) String() string { return proto.CompactTextString(m) }
func (*Call_Message) ProtoMessage()    {}

func (m *Call_Message) GetSlaveId() *SlaveID {
	if m != nil {
		return m.SlaveId
	}
	return nil
}

func (m *Call_Message) GetExecutorId() *ExecutorID {
	if m != nil {
		return m.ExecutorId
	}
	return nil
}

func (m *Call_Message) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterEnum("mesosproto.Event_Type", Event_Type_name, Event_Type_value)
	proto.RegisterEnum("mesosproto.Call_Type", Call_Type_name, Call_Type_value)
}

package igntransport

/*
// Find the ignition-transport7 package
#cgo pkg-config: ignition-transport7

// needed by unsafe
#include <stdlib.h>

// The Ignition Transport C header file
#include <ignition/transport/CIface.h>

// Callbacks proxying logic
// Unfortunately, we cannot use "const" modifier for arguments.
void ign_cb_proxy(char *_msg, size_t _size, char *_msgType, void *user_data);

static void _ign_subscribe(IgnTransportNode *n, char *_topic, void *cb_struct) {
	ignTransportSubscribeNonConst(n, _topic, ign_cb_proxy, cb_struct);
}

*/
import "C"
import (
	"gitlab.com/ignitionrobotics/web/ign-go"
	msgs "bitbucket.org/ignitionrobotics/web-cloudsim/ign-transport/proto/ignition/msgs"
	proto "github.com/golang/protobuf/proto"
	"github.com/mattn/go-pointer"
	"github.com/pkg/errors"
	"sync"
	"time"
	"unsafe"
)

/*
Main idea taken from here: https://dev.to/mattn/call-go-function-from-c-function-1n3

WARNING: threads in C++/C can break stuff in Go. Test it!

Older: https://stackoverflow.com/questions/6125683/call-go-functions-from-c
More: https://stackoverflow.com/questions/32215509/using-go-code-in-an-existing-c-project

some official doc:
https://github.com/golang/go/wiki/cgo#function-variables
https://golang.org/cmd/cgo/#hdr-Passing_pointers

Note: to compile protos, from the proto folder do: protoc --proto_path=. --go_out=. ignition/msgs/*.proto

*/

// IMPORTANT Protobuf: to compile ign-transport protobuf files, you need to run `protoc --proto_path=. --go_out=. ignition/msgs/*.proto`
// from the `web-cloudsim/ign-transport/proto/` folder.

// GoIgnTransportNode is the Go type we use to interact with our C api. It
// represents an ign-transport node.
type GoIgnTransportNode struct {
	// n is a pointer to an Ignition Transport Node.
	n *C.IgnTransportNode
	// partition is the ign_partition to which the ign-transport node will belong
	partition *string
	// topicsToPointers is a map of topics to callback pointers. Each subscribed topic
	// will have a callback function invoked from C to Go code. We cannot pass functions
	// directly to C. Instead, an 'unsafe.Pointer' pointing to that callback is passed to C.
	topicsToPointers map[string][]unsafe.Pointer
	// mutex needed to protect map from concurrent access
	lockTopicsToPointers sync.RWMutex
	// initializedPubTopics is a map used to track which topics -- used to publish -- were already
	// initialized, ie. the advertiser node was created and is ready to publish messages.
	initializedPubTopics map[string]bool
	// mutex needed to protect map from concurrent access
	lockIniatilizedPubTopics sync.RWMutex
}

// TODO: add loggers from Contexts

// NewIgnTransportNode creates a new ignition transport node.
func NewIgnTransportNode(partition *string) (*GoIgnTransportNode, error) {
	var ign GoIgnTransportNode

	if partition == nil {
		// ign-transport will use the partition from IGN_PARTITION env var.
		ign.n = C.ignTransportNodeCreate(nil)
	} else {
		partitionC := C.CString(*partition)
		defer C.free(unsafe.Pointer(partitionC))
		ign.n = C.ignTransportNodeCreate(partitionC)
	}
	ign.partition = partition
	ign.lockTopicsToPointers = sync.RWMutex{}
	ign.topicsToPointers = make(map[string][]unsafe.Pointer, 0)
	ign.initializedPubTopics = make(map[string]bool, 0)
	ign.lockIniatilizedPubTopics = sync.RWMutex{}
	return &ign, nil
}

// Free helps freeing C memory
// Note: you should have unsubscribed from all topics before invoking this func.
func (ign *GoIgnTransportNode) Free() {
	C.ignTransportNodeDestroy(&ign.n)
	ign.n = nil

	// Free the callback function pointers
	// First create a shallow copy of the map keys to avoid removing map keys while
	// iterating the map.
	ign.lockTopicsToPointers.RLock()
	keys := make([]string, 0)
	for k := range ign.topicsToPointers {
		keys = append(keys, k)
	}
	ign.lockTopicsToPointers.RUnlock()

	// Now free and remove map items
	for _, topic := range keys {
		ign.freeCallbackPointers(topic)
	}
}

// CallbackIgnT is a wrapper type for callback functions for ign-transport C interface.
type CallbackIgnT struct {
	Func func([]byte, string)
}

//export ign_cb_proxy
func ign_cb_proxy(msg *C.char, size C.size_t, msgType *C.char, fnp unsafe.Pointer) {

	defer func() {
		// check for panic and catch it. We don't want a server restart due to an error here
		if p := recover(); p != nil {
			// We should not be here! Just log it
			logger := ign.NewLogger("ign_cb_proxy", true, ign.VerbosityDebug)
			logger.Critical("Panic while running ign_transport's ign_cb_proxy function", p)
		}
	}()

	// receiving C strings: https://stackoverflow.com/questions/39708874/call-go-function-with-string-parameter-from-c
	fn := pointer.Restore(fnp).(*CallbackIgnT)
	// Converting C strings to Go Bytes: https://gist.github.com/zchee/b9c99695463d8902cd33#string
	bytes := C.GoBytes(unsafe.Pointer(msg), C.int(size))
	// Invoke the GO callback
	fn.Func(bytes, C.GoString(msgType))
}

// IgnTransportSubscribe registers a callback in C. It does this by wrapping the cb function
// and passes userdata into a Callback struct, and actually passing a pointer to that
// struct to C. It is expected that the C code will run a callback of type *void.
func (ign *GoIgnTransportNode) IgnTransportSubscribe(topic string, cb func([]byte, string)) error {

	if ign.n == nil {
		return errors.New("Internal ign.n node is nil")
	}

	ptr := pointer.Save(&CallbackIgnT{
		Func: cb,
	})

	ign.lockTopicsToPointers.Lock()
	ign.topicsToPointers[topic] = append(ign.topicsToPointers[topic], ptr)
	ign.lockTopicsToPointers.Unlock()

	topicC := C.CString(topic)
	defer C.free(unsafe.Pointer(topicC))

	C._ign_subscribe(ign.n, topicC, ptr)
	return nil
}

// IgnTransportUnsubscribe unsubscribes from a topic
func (ign *GoIgnTransportNode) IgnTransportUnsubscribe(topic string) error {

	if ign.n == nil {
		return errors.New("Internal ign.n node is nil")
	}

	var err error
	topicC := C.CString(topic)
	defer C.free(unsafe.Pointer(topicC))

	if i := C.ignTransportUnsubscribe(ign.n, topicC); i != 0 {
		err = errors.New("error while calling C.ignTransportUnsubscribe")
	}
	ign.freeCallbackPointers(topic)

	return err
}

func (ign *GoIgnTransportNode) freeCallbackPointers(topic string) {
	ign.lockTopicsToPointers.Lock()
	defer ign.lockTopicsToPointers.Unlock()

	// clear our internal map
	for _, ptr := range ign.topicsToPointers[topic] {
		ign.unregisterCallback(ptr)
	}
	delete(ign.topicsToPointers, topic)
}

func (ign *GoIgnTransportNode) internalPublishMsg(topicC, msgTypeC *C.char, data string) error {
	if ign.n == nil {
		return errors.New("Internal ign.n node is nil")
	}

	sm := msgs.StringMsg{
		Data: data,
	}
	// We marshal the proto to get a byte array (expected by C interface)
	out, _ := proto.Marshal(&sm)
	// CGO trick: when sending pointer to an slice, you need to send pointer to slice[0]
	// See: https://coderwall.com/p/m_ma7q/pass-go-slices-as-c-array-parameters
	p := unsafe.Pointer(&out[0])
	res := C.ignTransportPublish(ign.n, topicC, p, msgTypeC)
	if res != 0 {
		return errors.New("error invoking ignTransportPublish")
	}
	return nil
}

// IgnTransportPublishStringMsg publishes a message to a topic.
func (ign *GoIgnTransportNode) IgnTransportPublishStringMsg(topic, msg string) error {

	if ign.n == nil {
		return errors.New("Internal ign.n node is nil")
	}

	topicC := C.CString(topic)
	defer C.free(unsafe.Pointer(topicC))

	msgTypeC := C.CString("ignition.msgs.StringMsg")
	defer C.free(unsafe.Pointer(msgTypeC))

	// Is this the first time we send a message to this topic?
	ign.lockIniatilizedPubTopics.RLock()
	val, ok := ign.initializedPubTopics[topic]
	ign.lockIniatilizedPubTopics.RUnlock()

	if !ok || !val {
		// Create the advertise node.
		if res := C.ignTransportAdvertise(ign.n, topicC, msgTypeC); res != 0 {
			return errors.New("error invoking ignTransportAdvertise")
		}
		// Important: we need to sleep 1 second to let ing-transport
		// create the advertiser node. Otherwise messages won't be sent.
		time.Sleep(1000 * time.Millisecond)
		ign.lockIniatilizedPubTopics.Lock()
		ign.initializedPubTopics[topic] = true
		ign.lockIniatilizedPubTopics.Unlock()
	}

	// send the actual msg
	if err := ign.internalPublishMsg(topicC, msgTypeC, msg); err != nil {
		return errors.New("error invoking ignTransportPublish")
	}

	return nil
}

// unregisterCallback frees memory from a given callback
func (ign *GoIgnTransportNode) unregisterCallback(ptr unsafe.Pointer) error {
	pointer.Unref(ptr)
	return nil
}

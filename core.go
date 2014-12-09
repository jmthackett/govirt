package govirt

// #cgo pkg-config: libvirt
// #include <stdlib.h>
// #include <libvirt/libvirt.h>
import "C"
import "unsafe"
import "errors"

// Represents a libvirt node.
// Exposes methods for controlling the hypervisor.
type Node struct {
	Url  string
	conn Connection
	domain Domain
}

type Connection struct {
	Pointer C.virConnectPtr
}

type Domain struct {
	Pointer C.virDomainPtr
}

type Network struct {
}

type StorageVol struct {
}

type StoragePool struct {
}

func (p *Node) Connect() error {
	cUrl := C.CString(p.Url)
	defer C.free(unsafe.Pointer(cUrl))
	conn := C.virConnectOpenAuth(cUrl, C.virConnectAuthPtrDefault, 0) // TODO: Implement alias override // TODO: work out what the previous commenter meant
	if conn == nil {
		return errors.New("Unable to establish connection")
	}
	p.conn = Connection{conn}

	return nil
}

func (p *Node) ConnectGetCapabilities() string {
	retval := C.virConnectGetCapabilities(p.conn.Pointer)
	return C.GoString(retval)
}

func (p *Node) NodeGetFreeMemory() uint64 {
	retval := C.virNodeGetFreeMemory(p.conn.Pointer)
	return uint64(retval)
}

func (p *Node) Disconnect() int {
	if p.conn.Pointer != nil {
		return int(C.virConnectClose(p.conn.Pointer))
	}
	return -1
}

func (p *Node) DomainDefineXML(xml string) error  {
	domain := C.virDomainDefineXML(p.conn.Pointer, C.CString(xml))
	p.domain = Domain{domain}
	if p.domain.Pointer != nil {
		return errors.New("Unable to define domain!")
	}
	return nil
}

func Connect(url string) (Node, error) {
	var retval Node
	retval = Node{Url: url}
	err := retval.Connect()
	return retval, err
}

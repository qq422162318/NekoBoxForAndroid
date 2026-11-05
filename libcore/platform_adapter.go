package libcore

import (
	"context"
	"net/netip"
	"syscall"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/option"
	tun "github.com/sagernet/sing-tun"
	E "github.com/sagernet/sing/common/exceptions"
)

// boxAdapterPlatformInterface adapts the legacy libbox platform.Interface
// (boxPlatformInterfaceWrapper) to the sing-box 1.14+ adapter.PlatformInterface.
//
// Since sing-box 1.11 the box looks up adapter.PlatformInterface (which includes
// UsePlatformDefaultInterfaceMonitor) to decide whether the app handles default
// interface monitoring itself. Without it sing-box falls back to a raw netlink
// socket, which is banned on Android, causing:
//
//	create network monitor: netlink socket in Android is banned by Google
//
// The embedded boxPlatformInterfaceWrapper already provides UsePlatformDefaultInterfaceMonitor
// (true) and CreateDefaultInterfaceMonitor (stub); this type only fills in the
// remaining adapter.PlatformInterface methods so the box can find it.
var boxAdapterPlatformInterfaceInstance adapter.PlatformInterface = &boxAdapterPlatformInterface{boxPlatformInterfaceWrapperInstance}

type boxAdapterPlatformInterface struct {
	*boxPlatformInterfaceWrapper
}

func (w *boxAdapterPlatformInterface) UsePlatformInterface() bool {
	return true
}

func (w *boxAdapterPlatformInterface) OpenInterface(options *tun.Options, platformOptions option.TunPlatformOptions) (tun.Tun, error) {
	// delegate to the existing tun open path (fd dup + tun.New handled there)
	return w.OpenTun(options, platformOptions)
}

func (w *boxAdapterPlatformInterface) ProcessPlatformOptions(options option.TunPlatformOptions) error {
	return nil
}

func (w *boxAdapterPlatformInterface) UsePlatformDefaultInterfaceMonitor() bool {
	return true
}

func (w *boxAdapterPlatformInterface) UsePlatformNetworkInterfaces() bool {
	return false
}

func (w *boxAdapterPlatformInterface) NetworkInterfaces() ([]adapter.NetworkInterface, error) {
	return nil, nil
}

func (w *boxAdapterPlatformInterface) NetworkExtensionIncludeAllNetworks() bool {
	return w.IncludeAllNetworks()
}

func (w *boxAdapterPlatformInterface) RequestPermissionForWIFIState() error {
	return nil
}

func (w *boxAdapterPlatformInterface) ReadWIFIState(ctx context.Context) adapter.WIFIState {
	// boxPlatformInterfaceWrapper.ReadWIFIState returns adapter.WIFIState by value.
	return w.boxPlatformInterfaceWrapper.ReadWIFIState()
}

func (w *boxAdapterPlatformInterface) UsePlatformConnectionOwnerFinder() bool {
	return true
}

func (w *boxAdapterPlatformInterface) FindConnectionOwner(request *adapter.FindConnectionOwnerRequest) (*adapter.ConnectionOwner, error) {
	var network string
	switch request.IpProtocol {
	case int32(syscall.IPPROTO_TCP):
		network = "tcp"
	case int32(syscall.IPPROTO_UDP):
		network = "udp"
	default:
		return nil, E.New("unknown protocol: ", request.IpProtocol)
	}
	sourceAddr, err := netip.ParseAddr(request.SourceAddress)
	if err != nil {
		return nil, err
	}
	destAddr, err := netip.ParseAddr(request.DestinationAddress)
	if err != nil {
		return nil, err
	}
	source := netip.AddrPortFrom(sourceAddr, uint16(request.SourcePort))
	destination := netip.AddrPortFrom(destAddr, uint16(request.DestinationPort))
	return w.boxPlatformInterfaceWrapper.FindProcessInfo(context.Background(), network, source, destination)
}

func (w *boxAdapterPlatformInterface) UsePlatformWIFIMonitor() bool {
	return false
}

func (w *boxAdapterPlatformInterface) UsePlatformNotification() bool {
	return false
}

func (w *boxAdapterPlatformInterface) SendNotification(notification *adapter.Notification) error {
	return nil
}

func (w *boxAdapterPlatformInterface) CancelNotification(identifier string, typeID int32) error {
	return nil
}

func (w *boxAdapterPlatformInterface) MyInterfaceAddress() []netip.Addr {
	return nil
}

func (w *boxAdapterPlatformInterface) UsePlatformNeighborResolver() bool {
	return false
}

func (w *boxAdapterPlatformInterface) StartNeighborMonitor(listener adapter.NeighborUpdateListener) error {
	return nil
}

func (w *boxAdapterPlatformInterface) CloseNeighborMonitor(listener adapter.NeighborUpdateListener) error {
	return nil
}

func (w *boxAdapterPlatformInterface) UsePlatformShell() bool {
	return false
}

func (w *boxAdapterPlatformInterface) CheckPlatformShell() error {
	return nil
}

func (w *boxAdapterPlatformInterface) OpenShellSession(user *adapter.PlatformUser, command string, env []string, term string, rows int32, cols int32) (adapter.ShellSession, error) {
	return nil, E.New("platform shell not supported")
}

func (w *boxAdapterPlatformInterface) LookupUser(username string) (*adapter.PlatformUser, error) {
	return nil, E.New("platform user lookup not supported")
}

func (w *boxAdapterPlatformInterface) LookupSFTPServer() (string, error) {
	return "", nil
}

func (w *boxAdapterPlatformInterface) ReadSystemSSHHostKey() ([]byte, error) {
	return nil, nil
}

func (w *boxAdapterPlatformInterface) TailscaleHostname() string {
	return ""
}

func (w *boxAdapterPlatformInterface) UsePlatformBridge() bool {
	return false
}

func (w *boxAdapterPlatformInterface) CreateBridge(options adapter.BridgeOptions) (adapter.BridgeSession, error) {
	return nil, E.New("platform bridge not supported")
}

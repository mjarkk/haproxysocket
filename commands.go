package haproxysocket

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ShowErrors report last request and response errors for each proxy
func (h *HaproxyInstace) ShowErrors() error {
	return errors.New("Not yet available")
}

// ClearCounters clear max statistics counters (add 'all' for all counters)
func (h *HaproxyInstace) ClearCounters(all bool) error {
	toSend := "clear counters"
	if all {
		toSend = toSend + " all"
	}
	out, err := h.q(toSend)
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}
	return errors.New(out)
}

// ShowInfo report information about the running process
func (h *HaproxyInstace) ShowInfo() (map[string]string, error) {
	toReturn := map[string]string{}
	out, err := h.q("show info")
	if err != nil {
		return toReturn, err
	}

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			continue
		}
		toReturn[parts[0]] = parts[1]
	}

	return toReturn, nil
}

// ShowStat report counters for each proxy and server
func (h *HaproxyInstace) ShowStat() ([]map[string]string, error) {
	return h.qMap("show stat")
}

// ShowSchemaJSON report schema used for stats
func (h *HaproxyInstace) ShowSchemaJSON() (string, error) {
	return h.q("show schema json")
}

// DisableAgent disable agent checks
// The haproxy docs note that this function is depricated but that's not the case here
// Here we just use the recommended function
func (h *HaproxyInstace) DisableAgent(backend, server string) error {
	return h.Server(backend, server).Agent(true)
}

// DisableHealth disable health checks
// The haproxy docs note that this function is depricated but that's not the case here
// Here we just use the recommended function
func (h *HaproxyInstace) DisableHealth(backend, server string) error {
	return h.Server(backend, server).Health("down")
}

// DisableServer disable a server for maintenance
// The haproxy docs note that this function is depricated but that's not the case here
// Here we just use the recommended function
func (h *HaproxyInstace) DisableServer(backend, server string) error {
	return h.Server(backend, server).State("maint")
}

// EnableAgent enable agent checks
// The haproxy docs note that this function is depricated but that's not the case here
// Here we just use the recommended function
func (h *HaproxyInstace) EnableAgent(backend, server string) error {
	return h.Server(backend, server).Agent(true)
}

// EnableHealth enable health checks
// The haproxy docs note that this function is depricated but that's not the case here
// Here we just use the recommended function
func (h *HaproxyInstace) EnableHealth(backend, server, health string) error {
	return h.Server(backend, server).Health(health)
}

// EnableServer enable a disabled server
// The haproxy docs note that this function is depricated but that's not the case here
// Here we just use the recommended function
func (h *HaproxyInstace) EnableServer(backend, server string) error {
	return h.Server(backend, server).State("ready")
}

// SetMaxconnServer change a server's maxconn setting
func (h *HaproxyInstace) SetMaxconnServer(backend, server string, maxConn uint) error {
	if server == "" {
		return errors.New("server can't be an empty string")
	}
	out, err := h.q("set maxconn server " + h.Server(backend, server).server + " " + fmt.Sprintf("%v", maxConn))
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}
	return errors.New(out)
}

// ServerT is the response type for SetServer
type ServerT struct {
	q      func(subQuery string) (string, error) // run a "set server ..." command
	server string                                // just the backend/server string
}

// Server creates a instance of ServerT from where most things can be changed
// This returns a ServerT type that can be used for executing other queries, for example:
// h.Server("test-backend", "serv1").Addr("0.0.0.0")
// h.Server("test-backend", "serv1").Agent(true)
func (h *HaproxyInstace) Server(backend, server string) *ServerT {
	serv := backend + "/" + server
	toReturn := ServerT{
		q: func(subQuery string) (string, error) {
			return h.q("set server " + serv + " " + subQuery)
		},
		server: serv,
	}
	return &toReturn
}

// Addr <ip4 or ip6 address> [port <port>]
// Replace the current IP address of a server by the one provided.
// Optionnaly, the port can be changed using the 'port' parameter.
// Note that changing the port also support switching from/to port mapping
// (notation with +X or -Y), only if a port is configured for the health check.
func (s *ServerT) Addr(addr string, port ...string) error {
	query := "addr " + addr
	switch len(port) {
	case 0:
	case 1:
		query = query + " " + port[0]
	default:
		return errors.New("There can only be 0 or 1 ports defined")
	}

	out, err := s.q(query)
	if err != nil {
		return err
	}

	if strings.Contains(out, "changed from") {
		return nil
	}

	return errors.New(out)
}

// Agent [ up (true) | down (false) ]
// Force a server's agent to a new state. This can be useful to immediately
// switch a server's state regardless of some slow agent checks for example.
// Note that the change is propagated to tracking servers if any.
func (s *ServerT) Agent(upOrDown bool) error {
	toSend := "down"
	if upOrDown {
		toSend = "up"
	}

	out, err := s.q("agent " + toSend)
	if err != nil {
		return err
	}

	if strings.Contains(out, "changed from") {
		return nil
	}

	return errors.New(out)
}

// AgentAddr <addr>
// Change addr for servers agent checks. Allows to migrate agent-checks to
// another address at runtime. You can specify both IP and hostname, it will be
// resolved.
func (s *ServerT) AgentAddr(addr string) error {
	if addr == "" {
		return errors.New("addr can't be empty")
	}
	out, err := s.q("agent-addr " + addr)
	if err != nil {
		return err
	}

	if strings.Contains(out, "not enabled") || strings.Contains(out, "incorrect") {
		return errors.New(out)
	}

	return nil
}

// AgentSend <value>
// Change agent string sent to agent check target. Allows to update string while
// changing server address to keep those two matching.
func (s *ServerT) AgentSend(value string) error {
	if value == "" {
		return errors.New("value can't be empty")
	}

	out, err := s.q("agent-send " + value)
	if err != nil {
		return err
	}

	if strings.Contains(out, "not enabled") || strings.Contains(out, "cannot") {
		return errors.New(out)
	}

	return nil
}

// Health [ up | stopping | down ]
// Force a server's health to a new state. This can be useful to immediately
// switch a server's state regardless of some slow health checks for example.
// Note that the change is propagated to tracking servers if any.
func (s *ServerT) Health(health string) error {
	switch health {
	case "up", "stopping", "down":
	default:
		return errors.New("health has wrong value, must be \"up\", \"stopping\" or \"down\"")
	}

	out, err := s.q("health " + health)
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}
	return errors.New(out)
}

// CheckPort <port>
// Change the port used for health checking to <port>
func (s *ServerT) CheckPort(port string) error {
	out, err := s.q("check-port " + port)
	if err != nil {
		return err
	}
	if strings.Contains(out, "port updated") {
		return nil
	}
	return errors.New(out)
}

// State [ ready | drain | maint ]
// Force a server's administrative state to a new state. This can be useful to
// disable load balancing and/or any traffic to a server. Setting the state to
// "ready" puts the server in normal mode, and the command is the equivalent of
// the "enable server" command. Setting the state to "maint" disables any traffic
// to the server as well as any health checks. This is the equivalent of the
// "disable server" command. Setting the mode to "drain" only removes the server
// from load balancing but still allows it to be checked and to accept new
// persistent connections. Changes are propagated to tracking servers if any.
func (s *ServerT) State(state string) error {
	switch state {
	case "ready", "drain", "maint":
	default:
		return errors.New("state has wrong value, must be \"ready\", \"drain\" or \"maint\"")
	}

	out, err := s.q("state " + state)
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}
	return errors.New(out)
}

// Weight <weight>[%]
// Change a server's weight to the value passed in argument. This is the exact
// equivalent of the SetWeight function.
func (s *ServerT) Weight(newWeight string) error {
	if newWeight == "" {
		return errors.New("newWeight can't be an empty string")
	}
	out, err := s.q("weight " + newWeight)
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}
	return errors.New(out)
}

// FQDN <FQDN>
// Change a server's FQDN to the value passed in argument. This requires the
// internal run-time DNS resolver to be configured and enabled for this server.
func (s *ServerT) FQDN(fqdn string) error {
	if fqdn == "" {
		return errors.New("fqdn can't be empty")
	}

	out, err := s.q("fqdn " + fqdn)
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}

	return errors.New(out)
}

// GetWeight report a server's current weight
func (h *HaproxyInstace) GetWeight(backend, server string) (string, error) {
	if backend == "" || server == "" {
		return "", errors.New("Input can't be empty")
	}

	out, err := h.q("get weight " + backend + "/" + server)
	if err != nil {
		return "", err
	}

	return out, nil
}

// SetWeight change a server's weight
// The offical docs notes "set weight" as deprecated but thats not the case here
// Because this uses the newer "set server ..." instiad of "set weight ..."
func (h *HaproxyInstace) SetWeight(backend, server, setTo string) error {
	return h.Server(backend, server).Weight(setTo)
}

// SessionT is the response data when asking for the sessions
type SessionT struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Source  string `json:"source"`
	CPU     string `json:"cpu"`
	Latency string `json:"latency"`
	Age     string `json:"age"`
	RawRes  string `json:"rawRes"`
	Calls   string `json:"calls"`
	Expire  string `json:"expire"`
}

// ShowSess report the list of current sessions or dump this session
func (h *HaproxyInstace) ShowSess() ([]SessionT, error) {
	toReturn := []SessionT{}
	out, err := h.q("show sess")
	if err != nil {
		return toReturn, err
	}
	out = strings.TrimSpace(out)

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		toAdd := SessionT{RawRes: line}
		items := strings.Split(line, " ")
		for i, item := range items {
			if i == 0 && len(item) > 0 {
				toAdd.ID = strings.Replace(item, ":", "", 1)
				continue
			}

			nameAndVal := strings.Split(item, "=")
			if len(nameAndVal) < 2 {
				continue
			}
			name := nameAndVal[0]
			v := nameAndVal[1]

			switch name {
			case "proto":
				toAdd.Type = v
			case "src":
				toAdd.Source = v
			case "age":
				toAdd.Age = v
			case "calls":
				toAdd.Calls = v
			case "cpu":
				toAdd.CPU = v
			case "lat":
				toAdd.Latency = v
			case "exp":
				toAdd.Expire = v
			}
		}
		toReturn = append(toReturn, toAdd)
	}

	return toReturn, nil
}

// ShutdownSession kill a specific session
func (h *HaproxyInstace) ShutdownSession(id string) error {
	if id == "" {
		return errors.New("ID can't be empty")
	}

	out, err := h.q("shutdown session " + id)
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}

	return errors.New(out)
}

// ShutdownSessionsServer kill all sessions on a server
func (h *HaproxyInstace) ShutdownSessionsServer(backend, server string) error {
	out, err := h.q("shutdown sessions server " + h.Server(backend, server).server)
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}
	return errors.New(out)
}

// ClearTable remove an entry from a table
func (h *HaproxyInstace) ClearTable() error {
	return nil
}

// SetTable update or create a table entry's data
func (h *HaproxyInstace) SetTable(id string) error {
	return nil
}

// ShowTable report table usage stats or dump this table's contents
func (h *HaproxyInstace) ShowTable(id string) error {
	return nil
}

// DisableFrontend temporarily disable specific frontend
func (h *HaproxyInstace) DisableFrontend(frontend string) error {
	if frontend == "" {
		return errors.New("fontend can't be an empty string")
	}
	out, err := h.q("disable frontend " + frontend)
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}
	return errors.New(out)
}

// EnableFrontend re-enable specific frontend
func (h *HaproxyInstace) EnableFrontend(frontend string) error {
	if frontend == "" {
		return errors.New("fontend can't be an empty string")
	}
	out, err := h.q("enable frontend " + frontend)
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}
	return errors.New(out)
}

// SetMaxconnFrontend change a frontend's maxconn setting
func (h *HaproxyInstace) SetMaxconnFrontend(frontend string, maxConn uint) error {
	if frontend == "" {
		return errors.New("frontend can't be an empty string")
	}
	out, err := h.q("set maxconn frontend " + frontend + " " + fmt.Sprintf("%v", maxConn))
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}
	return errors.New(out)
}

// ShowServersState dump volatile server information (for backend)
func (h *HaproxyInstace) ShowServersState(backend string) ([]map[string]string, error) {
	toReturn := []map[string]string{}
	if backend == "" {
		return toReturn, errors.New("backend can't be empty")
	}
	out, err := h.q("show servers state " + backend)
	if err != nil {
		return nil, err
	}

	if strings.Contains(out, "Can't find backend") {
		return toReturn, errors.New(out)
	}

	lines := strings.Split(out, "\n")
	if len(lines) == 0 {
		return toReturn, errors.New("No output data")
	}
	if !strings.HasPrefix(lines[0], "1") {
		return toReturn, errors.New("unsupported \"show servers state ...\" output, only support version 1, recieved version " + lines[0])
	}

	out = strings.Join(lines[1:], "\n")
	toReturn, err = csvToArrMap(out, " ")
	if err != nil {
		return toReturn, err
	}

	return toReturn, nil
}

// ShowBackend list backends in the current running config
func (h *HaproxyInstace) ShowBackend() ([]map[string]string, error) {
	return h.qMap("show backend")
}

// ShutdownFrontend stop a specific frontend
func (h *HaproxyInstace) ShutdownFrontend(frontend string) error {
	if frontend == "" {
		return errors.New("frontend can't be an empty string")
	}

	out, err := h.q("shutdown frontend " + frontend)
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}

	return errors.New(out)
}

// SetDynamicCookieKeyBackend change a backend secret key for dynamic cookies
func (h *HaproxyInstace) SetDynamicCookieKeyBackend(backend, value string) error {
	if backend == "" {
		return errors.New("backend can't be an empty string")
	}
	out, err := h.q("set dynamic-cookie-key backend " + backend + " " + value)
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}
	return errors.New(out)
}

// DynamicCookieBackend enables (true) or disabled (false) dynamic cookies on a specific backend
func (h *HaproxyInstace) DynamicCookieBackend(backend string, setTo bool) error {
	if backend == "" {
		return errors.New("backend can't be an empty string")
	}

	prefix := "disable"
	if setTo {
		prefix = "enable"
	}
	out, err := h.q(prefix + " dynamic-cookie backend " + backend)
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}
	return errors.New(out)
}

// ShowStatResolvers dumps counters from all resolvers section and associated name servers
func (h *HaproxyInstace) ShowStatResolvers(id ...string) ([]map[string]string, error) {
	toEx := "show stat resolvers"
	switch len(id) {
	case 0:
	case 1:
		toEx = toEx + " " + id[0]
	default:
		return []map[string]string{}, errors.New("There can't be more than 1 IDs")
	}
	return h.qMap(toEx)
}

// SetMaxconnGlobal change the per-process maxconn setting
func (h *HaproxyInstace) SetMaxconnGlobal(maxConn uint) error {
	out, err := h.q("set maxconn global " + fmt.Sprintf("%v", maxConn))
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}
	return errors.New(out)
}

// SetRateLimit change a rate limiting value
func (h *HaproxyInstace) SetRateLimit(what string, value uint) error {
	switch what {
	case "connections", "http-compression", "sessions", "ssl-sessions":
	default:
		return errors.New("Unsupported \"what\", supported values: \"connections\", \"http-compression\", \"sessions\", \"ssl-sessions\"")
	}

	out, err := h.q("set rate-limit " + what + " global " + fmt.Sprintf("%v", value))
	if err != nil {
		return err
	}

	fmt.Println(out)

	return nil
}

// ShowEnv dump environment variables known to the process
func (h *HaproxyInstace) ShowEnv(name ...string) (map[string]string, error) {
	toReturn := map[string]string{}

	toEx := "show env"
	switch len(name) {
	case 0:
	case 1:
		toEx = toEx + " " + name[0]
	default:
		return toReturn, errors.New("Name can't have more than 1 entries")
	}
	out, err := h.q(toEx)
	if err != nil {
		return toReturn, err
	}

	failedLines := 0
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			failedLines++
			continue
		}
		toReturn[parts[0]] = parts[1]
	}

	if len(lines) < 5 || failedLines > 2 {
		return map[string]string{}, errors.New(out)
	}

	return toReturn, nil
}

// ShowCliSockets dump list of cli sockets
func (h *HaproxyInstace) ShowCliSockets() ([]map[string]string, error) {
	return h.qMap("show cli sockets", " ")
}

// AddACL add acl entry
func (h *HaproxyInstace) AddACL() error {
	return nil
}

// ClearACL clear the content of this acl
func (h *HaproxyInstace) ClearACL() error {
	return nil
}

// DelACL delete acl entry
func (h *HaproxyInstace) DelACL() error {
	return nil
}

// GetACL report the patterns matching a sample for an ACL
func (h *HaproxyInstace) GetACL() error {
	return nil
}

// ShowACL report available acls or dump an acl's contents
func (h *HaproxyInstace) ShowACL(id string) error {
	return nil
}

// AddMap add map entry
func (h *HaproxyInstace) AddMap() error {
	return nil
}

// ClearMap clear the content of this map
func (h *HaproxyInstace) ClearMap() error {
	return nil
}

// DelMap delete map entry
func (h *HaproxyInstace) DelMap() error {
	return nil
}

// GetMap report the keys and values matching a sample for a map
func (h *HaproxyInstace) GetMap() error {
	return nil
}

// SetMap modify map entry
func (h *HaproxyInstace) SetMap() error {
	return nil
}

// ShowMap report available maps or dump a map's contents
func (h *HaproxyInstace) ShowMap(id string) error {
	return nil
}

// PoolT is the data from 1 pool
type PoolT struct {
	Name      string `json:"name"`
	ID        string `json:"id"`        // TODO: find out what this value is
	Size      string `json:"size"`      // The size of 1 instance
	TotalSize string `json:"totalSize"` // The totalsize, .Used * .Size
	Used      uint   `json:"used"`
	Failures  uint   `json:"failures"`
	Users     uint   `json:"users"`
	Shared    bool   `json:"shared"`
	Raw       string `json:"raw"`
}

// ShowPools report information about the memory pools usage
func (h *HaproxyInstace) ShowPools() ([]PoolT, error) {
	toReturn := []PoolT{}

	out, err := h.q("show pools")
	if err != nil {
		return toReturn, err
	}

	if !strings.HasPrefix(out, "Dumping pools usage") {
		return toReturn, errors.New(out)
	}

	lines := strings.Split(out, "\n")
	preFix := "  - Pool "
	for _, line := range lines {
		if !strings.HasPrefix(line, preFix) {
			continue
		} else {
			line = strings.Replace(line, preFix, "", 1)
		}
		toAdd := PoolT{
			Raw: line,
		}

		parts := strings.Split(line, " ")
		currentPart := "procName"
		lastVal := ""
		for _, part := range parts {
			switch currentPart {
			case "procName":
				toAdd.Name = part
				currentPart = "size1"
			case "size1":
				toAdd.Size = strings.TrimLeft(part, "(")
				currentPart = "size2"
			case "size2":
				toAdd.Size = toAdd.Size + " " + strings.TrimRight(strings.TrimRight(part, "):"), ")")
				currentPart = ""
			case "totalSize1":
				toAdd.TotalSize = strings.TrimLeft(part, "(")
				currentPart = "totalSize2"
			case "totalSize2":
				toAdd.TotalSize = toAdd.TotalSize + " " + strings.TrimRight(strings.TrimRight(part, ")"), "),")
				currentPart = "used"
			case "used":
				i, err := strconv.Atoi(part)
				if err == nil {
					toAdd.Used = uint(i)
				}
				currentPart = ""
			case "id":
				if part != "," {
					toAdd.ID = part
					currentPart = ""
				}
			default:
				switch part {
				case "allocated":
					currentPart = "totalSize1"
				case "failures,", "failures":
					i, err := strconv.Atoi(lastVal)
					if err == nil {
						toAdd.Failures = uint(i)
					}
				case "users,", "users":
					i, err := strconv.Atoi(lastVal)
					if err == nil {
						toAdd.Users = uint(i)
					}
					currentPart = "id"
				case "[SHARED]":
					toAdd.Shared = true
				}
			}
			lastVal = part
		}

		toReturn = append(toReturn, toAdd)
	}

	return toReturn, nil
}

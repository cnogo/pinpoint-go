package agent

import "log"

const (
	SOCKSTAT_NONE                       byte = 1
	SOCKSTAT_BEING_CONNECT              byte = 2
	SOCKSTAT_CONNECTED                  byte = 3
	SOCKSTAT_CONNECT_FAILED             byte = 6
	SOCKSTAT_IGNORE                     byte = 9
	SOCKSTAT_RUN_WITHOUT_HANDSHAKE      byte = 10
	SOCKSTAT_RUN_SIMPLEX                byte = 11
	SOCKSTAT_RUN_DUPLEX                 byte = 12
	SOCKSTAT_BEING_CLOSE_BY_CLIENT      byte = 20
	SOCKSTAT_CLOSED_BY_CLIENT           byte = 22
	SOCKSTAT_UNEXPECTED_CLOSE_BY_CLIENT byte = 26
	SOCKSTAT_BEING_CLOSE_BY_SERVER      byte = 30
	SOCKSTAT_CLOSED_BY_SERVER           byte = 32
	SOCKSTAT_UNEXPECTED_CLOSE_BY_SERVER byte = 36
	SOCKSTAT_ERROR_UNKNOWN              byte = 40
	SOCKSTAT_ILLEGAL_STATE_CHANGE       byte = 41
	SOCKSTAT_ERROR_SYNC_STATE_SESSION   byte = 42
)

var sockStatCodeOp SocketStateCodeOp

type SocketStateCodeOp struct{}

func (p *SocketStateCodeOp) canChangeState(state byte, nextState byte) bool {
	switch nextState {
	case SOCKSTAT_BEING_CONNECT:
		return state == SOCKSTAT_NONE

	case SOCKSTAT_CONNECTED:
		return state == SOCKSTAT_NONE || state == SOCKSTAT_BEING_CONNECT ||
			state == SOCKSTAT_CONNECT_FAILED

	case SOCKSTAT_CONNECT_FAILED:
		return state == SOCKSTAT_BEING_CONNECT || state == SOCKSTAT_RUN_WITHOUT_HANDSHAKE ||
			state == SOCKSTAT_RUN_SIMPLEX || state == SOCKSTAT_RUN_DUPLEX

	case SOCKSTAT_IGNORE:
		return state == SOCKSTAT_CONNECTED

	case SOCKSTAT_RUN_WITHOUT_HANDSHAKE:
		return state == SOCKSTAT_CONNECTED

	case SOCKSTAT_RUN_SIMPLEX:
		return state == SOCKSTAT_RUN_WITHOUT_HANDSHAKE

	case SOCKSTAT_RUN_DUPLEX:
		return state == SOCKSTAT_RUN_WITHOUT_HANDSHAKE

	case SOCKSTAT_BEING_CLOSE_BY_CLIENT:
		return state == SOCKSTAT_RUN_WITHOUT_HANDSHAKE || state == SOCKSTAT_RUN_SIMPLEX ||
			state == SOCKSTAT_RUN_DUPLEX

	case SOCKSTAT_CLOSED_BY_CLIENT:
		return state == SOCKSTAT_BEING_CLOSE_BY_CLIENT

	case SOCKSTAT_UNEXPECTED_CLOSE_BY_CLIENT:
		return state == SOCKSTAT_CONNECTED || state == SOCKSTAT_RUN_WITHOUT_HANDSHAKE ||
			state == SOCKSTAT_RUN_SIMPLEX || state == SOCKSTAT_RUN_DUPLEX

	case SOCKSTAT_BEING_CLOSE_BY_SERVER:
		return state == SOCKSTAT_RUN_WITHOUT_HANDSHAKE || state == SOCKSTAT_RUN_SIMPLEX ||
			state == SOCKSTAT_RUN_DUPLEX

	case SOCKSTAT_CLOSED_BY_SERVER:
		return state == SOCKSTAT_BEING_CLOSE_BY_SERVER

	case SOCKSTAT_UNEXPECTED_CLOSE_BY_SERVER:
		return state == SOCKSTAT_RUN_WITHOUT_HANDSHAKE || state == SOCKSTAT_RUN_SIMPLEX ||
			state == SOCKSTAT_RUN_DUPLEX

	default:
		log.Println("error state change: ", state, " ", nextState)
		return false
	}
}

func (p *SocketStateCodeOp) isBeforeConnected(state byte) bool {
	switch state {
	case SOCKSTAT_NONE:
		return true
	case SOCKSTAT_BEING_CONNECT:
		return true
	}

	return false
}

func (p *SocketStateCodeOp) isRun(state byte) bool {
	switch state {
	case SOCKSTAT_RUN_WITHOUT_HANDSHAKE:
		return true
	case SOCKSTAT_RUN_SIMPLEX:
		return true
	case SOCKSTAT_RUN_DUPLEX:
		return true
	}

	return false
}

func (p *SocketStateCodeOp) isRunDuplex(state byte) bool {
	return state == SOCKSTAT_RUN_DUPLEX
}

func (p *SocketStateCodeOp) onClose(state byte) bool {
	switch state {
	case SOCKSTAT_BEING_CLOSE_BY_CLIENT:
		return true
	case SOCKSTAT_BEING_CLOSE_BY_SERVER:
		return true
	}

	return false
}

func (p *SocketStateCodeOp) isClosed(state byte) bool {
	switch state {
	case SOCKSTAT_CLOSED_BY_CLIENT:
		return true
	case SOCKSTAT_CLOSED_BY_SERVER:
		return true
	case SOCKSTAT_UNEXPECTED_CLOSE_BY_SERVER:
		return true
	case SOCKSTAT_ERROR_UNKNOWN:
		return true
	case SOCKSTAT_ILLEGAL_STATE_CHANGE:
		return true
	case SOCKSTAT_ERROR_SYNC_STATE_SESSION:
		return true
	}

	return false
}

func (p *SocketStateCodeOp) isError(state byte) bool {
	switch state {
	case SOCKSTAT_ILLEGAL_STATE_CHANGE:
		return true
	case SOCKSTAT_ERROR_UNKNOWN:
		return true
	}

	return false
}

type SocketStateChangeResult struct {
	changed           bool
	beforeState       byte
	currentState      byte
	updateWantedState byte
}

func NewSocketStateChange(changed bool, beforeState, currentState, updateWantedState byte) *SocketStateChangeResult {
	return &SocketStateChangeResult{
		changed:           changed,
		beforeState:       beforeState,
		currentState:      currentState,
		updateWantedState: updateWantedState,
	}
}

func (p *SocketStateChangeResult) isChange() bool {
	return p.changed
}

func (p *SocketStateChangeResult) getBeforeState() byte {
	return p.beforeState
}

func (p *SocketStateChangeResult) getCurrentState() byte {
	return p.currentState
}

func (p *SocketStateChangeResult) getUpdateWantedState() byte {
	return p.updateWantedState
}

type SocketState struct {
	beforeState  byte
	currentState byte
}

func NewSocketState() *SocketState {
	return &SocketState{
		beforeState:  SOCKSTAT_NONE,
		currentState: SOCKSTAT_NONE,
	}
}

func (p *SocketState) to(nextState byte) SocketStateChangeResult {
	enable := sockStatCodeOp.canChangeState(p.currentState, nextState)
	changed := false
	if enable {
		p.beforeState = p.currentState
		p.currentState = nextState
		changed = true
	}

	return SocketStateChangeResult{
		changed:           changed,
		beforeState:       p.beforeState,
		currentState:      p.currentState,
		updateWantedState: nextState,
	}
}

func (p *SocketState) toBeingConnect() SocketStateChangeResult {
	return p.to(SOCKSTAT_BEING_CONNECT)
}

func (p *SocketState) toConnected() SocketStateChangeResult {
	return p.to(SOCKSTAT_CONNECTED)
}

func (p *SocketState) toConnectFailed() SocketStateChangeResult {
	return p.to(SOCKSTAT_CONNECT_FAILED)
}

func (p *SocketState) toIgnore() SocketStateChangeResult {
	return p.to(SOCKSTAT_IGNORE)
}

func (p *SocketState) toRunWithoutHandShake() SocketStateChangeResult {
	return p.to(SOCKSTAT_RUN_WITHOUT_HANDSHAKE)
}

func (p *SocketState) toRunSimplex() SocketStateChangeResult {
	return p.to(SOCKSTAT_RUN_SIMPLEX)
}

func (p *SocketState) toRunDuplex() SocketStateChangeResult {
	return p.to(SOCKSTAT_RUN_DUPLEX)
}

func (p *SocketState) toBeingCloseByClient() SocketStateChangeResult {
	return p.to(SOCKSTAT_BEING_CLOSE_BY_CLIENT)
}

func (p *SocketState) toClosedByClient() SocketStateChangeResult {
	return p.to(SOCKSTAT_CLOSED_BY_CLIENT)
}

func (p *SocketState) toUnexpectedCloseByClient() SocketStateChangeResult {
	return p.to(SOCKSTAT_UNEXPECTED_CLOSE_BY_CLIENT)
}

func (p *SocketState) toBeingCloseByServer() SocketStateChangeResult {
	return p.to(SOCKSTAT_BEING_CLOSE_BY_SERVER)
}

func (p *SocketState) toClosedByServer() SocketStateChangeResult {
	return p.to(SOCKSTAT_CLOSED_BY_SERVER)
}

func (p *SocketState) toUnexpectedCloseByServer() SocketStateChangeResult {
	return p.to(SOCKSTAT_UNEXPECTED_CLOSE_BY_SERVER)
}

func (p *SocketState) toUnknownError() SocketStateChangeResult {
	return p.to(SOCKSTAT_ERROR_UNKNOWN)
}

func (p *SocketState) toSyncStateSessionError() SocketStateChangeResult {
	return p.to(SOCKSTAT_ERROR_SYNC_STATE_SESSION)
}

func (p *SocketState) getCurrentState() byte {
	return p.currentState
}

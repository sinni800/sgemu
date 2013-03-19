package GameServer

const (
	//CM - Client Message
	//SM - Server Message
	//CSM - Client and Server Message

	CSM_GAME_ENTER      = 0xE
	CM_DISCONNECT       = 0x15
	SM_PLAYER_STATS     = 0x15
	SM_UNIT_STAT        = 0x16
	SM_MAP_LOAD         = 0x17
	SM_PLAYER_LEAVE     = 0x19
	SM_INVENTORY_UPDATE = 0x1F
	SM_PLAYER_APPEAR    = 0x20
	CM_UNIT_EDIT        = 0x22
	CM_PROFILE          = 0x23
	CM_LEAVE_PROFILE    = 0x25
	SM_PROFILE          = 0x27

	CSM_CHAT             = 0x2A
	CM_SHOP_REQUEST      = 0x2D
	CM_MAPCHANGE_REQUEST = 0x2E
	SM_SHOP_RESPONSE     = 0x2F

	CSM_LAB_ENTER = 0x33

	CSM_PLAYER_NAME       = 0x47
	SM_PLAYER_NAME_BATTLE = 0x48

	CM_PING = 0x58
	SM_PONG = 0x56

	CSM_MOVE = 0x80
)

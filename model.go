package discordmarket

type Config struct {
	AppID               string
	Token               string
	DatabaseURL         string
	AlertChannel        string
	GuildID             string
	LotteryChannelID    string
	Currency            string
	EmojiReaction       string
	EmojiRole           string
	EmojiLottery        string
	EmojiLootbox        string
	EmojiPrivateChannel string
	ReactionPrice       int
	StatusPrice         int
	LootboxPrice        int
	LotteryPrice        int
	FarmRate            int
	VoiceRate           int
	VoiceDelta          int
}
type product struct {
	ID     int
	Dial   int
	Name   string
	Type   string
	RoleID string
	Price  int
}

type spender struct {
	UserID string
	Spent  int
}

type voicer struct {
	UserID string
	Spent  int64
}

type expiredRole struct {
	UserID string
	RoleID string
}

type order struct {
	OrderID    int
	CustomerID string
	Product    *product
	Expires    int64
	IsHidden   bool
}

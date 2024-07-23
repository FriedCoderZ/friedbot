package friedbot

const (
	SegmentTypeOther = iota
	SegmentTypeText
	SegmentTypeFace
	SegmentTypeImage
	SegmentTypeRecord
	SegmentTypeVideo
	SegmentTypeAt
	SegmentTypeRps
	SegmentTypeDice
	SegmentTypeShake
	SegmentTypePoke
	SegmentTypeAnonymous
	SegmentTypeShare
	SegmentTypeContact
	SegmentTypeLocation
	SegmentTypeMusic
	SegmentTypeReply
	SegmentTypeForward
	SegmentTypeNode
	SegmentTypeXML
	SegmentTypeJSON
)

const (
	MsgTypeOther = iota
	MsgTypePrivate
	MsgTypeGroup
)

type Message struct {
	ID       int64
	Type     int
	Segments []Segment
	User     *User
	GroupID  int64
	Time     int64
}

type User struct {
	ID              int64
	Nickname        string
	Card            string
	Sex             string
	Age             int
	Area            string
	JoinTime        int32
	LastSentTime    int32
	Level           int
	Role            string
	Unfriendly      bool
	Title           string
	TitleExpireTime int32
}

type Group struct {
	ID             int64
	Name           string
	MemberCount    int32
	MaxMemberCount int32
}

type Segment struct {
	Type int
	Data any
}

type SegmentText struct {
	Text string
}

type SegmentFace struct {
	ID string
}

type SegmentImage struct {
	File string
}

type SegmentRecord struct {
	File string
}

type SegmentVideo struct {
	File string
}

type SegmentAt struct {
	QQ string
}

type SegmentRPS struct{}

type SegmentDice struct{}

type SegmentShake struct{}

type SegmentPoke struct {
	Type string
	ID   string
}

type SegmentAnonymous struct{}

type SegmentShare struct {
	URL   string
	Title string
}

type SegmentContact struct {
	Type string
	ID   string
}

type SegmentLocation struct {
	Lat string
	Lon string
}

type SegmentMusic struct {
	Type  string
	ID    string
	URL   string
	Audio string
	Title string
}

type SegmentReply struct {
	ID string
}

type SegmentForward struct {
	ID string
}

type SegmentNode struct {
	ID string
}

type SegmentXML struct {
	Data string
}

type SegmentJSON struct {
	Data string
}

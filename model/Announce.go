package model

type AnnounceRequest struct {
	Agent           string `header:"User-Agent"`                 // 客户端
	OriInfoHash     string `form:"info_hash" binding:"required"` // 客户端上报的种子哈希码
	OriPeerID       string `form:"peer_id" binding:"required"`   // 客户端上报的同伴ID
	InfoHash        string ``                                    // 转换过后的哈希码
	PeerID          string ``                                    // 转换过后的同伴ID
	Port            uint16 `form:"port" binding:"min=1,max=65535"`
	DownloadedBytes uint   `form:"downloaded"`
	UploadedBytes   uint   `form:"uploaded"`
	LeftBytes       uint   `form:"left"`
	Event           string `form:"event"`
	IP              string `form:"ip" binding:"omitempty,ip"`
	IPv4            string `form:"ipv4" binding:"omitempty,ip"`
	IPv6            string `form:"ipv6" binding:"omitempty,ip"`
	Compact         uint8  `form:"compact" binding:"omitempty,min=0,max=1"`
	NoPeerID        uint8  `form:"no_peer_id" binding:"omitempty,min=0,max=1"`
	NumWanted       uint8  `form:"numwant" binding:"omitempty"`
}

type PeerItem struct {
	IP     string `bencode:"ip"`
	Port   uint   `bencode:"port"`
	PeerID string `bencode:"peer id,omitempty"`
}

type AnnounceResult struct {
	Interval     uint        `bencode:"interval"`     // 请求间隔（单位：秒）
	MinInterval  uint        `bencode:"min interval"` // 已完成下载总数
	SeederCount  uint        `bencode:"complete"`     // 当前做种数量
	LeecherCount uint        `bencode:"incomplete"`   // 正在下载数量
	Peers        []*PeerItem `bencode:"peers"`        // 同伴
}

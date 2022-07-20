package handler

import (
	"MoeTracker/model"
	"MoeTracker/redis"
	"MoeTracker/utils"
	"fmt"
	"github.com/anacrolix/torrent/bencode"
	"github.com/gin-gonic/gin"
	"net"
	"strconv"
	"strings"
)

func ReceiveAnnounce(c *gin.Context) {
	redisDb := redis.RedisDb

	var announceRequest model.AnnounceRequest
	err := c.ShouldBindQuery(&announceRequest)
	if err != nil {
		fmt.Println("解析失败:" + err.Error())
	}
	if err := fillAgent(c, &announceRequest); err != nil {
	}
	if err := seedIP(c, &announceRequest); err != nil {
	}
	announceRequest.InfoHash = utils.RestoreToHexString(announceRequest.OriInfoHash)
	announceRequest.PeerID = utils.RestoreToHexString(announceRequest.OriPeerID)

	if announceRequest.Event == "stopped" {
		redisDb.Del(announceRequest.InfoHash, announceRequest.PeerID+":"+announceRequest.IPv4+":"+strconv.Itoa(int(announceRequest.Port)))
		return
	}

	// 将自己加入Redis
	redisDb.SAdd(announceRequest.InfoHash, announceRequest.PeerID+":"+announceRequest.IPv4+":"+strconv.Itoa(int(announceRequest.Port)))
	// 从Redis获取其他用户
	ret, _ := redisDb.SMembers(announceRequest.InfoHash).Result()
	var peers []*model.PeerItem
	for _, v := range ret {
		peerID := strings.Split(v, ":")[0]
		ip := strings.Split(v, ":")[1]
		port, _ := strconv.ParseUint(strings.Split(v, ":")[2], 10, 32)
		if announceRequest.PeerID != peerID {
			peers = append(peers, &model.PeerItem{
				PeerID: utils.RestoreToByteString(peerID),
				IP:     ip,
				Port:   uint(port),
			})
		}
	}

	c.String(200, string(bencode.MustMarshal(model.AnnounceResult{
		Interval:     30,
		MinInterval:  2,
		SeederCount:  1,
		LeecherCount: 1,
		Peers:        peers,
	})))
}

// 处理 IP
// 优先级：IP > IPv4 > IPv6 > ctx.ip
func seedIP(ctx *gin.Context, request *model.AnnounceRequest) error {
	// 如果上报的 IPv4 地址有误，则清空
	if !IsIPv4(request.IPv4) {
		request.IPv4 = ""
	}

	// 如果上报的 IPv6 地址有误，则清空
	if !IsIPv6(request.IPv6) {
		request.IPv6 = ""
	}

	// 如果上报的 IP 地址有效，则覆盖对应的 IPv4/IPv6 地址
	if IsIPv4(request.IP) {
		request.IPv4 = request.IP
	}
	if IsIPv6(request.IP) {
		request.IPv6 = request.IP
	}

	// 如果均为空，则使用客户端地址填充
	if request.IPv4 == "" && request.IPv6 == "" {
		clientIP := ctx.ClientIP()
		if IsIPv4(clientIP) {
			request.IPv4 = clientIP
		} else {
			request.IPv6 = clientIP
		}
	}

	return nil
}

func IsIPv4(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && strings.Contains(ipAddr, ".")
}

func IsIPv6(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && strings.Contains(ipAddr, ":")
}

func fillAgent(ctx *gin.Context, request *model.AnnounceRequest) error {
	// Gin Bind 可能失效
	if request.Agent == "" {
		request.Agent = ctx.Request.Header.Get("User-Agent")
	}

	return nil
}

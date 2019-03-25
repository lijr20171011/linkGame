package linkGame

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Chess ...
// 棋子
type Chess struct {
	X      int // 对应列 col
	Y      int // 对应行 row
	Status int // 0代表消除
}

type ChessPanel [][]int

var cRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// NewChessPanel ...
// 创建棋盘 row:行数 col:列数 kinds: 棋子种类
func NewChessPanel(row, col, kinds int) *ChessPanel {
	cp := initChessPanel(row+2, col+2)
	totalCount := row * col
	for i := 0; i < totalCount; {
		chessKind := cRand.Intn(kinds) + 1
		for cnt := 0; cnt < 2; cnt++ {
			for {
				x := cRand.Intn(col)
				y := cRand.Intn(row)
				if cp[y+1][x+1] == 0 {
					cp[y+1][x+1] = chessKind
					i++
					break
				}
			}
		}
	}
	return &cp
}

// 初始化棋盘矩阵
func initChessPanel(row, col int) ChessPanel {
	cp := make([][]int, row)
	for i := range cp {
		cp[i] = make([]int, col)
	}
	return cp
}

// 判断坐标是否存在
func (cp *ChessPanel) IsPointExist(x, y int) bool {
	// 校验cp
	if cp == nil {
		return false
	}
	// 校验cp行
	row := len(*cp)
	if row <= 2 {
		return false
	}
	// 校验cp列
	col := len((*cp)[0])
	if col <= 2 {
		return false
	}
	// 校验 x 坐标
	if x < 1 || x > col-1-1 {
		return false
	}
	// 校验 y 坐标
	if y < 1 || y > row-1-1 {
		return false
	}
	return true
}

// 判断是否为棋子
func (cp *ChessPanel) IsChessExist(x, y int) (int, bool) {
	// 判断位置是否存在
	if !cp.IsPointExist(x, y) {
		return 0, false
	}
	// 判断是否为空
	if (*cp)[y][x] == 0 {
		return 0, false
	}
	return (*cp)[y][x], true
}

// 从扫描字符串中解析Chess
func (cp *ChessPanel) ParseChess(input string) (Chess, error) {
	var c Chess
	var err error
	// 解析point
	positions := strings.Split(input, ",")
	if len(positions) != 2 {
		return c, errors.New("x,y 需以英文逗号(,)隔开")
	}
	// 解析x
	c.X, err = strconv.Atoi(positions[0])
	if err != nil {
		return c, errors.New("解析x坐标异常;" + err.Error())
	}
	// 解析y
	c.Y, err = strconv.Atoi(positions[1])
	if err != nil {
		return c, errors.New("解析y坐标异常;" + err.Error())
	}
	// 判断point是否有效
	status, exist := cp.IsChessExist(c.X, c.Y)
	if !exist {
		return c, errors.New("节点异常")
	}
	c.Status = status
	return c, nil
}

// 消除
func (cp *ChessPanel) Offset(x1, y1, x2, y2 int) bool {
	// 判断point1是否有效
	p1, exist := cp.IsChessExist(x1, y1)
	if !exist {
		return false
	}
	// 判断point2是否有效
	p2, exist := cp.IsChessExist(x2, y2)
	if !exist {
		return false
	}
	// 判断是否相等
	if p1 != p2 {
		return false
	}
	// 消除
	cp.setChess(x1, y1, 0)
	cp.setChess(x2, y2, 0)
	return true
}

// 设置
func (cp *ChessPanel) setChess(x, y, value int) {
	(*cp)[y][x] = value
}

// 判断是否连通
func (cp *ChessPanel) IsLinked(c1, c2 Chess) bool {
	// 是否0折连通
	if cp.isLinked0(c1.X, c1.Y, c2.X, c2.Y) {
		return true
	}
	// 是否一折连通
	if cp.isLinked1(c1.X, c1.Y, c2.X, c2.Y) {
		return true
	}
	// 是否二折连通
	if cp.isLinked2(c1.X, c1.Y, c2.X, c2.Y) {
		return true
	}
	return false
}

// 0折连通
func (cp *ChessPanel) isLinked0(x1, y1, x2, y2 int) bool {
	// 两点是否在一条直线上
	if x1 != x2 && y1 != y2 {
		return false
	}
	// 判断x是否相同
	if x1 == x2 {
		// 在同一列,遍历y轴
		if y1 < y2 {
			for y := y1 + 1; y <= y2; y++ {
				// 找到y2返回
				if y == y2 {
					return true
				}
				// 遇到障碍结束遍历
				if !cp.IsEmpty(x1, y) {
					break
				}
			}
		} else {
			for y := y1 - 1; y >= y2; y-- {
				// 找到y2返回
				if y == y2 {
					return true
				}
				// 遇到障碍结束遍历
				if !cp.IsEmpty(x1, y) {
					break
				}
			}
		}
	} else if y1 == y2 {
		// 在同一行,遍历x轴
		if x1 < x2 {
			for x := x1 + 1; x <= x2; x++ {
				// 找到x2返回
				if x == x2 {
					return true
				}
				// 遇到障碍结束遍历
				if !cp.IsEmpty(x, y1) {
					break
				}
			}
		} else {
			for x := x1 - 1; x >= x2; x-- {
				// 找到x2返回
				if x == x2 {
					return true
				}
				// 遇到障碍结束遍历
				if !cp.IsEmpty(x, y1) {
					break
				}
			}
		}
	}
	return false
}

// 1折连通
func (cp *ChessPanel) isLinked1(x1, y1, x2, y2 int) bool {
	// 判断拐点1(x1,y2)是否0折连通起点终点
	if cp.IsEmpty(x1, y2) && cp.isLinked0(x1, y2, x1, y1) && cp.isLinked0(x1, y2, x2, y2) {
		return true
	}
	//  判断拐点1(x2,y1)是否0折连通起点终点
	if cp.IsEmpty(x2, y1) && cp.isLinked0(x2, y1, x1, y1) && cp.isLinked0(x2, y1, x2, y2) {
		return true
	}
	return false
}

// 2折连通
func (cp *ChessPanel) isLinked2(x1, y1, x2, y2 int) bool {
	// 以(x1,1)为起点遍历的上下左右节点 寻找可以一折连通的路径
	// 向左遍历(x 减少)
	for x := x1 - 1; x >= 0; x-- {
		// 遇到障碍放弃这条路
		if !cp.IsEmpty(x, y1) {
			break
		}
		// 判断当前节点是否可以一折连通终点
		if cp.isLinked1(x, y1, x2, y2) {
			return true
		}
	}
	// 向右遍历(x 增加)
	for x := x1 + 1; x <= len((*cp)[0])-1; x++ {
		// 遇到障碍放弃这条路
		if !cp.IsEmpty(x, y1) {
			break
		}
		// 判断当前节点是否可以一折连通终点
		if cp.isLinked1(x, y1, x2, y2) {
			return true
		}
	}
	// 向上遍历(y 减少)
	for y := y1 - 1; y >= 0; y-- {
		// 遇到障碍放弃这条路
		if !cp.IsEmpty(x1, y) {
			break
		}
		// 判断当前节点是否可以一折连通终点
		if cp.isLinked1(x1, y, x2, y2) {
			return true
		}
	}
	// 向下遍历(y 增加)
	for y := y1 + 1; y <= len(*cp)-1; y++ {
		// 遇到障碍放弃这条路
		if !cp.IsEmpty(x1, y) {
			break
		}
		// 判断当前节点是否可以一折连通终点
		if cp.isLinked1(x1, y, x2, y2) {
			return true
		}
	}
	return false
}

// 获取节点值
func (cp *ChessPanel) GetValue(x, y int) int {
	if x < 0 || x >= len((*cp)[0]) {
		panic("x out of index")
	}
	if y < 0 || y >= len(*cp) {
		panic("y out of index")
	}
	return (*cp)[y][x]
}

// 是否为空
func (cp *ChessPanel) IsEmpty(x, y int) bool {
	if cp.GetValue(x, y) == 0 {
		return true
	}
	return false
}

// 打印棋盘信息
func (cp ChessPanel) String() string {
	// 校验行
	row2 := len(cp)
	if row2 < 2 {
		return ""
	}
	// 校验列
	col2 := len(cp[0])
	if col2 < 2 {
		return ""
	}
	// 单起一行,整理列下标
	info := "    列    "
	for i := 1; i <= col2-2; i++ {
		info += fmt.Sprintf("% 4d ", i)
	}
	info += "\n\n"

	// 添加行下标和内容
	for i, v := range cp {
		if i > 0 && i < row2-1 {
			// 添加行下标
			info += fmt.Sprintf("% 3d 行    ", i)
			for j, vv := range v {
				if j > 0 && j < col2-1 {
					// 打印内容
					info += fmt.Sprintf("% 4d ", vv)
				}
			}
			// 一行结束
			info += "\n"
		}
	}

	return info
}

// 用时统计
func TimeCost(tStart, tEnd time.Time) string {
	timeCost := ""
	duration := tEnd.Sub(tStart)
	// 时
	tHour := int(duration / time.Hour)
	if tHour > 0 {
		timeCost += fmt.Sprintf("%dh", tHour)
		duration -= time.Duration(tHour) * time.Hour
	}
	// 分
	tMinute := int(duration / time.Minute)
	if tMinute > 0 {
		timeCost += fmt.Sprintf("%dm", tMinute)
		duration -= time.Duration(tMinute) * time.Minute
	}
	// 秒
	timeCost += fmt.Sprintf("%.4fs", duration.Seconds())
	return timeCost
}

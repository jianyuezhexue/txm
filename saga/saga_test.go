package saga

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test0(t *testing.T) {
	assert.Equal(t, "1", "1")
}

func TestSaga(t *testing.T) {

	// 0. 模拟入口函数传参
	req := "1"

	// 1. 实例化Saga事务管理器
	ctx := context.Background()
	saga := NewSaga(ctx, "user_001")

	// 2. 注册入库回传相关函数
	saga.Register("func1", func(txm *Saga) (any, error) {
		// 正向执行函数
		res, err := EntryReturn(req)
		if err != nil {
			return nil, err
		}

		// 预埋逆向函数使用到的值
		err = txm.SetVal("操作前的最近入库时间", res)
		if err != nil {
			return nil, err
		}

		return res, nil
	}, func(txm *Saga) (any, error) {

		// 获取正向函数预埋好的值
		rollBackTime, err := txm.GetVal("操作前的最近入库时间")
		if err != nil {
			return nil, err
		}

		// 逆向补偿函数
		res, err := EntryReturnRollBack(req, rollBackTime.(string))
		if err != nil {
			return nil, err
		}

		return res, nil
	})

	// 6. 执行提交
	saga.Commit()

	// 7. 结果断言
	assert.Equal(t, "StatusSucceed", "StatusSucceed")
}

// 函数1执行:入库回传逻辑处理
func EntryReturn(req1 string) (string, error) {
	return "", nil
}

// 函数1注册具柄函数

// 函数1回滚：入库回传处理补偿
func EntryReturnRollBack(req2 string, time string) (string, error) {
	return "", nil
}

// 函数2执行:发货单新建
func NewDeliver(req3 string) (string, error) {
	return "", nil
}

// 函数2回滚：删除发货单
func NewDeliverRollBack(req4 string) (string, error) {
	return "", nil
}

// 函数3执行: 销售单状态更新
func SaleOrderReturn(req5 string) (string, error) {
	return "", nil
}

// 函数3回滚：销售单撤销状态更改
func SaleOrderReturnRollBack(req5 string) (string, error) {
	return "", nil
}

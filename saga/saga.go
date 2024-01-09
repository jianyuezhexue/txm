package saga

import (
	"context"
	"fmt"
	"sync"

	"github.com/gomodule/redigo/redis"
	"github.com/jianyuezhexue/txm/util"
	"gorm.io/gorm"
)

// Saga 分布式事务接口定义
type SagaTxManager interface {
	singleLock() error                         // 单用户锁
	SetVal(string, any) error                  // 设置值
	GetVal(string) (any, error)                // 获取值
	Register(string, SagaFunc, SagaFunc) error // 注册函数
	checkAndRepair() error                     // 注册函数检查&基本信息补全
	Test() error                               // 调试测试
	Commit() error                             // 执行提交
	Transaction(SagaFunc) error                // 事务执行
}

// 注册函数
type SagaFunc func(txm *Saga) (any, error)

// 事务实例结构体
type Saga struct {
	ctx              context.Context        // 上下文
	opts             *Options               // 配置选项｜超时控制，轮询时间间隔
	db               *gorm.DB               // DB连接
	redis            redis.Conn             // Redis连接
	val              map[string]interface{} // 中间值传递
	executeFuncPool  map[string]SagaFunc    // 执行函数池
	roolBackFuncPool map[string]SagaFunc    // 补偿函数池
	errArr           []string               // 异常待回滚函数
	errChan          chan (error)           // 异常错误
}

// 实例化 Saga
func NewSaga(ctx context.Context, opts ...Option) SagaTxManager {
	txManager := &Saga{
		ctx:              ctx,
		opts:             &Options{},
		val:              map[string]interface{}{},
		executeFuncPool:  map[string]SagaFunc{},
		roolBackFuncPool: map[string]SagaFunc{},
		errArr:           []string{},
		errChan:          make(chan error),
	}

	// 配置参数
	for _, opt := range opts {
		opt(txManager.opts)
	}
	repair(txManager.opts)

	// 连接Db
	txManager.db = util.InitGorm()
	txManager.redis = util.GetRedisConn()

	return txManager
}

// 单用户锁
func (s *Saga) singleLock() error {
	return nil
}

// 设置值
func (s *Saga) SetVal(key string, val any) error {
	_, ok := s.val[key]
	if ok {
		return fmt.Errorf("key[%v],已经设置了值,请检查！", key)
	}
	s.val[key] = val
	return nil
}

// 获取值
func (s *Saga) GetVal(key string) (any, error) {
	val, ok := s.val[key]
	if !ok {
		return nil, fmt.Errorf("key[%v],不存在,请检查！", key)
	}
	return val, nil
}

// 注册组件方法
func (s *Saga) Register(funcName string, executeFunc SagaFunc, roolBackFunc SagaFunc) error {
	_, ok := s.executeFuncPool[funcName]
	if ok {
		return fmt.Errorf("正向函数[%v]重复注册,请检查", funcName)
	}
	_, ok2 := s.roolBackFuncPool[funcName]
	if ok2 {
		return fmt.Errorf("补偿函数[%v]重复注册,请检查", funcName)
	}
	s.executeFuncPool[funcName] = executeFunc
	s.roolBackFuncPool[funcName] = roolBackFunc
	return nil
}

// 校验&修复补全
func (s *Saga) checkAndRepair() error {
	return nil
}

// 调试测试
func (s *Saga) Test() error {
	return nil
}

// 提交事务
func (s *Saga) Commit() error {
	// 释放资源
	defer close(s.errChan)
	defer s.redis.Close()

	// 用户+函数名 单用户加锁

	// 创建唯一事务Id

	// 开启本地事务
	s.db.Transaction(func(tx *gorm.DB) error {

		// 超时控制

		// 并发开始
		var wg sync.WaitGroup

		// 并发执行正向函数
		for funcName, itemFunc := range s.executeFuncPool {
			wg.Add(1)

			// 悬挂校验
			// todo

			// 正常执行正向操作
			res, err := itemFunc(s)
			wg.Done()
			fmt.Println(res)
			if err != nil {
				s.errChan <- err
				s.errArr = append(s.errArr, funcName)
				break
			}
		}

		wg.Wait()

		//  全部执行成功
		if len(s.errChan) == 0 {
			return nil
		}

		// 有异常执行逆向函数
		for _, funcName := range s.errArr {
			fn := s.roolBackFuncPool[funcName]
			// 空补偿校验
			// todo

			// 正常执行补偿
			res, err := fn(s)
			if err != nil {
				// 记录日志 & 记录到重试数组中
				fmt.Println(err.Error())
				fmt.Println(res)
			}
		}

		return nil
	})

	// 返回异常信息
	err := <-s.errChan
	return err
}

// 事务执行
func (s *Saga) Transaction(SagaFunc) error {
	return nil
}

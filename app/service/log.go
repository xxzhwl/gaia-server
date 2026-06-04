package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/xxzhwl/gaia"
	esclient "github.com/xxzhwl/gaia/components/es"
)

// LogService 日志查询服务
type LogService struct {
	esClient *esclient.Client
	initOnce sync.Once
	initErr  error
}

var logService *LogService

func init() {
	logService = &LogService{}
}

func GetLogService() *LogService {
	return logService
}

// initClient 延迟初始化 ES 客户端
func (s *LogService) initClient() (*esclient.Client, error) {
	s.initOnce.Do(func() {
		// 检查是否配置了ES
		esAddress := gaia.GetSafeConfString("Framework.ES.Address")
		if esAddress == "" {
			s.initErr = fmt.Errorf("ES未配置")
			return
		}

		client, err := esclient.NewFrameWorkEs()
		if err != nil {
			s.initErr = fmt.Errorf("初始化ES客户端失败: %w", err)
			return
		}

		s.esClient = client
	})

	if s.initErr != nil {
		return nil, s.initErr
	}

	return s.esClient, nil
}

// SortKv 排序条件
type SortKv struct {
	Name string `json:"name"`
	Desc bool   `json:"desc"`
}

// OpArg 查询条件
type OpArg struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
	Op    string `json:"op"`
}

// QueryES 查询ES日志
func (s *LogService) QueryES(index string, from, size int, must, should, not []OpArg, sorts []SortKv) (map[string]any, error) {
	client, err := s.initClient()
	if err != nil {
		return nil, err
	}

	// 转换为es包中的类型
	esSorts := make([]esclient.SortKv, 0, len(sorts))
	for _, sort := range sorts {
		esSorts = append(esSorts, esclient.SortKv{
			Name: sort.Name,
			Desc: sort.Desc,
		})
	}

	esMust := make([]esclient.OpArg, 0, len(must))
	for _, arg := range must {
		esMust = append(esMust, esclient.OpArg{
			Key:   arg.Key,
			Value: arg.Value,
			Op:    arg.Op,
		})
	}

	esShould := make([]esclient.OpArg, 0, len(should))
	for _, arg := range should {
		esShould = append(esShould, esclient.OpArg{
			Key:   arg.Key,
			Value: arg.Value,
			Op:    arg.Op,
		})
	}

	esNot := make([]esclient.OpArg, 0, len(not))
	for _, arg := range not {
		esNot = append(esNot, esclient.OpArg{
			Key:   arg.Key,
			Value: arg.Value,
			Op:    arg.Op,
		})
	}

	// 构建查询参数
	arg := esclient.SimpleSearchArg{
		Index:  index,
		From:   int64(from),
		Size:   int64(size),
		Sorts:  esSorts,
		Must:   esMust,
		Should: esShould,
		Not:    esNot,
	}

	// 执行查询
	result, err := client.SimpleSearch(arg)
	if err != nil {
		return nil, fmt.Errorf("ES查询失败: %w", err)
	}

	// 构建符合前端期望的返回格式
	hits := make([]map[string]any, 0, len(result.Hits.Data))
	for _, hit := range result.Hits.Data {
		hits = append(hits, map[string]any{
			"_id":     hit.Id,
			"_score":  hit.Score,
			"_source": hit.Detail,
		})
	}

	resultMap := map[string]any{
		"hits": map[string]any{
			"total": map[string]any{
				"value": result.Hits.Total.Value,
			},
			"hits": hits,
		},
	}

	return resultMap, nil
}

// GetIndices 获取ES索引列表
func (s *LogService) GetIndices() ([]string, error) {
	client, err := s.initClient()
	if err != nil {
		return nil, err
	}

	// 获取所有索引
	// 使用ES客户端获取索引列表
	cli := client.GetCli()
	indicesResp, err := cli.Cat.Indices().Do(gaia.NewContextTrace().GetParentCtx())
	if err != nil {
		return nil, fmt.Errorf("获取索引列表失败: %w", err)
	}

	indices := make([]string, 0)
	for _, index := range indicesResp {
		if index.Index != nil {
			indices = append(indices, *index.Index)
		}
	}

	return indices, nil
}

// GetMapping 获取ES字段映射
func (s *LogService) GetMapping(index string) (map[string]any, error) {
	client, err := s.initClient()
	if err != nil {
		return nil, err
	}

	// 获取字段映射
	cli := client.GetCli()
	mappingResp, err := cli.Indices.GetMapping().Index(index).Do(gaia.NewContextTrace().GetParentCtx())
	if err != nil {
		return nil, fmt.Errorf("获取字段映射失败: %w", err)
	}

	// 转换为map
	result := make(map[string]any)
	for idx, mapping := range mappingResp {
		result[idx] = mapping.Mappings
	}

	return result, nil
}

// ExportLogs 导出日志
func (s *LogService) ExportLogs(index string, must, should, not []OpArg) ([]byte, error) {
	client, err := s.initClient()
	if err != nil {
		return nil, err
	}

	// 转换为es包中的类型
	esMust := make([]esclient.OpArg, 0, len(must))
	for _, arg := range must {
		esMust = append(esMust, esclient.OpArg{
			Key:   arg.Key,
			Value: arg.Value,
			Op:    arg.Op,
		})
	}

	esShould := make([]esclient.OpArg, 0, len(should))
	for _, arg := range should {
		esShould = append(esShould, esclient.OpArg{
			Key:   arg.Key,
			Value: arg.Value,
			Op:    arg.Op,
		})
	}

	esNot := make([]esclient.OpArg, 0, len(not))
	for _, arg := range not {
		esNot = append(esNot, esclient.OpArg{
			Key:   arg.Key,
			Value: arg.Value,
			Op:    arg.Op,
		})
	}

	// 构建查询参数，不限制返回数量
	arg := esclient.SimpleSearchArg{
		Index:  index,
		From:   0,
		Size:   10000,
		Must:   esMust,
		Should: esShould,
		Not:    esNot,
	}

	// 执行查询
	result, err := client.SimpleSearch(arg)
	if err != nil {
		return nil, fmt.Errorf("ES查询失败: %w", err)
	}

	// 构建符合前端期望的返回格式
	hits := make([]map[string]any, 0, len(result.Hits.Data))
	for _, hit := range result.Hits.Data {
		hits = append(hits, map[string]any{
			"_id":     hit.Id,
			"_score":  hit.Score,
			"_source": hit.Detail,
		})
	}

	exportResult := map[string]any{
		"hits": map[string]any{
			"total": map[string]any{
				"value": result.Hits.Total.Value,
			},
			"hits": hits,
		},
	}

	// 将结果转换为JSON格式
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(exportResult); err != nil {
		return nil, fmt.Errorf("编码结果失败: %w", err)
	}

	return buf.Bytes(), nil
}

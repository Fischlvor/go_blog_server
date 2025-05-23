package task

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"server/global"
	"server/model/elasticsearch"
	"server/service"
	"strconv"
)

// UpdateArticleViewsSyncTask 将 Redis 中的文章浏览量（增量），同步到 Elasticsearch
func UpdateArticleViewsSyncTask() error {
	// 获取redis中的缓存数据
	articleView := service.ServiceGroupApp.ArticleService.NewArticleView()

	viewsInfo := articleView.GetAllInfo()
	for id, num := range viewsInfo {
		// 无变化就跳过
		if num == 0 {
			continue
		}

		// 更新数据 之前的数据+缓存中的数据
		source := "ctx._source.views += " + strconv.Itoa(num)
		script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
		if _, err := global.ESClient.Update(elasticsearch.ArticleIndex(), id).Script(&script).Do(context.TODO()); err != nil {
			return err
		}
	}

	// 清除redis中的数据
	articleView.Clear()
	return nil
}

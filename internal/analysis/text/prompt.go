package text

import (
	"fmt"
	"strings"

	"github.com/axpz/video-publisher/internal/analysis"
)

func buildPrompt(description string, opts analysis.AnalyzeOptions) string {
	var b strings.Builder

	b.WriteString("作为平台深度思考专家，根据下面这段一句话描述，生成适合发布到视频网站的元数据，输出严格的 JSON，不要包含任何解释文字或代码块标记。")
	b.WriteString("字段包括：title, desc, tags, category。")
	if opts.Platform != "" {
		b.WriteString("平台是 ")
		b.WriteString(opts.Platform)
		b.WriteString("。我们是专业的暖宝宝OEM/ODM源头生产厂商, 质量符合全球标准")
	}
	b.WriteString("title 是简洁有吸引力的标题；")
	b.WriteString("desc 是简短描述；")
	b.WriteString("tags 是相关标签数组，生成的tags最好少于6个；")
	b.WriteString(fmt.Sprintf("category 是 %s 视频分类的数字字符串，如果不确定使用\"28\"；", opts.Platform))
	b.WriteString("使用语言: ")
	b.WriteString(opts.Language)
	b.WriteString("。")
	b.WriteString("下面是视频内容的一句话描述：")
	b.WriteString(description)
	return b.String()
}

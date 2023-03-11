package internet

func AddDraft(ctx *MsgContext) {
	ResponseText(ctx.ResponseWriter, "发布成功")
}

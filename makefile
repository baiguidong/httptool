
.PHONY:all

all:$(OBJS)
	go build  -o tools -v
	@echo "****项目tools配置构建成功****"

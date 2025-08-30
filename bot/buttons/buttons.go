package buttons

type Button struct {
	Name  string
	Value string
}

func newButton(name string, value string) Button {
	return Button{
		Name:  name,
		Value: value,
	}
}

var Return = newButton("🔙 Назад", "return")
var Cancel = newButton("❌ Отменить", "cancel")
var Info = newButton("ℹ️ Информация", "info")

var Channels = newButton("📋 Каналы", "channels")
var AddPost = newButton("🕑 Запланировать публикацию", "add_post")
var Settings = newButton("⚙️ Настройки", "settings")

var AddChannel = newButton("➕ Добавить канал", "add_channel")
var Next = newButton("➡️", "next")
var Previous = newButton("⬅️", "previous")

var UpdateChannel = newButton("🔄Обновить имя", "update_channel_name")
var RemoveChannel = newButton("Удалить канал", "remove_channel")

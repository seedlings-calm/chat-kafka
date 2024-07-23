package common

const (
	Private   = "private"
	Group     = "group"
	Broadcast = "broadcast"
)

// type ChatKafkaAny struct {
// 	Any   *anypb.Any
// 	Value interface{}
// }

// // 添加其他方法
// func (w *ChatKafkaAny) GetTypeUrl() string {
// 	return w.Any.TypeUrl
// }

// func (w *ChatKafkaAny) DecodeAnyMessage() {
// 	switch w.Any.TypeUrl {
// 	case "chat.kafka.proto.text":
// 		var msg types.TextMsg
// 		if err := w.Any.UnmarshalTo(&msg); err != nil {
// 			log.Fatalf("Failed to unmarshal : %v", err)
// 		}
// 	case "chat.kafka.proto.face":
// 		var msg types.TextMsg
// 		if err := w.Any.UnmarshalTo(&msg); err != nil {
// 			log.Fatalf("Failed to unmarshal : %v", err)
// 		}
// 	case "chat.kafka.proto.sound":
// 		var msg types.TextMsg
// 		if err := w.Any.UnmarshalTo(&msg); err != nil {
// 			log.Fatalf("Failed to unmarshal : %v", err)
// 		}
// 	default:
// 		log.Printf("Unknown TypeUrl: %s", w.Any.TypeUrl)
// 	}

// }

// func TypeUrlHeader(name string) string {
// 	return fmt.Sprintf("chat.kafka.proto.%s", name)
// }
// func CreateTextMessage(value string) (any *anypb.Any, err error) {
// 	msg := &types.TextMsg{
// 		Value: value,
// 	}
// 	any, err = anypb.New(msg)
// 	any.TypeUrl = TypeUrlHeader("text")
// 	return
// }

// func CreateFaceMessage(index int64, data string) (any *anypb.Any, err error) {
// 	msg := &types.FaceMsg{
// 		Index: index,
// 		Data:  data,
// 	}
// 	any, err = anypb.New(msg)
// 	any.TypeUrl = TypeUrlHeader("face")
// 	return
// }

// func CreateSoundMessage(url string) (any *anypb.Any, err error) {
// 	msg := &types.SoundMsg{
// 		Url:          url,
// 		Uuid:         "",
// 		Size:         0,
// 		Second:       0,
// 		DownloadFlag: 1,
// 	}
// 	any, err = anypb.New(msg)
// 	any.TypeUrl = TypeUrlHeader("sound")
// 	return
// }

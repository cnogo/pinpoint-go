package agent

import "github.com/cnogo/pinpoint-go/protocol/thrift/trace"

func NewByteAnnotation(key int32, v int8) *trace.TAnnotation {
	value := trace.NewTAnnotationValue()
	value.ByteValue = &v

	annotation := trace.NewTAnnotation()
	annotation.Key = key
	annotation.Value = value
	return annotation
}

func NewBoolAnnotation(key int32, v bool) *trace.TAnnotation {
	value := trace.NewTAnnotationValue()
	value.BoolValue = &v

	annotation := trace.NewTAnnotation()
	annotation.Key = key
	annotation.Value = value
	return annotation
}

func NewInt32Annotation(key int32, v int32) *trace.TAnnotation {
	value := trace.NewTAnnotationValue()
	value.IntValue = &v

	annotation := trace.NewTAnnotation()
	annotation.Key = key
	annotation.Value = value
	return annotation
}

func NewBinaryAnnotation(key int32, v []byte) *trace.TAnnotation {
	value := trace.NewTAnnotationValue()
	value.BinaryValue = v

	annotation := trace.NewTAnnotation()
	annotation.Key = key
	annotation.Value = value
	return annotation
}


func NewDoubleAnnotation(key int32, v float64) *trace.TAnnotation {
	value := trace.NewTAnnotationValue()
	value.DoubleValue = &v

	annotation := trace.NewTAnnotation()
	annotation.Key = key
	annotation.Value = value
	return annotation
}

func NewLongAnnotation(key int32, v int64) *trace.TAnnotation {
	value := trace.NewTAnnotationValue()
	value.LongValue = &v

	annotation := trace.NewTAnnotation()
	annotation.Key = key
	annotation.Value = value
	return annotation
}

func NewShortAnnotation(key int32, v int16) *trace.TAnnotation {
	value := trace.NewTAnnotationValue()
	value.ShortValue = &v
	annotation := trace.NewTAnnotation()
	annotation.Key = key
	annotation.Value = value
	return annotation
}

func NewStringAnnotation(key int32, v string) *trace.TAnnotation {
	value := trace.NewTAnnotationValue()
	value.StringValue = &v
	annotation := trace.NewTAnnotation()
	annotation.Key = key
	annotation.Value = value
	return annotation
}

func NewIntStringStringAnnotation(key, v int32, s1, s2 string) *trace.TAnnotation {
	value := trace.NewTAnnotationValue()
	issValue := trace.NewTIntStringStringValue()
	issValue.IntValue = v
	if s1 != "" {
		issValue.StringValue1 = &s1
	}

	if s2 != "" {
		issValue.StringValue2 = &s2
	}
	value.IntStringStringValue = issValue

	annotation := trace.NewTAnnotation()
	annotation.Key = key
	annotation.Value = value

	return annotation
}

func NewIntStringAnnotation(key, v int32, s string) *trace.TAnnotation {
	value := trace.NewTAnnotationValue()
	isValue := trace.NewTIntStringValue()
	isValue.IntValue = v
	if s != "" {
		isValue.StringValue = &s
	}

	value.IntStringValue = isValue
	annotation := trace.NewTAnnotation()
	annotation.Key = key
	annotation.Value = value

	return annotation
}

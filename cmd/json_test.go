package cmd

import (
	"reflect"
	"testing"
)

func Test_jsonValue_Get(t *testing.T) {
	type fields struct {
		kind    jsonValueKind
		array   []any
		object  map[string]any
		str     string
		num     float64
		boolean bool
	}
	tests := []struct {
		name   string
		fields fields
		want   any
	}{
		{
			name: "Null",
			fields: fields{
				kind: Null,
			},
			want: nil,
		},
		{
			name: "Boolean",
			fields: fields{
				kind:    Boolean,
				boolean: true,
			},
			want: true,
		},
		{
			name: "Number",
			fields: fields{
				kind: Number,
				num:  123.45,
			},
			want: 123.45,
		},
		{
			name: "String",
			fields: fields{
				kind: String,
				str:  "this is a test",
			},
			want: "this is a test",
		},
		{
			name: "Array",
			fields: fields{
				kind:  Array,
				array: []any{"a", "b", "c"},
			},
			want: []any{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonValue{
				kind:    tt.fields.kind,
				array:   tt.fields.array,
				object:  tt.fields.object,
				str:     tt.fields.str,
				num:     tt.fields.num,
				boolean: tt.fields.boolean,
			}
			if got := j.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jsonValue.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jsonValue_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantKind jsonValueKind
		want     any
	}{
		{
			name: "Null",
			args: args{
				data: []byte("null"),
			},
			wantKind: Null,
			want:     nil,
		},
		{
			name: "True",
			args: args{
				data: []byte("true"),
			},
			wantKind: Boolean,
			want:     true,
		},
		{
			name: "False",
			args: args{
				data: []byte("false"),
			},
			wantKind: Boolean,
			want:     false,
		},
		{
			name: "Array",
			args: args{
				data: []byte("[1,2,3]"),
			},
			wantKind: Array,
			want:     []any{1.0, 2.0, 3.0},
		},
		{
			name: "String",
			args: args{
				data: []byte(`"hello, world"`),
			},
			wantKind: String,
			want:     "hello, world",
		},
		{
			name: "Number",
			args: args{
				data: []byte(`123`),
			},
			wantKind: Number,
			want:     123.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonValue{}
			if err := j.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("jsonValue.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if j.kind != tt.wantKind {
				t.Errorf("got kind %v, want %v", j.kind, tt.wantKind)
			}
			if got := j.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jsonValue.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

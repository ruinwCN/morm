package entity

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestGetTableName(t *testing.T) {
	type Model1 struct {
		Id   int64  `json:"id" borm:"id" borm_table:"model1"`
		Name string `json:"name" borm:"name"`
	}
	type Model2 struct {
		Id   int64  `json:"id" borm:"id"`
		Name string `json:"name" borm:"name"`
	}

	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				&Model1{1, "model1"},
			},
			want:    "model1",
			wantErr: false,
		},
		{
			name: "",
			args: args{
				&Model2{1, "model2"},
			},
			want:    "model2",
			wantErr: false,
		},
		{
			name: "",
			args: args{
				Model2{1, "model2"},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTableName(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTableName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTableName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMapWithBORM(t *testing.T) {

	type Model struct {
		Id       int64         `json:"id" borm:"id" borm_table:"model1"`
		Name     string        `json:"name" borm:"name"`
		Ignore   string        `json:"ignore" borm:"ignore" borm_ignore:"true"`
		Age      int32         `json:"age" borm:"age"`
		Duration time.Duration `json:"duration" borm:"duration"`
	}
	type ModelSlice struct {
		Id    int64    `json:"id" borm:"id"`
		Name  string   `json:"name" borm:"name"`
		Slice []string `json:"slice" borm:"slice"`
	}
	type ModelMap struct {
		Id   int64             `json:"id" borm:"id"`
		Name string            `json:"name" borm:"name"`
		Map  map[string]string `json:"map" borm:"map"`
	}
	type ModelObject struct {
		Id   int64  `json:"id" borm:"id"`
		Name string `json:"name" borm:"name"`
		Obj  Model  `json:"obj" borm:"obj"`
	}
	type ModelMap2 struct {
		Id   int64            `json:"id" borm:"id"`
		Name string           `json:"name" borm:"name"`
		Map  map[string]Model `json:"map" borm:"map"`
	}

	mp := make(map[string]string)
	mp["key"] = "value"

	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "baseModel",
			args: args{
				&Model{1, "baseModel", "ig", 30, 12345678},
			},
			want: map[string]interface{}{"id": int64(1), "name": "baseModel", "age": int64(30), "duration": int64(12345678)},
		},
		{
			name: "modelSlice",
			args: args{
				&ModelSlice{2, "modelSlice", []string{"1", "2"}},
			},
			want: map[string]interface{}{"id": int64(2), "name": "modelSlice", "slice": "[\"1\",\"2\"]"},
		},
		{
			name: "modelMap",
			args: args{
				&ModelMap{1, "modelMap", mp},
			},
			want: map[string]interface{}{"id": int64(1), "name": "modelMap", "map": "{\"key\":\"value\"}"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMapWithMORM(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMapWithBORM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMapWithBORM() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeParse(t *testing.T) {
	T, err := TimeParse("2020-11-18 17:05:04")
	if err != nil {
	}
	fmt.Println(T)
}

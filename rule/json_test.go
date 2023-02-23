package rule

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func Test_valueTypeConversion(t *testing.T) {
	vt, err := valueTypeConversion("int")
	assert.NoError(t, err)
	assert.Equal(t, IIntType, vt)
	vt, err = valueTypeConversion("nnn")
	assert.Error(t, err)
}

func Test_verifyTimeType(t *testing.T) {
	err := verifyTimeType(&Field{
		ValueType: IntType,
	})
	assert.Error(t, err)
	err = verifyTimeType(&Field{ValueType: TimeString})
	assert.Error(t, err)
	err = verifyTimeType(&Field{ValueType: TimeString, TimeLayout: time.RFC3339Nano})
	assert.NoError(t, err)
}

func Test_verifySTableName(t *testing.T) {
	type args struct {
		sTableName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "start with an illegal character",
			args: args{
				sTableName: "1stb",
			},
			wantErr: assert.Error,
		},
		{
			name: "empty",
			args: args{
				sTableName: "",
			},
			wantErr: assert.Error,
		},
		{
			name: "contains illegal character",
			args: args{
				sTableName: "a-b",
			},
			wantErr: assert.Error,
		},
		{
			name: "normal",
			args: args{
				sTableName: "normal",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, verifySTableName(tt.args.sTableName), fmt.Sprintf("verifySTableName(%v)", tt.args.sTableName))
		})
	}
}

func Test_generateFieldSql(t *testing.T) {
	type args struct {
		fields []*Field
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "int",
			args: args{
				fields: []*Field{
					{
						Name:      "c1",
						ValueType: IntType,
					},
				},
			},
			want: []string{"c1 bigint"},
		},
		{
			name: "float",
			args: args{
				fields: []*Field{
					{
						Name:      "c1",
						ValueType: FloatType,
					},
				},
			},
			want: []string{"c1 double"},
		},
		{
			name: "bool",
			args: args{
				fields: []*Field{
					{
						Name:      "c1",
						ValueType: BoolType,
					},
				},
			},
			want: []string{"c1 bool"},
		},
		{
			name: "binary",
			args: args{
				fields: []*Field{
					{
						Name:      "c1",
						ValueType: StringType,
						Length:    30,
					},
				},
			},
			want: []string{"c1 binary(30)"},
		},
		{
			name: "timestamp",
			args: args{
				fields: []*Field{
					{
						Name:      "c1",
						ValueType: TimeString,
					},
				},
			},
			want: []string{"c1 timestamp"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, generateFieldSql(tt.args.fields), "generateFieldSql(%v)", tt.args.fields)
		})
	}
}

func Test_generateInsertSql(t *testing.T) {
	type args struct {
		b      *bytes.Buffer
		values []interface{}
		want   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "all",
			args: args{
				b: &bytes.Buffer{},
				values: []interface{}{
					nil,
					int64(1),
					float64(2.3),
					true,
					time.Unix(0, 0).UTC(),
					"test",
					int8(4),
				},
				want: "null,1,2.300000,true,'1970-01-01T00:00:00Z','test',4",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generateInsertSql(tt.args.b, tt.args.values)
			assert.Equal(t, tt.args.want, tt.args.b.String())
		})
	}
}

func TestEscapeString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{
				s: "a",
			},
			want: "a",
		},
		{
			name: "escape",
			args: args{"a'b"},
			want: "a\\'b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, EscapeString(tt.args.s), "EscapeString(%v)", tt.args.s)
		})
	}
}

func TestPathExist(t *testing.T) {
	f, err := os.CreateTemp("", "*")
	assert.NoError(t, err)
	exists := PathExist(f.Name())
	assert.Equal(t, true, exists)
	f.Close()
	err = os.Remove(f.Name())
	assert.NoError(t, err)
}

func Test_parseValue(t *testing.T) {
	type args struct {
		column *Column
		value  gjson.Result
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "int",
			args: args{
				column: &Column{
					ValueType: IIntType,
				},
				value: gjson.Result{
					Type: gjson.Number,
					Num:  123,
				},
			},
			want:    int64(123),
			wantErr: assert.NoError,
		},
		{
			name: "float",
			args: args{
				column: &Column{
					ValueType: IDoubleType,
				},
				value: gjson.Result{
					Type: gjson.Number,
					Num:  123.123,
				},
			},
			want:    float64(123.123),
			wantErr: assert.NoError,
		},
		{
			name: "bool",
			args: args{
				column: &Column{
					ValueType: IBoolType,
				},
				value: gjson.Result{
					Type: gjson.True,
				},
			},
			want:    true,
			wantErr: assert.NoError,
		},
		{
			name: "string",
			args: args{
				column: &Column{
					ValueType: IStringType,
				},
				value: gjson.Result{
					Type: gjson.String,
					Str:  "test",
				},
			},
			want:    "test",
			wantErr: assert.NoError,
		},
		{
			name: "time string",
			args: args{
				column: &Column{
					ValueType:  ITimeString,
					TimeLayout: time.RFC3339Nano,
				},
				value: gjson.Result{
					Type: gjson.String,
					Str:  "1970-01-01T00:00:00Z",
				},
			},
			want:    time.Unix(0, 0).UTC(),
			wantErr: assert.NoError,
		},
		{
			name: "time string wrong",
			args: args{
				column: &Column{
					ValueType:  ITimeString,
					TimeLayout: time.RFC3339Nano,
				},
				value: gjson.Result{
					Type: gjson.String,
					Str:  "xxx",
				},
			},
			want:    time.Time{},
			wantErr: assert.Error,
		},
		{
			name: "time second",
			args: args{
				column: &Column{
					ValueType: ITimeSecond,
				},
				value: gjson.Result{
					Type: gjson.Number,
					Num:  1,
				},
			},
			want:    time.Unix(1, 0),
			wantErr: assert.NoError,
		},
		{
			name: "time millisecond",
			args: args{
				column: &Column{
					ValueType: ITimeMillisecond,
				},
				value: gjson.Result{
					Type: gjson.Number,
					Num:  1,
				},
			},
			want:    time.Unix(0, 1000000),
			wantErr: assert.NoError,
		},
		{
			name: "time microsecond",
			args: args{
				column: &Column{
					ValueType: ITimeMicrosecond,
				},
				value: gjson.Result{
					Type: gjson.Number,
					Num:  1,
				},
			},
			want:    time.Unix(0, 1000),
			wantErr: assert.NoError,
		},
		{
			name: "time nanosecond",
			args: args{
				column: &Column{
					ValueType: ITimeNanoSecond,
				},
				value: gjson.Result{
					Type: gjson.Number,
					Num:  1,
				},
			},
			want:    time.Unix(0, 1),
			wantErr: assert.NoError,
		},
		{
			name: "wrong type",
			args: args{
				column: &Column{
					ValueType: 0,
				},
				value: gjson.Result{
					Type: gjson.Number,
					Num:  1,
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseValue(tt.args.column, tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("parseValue(%v, %v)", tt.args.column, tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "parseValue(%v, %v)", tt.args.column, tt.args.value)
		})
	}
}

const ExampleRule = `[
  {
    "ruleName": "test_1",
    "topic": "device_type_01",
    "rule": {
      "sTable": "device_type_01",
      "table": {
        "path": "info.id"
      },
      "tags": [
        {
          "name": "area",
          "valueType": "string",
          "length": 30,
          "path": "info.zone"
        }
      ],
      "columns": [
        {
          "name": "ts",
          "valueType": "timeString",
          "timeLayout": "2006-01-02 15:04:05",
          "path": "time"
        },
        {
          "name": "value",
          "valueType": "float",
          "path": "value"
        },
        {
          "name": "test_default",
          "valueType": "int",
          "path": "",
          "defaultValue": null
        }
      ]
    }
  }
]`

func TestNewRuleManage(t *testing.T) {
	_, err := NewRuleManage("")
	assert.Error(t, err)
	f, err := os.CreateTemp("", "*")
	assert.NoError(t, err)
	_, err = NewRuleManage(f.Name())
	assert.Error(t, err)
	f.Close()
	os.Remove(f.Name())
	f, err = os.CreateTemp("", "*")
	assert.NoError(t, err)
	_, err = f.Write([]byte(ExampleRule))
	assert.NoError(t, err)
	f.Close()
	manage, err := NewRuleManage(f.Name())
	assert.NoError(t, err)
	sqls := manage.GenerateCreateSql()
	assert.Equal(t, []string{"create stable if not exists device_type_01 (ts timestamp,value double,test_default bigint) tags(area binary(30))"}, sqls)
	exist := manage.RuleExist("device_type_01")
	assert.True(t, exist)
	exist = manage.RuleExist("xxx")
	assert.False(t, exist)
	m := manage.GetPathMap("device_type_01")
	assert.Equal(t, map[string]*Column{
		"time": {
			Name:       "ts",
			Index:      0,
			ValueType:  ITimeString,
			FieldType:  ColumnType,
			TimeLayout: "2006-01-02 15:04:05",
			Path:       "time",
		},
		"value": {
			Name:      "value",
			Index:     1,
			ValueType: IDoubleType,
			FieldType: ColumnType,
			Path:      "value",
		},
		"": {
			Name:         "test_default",
			Index:        2,
			ValueType:    IIntType,
			FieldType:    ColumnType,
			DefaultValue: nil,
			Path:         "",
		},
		"info.zone": {
			Name:      "area",
			Index:     0,
			ValueType: IStringType,
			FieldType: TagType,
			Path:      "info.zone",
		},
	}, m)
	result, err := manage.Parse("device_type_01", []byte(`{
    "info": {
        "zone": "zone1",
        "id": "device1"
    },
    "time": "2023-02-22 03:30:53",
    "value": 12
}`))
	assert.NoError(t, err)
	sql := result.ToSql()
	assert.Equal(t, "insert into _device1 using device_type_01 tags('zone1') values('2023-02-22T03:30:53Z',12.000000,null)", sql)
	result, err = manage.Parse("xxx", nil)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

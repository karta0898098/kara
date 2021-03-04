package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginationGetLimitAndOffset(t *testing.T) {
	type fields struct {
		Page       int
		PerPage    int
		TotalCount int64
		TotalPage  int64
	}
	tests := []struct {
		name   string
		fields fields
		limit  int
		offset int
	}{
		{
			name: "Success",
			fields: fields{
				Page:       1,
				PerPage:    10,
				TotalCount: 0,
				TotalPage:  0,
			},
			limit:  10,
			offset: 0,
		},
		{
			name: "Success",
			fields: fields{
				Page:       2,
				PerPage:    10,
				TotalCount: 0,
				TotalPage:  0,
			},
			limit:  10,
			offset: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pagination{
				Page:       tt.fields.Page,
				PerPage:    tt.fields.PerPage,
				TotalCount: tt.fields.TotalCount,
				TotalPage:  tt.fields.TotalPage,
			}
			limit, offset := p.GetLimitAndOffset()
			assert.Equal(t, tt.limit, limit)
			assert.Equal(t, tt.offset, offset)
		})
	}
}

func TestPaginationSetTotalCountAndPage(t *testing.T) {
	type fields struct {
		Page    int
		PerPage int
	}
	type args struct {
		total int64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		totalPage int64
	}{
		{
			name: "Success",
			fields: fields{
				Page:    1,
				PerPage: 20,
			},
			args: args{
				total: 30,
			},
			totalPage: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pagination{
				Page:    tt.fields.Page,
				PerPage: tt.fields.PerPage,
			}
			p.SetTotalCountAndPage(tt.args.total)

			assert.Equal(t, tt.totalPage, p.TotalPage)
		})
	}
}

func TestPaginationCheckOrSetDefault(t *testing.T) {
	type fields struct {
		Page       int
		PerPage    int
		TotalCount int64
		TotalPage  int64
	}
	type args struct {
		params []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Pagination
	}{
		{
			name: "Success",
			fields: fields{
				Page:    0,
				PerPage: 0,
			},
			args: args{
				params: []int{10},
			},
			want: &Pagination{
				Page:    1,
				PerPage: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pagination{
				Page:       tt.fields.Page,
				PerPage:    tt.fields.PerPage,
				TotalCount: tt.fields.TotalCount,
				TotalPage:  tt.fields.TotalPage,
			}
			paging := p.CheckOrSetDefault(tt.args.params...)
			assert.Equal(t, tt.want, paging)
		})
	}
}

func TestPaginationOffset(t *testing.T) {
	type fields struct {
		Page       int
		PerPage    int
		TotalCount int64
		TotalPage  int64
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "ReturnZeroOffset",
			fields: fields{
				Page:       1,
				PerPage:    10,
				TotalCount: 0,
				TotalPage:  0,
			},
			want: 0,
		},
		{
			name:   "ReturnTenOffset",
			fields: fields{
				Page:       2,
				PerPage:    10,
				TotalCount: 0,
				TotalPage:  0,
			},
			want:   10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pagination{
				Page:       tt.fields.Page,
				PerPage:    tt.fields.PerPage,
				TotalCount: tt.fields.TotalCount,
				TotalPage:  tt.fields.TotalPage,
			}
			offset := p.offset()
			assert.Equal(t, tt.want, offset)
		})
	}
}

package webhook_creator

import (
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/git_provider"
	"github.com/rookout/piper/pkg/utils"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestWebhookCreatorImpl_RunDiagnosis(t *testing.T) {
	type fields struct {
		clients *clients.Clients
		cfg     *conf.GlobalConfig
		hooks   map[int64]*git_provider.HookWithStatus
	}
	type args struct {
		ctx *context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WebhookCreatorImpl{
				clients: tt.fields.clients,
				cfg:     tt.fields.cfg,
				hooks:   tt.fields.hooks,
			}
			if err := wc.RunDiagnosis(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("RunDiagnosis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWebhookCreatorImpl_SetHealth(t *testing.T) {
	type fields struct {
		clients *clients.Clients
		cfg     *conf.GlobalConfig
		hooks   map[int64]*git_provider.HookWithStatus
	}
	type args struct {
		status bool
		hookID *int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "empty hook list",
			fields: fields{
				hooks: make(map[int64]*git_provider.HookWithStatus, 0),
			},
			args: args{status: false,
				hookID: utils.IPtr(123)},
			wantErr: true,
		},
		{
			name: "List with one org hook",
			fields: fields{
				hooks: map[int64]*git_provider.HookWithStatus{
					123: &git_provider.HookWithStatus{HookID: utils.IPtr(123),
						HealthStatus: false,
						RepoName:     nil},
				},
			},
			args: args{
				status: true,
				hookID: utils.IPtr(123),
			},
			wantErr: false,
		},
		{
			name: "List with one multiple repo hooks",
			fields: fields{
				hooks: map[int64]*git_provider.HookWithStatus{
					123: &git_provider.HookWithStatus{HookID: utils.IPtr(123),
						HealthStatus: false,
						RepoName:     utils.SPtr("repo1")},
					234: &git_provider.HookWithStatus{HookID: utils.IPtr(234),
						HealthStatus: false,
						RepoName:     utils.SPtr("repo2")},
				},
			},
			args: args{
				status: true,
				hookID: utils.IPtr(234),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WebhookCreatorImpl{
				hooks: tt.fields.hooks,
			}
			if err := wc.SetHealth(tt.args.status, tt.args.hookID); (err != nil) != tt.wantErr {
				t.Errorf("SetHealth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWebhookCreatorImpl_Start(t *testing.T) {
	type fields struct {
		clients *clients.Clients
		cfg     *conf.GlobalConfig
		hooks   map[int64]*git_provider.HookWithStatus
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WebhookCreatorImpl{
				clients: tt.fields.clients,
				cfg:     tt.fields.cfg,
				hooks:   tt.fields.hooks,
			}
			wc.Start()
		})
	}
}

func TestWebhookCreatorImpl_Stop(t *testing.T) {
	type fields struct {
		clients *clients.Clients
		cfg     *conf.GlobalConfig
		hooks   map[int64]*git_provider.HookWithStatus
	}
	type args struct {
		ctx *context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WebhookCreatorImpl{
				clients: tt.fields.clients,
				cfg:     tt.fields.cfg,
				hooks:   tt.fields.hooks,
			}
			wc.Stop(tt.args.ctx)
		})
	}
}

func TestWebhookCreatorImpl_checkHooksHealth(t *testing.T) {
	type fields struct {
		clients *clients.Clients
		cfg     *conf.GlobalConfig
		hooks   map[int64]*git_provider.HookWithStatus
	}
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WebhookCreatorImpl{
				clients: tt.fields.clients,
				cfg:     tt.fields.cfg,
				hooks:   tt.fields.hooks,
			}
			if got := wc.checkHooksHealth(tt.args.timeout); got != tt.want {
				t.Errorf("checkHooksHealth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebhookCreatorImpl_pingHooks(t *testing.T) {
	type fields struct {
		clients *clients.Clients
		cfg     *conf.GlobalConfig
		hooks   map[int64]*git_provider.HookWithStatus
	}
	type args struct {
		ctx *context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WebhookCreatorImpl{
				clients: tt.fields.clients,
				cfg:     tt.fields.cfg,
				hooks:   tt.fields.hooks,
			}
			wc.pingHooks(tt.args.ctx)
		})
	}
}

func TestWebhookCreatorImpl_recoverHook(t *testing.T) {
	type fields struct {
		clients *clients.Clients
		cfg     *conf.GlobalConfig
		hooks   map[int64]*git_provider.HookWithStatus
	}
	type args struct {
		ctx    *context.Context
		hookID *int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WebhookCreatorImpl{
				clients: tt.fields.clients,
				cfg:     tt.fields.cfg,
				hooks:   tt.fields.hooks,
			}
			if err := wc.recoverHook(tt.args.ctx, tt.args.hookID); (err != nil) != tt.wantErr {
				t.Errorf("recoverHook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWebhookCreatorImpl_setAllHooksHealth(t *testing.T) {
	type fields struct {
		clients *clients.Clients
		cfg     *conf.GlobalConfig
		hooks   map[int64]*git_provider.HookWithStatus
	}
	type args struct {
		status bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WebhookCreatorImpl{
				clients: tt.fields.clients,
				cfg:     tt.fields.cfg,
				hooks:   tt.fields.hooks,
			}
			wc.setAllHooksHealth(tt.args.status)
		})
	}
}

func TestWebhookCreatorImpl_setWebhooks(t *testing.T) {
	type fields struct {
		clients *clients.Clients
		cfg     *conf.GlobalConfig
		hooks   map[int64]*git_provider.HookWithStatus
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WebhookCreatorImpl{
				clients: tt.fields.clients,
				cfg:     tt.fields.cfg,
				hooks:   tt.fields.hooks,
			}
			if err := wc.setWebhooks(); (err != nil) != tt.wantErr {
				t.Errorf("setWebhooks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWebhookCreatorImpl_unsetWebhooks(t *testing.T) {
	type fields struct {
		clients *clients.Clients
		cfg     *conf.GlobalConfig
		hooks   map[int64]*git_provider.HookWithStatus
	}
	type args struct {
		ctx *context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WebhookCreatorImpl{
				clients: tt.fields.clients,
				cfg:     tt.fields.cfg,
				hooks:   tt.fields.hooks,
			}
			if err := wc.unsetWebhooks(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("unsetWebhooks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
